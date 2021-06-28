package main

import (
	"encoding/json"
	"log"
	"mime"
	"net/http"

	"github.com/gorilla/mux"

	"golang-mysql-guest-list/internal/gueststore"
)

type guestServer struct {
	store *gueststore.GuestStore
}

func NewGuestServer() *guestServer {
	store := gueststore.New()
	return &guestServer{store: store}
}

type ResponseName struct {
	Name string `json:"name"`
}

// renderJSON renders 'v' as JSON and writes it as a response into w.
func renderJSON(w http.ResponseWriter, v interface{}) {
	js, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

// addGuestListRecord creates a new table or updates existing table's capacity based on new guest list record.
// Handler for POST /guest_list/name
func (gs *guestServer) addGuestListRecord(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling guest create at %s\n", req.URL.Path)

	const MAXTABLESIZE = 10

	// Types used internally in this handler to (de-)serialize the request and
	// response from/to JSON.
	type RequestGuest struct {
		Table              int `json:"table"`
		AccompanyingGuests int `json:"accompanying_guests"`
	}

	// Enforce a JSON Content-Type.
	contentType := req.Header.Get("Content-Type")
	mediatype, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if mediatype != "application/json" {
		http.Error(w, "expect application/json Content-Type", http.StatusUnsupportedMediaType)
		return
	}

	name, _ := mux.Vars(req)["name"]
	log.Printf("received name=%s\n", name)
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	var rg RequestGuest
	if err := dec.Decode(&rg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	table, err := gs.store.GetTable(rg.Table)
	if err != nil {
		// create new, if table not found
		id := gs.store.CreateTable(rg.Table, rg.AccompanyingGuests+1)
		log.Printf("created new table with id=%d\n", id)
	} else {
		// check if new the arrangement will exceed the maximum expected table capacity
		newCapacity := table.AvailableSeats + rg.AccompanyingGuests + 1
		if newCapacity > MAXTABLESIZE {
			http.Error(w, "not enough capacity on the requested table", http.StatusBadRequest)
			return
		}
		// update available seats on the table
		id, err := gs.store.UpdateTableCapacity(table.Id, newCapacity)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("updated table with id=%d, new capacity=%d\n", id, newCapacity)
	}
	gs.store.AddGuestToTableList(name, rg.Table, rg.AccompanyingGuests)
	renderJSON(w, ResponseName{Name: name})
}

// getGuestList returns all the guests registered on the list.
// Handler for GET /guest_list/
func (gs *guestServer) getGuestList(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get guest list %s\n", req.URL.Path)

	type all struct {
		Guests []gueststore.ListedGuest `json:"guests"`
	}

	allGuests := gs.store.GetAllRegisteredGuests()
	renderJSON(w, all{Guests: allGuests})
}

func (gs *guestServer) getGuestHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get guest at %s\n", req.URL.Path)

	name, _ := mux.Vars(req)["name"]
	guest, err := gs.store.GetGuest(name)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	renderJSON(w, guest)
}

// seatArrivingGuest attempts to seat an arriving guest, after checking table availability.
// Handler for PUT /guests/{name}
func (gs *guestServer) seatArrivingGuest(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling arriving guest at %s\n", req.URL.Path)

	type RequestGuest struct {
		AccompanyingGuests int `json:"accompanying_guests"`
	}

	name, _ := mux.Vars(req)["name"]
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	var rg RequestGuest
	if err := dec.Decode(&rg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// find guest's table
	tableID, err := gs.store.GetSeatingMap(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// get toal avaible seats
	t, err := gs.store.GetTable(tableID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check if table is available
	log.Printf("table=%d has available seats=%d\n", t.Id, t.AvailableSeats-t.SeatedCount)
	if t.SeatedCount+rg.AccompanyingGuests+1 > t.AvailableSeats {
		http.Error(w, "not enough available seats left on the table", http.StatusBadRequest)
		return
	}

	// seat the guest
	n := gs.store.CreateGuest(name, tableID, rg.AccompanyingGuests)
	renderJSON(w, ResponseName{Name: n})
}

// deleteGuests removes guest and accompanying guests from the table
// Handler for DELETE /guests/name
func (gs *guestServer) deleteGuests(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling delete guests at %s\n", req.URL.Path)

	name, _ := mux.Vars(req)["name"]
	err := gs.store.DeleteGuests(name)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
}

// getGuests returns all currently seated guests.
// Handler for GET /guests
func (gs *guestServer) getGuests(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get all seated guests at %s\n", req.URL.Path)

	type all struct {
		Guests []gueststore.ArrivedGuest `json:"guests"`
	}

	allGuests := gs.store.GetAllSeatedGuests()
	renderJSON(w, all{Guests: allGuests})
}

// getGuests returns total number of empty seats.
// Handler for GET /seats_empty
func (gs *guestServer) getEmptySeats(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get empty seats at %s\n", req.URL.Path)

	type seatsEmpty struct {
		SeatsEmpty int `json:"seats_empty"`
	}

	renderJSON(w, seatsEmpty{SeatsEmpty: gs.store.GetEmptySeats()})
}

func main() {
	router := mux.NewRouter()
	router.StrictSlash(true)
	server := NewGuestServer()

	router.HandleFunc("/guest_list/{name:[[:alpha:]]+}/", server.addGuestListRecord).Methods("POST")
	router.HandleFunc("/guest_list/", server.getGuestList).Methods("GET")
	router.HandleFunc("/guests/{name:[[:alpha:]]+}/", server.seatArrivingGuest).Methods("PUT")
	router.HandleFunc("/guests/{name:[[:alpha:]]+}/", server.deleteGuests).Methods("DELETE")
	router.HandleFunc("/guests/", server.getGuests).Methods("GET")
	router.HandleFunc("/seats_empty/", server.getEmptySeats).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}

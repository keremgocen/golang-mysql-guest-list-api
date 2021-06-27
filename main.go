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
	log.Printf("handling get all guests at %s\n", req.URL.Path)

	type all struct {
		Guests []gueststore.Guest `json:"guests"`
	}

	allGuests := gs.store.GetAllGuests()
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

// Guest Arrives
// A guest may arrive with an entourage that is not the size indicated at the guest list.
// If the table is expected to have space for the extras, allow them to come. Otherwise, this method should throw an error.
// ```
// PUT /guests/name
// body:
// {
//     "accompanying_guests": int
// }
// response:
// {
//     "name": "string"
// }
// ```
func (gs *guestServer) seatArrivingGuest(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling arriving guest at %s\n", req.URL.Path)

	type RequestGuest struct {
		AccompanyingGuests int `json:"accompanyingGuests"`
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
	// tableID, err := gs.store.GetGuestTable(name)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// } else {
	// 	// check available seats and update table
	// 	id, err := gs.store.UpdateTable(tableID, rg.AccompanyingGuests+1)
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusBadRequest)
	// 	}
	// 	log.Printf("updated table with id=%d\n", id)
	// }
	// // seat the guest
	// n := gs.store.CreateGuest(name, tableID, rg.AccompanyingGuests)
	renderJSON(w, ResponseName{Name: name})
}

// func (gs *guestServer) removeGuestsHandler(w http.ResponseWriter, req *http.Request) {
// 	log.Printf("handling delete guests at %s\n", req.URL.Path)

// 	name, _ := mux.Vars(req)["name"]
// 	err := gs.store.DeleteGuests(name)

// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 	}
// }

func main() {
	router := mux.NewRouter()
	router.StrictSlash(true)
	server := NewGuestServer()

	router.HandleFunc("/guest_list/{name:[[:alpha:]]+}", server.addGuestListRecord).Methods("POST")
	router.HandleFunc("/guest_list/", server.getGuestList).Methods("GET")
	// PUT    /guests/name        :  seat arriving guest
	router.HandleFunc("/guests/{name:[[:alpha:]]+}/", server.seatArrivingGuest).Methods("PUT")

	// DELETE /guests/name        :  remove guest and accompanying guests from the table
	// router.HandleFunc("/guests/{name:[[:alpha:]]+}/", server.removeGuestsHandler).Methods("DELETE")
	// GET    /guests             :  get the guest list of arrived guests
	// GET    /seats_empty        :  get number of empty seats

	log.Fatal(http.ListenAndServe(":8080", router))
}

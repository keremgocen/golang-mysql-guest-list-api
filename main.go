// package main

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"

// 	"encoding/json"

// 	"github.com/gorilla/mux"
// )

// type Person struct {
// 	UserRef     string `json:"userRef"`
// 	Name        string `json:"name"`
// 	DOB         string `json:"dob"`
// 	PhoneNumber string `json:"phoneNumber"`
// 	Address     string `json:"address"`
// }

// // let's declare a global Persons array
// // that we can then populate in our main function
// // to simulate a database
// var Persons []Person

// func homePage(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Welcome home!")
// }

// func returnAllPersons(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("Endpoint Hit: returnAllPersons")
// 	json.NewEncoder(w).Encode(Persons)
// }

// func returnSinglePerson(w http.ResponseWriter, r *http.Request) {
// 	vars := mux.Vars(r)
// 	key := vars["userRef"]

// 	// Loop over all of our Persons
// 	// if the Person.Id equals the key we pass in
// 	// return the Person encoded as JSON
// 	for _, person := range Persons {
// 		if person.UserRef == key {
// 			json.NewEncoder(w).Encode(person)
// 		}
// 	}
// }

// func createNewPerson(w http.ResponseWriter, r *http.Request) {
// 	// get the body of our POST request
// 	// unmarshal this into a new Person struct
// 	// append this to our Persons array.
// 	reqBody, _ := ioutil.ReadAll(r.Body)
// 	var person Person
// 	json.Unmarshal(reqBody, &person)
// 	// update our global Persons array to include
// 	// our new Person
// 	Persons = append(Persons, person)

// 	json.NewEncoder(w).Encode(person)
// }

// func deletePerson(w http.ResponseWriter, r *http.Request) {
// 	// once again, we will need to parse the path parameters
// 	vars := mux.Vars(r)
// 	// we will need to extract the `id` of the person we
// 	// wish to delete
// 	id := vars["userRef"]

// 	// we then need to loop through all our persons
// 	for index, person := range Persons {
// 		// if our id path parameter matches one of our
// 		// persons
// 		if person.UserRef == id {
// 			// updates our Persons array to remove the
// 			// person
// 			Persons = append(Persons[:index], Persons[index+1:]...)
// 		}
// 	}

// }

// func handleRequests() {
// 	// creates a new instance of a mux router
// 	myRouter := mux.NewRouter().StrictSlash(true)
// 	// replace http.HandleFunc with myRouter.HandleFunc
// 	myRouter.HandleFunc("/", homePage)
// 	myRouter.HandleFunc("/all", returnAllPersons)
// 	myRouter.HandleFunc("/person/{id}", returnSinglePerson)
// 	myRouter.HandleFunc("/person", createNewPerson).Methods("POST")
// 	myRouter.HandleFunc("/person/{id}", deletePerson).Methods("DELETE")

// 	// finally, instead of passing in nil, we want
// 	// to pass in our newly created router as the second
// 	// argument
// 	log.Fatal(http.ListenAndServe(":8080", myRouter))
// }

// func main() {
// 	// client := redis.NewClient(&redis.Options{
// 	// 	Addr:     "localhost:6379",
// 	// 	Password: "",
// 	// 	DB:       0,
// 	// })

// 	_, err := json.Marshal(Person{UserRef: "KFG-734", Name: "Kerem", DOB: "18.06.1985", PhoneNumber: "+901111111111", Address: "Woodbine Close TW2 something"})
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// 	// err = client.Set("KFG-734", json, 0).Err()
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// }
// 	// val, err := client.Get("KFG-734").Result()
// 	// if err != nil {
// 	// 	fmt.Println(err)
// 	// }
// 	// fmt.Println(val)

// 	Persons = []Person{
// 		{UserRef: "KFG-734", Name: "Kerem", DOB: "18.06.1985", PhoneNumber: "+441111111111", Address: "Woodbine Close TW2"},
// 		{UserRef: "XRT-251", Name: "John", DOB: "11.11.1111", PhoneNumber: "+442222222222", Address: "Somewhere in London"},
// 	}

// 	handleRequests()
// }
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

// 	```
// POST /guest_list/name
// body:
// {
//     "table": int,
//     "accompanying_guests": int
// }
// response:
// {
//     "name": "string"
// }
// ```
func (gs *guestServer) createGuestRecordHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling guest create at %s\n", req.URL.Path)

	const TABLESIZE = 10

	// Types used internally in this handler to (de-)serialize the request and
	// response from/to JSON.
	type RequestGuest struct {
		Table              int `json:"table"`
		AccompanyingGuests int `json:"accompanyingGuests"`
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
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields()
	var rg RequestGuest
	if err := dec.Decode(&rg); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	table, err := gs.store.GetTable(rg.Table)
	if err != nil {
		// create new table if not found
		id := gs.store.CreateTable(rg.Table, rg.AccompanyingGuests+1, TABLESIZE)
		log.Printf("created new table with id %d\n", id)
	} else {
		// update table
		id, err := gs.store.UpdateTable(table.Id, rg.AccompanyingGuests+1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		log.Printf("updated table with id %d\n", id)
	}
	gs.store.SetGuestTable(name, table.Id)
	renderJSON(w, ResponseName{Name: name})
}

func (gs *guestServer) getAllGuestsHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling get all guests at %s\n", req.URL.Path)

	allGuests := gs.store.GetAllGuests()
	renderJSON(w, allGuests)
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
func (gs *guestServer) arrivingGuestsHandler(w http.ResponseWriter, req *http.Request) {
	log.Printf("handling arriving guest at %s\n", req.URL.Path)

	// Types used internally in this handler to (de-)serialize the request and
	// response from/to JSON.
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
	tableID, err := gs.store.GetGuestTable(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		// check available seats and update table
		id, err := gs.store.UpdateTable(tableID, rg.AccompanyingGuests+1)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		log.Printf("updated table with id %d\n", id)
	}
	// seat the guest
	n := gs.store.CreateGuest(name, tableID, rg.AccompanyingGuests)
	renderJSON(w, ResponseName{Name: n})
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

	// POST   /guest_list/name    :  add a guest to the guestlist
	router.HandleFunc("/guest_list/{name:[[:alpha:]]+}/", server.createGuestRecordHandler).Methods("POST")
	// GET    /guest_list         :  get the guest list
	router.HandleFunc("/guest_list", server.getAllGuestsHandler).Methods("GET")
	// PUT    /guests/name        :  seat arriving guest
	router.HandleFunc("/guests/{name:[:alpha:]+}/", server.arrivingGuestsHandler).Methods("PUT")

	// DELETE /guests/name        :  remove guest and accompanying guests from the table
	// router.HandleFunc("/guests/{name:[[:alpha:]]+}/", server.removeGuestsHandler).Methods("DELETE")
	// GET    /guests             :  get the guest list of arrived guests
	// GET    /seats_empty        :  get number of empty seats

	log.Fatal(http.ListenAndServe(":8080", router))
}

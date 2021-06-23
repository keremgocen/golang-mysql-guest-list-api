// Package classification Account API.
//
// this is to show how to write RESTful APIs in golang.
// that is to provide a detailed overview of the language specs
//
// Terms Of Service:
//
//     Schemes: http, https
//     Host: localhost:8080
//     Version: 1.0.0
//     Contact: Kerem Gocen<my@email.com>
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Security:
//     - api_key:
//
//     SecurityDefinitions:
//     api_key:
//          type: apiKey
//          name: KEY
//          in: header
//
// swagger:meta
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"encoding/json"

	"github.com/gorilla/mux"
)

type Person struct {
	UserRef     string `json:"userRef"`
	Name        string `json:"name"`
	DOB         string `json:"dob"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address"`
}

// let's declare a global Persons array
// that we can then populate in our main function
// to simulate a database
var Persons []Person

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func returnAllPersons(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllPersons")
	json.NewEncoder(w).Encode(Persons)
}

func returnSinglePerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["userRef"]

	// Loop over all of our Persons
	// if the Person.Id equals the key we pass in
	// return the Person encoded as JSON
	for _, person := range Persons {
		if person.UserRef == key {
			json.NewEncoder(w).Encode(person)
		}
	}
}

func createNewPerson(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// unmarshal this into a new Person struct
	// append this to our Persons array.
	reqBody, _ := ioutil.ReadAll(r.Body)
	var person Person
	json.Unmarshal(reqBody, &person)
	// update our global Persons array to include
	// our new Person
	Persons = append(Persons, person)

	json.NewEncoder(w).Encode(person)
}

func deletePerson(w http.ResponseWriter, r *http.Request) {
	// once again, we will need to parse the path parameters
	vars := mux.Vars(r)
	// we will need to extract the `id` of the person we
	// wish to delete
	id := vars["userRef"]

	// we then need to loop through all our persons
	for index, person := range Persons {
		// if our id path parameter matches one of our
		// persons
		if person.UserRef == id {
			// updates our Persons array to remove the
			// person
			Persons = append(Persons[:index], Persons[index+1:]...)
		}
	}

}

func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	// replace http.HandleFunc with myRouter.HandleFunc
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/all", returnAllPersons)
	myRouter.HandleFunc("/person/{id}", returnSinglePerson)
	myRouter.HandleFunc("/person", createNewPerson).Methods("POST")
	myRouter.HandleFunc("/person/{id}", deletePerson).Methods("DELETE")

	// finally, instead of passing in nil, we want
	// to pass in our newly created router as the second
	// argument
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	// client := redis.NewClient(&redis.Options{
	// 	Addr:     "localhost:6379",
	// 	Password: "",
	// 	DB:       0,
	// })

	_, err := json.Marshal(Person{UserRef: "KFG-734", Name: "Kerem", DOB: "18.06.1985", PhoneNumber: "+901111111111", Address: "Woodbine Close TW2 something"})
	if err != nil {
		fmt.Println(err)
	}

	// err = client.Set("KFG-734", json, 0).Err()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// val, err := client.Get("KFG-734").Result()
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println(val)

	Persons = []Person{
		{UserRef: "KFG-734", Name: "Kerem", DOB: "18.06.1985", PhoneNumber: "+441111111111", Address: "Woodbine Close TW2"},
		{UserRef: "XRT-251", Name: "John", DOB: "11.11.1111", PhoneNumber: "+442222222222", Address: "Somewhere in London"},
	}

	handleRequests()
}

// package gueststore provides a simple in-memory "data store" for guests.
// Guests are uniquely identified by their names.
//
// Created by Kerem Gocen, inspired by the (taskstore implementation)[https://github.com/eliben/code-for-blog/blob/master/2021/go-rest-servers/gorilla/internal/taskstore/taskstore.go] by Eli Bendersky [https://eli.thegreenplace.net]
package gueststore

import (
	"fmt"
	"log"
	"sync"
)

type Guest struct {
	Table              int    `json:"table"`
	Name               string `json:"name"`
	AccompanyingGuests int    `json:"accompanying_guests"`
}

type Table struct {
	Id             int `json:"id"`
	AvailableSeats int `json:"available_seats"`
}

type GuestList struct {
	Guests []Guest `json:"guests"`
}

// GuestStore is a simple in-memory database of guests; GuestStore methods are
// safe to call concurrently.
type GuestStore struct {
	sync.Mutex
	guests    map[string]Guest
	tables    map[int]Table
	guestList map[int]GuestList
}

func New() *GuestStore {
	gs := &GuestStore{}
	gs.guests = make(map[string]Guest)
	gs.tables = make(map[int]Table)
	gs.guestList = make(map[int]GuestList)
	return gs
}

// AddGuestToTableList assigns a guest to a table for seating lookup.
func (gs *GuestStore) AddGuestToTableList(name string, id int, accompanying int) {
	gs.Lock()
	defer gs.Unlock()

	gl := gs.guestList[id]
	log.Printf("obtained guest list for table=%d, %#v\n", id, gl.Guests)
	newGuest := Guest{
		Table:              id,
		Name:               name,
		AccompanyingGuests: accompanying,
	}
	gl.Guests = append(gl.Guests, newGuest)
	gs.guestList[id] = gl
	log.Printf("updated guest list for table=%d, %#v\n", id, gl.Guests)
}

// GetTable retrieves a table from the store, by id. If no such id exists, an
// error is returned.
func (gs *GuestStore) GetTable(id int) (Table, error) {
	gs.Lock()
	defer gs.Unlock()

	t, ok := gs.tables[id]
	if ok {
		return t, nil
	} else {
		return Table{}, fmt.Errorf("table with id=%d not found", id)
	}
}

// UpdateTableCapacity updates an existing table's available seat count in the store.
func (gs *GuestStore) UpdateTableCapacity(id int, newGuests int) (int, error) {
	gs.Lock()
	defer gs.Unlock()

	t, ok := gs.tables[id]
	if ok {
		newTable := Table{
			Id:             id,
			AvailableSeats: t.AvailableSeats + newGuests,
		}
		gs.tables[id] = newTable
		return newTable.Id, nil
	} else {
		return -1, fmt.Errorf("table with id=%d not found", id)
	}
}

// CreateTable creates a new table in the store.
func (gs *GuestStore) CreateTable(id int, capacity int) int {
	gs.Lock()
	defer gs.Unlock()

	table := Table{
		Id:             id,
		AvailableSeats: capacity,
	}

	log.Printf("created table %#v", table)

	gs.tables[id] = table
	return table.Id
}

// CreateGuest creates a new guest in the store.
func (gs *GuestStore) CreateGuest(name string, table int, accompanyingGuests int) string {
	gs.Lock()
	defer gs.Unlock()

	guest := Guest{
		Name:               name,
		Table:              table,
		AccompanyingGuests: accompanyingGuests,
	}

	gs.guests[name] = guest
	return guest.Name
}

// GetGuest retrieves a guest from the store, by name. If no such name exists, an
// error is returned.
func (gs *GuestStore) GetGuest(name string) (Guest, error) {
	gs.Lock()
	defer gs.Unlock()

	g, ok := gs.guests[name]
	if ok {
		return g, nil
	} else {
		return Guest{}, fmt.Errorf("guest with name=%s not found", name)
	}
}

// // DeleteTask deletes the task with the given id. If no such id exists, an error
// // is returned.
// func (ts *TaskStore) DeleteTask(id int) error {
// 	ts.Lock()
// 	defer ts.Unlock()

// 	if _, ok := ts.tasks[id]; !ok {
// 		return fmt.Errorf("task with id=%d not found", id)
// 	}

// 	delete(ts.tasks, id)
// 	return nil
// }

// // DeleteAllTasks deletes all tasks in the store.
// func (ts *TaskStore) DeleteAllTasks() error {
// 	ts.Lock()
// 	defer ts.Unlock()

// 	ts.tasks = make(map[int]Task)
// 	return nil
// }

// GetAllGuests returns all the guests in the store, in arbitrary order.
func (gs *GuestStore) GetAllGuests() []Guest {
	gs.Lock()
	defer gs.Unlock()

	allGuests := make([]Guest, 0)
	for _, guestList := range gs.guestList {
		allGuests = append(allGuests, guestList.Guests...)
	}
	return allGuests
}

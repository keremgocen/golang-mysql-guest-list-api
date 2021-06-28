// package gueststore provides a simple in-memory "data store" for guests.
// Guests are uniquely identified by their names.
//
// Created by Kerem Gocen, inspired by the (taskstore implementation)[https://github.com/eliben/code-for-blog/blob/master/2021/go-rest-servers/gorilla/internal/taskstore/taskstore.go] by Eli Bendersky [https://eli.thegreenplace.net]
package gueststore

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type ListedGuest struct {
	Table              int    `json:"table"`
	Name               string `json:"name"`
	AccompanyingGuests int    `json:"accompanying_guests"`
}

type ArrivedGuest struct {
	Name               string    `json:"name"`
	AccompanyingGuests int       `json:"accompanying_guests"`
	TimeArrived        time.Time `json:"time_arrived"`
}

type Table struct {
	Id             int `json:"id"`
	AvailableSeats int `json:"available_seats"`
	SeatedCount    int `json:"seated_count"`
}

type GuestList struct {
	Guests []ListedGuest `json:"guests"`
}

// GuestStore is a simple in-memory database of guests; GuestStore methods are
// safe to call concurrently.
type GuestStore struct {
	sync.Mutex
	guests     map[string]ArrivedGuest
	tables     map[int]Table
	guestList  map[int]GuestList
	seatingMap map[string]int
}

func New() *GuestStore {
	gs := &GuestStore{}
	gs.guests = make(map[string]ArrivedGuest)
	gs.tables = make(map[int]Table)
	gs.guestList = make(map[int]GuestList)
	gs.seatingMap = make(map[string]int)
	return gs
}

// AddGuestToTableList assigns a guest to a table for seating lookup.
func (gs *GuestStore) AddGuestToTableList(name string, id int, accompanying int) {
	gs.Lock()
	defer gs.Unlock()

	gl := gs.guestList[id]
	log.Printf("obtained guest list for table=%d, %#v\n", id, gl.Guests)
	newGuest := ListedGuest{
		Table:              id,
		Name:               name,
		AccompanyingGuests: accompanying,
	}
	gl.Guests = append(gl.Guests, newGuest)
	gs.guestList[id] = gl

	gs.seatingMap[name] = id
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
func (gs *GuestStore) UpdateTableCapacity(id int, newCapacity int) (int, error) {
	gs.Lock()
	defer gs.Unlock()

	_, ok := gs.tables[id]
	if ok {
		newTable := Table{
			Id:             id,
			AvailableSeats: newCapacity,
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
		SeatedCount:    0,
	}

	log.Printf("created table %#v", table)

	gs.tables[id] = table
	return table.Id
}

// CreateGuest creates a new guest in the store.
func (gs *GuestStore) CreateGuest(name string, id int, accompanyingGuests int) string {
	gs.Lock()
	defer gs.Unlock()

	guest := ArrivedGuest{
		Name:               name,
		AccompanyingGuests: accompanyingGuests,
		TimeArrived:        time.Now().UTC(),
	}

	// updated seated count at the table
	t := gs.tables[id]
	newTable := Table{
		Id:             id,
		AvailableSeats: t.AvailableSeats,
		SeatedCount:    t.SeatedCount + accompanyingGuests + 1,
	}
	gs.tables[id] = newTable
	log.Printf("CreateGuest seated count at table=%d is updated as %d, from %d", id, newTable.SeatedCount, t.SeatedCount)

	gs.guests[name] = guest
	return guest.Name
}

// GetSeatingMap returns table number for the given guest name.
func (gs *GuestStore) GetSeatingMap(name string) (int, error) {
	gs.Lock()
	defer gs.Unlock()

	s, ok := gs.seatingMap[name]
	if ok {
		return s, nil
	} else {
		return -1, fmt.Errorf("seating map for guest=%s not found", name)
	}
}

// GetGuest retrieves a guest from the store, by name. If no such name exists, an
// error is returned.
func (gs *GuestStore) GetGuest(name string) (ArrivedGuest, error) {
	gs.Lock()
	defer gs.Unlock()

	g, ok := gs.guests[name]
	if ok {
		return g, nil
	} else {
		return ArrivedGuest{}, fmt.Errorf("guest with name=%s not found", name)
	}
}

// DeleteGuests removes the guest with the given name. If no such name exists, an error
// is returned. Guest's assigned table seated count is also updated.
func (gs *GuestStore) DeleteGuests(name string) error {
	gs.Lock()
	defer gs.Unlock()

	if _, ok := gs.guests[name]; !ok {
		return fmt.Errorf("guest with name=%s not found", name)
	}

	extraSeats := gs.guests[name].AccompanyingGuests + 1
	delete(gs.guests, name)

	tID := gs.seatingMap[name]
	if _, ok := gs.tables[tID]; !ok {
		return fmt.Errorf("table with id=%d not found", tID)
	}
	t := gs.tables[tID]
	newTable := Table{
		Id:             tID,
		AvailableSeats: t.AvailableSeats,
		SeatedCount:    t.SeatedCount - extraSeats,
	}
	gs.tables[tID] = newTable
	log.Printf("DeleteGuests seated count at table=%#v is updated as %d, previous was %d", newTable, newTable.SeatedCount, t.SeatedCount)

	return nil
}

// GetAllRegisteredGuests returns all the guests in the store, in arbitrary order.
func (gs *GuestStore) GetAllRegisteredGuests() []ListedGuest {
	gs.Lock()
	defer gs.Unlock()

	allGuests := make([]ListedGuest, 0)
	for _, guestList := range gs.guestList {
		allGuests = append(allGuests, guestList.Guests...)
	}
	return allGuests
}

// GetAllSeatedGuests returns all the seated guests in the store, in arbitrary order.
func (gs *GuestStore) GetAllSeatedGuests() []ArrivedGuest {
	gs.Lock()
	defer gs.Unlock()

	allGuests := make([]ArrivedGuest, 0)
	for _, guest := range gs.guests {
		allGuests = append(allGuests, guest)
	}
	return allGuests
}

// GetEmptySeats returns total number of empty seats.
func (gs *GuestStore) GetEmptySeats() int {
	gs.Lock()
	defer gs.Unlock()

	count := 0
	for _, t := range gs.tables {
		count += (t.AvailableSeats - t.SeatedCount)
		log.Printf("GetEmptySeats table=%#v has available seats=%d, total count=%d\n", t, t.AvailableSeats-t.SeatedCount, count)
	}
	return count
}

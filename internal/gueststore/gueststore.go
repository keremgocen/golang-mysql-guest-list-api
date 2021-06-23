// package gueststore provides a simple in-memory "data store" for guests.
// Guests are uniquely identified by their names.
//
// Created by Kerem Gocen, inspired by the (taskstore implementation)[https://github.com/eliben/code-for-blog/blob/master/2021/go-rest-servers/gorilla/internal/taskstore/taskstore.go] by Eli Bendersky [https://eli.thegreenplace.net]
package gueststore

import (
	"fmt"
	"sync"
)

// "table": int,
// "accompanying_guests": int
type Guest struct {
	Table              int    `json:"table"`
	Name               string `json:"name"`
	AccompanyingGuests int    `json:"accompanyingGuests"`
}

// GuestStore is a simple in-memory database of guests; GuestStore methods are
// safe to call concurrently.
type GuestStore struct {
	sync.Mutex
	guests map[string]Guest
	name   string
}

func New() *GuestStore {
	gs := &GuestStore{}
	gs.guests = make(map[string]Guest)
	return gs
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

	allGuests := make([]Guest, 0, len(gs.guests))
	for _, guest := range gs.guests {
		allGuests = append(allGuests, guest)
	}
	return allGuests
}

// // GetTasksByTag returns all the tasks that have the given tag, in arbitrary
// // order.
// func (ts *TaskStore) GetTasksByTag(tag string) []Task {
// 	ts.Lock()
// 	defer ts.Unlock()

// 	var tasks []Task

// taskloop:
// 	for _, task := range ts.tasks {
// 		for _, taskTag := range task.Tags {
// 			if taskTag == tag {
// 				tasks = append(tasks, task)
// 				continue taskloop
// 			}
// 		}
// 	}
// 	return tasks
// }

// // GetTasksByDueDate returns all the tasks that have the given due date, in
// // arbitrary order.
// func (ts *TaskStore) GetTasksByDueDate(year int, month time.Month, day int) []Task {
// 	ts.Lock()
// 	defer ts.Unlock()

// 	var tasks []Task

// 	for _, task := range ts.tasks {
// 		y, m, d := task.Due.Date()
// 		if y == year && m == month && d == day {
// 			tasks = append(tasks, task)
// 		}
// 	}

// 	return tasks
// }

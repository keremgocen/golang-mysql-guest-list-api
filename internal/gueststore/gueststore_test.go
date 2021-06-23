package gueststore

import (
	"testing"
)

func TestCreateAndGet(t *testing.T) {
	// Create a store and a single guest.
	gs := New()
	name := gs.CreateGuest("John", 1, 3)

	// We should be able to retrieve this guest by name, but nothing with other names.
	guest, err := gs.GetGuest(name)
	if err != nil {
		t.Fatal(err)
	}

	if guest.Name != name {
		t.Errorf("got guest.Name=%v, name=%v", guest.Name, name)
	}
	if guest.Table != 1 {
		t.Errorf("got Table=%d, want %d", guest.Table, 1)
	}

	// Asking for all guests, we only get the one we put in.
	allGuests := gs.GetAllGuests()
	if len(allGuests) != 1 || allGuests[0].Name != name {
		t.Errorf("got len(allGuests)=%d, allGuests[0].Name=%v; want 1, %v", len(allGuests), allGuests[0].Name, name)
	}

	_, err = gs.GetGuest("Noone")
	if err == nil {
		t.Fatal("got nil, want error")
	}

	// Add another guest. Expect to find two guests in the store.
	gs.CreateGuest("Joe", 2, 2)
	allGuests2 := gs.GetAllGuests()
	if len(allGuests2) != 2 {
		t.Errorf("got len(allGuests2)=%d; want 2", len(allGuests2))
	}
}

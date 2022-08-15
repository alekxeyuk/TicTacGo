package main

import (
	"github.com/dustin/go-broadcast"
	"github.com/google/uuid"
)

// Creating a map of room IDs to broadcasters.
var roomChannels = make(map[uuid.UUID]broadcast.Broadcaster)
var roomCounter uint64 = 0

// It creates a channel, registers it with the room, and returns it.
func openListener(roomid uuid.UUID) chan interface{} {
	listener := make(chan interface{})
	getRoom(roomid).Register(listener)
	return listener
}

// It closes the listener channel and unregisters it from the room.
func closeListener(roomid uuid.UUID, listener chan interface{}) {
	getRoom(roomid).Unregister(listener)
	close(listener)
}

// It creates a new room, assigns it a unique ID, creates a new broadcaster for that room, and then
// adds the broadcaster to the roomChannels map.
func newRoom() (roomid uuid.UUID, b broadcast.Broadcaster) {
	roomid = uuid.New()
	b = broadcast.NewBroadcaster(10)
	roomChannels[roomid] = b
	roomCounter++
	return
}

func deleteRoom(roomid uuid.UUID) bool {
	r, ok := roomChannels[roomid]
	if ok {
		r.Close()
		delete(roomChannels, roomid)
		roomCounter--
	}
	return ok
}

// Returns a room channel for the given room id. If the room does not exist,
// returns nil.
func getRoom(roomid uuid.UUID) broadcast.Broadcaster {
	return roomChannels[roomid]
}

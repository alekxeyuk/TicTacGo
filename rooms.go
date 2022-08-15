package main

import (
	"github.com/dustin/go-broadcast"
	"github.com/google/uuid"
)

// Creating a map of room IDs to broadcasters.
var roomChannels = make(map[string]broadcast.Broadcaster)
var roomCounter uint64 = 0
var globalRoomID, _ = newRoom()

// It creates a channel, registers it with the room, and returns it.
func openListener(roomid string) chan interface{} {
	listener := make(chan interface{})
	getRoom(roomid).Register(listener)
	return listener
}

// It closes the listener channel and unregisters it from the room.
func closeListener(roomid string, listener chan interface{}) {
	if room := getRoom(roomid); room != nil {
		room.Unregister(listener)
	}
	close(listener)
}

// It creates a new room, assigns it a unique ID, creates a new broadcaster for that room, and then
// adds the broadcaster to the roomChannels map.
func newRoom() (roomid string, b broadcast.Broadcaster) {
	roomid = uuid.NewString()
	b = broadcast.NewBroadcaster(10)
	roomChannels[roomid] = b
	roomCounter++
	return
}

func deleteRoom(roomid string) bool {
	r, ok := roomChannels[roomid]
	if ok {
		r.Submit("stop")
		r.Close()
		delete(roomChannels, roomid)
		roomCounter--
	}
	return ok
}

// Returns a room channel for the given room id. If the room does not exist,
// returns nil.
func getRoom(roomid string) broadcast.Broadcaster {
	return roomChannels[roomid]
}

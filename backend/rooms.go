package main

import (
	"sync"

	"github.com/dustin/go-broadcast"
	"github.com/google/uuid"
)

type Message struct {
	Type string
	Body string
}

type RoomState byte
type PlayerSign byte

const (
	EMPTY_CELL PlayerSign = iota
	PLAYER_X
	PLAYER_O
)

const (
	EMPTY_ROOM RoomState = iota
	GAME_IN_PROGRESS
	GAME_FINISHED
)

//go:generate stringer -type=RoomState,PlayerSign

type Room struct {
	id            string
	b             broadcast.Broadcaster
	board         [9]PlayerSign
	boardLock     sync.RWMutex
	currentPlayer PlayerSign
	users         [2]*User
	state         RoomState
}

// Creating a map of room IDs to broadcasters.
var rooms = make(map[string]*Room)
var roomsCounter uint64 = 0
var globalRoom, _ = newRoom()

// It creates a channel, registers it with the room, and returns it.
func openListener(roomid string) chan interface{} {
	listener := make(chan interface{})
	getRoom(roomid).b.Register(listener)
	return listener
}

// It closes the listener channel and unregisters it from the room.
func closeListener(roomid string, listener chan interface{}) {
	if room := getRoom(roomid); room.id != "" {
		room.b.Unregister(listener)
	}
	close(listener)
}

func joinableRooms() []*Room {
	rs := make([]*Room, 0)
	for _, room := range rooms {
		if room.id != globalRoom.id && (room.users[0] == nil || room.users[1] == nil) {
			rs = append(rs, room)
		}
	}
	return rs
}

func (r *Room) addUser(user *User) {
	r.boardLock.Lock()
	defer r.boardLock.Unlock()
	if r.users[0] == nil {
		r.users[0] = user
	} else {
		r.users[1] = user
	}
}

// It creates a new room, assigns it a unique ID, creates a new broadcaster for that room, and then
// adds the broadcaster to the roomChannels map.
func newRoom() (r *Room, b broadcast.Broadcaster) {
	roomid := uuid.NewString()
	b = broadcast.NewBroadcaster(10)
	r = &Room{roomid, b, [9]PlayerSign{}, sync.RWMutex{}, PLAYER_X, [2]*User{}, EMPTY_ROOM}
	rooms[roomid] = r
	roomsCounter++
	return
}

func deleteRoom(roomid string) bool {
	r, ok := rooms[roomid]
	if ok {
		r.b.Submit(Message{"stop", ""})
		r.b.Close()
		delete(rooms, roomid)
		roomsCounter--
	}
	return ok
}

// Returns a room channel for the given room id. If the room does not exist,
// returns nil.
func getRoom(roomid string) *Room {
	return rooms[roomid]
}

func (r *Room) getBoard() [9]PlayerSign {
	r.boardLock.RLock()
	defer r.boardLock.RUnlock()
	return r.board
}

func (r *Room) setBoard(board [9]PlayerSign) {
	r.boardLock.Lock()
	defer r.boardLock.Unlock()
	r.board = board
}

func (r *Room) move(index int64) {
	r.boardLock.Lock()
	defer r.boardLock.Unlock()
	r.board[index] = r.currentPlayer
}

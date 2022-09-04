package main

import (
	"sync"
	"time"

	"github.com/dustin/go-broadcast"
	"github.com/google/uuid"
)

type Message struct {
	Type string
	Body interface{}
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
	currentPlayer *User
	users         [2]*User
	state         RoomState
}

var winCheckLines = [8][3]int64{
	{0, 1, 2},
	{3, 4, 5},
	{6, 7, 8},
	{0, 3, 6},
	{1, 4, 7},
	{2, 5, 8},
	{0, 4, 8},
	{2, 4, 6},
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
	if room := getRoom(roomid); room != nil {
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
		user.sign = PLAYER_X
		r.currentPlayer = user
		r.users[0] = user
	} else {
		user.sign = PLAYER_O
		r.users[1] = user
	}
}

// It creates a new room, assigns it a unique ID, creates a new broadcaster for that room, and then
// adds the broadcaster to the roomChannels map.
func newRoom() (r *Room, b broadcast.Broadcaster) {
	roomid := uuid.NewString()
	b = broadcast.NewBroadcaster(10)
	r = &Room{roomid, b, [9]PlayerSign{}, sync.RWMutex{}, nil, [2]*User{}, EMPTY_ROOM}
	rooms[roomid] = r
	roomsCounter++
	return
}

func deleteRoom(roomid string, mess interface{}) bool {
	r, ok := rooms[roomid]
	if ok {
		r.b.Submit(mess)
		go func() {
			time.Sleep(1 * time.Second)
			r.b.Close()
			delete(rooms, roomid)
			roomsCounter--
		}()
	}
	return ok
}

// Returns a room channel for the given room id. If the room does not exist,
// returns nil.
func getRoom(roomid string) *Room {
	return rooms[roomid]
}

func (r *Room) move(index int64) {
	r.boardLock.Lock()
	defer r.boardLock.Unlock()
	r.board[index] = r.currentPlayer.sign
	r.switchPlayer()
}

func (r *Room) switchPlayer() {
	if r.currentPlayer.sign == PLAYER_X {
		r.currentPlayer = r.users[1]
	} else {
		r.currentPlayer = r.users[0]
	}
}

func (r *Room) cellIsEmpty(index int64) bool {
	r.boardLock.RLock()
	defer r.boardLock.RUnlock()
	return r.board[index] == EMPTY_CELL
}

func (r *Room) checkWin() bool {
	r.boardLock.Lock()
	defer r.boardLock.Unlock()
	for _, line := range winCheckLines {
		a, b, c := line[0], line[1], line[2]
		if r.board[a] == r.board[b] && r.board[b] == r.board[c] && r.board[a] != EMPTY_CELL {
			return true
		}
	}
	return false
}

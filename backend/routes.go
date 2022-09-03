package main

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func index(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Nothing to see here",
	})
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
	globalRoom.b.Submit(Message{"ping", "pong"})
}

func sendCountEvent() {
	globalRoom.b.Submit(Message{"count", strconv.FormatUint(roomsCounter, 10)})
}

func roomMOVE(c *gin.Context) {
	roomid := c.Param("roomid")
	room := getRoom(roomid)
	ok, _ := authorized(c)
	if roomid == "global" || room == nil || !ok {
		c.Status(http.StatusBadRequest)
		return
	}

	cell, err := strconv.ParseInt(c.Query("cell"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	room.move(cell)
	room.b.Submit(Message{"move", c.Query("cell")})
}

func roomJoinOrCreate(c *gin.Context, u *User) {
	j := joinableRooms()
	var room *Room
	if len(j) == 0 {
		room, _ = newRoom()
	} else {
		room = j[0]
	}
	room.addUser(u)
	u.roomId = room.id
	c.JSON(http.StatusOK, gin.H{
		"uuid": room.id,
	})
}

func roomRANDOM(c *gin.Context) {
	ok, user := authorized(c)
	if !ok {
		c.Status(http.StatusUnauthorized)
		return
	}

	u := getUser(user)
	if u.roomId == "" {
		roomJoinOrCreate(c, u)
	} else {
		r := getRoom(u.roomId)
		if r != nil {
			c.JSON(http.StatusOK, gin.H{
				"uuid": r.id,
			})
		} else {
			roomJoinOrCreate(c, u)
		}
	}
	go sendCountEvent()
}

func roomCOUNT(c *gin.Context) {
	ok, user := authorized(c)
	c.JSON(http.StatusOK, gin.H{
		"room_count": roomsCounter,
		"rooms":      rooms,
		"user":       user,
		"user_count": userCounter,
		"users":      userMap,
		"ok":         ok,
	})
}

func roomLIST(c *gin.Context) {
	// return a list of all rooms as JSON
	roomsJSON := make([]gin.H, 0)
	for _, room := range rooms {
		roomsJSON = append(roomsJSON, gin.H{
			"uuid":           room.id,
			"board":          room.board,
			"current_player": room.currentPlayer.String(),
			"players":        room.users,
			"state":          room.state.String(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"rooms": roomsJSON,
	})
}

func roomDELETE(c *gin.Context) {
	roomid := c.Param("roomid")
	if roomid == globalRoom.id {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "cannot delete global room",
		})
		return
	}

	if ok := deleteRoom(roomid); ok {
		c.JSON(http.StatusOK, gin.H{
			"message": "deleted",
		})
		go sendCountEvent()
	} else {
		c.JSON(http.StatusGone, gin.H{
			"message": "not found",
		})
	}
}

func roomSTREAM(c *gin.Context) {
	roomid := c.Param("roomid")
	if roomid == "global" {
		roomid = globalRoom.id
	} else if getRoom(roomid) == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "not found",
		})
		return
	}

	listener := openListener(roomid)
	ticker := time.NewTicker(10 * time.Second)
	defer func() {
		closeListener(roomid, listener)
		ticker.Stop()
	}()

	c.Stream(func(w io.Writer) bool {
		select {
		case msg := <-listener:
			switch m := msg.(type) {
			case Message:
				c.SSEvent(m.Type, m.Body)
				if m.Type == "stop" {
					return false
				}
			default:
				c.SSEvent("message", msg)
			}
		case <-ticker.C:
			c.SSEvent("time", time.Now().String())
		}
		return true
	})
}

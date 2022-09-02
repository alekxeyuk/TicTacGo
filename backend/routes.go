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
	getRoom(globalRoomID).b.Submit(Message{"ping", "pong"})
}

func sendCountEvent() {
	getRoom(globalRoomID).b.Submit(Message{"count", strconv.FormatUint(roomsCounter, 10)})
}

func roomMOVE(c *gin.Context) {
	roomid := c.Param("roomid")
	room := getRoom(roomid)
	ok, _ := authorized(c)
	if roomid == "global" || room.id == "" || !ok {
		c.Status(http.StatusBadRequest)
		return
	}

	cell, err := strconv.ParseInt(c.Query("cell"), 10, 64)
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	time.Sleep(time.Duration(cell) * time.Second)
}

func roomRANDOM(c *gin.Context) {
	roomid, _ := newRoom()
	c.JSON(http.StatusCreated, gin.H{
		"uuid": roomid,
	})
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
			"uuid": room.id,
			"board": room.board,
			"current_player": room.currentPlayer.String(),
			"players": room.users,
			"state": room.state.String(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"rooms": roomsJSON,
	})
}

func roomDELETE(c *gin.Context) {
	roomid := c.Param("roomid")
	if roomid == globalRoomID {
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
		roomid = globalRoomID
	} else if getRoom(roomid).id == "" {
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

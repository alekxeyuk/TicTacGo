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
	getRoom(globalRoomID).Submit(Message{"ping", "pong"})
}

func sendCountEvent() {
	getRoom(globalRoomID).Submit(Message{"count", strconv.FormatUint(roomCounter, 10)})
}

func roomNEW(c *gin.Context) {
	roomid, _ := newRoom()
	c.JSON(http.StatusCreated, gin.H{
		"uuid": roomid,
	})
	go sendCountEvent()
}

func roomCOUNT(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"count": roomCounter,
	})
}

func roomLIST(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"rooms": "roomChannels",
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

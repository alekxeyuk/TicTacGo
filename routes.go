package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.templ.html", gin.H{
		"title":     "Main page",
		"roomCount": roomCounter,
		"rooms":     roomChannels,
	})
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
	getRoom(globalRoomID).Submit(Message{"ping", "pong"})
}

func roomNEW(c *gin.Context) {
	roomid, _ := newRoom()
	c.JSON(http.StatusCreated, gin.H{
		"uuid": roomid,
	})
	getRoom(globalRoomID).Submit(Message{"rooms", fmt.Sprintf("created room %s", roomid)})
}

func roomCOUNT(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"count": roomCounter,
	})
}

func roomLIST(c *gin.Context) {
	c.HTML(http.StatusOK, "index/list", gin.H{
		"rooms": roomChannels,
	})
}

func roomDELETE(c *gin.Context) {
	roomid := c.Param("roomid")
	if ok := deleteRoom(roomid); ok {
		c.JSON(http.StatusOK, gin.H{
			"message": "deleted",
		})
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
	ticker := time.NewTicker(1 * time.Second)
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
			// if msg == "stop" {
			// 	return false
			// }
		case <-ticker.C:
			c.SSEvent("time", time.Now().String())
		}
		return true
	})
}

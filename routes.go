package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.templ.html", gin.H{
		"title": "Main page",
		"roomCount": roomCounter,
	})
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func roomNEW(c *gin.Context) {
	roomid, _ := newRoom()
	c.JSON(http.StatusCreated, gin.H{
		"uuid": roomid.String(),
	})
}

func roomCOUNT(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"count": roomCounter,
	})
}

func roomDELETE(c *gin.Context) {
	roomid := c.Param("roomid")
	ok := deleteRoom(uuid.MustParse(roomid))
	if ok {
		c.JSON(http.StatusOK, gin.H{
			"message": "deleted",
		})
	} else {
		c.JSON(http.StatusGone, gin.H{
			"message": "not found",
		})
	}
}
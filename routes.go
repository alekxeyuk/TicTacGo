package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.templ.html", gin.H{
		"title": "Main page",
	})
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
package main

import (
	"fmt"
	"runtime"

	"github.com/alekxeyuk/TicTacGo/constant"
	"github.com/gin-gonic/gin"
)

func main() {
	ConfigRuntime()
	StartGin()
}

func ConfigRuntime() {
	nuCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(nuCPU)
	fmt.Printf("Running with %d CPUs\n", nuCPU)
	fmt.Printf("Version: %s\nBuilt at: %s\n", constant.Version, constant.BuildTime)
}

func setupRouter(router *gin.Engine) {
	router.GET("/", index)
	router.GET("/ping", ping)

	rooms := router.Group("/room")
	{
		rooms.POST("/new", roomNEW)
		rooms.GET("/count", roomCOUNT)
		rooms.GET("/list", roomLIST)
		rooms.DELETE("/:roomid", roomDELETE)
	}

	sse := router.Group("/stream")
	{
		sse.GET("/:roomid", roomSTREAM)
	}
}

func StartGin() {
	router := gin.Default()
	setupRouter(router)
	router.Run(":80")
}

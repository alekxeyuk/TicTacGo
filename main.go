package main

import (
	"fmt"
	"runtime"

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
}

func setupRouter(router *gin.Engine)  {
	router.Static("/static", "resources/static")
	router.GET("/", index)
	router.GET("/ping", ping)
}

func StartGin() {
	// gin.SetMode(gin.ReleaseMode)
	// router.Use(rateLimit, gin.Recovery())
	// router.GET("/room/:roomid", roomGET)
	// router.POST("/room-post/:roomid", roomPOST)
	// router.GET("/stream/:roomid", streamRoom)

	router := gin.Default()
	router.LoadHTMLGlob("resources/*.templ.html")
	setupRouter(router)
	router.Run(":80")
}
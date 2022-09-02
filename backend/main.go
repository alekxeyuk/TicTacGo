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
	router.Use(CORSMiddleware())
	router.Use(UserMiddleware())
	router.GET("/", index)
	router.GET("/ping", ping)

	rooms := router.Group("/room")
	{
		rooms.GET("/random", roomRANDOM)
		rooms.GET("/:roomid/stream", roomSTREAM)
		rooms.GET("/count", roomCOUNT)
		rooms.PUT("/:roomid/move", roomMOVE)
		rooms.GET("/list", roomLIST)
		// rooms.DELETE("/:roomid", roomDELETE)
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func StartGin() {
	router := gin.Default()
	setupRouter(router)
	router.Run(":80")
}

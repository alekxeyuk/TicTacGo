package main

import (
	"fmt"
	"html/template"
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

func setupRouter(router *gin.Engine)  {
	router.Static("/static", "resources/static")
	router.GET("/", index)
	router.GET("/ping", ping)

	rooms := router.Group("/room")
	{
		// rooms.GET("/:roomid", roomGET)
		// rooms.POST("/:roomid", roomPOST)
		rooms.POST("/new", roomNEW)
		rooms.GET("/count", roomCOUNT)
		rooms.DELETE("/:roomid", roomDELETE)
	}

	sse := router.Group("/stream")
	{
		sse.GET("/:roomid", roomSTREAM)
	}
}

func StartGin() {
	// gin.SetMode(gin.ReleaseMode)
	// router.Use(rateLimit, gin.Recovery())
	// router.GET("/room/:roomid", roomGET)
	// router.POST("/room-post/:roomid", roomPOST)
	// router.GET("/stream/:roomid", streamRoom)

	router := gin.Default()
	router.SetFuncMap(template.FuncMap{
        "plural": func (c uint64) (end string) {
			if c != 1 {
				end = "s"
			}
			return
		},
    })
	router.LoadHTMLGlob("resources/*.templ.html")
	setupRouter(router)
	router.Run(":80")
}
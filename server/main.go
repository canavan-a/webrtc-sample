package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

func main() {

	// No longer using webrtc, using WS streaming
	fmt.Println("Hello World")

	r := gin.Default()
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	r.Use(cors.New(config))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/relay", handleRelayServer)

	r.Run(":5000")

}

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	Mutex   = sync.Mutex{}
	Clients = make(map[*websocket.Conn]bool)
)

func handleRelayServer(c *gin.Context) {
	conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	Mutex.Lock()
	Clients[conn] = true
	Mutex.Unlock()

	defer func() {
		Mutex.Lock()
		Clients[conn] = false
		Mutex.Unlock()
	}()

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		Mutex.Lock()
		for client := range Clients {
			if Clients[client] && client != conn {
				client.WriteMessage(msgType, msg)
			}
		}
		Mutex.Unlock()
	}

}

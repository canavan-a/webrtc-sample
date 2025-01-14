package dialer

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Dialer struct {
	URL         string
	Conn        *websocket.Conn
	Channel     chan []byte
	sendChannel chan []byte
}

func NewDialer(url string) *Dialer {
	return &Dialer{
		URL:         url,
		Channel:     make(chan []byte, 100), // Buffered channel to prevent blocking
		sendChannel: make(chan []byte, 100), // Buffered channel for sending messages
	}
}

// Expose a method to access the Channel safely.
func (d *Dialer) GetChannel() chan []byte {
	return d.Channel
}

func (d *Dialer) Dial() (err error) {
	for {
		conn, _, err := websocket.DefaultDialer.Dial(d.URL, nil)
		if err != nil {
			return fmt.Errorf("dial failed: %w", err)
		}
		d.Conn = conn

		fmt.Println("connected")
		go d.listen() // Start listening for messages

		for {
			_, msg, err := d.Conn.ReadMessage()
			if err != nil {
				// If the connection is closed (normal or error), log it and break the loop
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					log.Println("Connection closed normally.")
				} else {
					log.Println("Error reading message:", err)
				}
				break
			}
			// Pass received message to the Channel
			d.Channel <- msg
		}

		// Connection lost or closed, reconnect after a delay
		d.Conn.Close()
		log.Println("Connection lost, reconnecting...")
	}
}

func (d *Dialer) Send(message string) {
	// Send message through the sendChannel
	d.sendChannel <- []byte(message)
}

func (d *Dialer) listen() {
	for {
		select {
		case msg := <-d.sendChannel:
			// Send binary message to WebSocket
			err := d.Conn.WriteMessage(websocket.BinaryMessage, msg)
			if err != nil {
				log.Println("Error sending message:", err)
				return
			}
		}
	}
}

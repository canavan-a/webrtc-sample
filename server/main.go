package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4"
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

	r.Run(":6789")

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
	Clients = make(map[string]*Connection)
)

type Connection struct {
	WebsocketConn *websocket.Conn
	WebrtcConn    *webrtc.PeerConnection
}

func handleRelayServer(c *gin.Context) {
	conn, err := Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	uuid, err := uuid.NewUUID()
	if err != nil {
		return
	}
	clientId := uuid.String()

	Mutex.Lock()
	myC := Connection{
		WebsocketConn: conn,
		WebrtcConn:    &webrtc.PeerConnection{},
	}
	Clients[clientId] = &myC
	Mutex.Unlock()

	defer func() {
		Mutex.Lock()
		delete(Clients, clientId)
		Mutex.Unlock()
	}()

	for {
		msgType, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		if msgType == websocket.TextMessage {
			var offer struct {
				Type          string `json:"type"`
				Sdp           string `json:"sdp,omitempty"`
				Candidate     string `json:"candidate,omitempty"`
				SdpMid        string `json:"sdpMid,omitempty"`
				SdpMLineIndex uint16 `json:"sdpMLineIndex,omitempty"`
			}
			err := json.Unmarshal(msg, &offer)
			if err != nil {
				fmt.Println("Error unmarshaling message:", err)
				continue
			}
			if offer.Type == "offer" {
				// Create a WebRTC offer object
				webrtcOffer := webrtc.SessionDescription{
					SDP:  offer.Sdp,
					Type: webrtc.SDPTypeOffer,
				}

				// Initialize PeerConnection (if not already done)\

				peerConnection, err := initPeerConnection(clientId, webrtcOffer)
				if err != nil {
					fmt.Println("Error initializing PeerConnection:", err)
					continue
				}

				Mutex.Lock()
				Clients[clientId].WebrtcConn = peerConnection
				Mutex.Unlock()

				// Process the offer (e.g., set it on a PeerConnection)
				fmt.Println("Received WebRTC offer:", webrtcOffer.SDP)
			} else if offer.Type == "candidate" {
				candidate := webrtc.ICECandidateInit{
					Candidate:     offer.Candidate,
					SDPMid:        &offer.SdpMid,
					SDPMLineIndex: &offer.SdpMLineIndex,
				}
				Mutex.Lock()
				err := Clients[clientId].WebrtcConn.AddICECandidate(candidate)
				Mutex.Unlock()
				if err != nil {
					fmt.Println("Error adding ice candidate:", err)
					continue
				}
			} else {
				fmt.Println("Received unsupported message type")
			}

		}

	}

}

func initPeerConnection(clientId string, offer webrtc.SessionDescription) (*webrtc.PeerConnection, error) {
	configuration := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"}, // Google's public STUN server
			},
		},
	}

	peerConnection, err := webrtc.NewPeerConnection(configuration)
	if err != nil {
		log.Fatal(err)
	}

	err = peerConnection.SetRemoteDescription(offer)
	if err != nil {
		return nil, err
	}

	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		// defer candidatesGathered.Done()
		if c == nil {
			return
		}

		candidateMessage := c.ToJSON() // This will serialize the ICECandidate

		Mutex.Lock()

		// Send the candidate as a WebSocket message to the client
		err := Clients[clientId].WebsocketConn.WriteJSON(candidateMessage)
		if err != nil {
			log.Fatal(err)
		}

		// Unlock the mutex after sending the candidate
		Mutex.Unlock()

		fmt.Println("CANDIDATE GATHERED")
		fmt.Println(candidateMessage)
	})

	addr := "127.0.0.1:5005"
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", addr, err)
	}
	defer conn.Close()

	videoTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: "video/H264"}, "video", "rtcVideoStream")
	if err != nil {
		log.Fatal(err)
	}

	audioTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: "audio/opus"}, "audio", "rtcAudiooStream")
	if err != nil {
		log.Fatal(err)
	}

	videoStreamer := CreateMediaStreamer(5005, "0.0.0.0", "video/H264", videoTrack)
	go videoStreamer.startReader()

	audioStreamer := CreateMediaStreamer(5006, "0.0.0.0", "audio/opus", audioTrack)
	go audioStreamer.startReader()

	_, err = peerConnection.AddTrack(videoTrack)
	_, err = peerConnection.AddTrack(audioTrack)

	var myAnswerOption webrtc.AnswerOptions
	mySDP, err := peerConnection.CreateAnswer(&myAnswerOption)
	if err != nil {
		log.Fatal(err)
	}
	err = peerConnection.SetLocalDescription(mySDP)
	if err != nil {
		log.Fatal(err)
	}
	Mutex.Lock()
	Clients[clientId].WebsocketConn.WriteMessage(websocket.TextMessage, []byte(mySDP.SDP))
	Mutex.Unlock()

	return peerConnection, nil

}

type RTPMediaStreamer struct {
	Port        uint16
	Hostname    string
	MimeType    string
	WebRTCTrack *webrtc.TrackLocalStaticRTP
}

func CreateMediaStreamer(Port uint16, Hostname string, MimeType string, WebRTCTrack *webrtc.TrackLocalStaticRTP) RTPMediaStreamer {
	return RTPMediaStreamer{
		Port: Port, Hostname: Hostname, MimeType: MimeType, WebRTCTrack: WebRTCTrack,
	}
}

func (s *RTPMediaStreamer) startReader() {
	address := fmt.Sprintf("%s:%d", s.Hostname, s.Port)

	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatal(err)

	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	buf := make([]byte, 1500)

	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("Error reading from UDP: %v", err)
			continue
		}

		var packet rtp.Packet

		err = packet.Unmarshal(buf[:n])
		if err != nil {
			log.Printf("Error parsing packet")
			continue
		}

		err = s.WebRTCTrack.WriteRTP(&packet)
		if err != nil {
			log.Printf("Error adding pcket to track")
			continue
		}

		fmt.Println(packet.SequenceNumber)
		fmt.Println(packet.Payload)

	}
}

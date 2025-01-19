package main

import (
	"fmt"
	"log"
	"net"
	"peer/dialer"

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4"
)

func main() {
	d := dialer.NewDialer("ws://localhost:6789/relay")
	go d.Dial()

	// for count < 10 {
	// 	d.Send("hello world")
	// 	time.Sleep(1 * time.Second)
	// 	count += 1
	// }

	// Use the exposed channel to listen for messages
	// channel := d.GetChannel()
	// for msg := range channel {
	// 	log.Printf("Received message: %s\n", msg)
	// }

	err := getICECandidates(doThing)
	if err != nil {
		panic(err)
	}

}
func doThing(b []byte) {
	fmt.Println(b)
}

func getICECandidates(onIceCandidateFunction func(b []byte)) error {
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

	// 3. Add media tracks (audio/video)
	// peerConnection.AddTrack(track) // Add your tracks here

	// 4. Gather ICE candidates (automatically handled)
	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		// defer candidatesGathered.Done()
		if c == nil {
			return
		}

		onIceCandidateFunction([]byte("hsas"))

		fmt.Println("CANDIDATE GATHERED")
		fmt.Println(c)
	})

	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		log.Fatal(err)
	}

	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		log.Fatal("Failed to set local description:", err)
	}

	// candidatesGathered.Wait()

	select {}

	return nil
}

func connectToStreams() {
	// Listen on UDP port 5004
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

	audioTrack, err := webrtc.NewTrackLocalStaticRTP(webrtc.RTPCodecCapability{MimeType: "audio/G722"}, "audio", "rtcAudiooStream")
	if err != nil {
		log.Fatal(err)
	}

	videoStreamer := CreateMediaStreamer(5005, "0.0.0.0", "video/H264", videoTrack)
	go videoStreamer.startReader()

	audioStreamer := CreateMediaStreamer(5006, "0.0.0.0", "audio/G722", audioTrack)
	go audioStreamer.startReader()

	select {}
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

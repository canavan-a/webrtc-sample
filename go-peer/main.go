package main

import (
	"fmt"
	"log"
	"net"

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v4"
)

func main() {
	// Listen on UDP port 5004
	addr := "127.0.0.1:5005"
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", addr, err)
	}
	defer conn.Close()

	// config := webrtc.Configuration{}

	// peerConnection, err := webrtc.NewPeerConnection(config)
	// if err != nil {
	// 	log.Fatal(err)
	// }

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

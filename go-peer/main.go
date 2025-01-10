package main

import (
	"fmt"
	"log"
	"net"

	"github.com/pion/rtp"
)

func main() {
	// Listen on UDP port 5004
	addr := "127.0.0.1:5005"
	conn, err := net.ListenPacket("udp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", addr, err)
	}
	defer conn.Close()

	buf := make([]byte, 2000) // Adjust buffer size as needed

	for {
		// Read a packet from the UDP connection
		n, addr, err := conn.ReadFrom(buf)
		if err != nil {
			log.Println("Error reading UDP packet:", err)
			continue
		}

		fmt.Printf("Received packet of size %d from %s\n", n, addr.String())

		// Decode the RTP packet
		packet := &rtp.Packet{}
		err = packet.Unmarshal(buf[:n])
		if err != nil {
			log.Println("Failed to unmarshal RTP packet:", err)
			continue
		}

		// Handle the RTP packet (e.g., print sequence number and timestamp)
		fmt.Printf("Received RTP packet - Seq: %d, Timestamp: %d\n", packet.SequenceNumber, packet.Timestamp)

		// You can process the packet further (e.g., decode video/audio data here)
	}

}

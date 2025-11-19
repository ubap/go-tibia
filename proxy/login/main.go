package main

import (
	"encoding/hex"
	"fmt"
	"goTibia/protocol"
	"io"
	"log"
	"net"
)

// The address and port for our dummy server to listen on.
const listenAddr = ":7171"

func main() {
	// Start listening for incoming TCP connections on the specified address.
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	// Ensure the listener is closed when the main function exits.
	defer listener.Close()

	log.Printf("Dummy server listening on %s. Waiting for Tibia client to connect...", listenAddr)

	// Loop forever, accepting new connections.
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	protoConn := protocol.NewConnection(conn)
	defer protoConn.Close()

	log.Printf("Accepted connection from %s", protoConn.RemoteAddr())

	messageBytes, err := protoConn.ReadMessage()
	if err != nil {
		if err == io.EOF {
			log.Printf("Client %s disconnected.", protoConn.RemoteAddr())
		} else {
			log.Printf("Error reading message from %s: %v", protoConn.RemoteAddr(), err)
		}
		return // End the handler for this connection.
	}

	packet, err := protocol.ParseLoginPacket(messageBytes)
	if err != nil {
		log.Printf("Error parsing login packet: %v", err)
		return
	}

	log.Printf("Packet received: %v", packet)

	// --- The most important part for reverse engineering ---
	// Print a detailed hex dump of the packet's body.
	fmt.Printf("\n--- Packet Received from %s ---\n", conn.RemoteAddr())
	fmt.Printf("%s", hex.Dump(messageBytes))
	fmt.Println("--- End of Packet ---")

	PrintAsGoSlice(messageBytes)
}

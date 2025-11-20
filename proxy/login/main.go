package main

import (
	"encoding/hex"
	"fmt"
	"goTibia/protocol"
	"io"
	"log"
	"net"
	"sync"
)

// The address and port for our dummy server to listen on.
const (
	listenAddr     = ":7171"
	realServerAddr = "world.fibula.app:7171"
)

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

func handleConnection(clientConn net.Conn) {
	protoConn := protocol.NewConnection(clientConn)
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
	log.Printf("Successfully decrypted packet: Account=%d, Password='%s', Version=%d",
		packet.AccountNumber, packet.Password, packet.ClientVersion)

	serverConn, err := net.Dial("tcp", realServerAddr)
	if err != nil {
		log.Printf("Failed to connect to real server at %s: %v", realServerAddr, err)
		return
	}
	protoServerConn := protocol.NewConnection(serverConn)
	defer protoServerConn.Close()
	log.Printf("Successfully connected to real server %s", realServerAddr)

	// 5. Re-serialize the packet, but this time encrypt it with the TARGET server's public key.
	outgoingMessageBytes, err := packet.Marshal()
	if err != nil {
		log.Printf("Failed to marshal outgoing packet: %v", err)
		return
	}

	// 6. Send the re-encrypted packet to the real server.
	if err := protoServerConn.WriteMessage(outgoingMessageBytes); err != nil {
		log.Printf("Failed to send login packet to real server: %v", err)
		return
	}
	log.Println("Sent re-encrypted login packet to real server.")

	// 7. Bridge the two connections to shuttle data back and forth.
	log.Println("Bridging connections...")
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		io.Copy(serverConn, clientConn)
		serverConn.Close()
	}()
	go func() {
		defer wg.Done()

		// Create our hex dumper with a clear label.
		dumper := &HexDumpWriter{Prefix: "SERVER -> CLIENT"}

		// Create a TeeReader. It reads from serverConn.
		// Everything it reads is also written to our dumper.
		teeReader := io.TeeReader(serverConn, dumper)

		// Now, copy from the teeReader to the client. The client gets the
		// exact same data, but we get to see it as it passes through.
		io.Copy(clientConn, teeReader)

		// Once copying is done, close the write-half of the client connection.
		clientConn.(*net.TCPConn).CloseWrite()
	}()

	wg.Wait()
	log.Printf("Connection bridge for %s closed.", clientConn.RemoteAddr())
}

type HexDumpWriter struct {
	// Prefix allows us to label the output, e.g., "SERVER->" or "CLIENT->"
	Prefix string
}

// Write is the only method needed to satisfy the io.Writer interface.
func (w *HexDumpWriter) Write(p []byte) (n int, err error) {
	fmt.Printf("\n--- Data Dump (%s) ---\n", w.Prefix)
	fmt.Printf("%s", hex.Dump(p))
	fmt.Println("--- End of Dump ---")
	// We return the number of bytes processed and no error.
	return len(p), nil
}

package login

import (
	"goTibia/protocol"
	"goTibia/proxy"
	"io"
	"log"
	"net"
	"sync"
)

// Server is a struct that manages the login proxy.
// It holds the configuration needed for the login process.
type Server struct {
	ListenAddr     string
	RealServerAddr string
	// You could add other dependencies here, like a specific logger.
}

// NewServer is a constructor for the login server.
func NewServer(listenAddr, realServerAddr string) *Server {
	return &Server{
		ListenAddr:     listenAddr,
		RealServerAddr: realServerAddr,
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	defer listener.Close()
	log.Printf("Login proxy listening on %s, forwarding to %s", s.ListenAddr, s.RealServerAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Login proxy failed to accept connection: %v", err)
			continue
		}
		go s.handleConnection(conn) // Call the method on our server instance
	}
}

func (p *Server) handleConnection(clientConn net.Conn) {
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

	serverConn, err := net.Dial("tcp", p.RealServerAddr)
	if err != nil {
		log.Printf("Failed to connect to real server at %s: %v", p.RealServerAddr, err)
		return
	}
	protoServerConn := protocol.NewConnection(serverConn)
	defer protoServerConn.Close()
	log.Printf("Successfully connected to real server %s", p.RealServerAddr)

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
		dumper := &proxy.HexDumpWriter{Prefix: "SERVER -> CLIENT"}

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

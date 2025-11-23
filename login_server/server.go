package login_server

import (
	"fmt"
	"goTibia/packets/login"
	"goTibia/protocol"
	"log"
	"net"
	"strconv"
	"time"
)

type Server struct {
	ListenAddr     string
	RealServerAddr string
}

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
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(clientConn net.Conn) {
	protoClientConn := protocol.NewConnection(clientConn)
	defer protoClientConn.Close()
	log.Printf("Login: Accepted connection from %s", protoClientConn.RemoteAddr())

	packetReader, err := protoClientConn.ReadMessage()
	if err != nil {
		log.Printf("error reading message from %s: %v", protoClientConn.RemoteAddr(), err)
		return
	}

	loginPacket, err := login.ParseCredentialsPacket(packetReader)
	if err != nil {
		log.Printf("Login: Failed to parse login packet: %v", err)
		return
	}

	protoServerConn, err := s.connectToServer()
	if err != nil {
		log.Printf("Login: Failed to connect to %s: %v", protoClientConn.RemoteAddr(), err)
		return
	}
	defer protoServerConn.Close()

	if err := protoServerConn.SendPacket(loginPacket); err != nil {
		log.Printf("Login: Failed to forward credentials to backend: %v", err)
		return
	}

	log.Println("Login: Credentials forwarded to backend.")

	protoServerConn.EnableXTEA(loginPacket.XTEAKey)
	protoClientConn.EnableXTEA(loginPacket.XTEAKey)

	message, err := protoServerConn.ReadMessage()
	if err != nil {
		log.Printf("Login: Failed to read server response for %s: %v", protoClientConn.RemoteAddr(), err)
		return
	}

	loginResultMessage, err := login.ParseLoginResultMessage(message)
	if err != nil {
		log.Printf("Login: Failed to receive login result message for %s: %v", protoClientConn.RemoteAddr(), err)
		return
	}

	s.injectMotd(loginResultMessage)

	err = protoClientConn.SendPacket(loginResultMessage)
	if err != nil {
		log.Printf("Login: Failed to send login result message for %s: %v", protoClientConn.RemoteAddr(), err)
		return
	}

	log.Printf("Login: Connection for %s finished.", protoClientConn.RemoteAddr())
}

func (s *Server) connectToServer() (*protocol.Connection, error) {
	conn, err := net.DialTimeout("tcp", s.RealServerAddr, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("backend unavailable at %s: %w", s.RealServerAddr, err)
	}

	return protocol.NewConnection(conn), nil
}

func (s *Server) injectMotd(message *login.LoginResultMessage) {
	message.Motd = &login.Motd{
		MotdId:  strconv.Itoa(int(time.Now().Unix())),
		Message: "Welcome to the go-tibia!\nImprove your go coding skills!",
	}
}

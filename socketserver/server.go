package socketserver

import (
	"backups/quit"
	"log"
	"net"
)

type TCPConnHandler func(*net.TCPConn)

type Server struct {
	ListenURL string
	Workers   int
	Handle    TCPConnHandler
	Control   *quit.Control
}

func (server *Server) Serve() error {
	log.Println("Server: Starting...")
	address, err := net.ResolveTCPAddr("tcp", server.ListenURL)
	if err != nil {
		return err
	}
	ln, err := net.ListenTCP("tcp", address)
	if err != nil {
		return err
	}

	connections := make(chan *net.TCPConn)
	for i := 0; i < server.Workers; i++ {
		go func() {
			quit, err := server.Control.Sub()
			if err != nil {
				log.Println("Server: terminated before starting.")
				return
			}
		loop:
			for {
				select {
				case conn := <-connections:
					log.Println("Server: Handling conn.")
					server.Handle(conn)
					conn.Close()
				case <-quit:
					log.Println("Server: Shuting down worker.", err)
					break loop
				}
			}
			quit <- true
		}()
	}

	go func() {
		for {
			conn, err := ln.AcceptTCP()
			if err != nil {
				log.Println("Server: Shuting down listener.", err)
				break
			}
			connections <- conn
		}
	}()

	quit, err := server.Control.Sub()
	if err != nil {
		return err
	}

	go func() {
		<-quit
		log.Printf("Server: Shutting Down...")
		ln.Close()
		quit <- true
	}()

	return nil
}

package socketserver

import (
	"log"
	"net"
	"backups/quit"
)


func ServeTCP(serverUrl string, workers int, handle func(*net.TCPConn)) error {
	log.Println("server: Starting Server.")
	address, err := net.ResolveTCPAddr("tcp", serverUrl)
	if err != nil {
		return err
	}
	ln, err := net.ListenTCP("tcp", address)
	if err != nil {
		return err
	}

	connections := make(chan *net.TCPConn)
	for i:=0; i< workers; i++ {
		go func() {
			quit := quit.Sub()
			for {
				select {
				case conn := <- connections:
					log.Println("server: Handling conn")
					handle(conn)
					conn.Close()
				case <-quit:
					log.Println("server: Shuting down socket worker", err)
					break
				}
			}
		}()
	}

	go func() {
		for {
			conn, err := ln.AcceptTCP()
			if err != nil {
				log.Println("server: Shuting down socket server", err)
				break
			}
			connections <- conn
		}
	}()

	quit := quit.Sub()
	go func() {
		<-quit
		ln.Close()
	}()
	return nil
}

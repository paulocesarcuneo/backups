package main

import (
	. "backups/director"
	"backups/quit"
	. "backups/reception"
	server "backups/socketserver"
	. "backups/storekeeper"
	. "backups/storeregistry"
	. "backups/synchronizer"
	. "backups/commands"
	"bufio"
	"log"
	"net"
	"io"
	// "time"
)

func tcpHandle(reception chan Request) func(*net.TCPConn){
	return func(conn*net.TCPConn) {
		reader := bufio.NewReader(conn)
		writer := bufio.NewWriter(conn)
		for {
			err, cmd := ReadCommand(reader)
			if err != nil {
				if err == io.EOF {
					break
				}
				WriteCommand(writer, Acknowledge{Err: err})
				continue
			}
			res:= make(chan Command)
			go func() {
				WriteCommand(writer, <-res)
				writer.Flush()
			}()
			reception <- Request{Input:cmd, Output: res}
		}
		log.Println("server: tcp closed")
	}
}

func main() {
	quit := quit.Sub()
	url:=":9000"
	workers:= 4
	storeRegistry := NewStoreRegistry()
	storekeeperChan := StoreKeeper(&storeRegistry)
	directorChan := Director(storekeeperChan, Synchronizer)
	receptionChan := Receptionist(directorChan, &storeRegistry)
	server.ServeTCP(url, workers, tcpHandle(receptionChan))
	<-quit
}

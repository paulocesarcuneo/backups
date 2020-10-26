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
	"backups/config"
	"bufio"
	"log"
	"net"
	"io"
	"os"
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
	if _, err := os.Stat(SyncFolder); os.IsNotExist(err) {
		os.Mkdir(SyncFolder, os.ModePerm)
	}
	storeRegistry := NewStoreRegistry()
	storekeeperChan := StoreKeeper(&storeRegistry)
	directorChan := Director(storekeeperChan, SynchronizeTask)
	receptionChan := Receptionist(directorChan, &storeRegistry)
	server.ServeTCP(config.Config.ServerUrl,  config.Config.Threads, tcpHandle(receptionChan))
	<-quit
}

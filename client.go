package main

import (
	. "backups/commands"
	"backups/quit"
	"backups/archiver"
	server "backups/socketserver"
	"fmt"
	"io"
	"log"
	"strings"
	"os"
	"bufio"
	"net"
)

func tcpHandler(conn *net.TCPConn) {
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
		switch v:=cmd.(type) {
		case Archive:
			md5, err := archiver.Tar(v.Path, "/tmp/archiver/" + strings.ReplaceAll(v.Path, "/", "_"))
			log.Printf("client: archived: ", md5, err)
			content:= md5
			WriteCommand(writer, Transfer{Size: len(content), Reader:strings.NewReader(content)})
		default :
			WriteCommand(writer, Acknowledge{Err: fmt.Errorf("Unhandled command")})
		}
		writer.Flush()
	}
	log.Println("client: tcp closed")
}

func main() {
	log.Println(string(1))
	url := ":9001"
	workers := 4
	path:= "/tmp/archiver"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
	server.ServeTCP(url, workers, tcpHandler)
	quit:= quit.Sub()
	<-quit
}

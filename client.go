package main

import (
	"backups/commands"
	"backups/quit"
	"backups/archiver"
	"backups/socketserver"
	"backups/config"
	"fmt"
	"io"
	"log"
	"os"
	"bufio"
	"net"
)

func archiveTCPHandler(conn *net.TCPConn) {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	archiveMan := archiver.NewArchiver()
	for {
		err, cmd := commands.ReadCommand(reader)
		switch {
		case err == io.EOF:
			break
		case err != nil:
			commands.WriteCommand(writer, commands.Acknowledge{Err: err})
		default:
			log.Println("client: handling ", cmd)
			switch archRequest := cmd.(type) {
			case commands.Archive:
				archiveMan.Transfer(archRequest, writer)
			default :
				commands.WriteCommand(writer, commands.Acknowledge{Err: fmt.Errorf("Unhandled command")})
			}
		}
		writer.Flush()
	}
}

func main() {
	path:= archiver.ArchivePath
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.ModePerm)
	}
	socketserver.ServeTCP(config.Config.ClientUrl,
		config.Config.Threads,
		archiveTCPHandler)
	quit:= quit.Sub()
	<-quit
}

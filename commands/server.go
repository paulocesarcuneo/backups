package commands

import (
	"bufio"
	"errors"
	"io"
	"log"
	"net"
	"syscall"
)

type CommandHandler func(Command) Command

func TCPCommandAdapter(commandHandler CommandHandler) func(*net.TCPConn) {
	return func(conn *net.TCPConn) {
		reader := bufio.NewReader(conn)
		writer := bufio.NewWriter(conn)
		for {

			var written int
			request, err := ReadCommand(reader)
			switch {
			case err == io.EOF:
				log.Println("TCPCommandAdapter: EOF")
				return
			case errors.Is(err, syscall.EPIPE):
				log.Println("TCPCommandAdapter: EPIPE")
				return
			case err != nil:
				log.Println("TCPCommandAdapter: Command Read Error", err)
				written, err = WriteCommand(writer, Error{Error: err.Error()})
			default:
				written, err = WriteCommand(writer, commandHandler(request))
			}
			if err != nil {
				log.Println("TCPCommandAdapter: Command Write Error", written, err)
			}
		}
	}
}

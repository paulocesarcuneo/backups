package commands

import (
	"bufio"
	"errors"
	"io"
	"net"
	"sync"
	"syscall"
)

type Client struct {
	URL    string
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
	lock   sync.Mutex
}

func (client *Client) Send(command Command) (Command, error) {
	client.lock.Lock()
	defer client.lock.Unlock()

	err := client.connect()
	if err != nil {
		return nil, err
	}
	_, err = WriteCommand(client.writer, command)
	if err != nil {
		client.resetOnBroken(err)
		return nil, errors.New("Client: Send " + err.Error())
	}
	var response Command
	response, err = ReadCommand(client.reader)
	if err != nil {
		client.resetOnBroken(err)
		err = errors.New("Client: Recv " + err.Error())
	}
	return response, err
}

func (client *Client) connect() (err error) {
	if client.conn == nil {
		client.conn, err = net.Dial("tcp", client.URL)
		if err != nil {
			return
		}
		client.reader = bufio.NewReader(client.conn)
		client.writer = bufio.NewWriter(client.conn)
	}
	return
}

func (client *Client) resetOnBroken(err error) {
	switch {
	case err == io.EOF:
		client.conn = nil
	case errors.Is(err, syscall.EPIPE):
		client.conn = nil
	default:
	}
}

func (client *Client) Close() error {
	client.lock.Lock()
	defer client.lock.Unlock()

	return client.conn.Close()
}

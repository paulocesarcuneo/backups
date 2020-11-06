package commands

import (
	"io"
	"os"
	"time"
)

type Command interface{}

type Register struct {
	Name     string
	Path     string
	Port     string
	Interval int
}

func (register *Register) URL() string {
	return register.Name + ":" + register.Port
}

type UnRegister struct {
	Name  string
	Path  string
	Cause string
}

type History struct {
	Name string
	Path string
}

type HistoryEntry struct {
	Size int64
	Date time.Time
}

type HistoryList struct {
	Events []HistoryEntry
}

type Archive struct {
	Path string
}

type Transfer struct {
	Size   int64
	Reader io.Reader
}

func (transfer *Transfer) WriteFile(location string) (int64, error) {
	file, err := os.Create(location)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	return transfer.WriteTo(file)
}

func (transfer *Transfer) WriteTo(writer io.Writer) (int64, error) {
	return io.CopyN(writer, transfer.Reader, transfer.Size)
}

type Store struct {
	Name     string
	Path     string
	Date     time.Time
	Transfer Transfer
}

type Acknowledge struct{}

type Error struct {
	Error string
}

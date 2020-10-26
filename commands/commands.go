package commands

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)


type Command interface{}

type Register struct {
	Name string
	Path string
	Port string
	Interval int
}

type UnRegister struct {
	Name string
	Path string
	Cause string
}

type History struct {
	Name string
	Path string
}

type HistoryEntry struct {
	Size int
	Date time.Time
}

type HistoryList struct {
	Events []HistoryEntry
}

type Archive struct {
	Path string
}

type Transfer struct {
	Size int
	Reader io.Reader
}

func (t* Transfer) WriteFile(location string) (string, error) {
	file, err := os.Create(location)
	if err != nil {
		return "failure", err
	}
	defer file.Close()
	_, err = io.CopyN(file, t.Reader, int64(t.Size))
	if err != nil {
		return "failure", err
	}
	return "ok", nil
}

type Acknowledge struct {
	Err error
}

func (t* Transfer) Write(path string) (int64, error) {
    f, err := os.Create(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	return io.CopyN(f, t.Reader, int64(t.Size))
}

var RegisterIntervalInvalid = errors.New("interval invalid format")
var RegisterInvalidFormat = errors.New("invalid format for register command")
var UnregisterInvalidFormat = errors.New("invalid format for unregister command")
var HistoryInvalidFormat = errors.New("invalid format for history command")
var UnknownCommand = errors.New("Unknown command")
var TransferInvalidSize = errors.New("invalid size value")
var UnableToHandle = errors.New("Unable to handle command")

func tokenize(reader *bufio.Reader) (error, []string) {
	data, err := reader.ReadString('.')
	if err != nil {
		return err, nil
	}
	tokens := strings.Split(strings.TrimSuffix(data,"."), ",")
	return nil, tokens
}

func ReadCommand(reader *bufio.Reader) (error, Command) {
	err, tokens := tokenize(reader)
	if err != nil {
		return err, nil
	}
	cmdName := strings.TrimSpace(strings.ToUpper(tokens[0]))
	tokens = tokens[1:]
	switch  cmdName {
	case "REGISTER":
		interval, err := strconv.Atoi(tokens[3])
		if err!= nil || interval <= 0 {
			return RegisterIntervalInvalid, nil
		}
		return nil, Register{
			Name:tokens[0],
			Path:tokens[1],
			Port:tokens[2],
			Interval: interval}
	case "UNREGISTER":
		cause:=""
		if len(tokens) == 3 {
			cause = tokens[2]
		}
		return nil, UnRegister{
			Name: tokens[0],
			Path: tokens[1],
			Cause: cause }
	case "HISTORY":
		return nil, History{
			Name:tokens[0],
			Path:tokens[1]}
	case "ARCHIVE":
		return nil, Archive{Path:tokens[0]}
	case "TRANSFER":
		size, err := strconv.Atoi(tokens[0])
		if err!= nil {
			return TransferInvalidSize, nil
		}
		return nil, Transfer{Size:size, Reader: reader}
	case "ACKNOWLEDGE":
		var errorToken error = nil
		if len(tokens) > 0 {
			errorToken = fmt.Errorf(tokens[0])
		}
		return nil, Acknowledge{Err: errorToken}
	default:
		return UnknownCommand, nil
	}
}

func  WriteCommand(w *bufio.Writer, acmd Command) (int, error) {
	switch cmd := acmd.(type) {
	case Register:
		return w.WriteString("register,"+cmd.Name+","+cmd.Path+","+ cmd.Port+","+strconv.Itoa(cmd.Interval)+".")
	case UnRegister:
		return w.WriteString("unregister,"+cmd.Name+","+cmd.Path+","+cmd.Cause+".")
	case History:
		return w.WriteString("history,"+cmd.Name+","+cmd.Path+".")
	case Archive:
		return w.WriteString("archive,"+cmd.Path+".")
	case Transfer:
		written, err := w.WriteString("transfer,"+strconv.Itoa(cmd.Size)+".")
		if err != nil {
			return written, err
		}
		// TODO Fix written type, should be int64
		writtenFile, err := io.CopyN(w, cmd.Reader, int64(cmd.Size))
		return written + int(writtenFile), err
	case Acknowledge:
		msg := ""
		if cmd.Err != nil {
			msg = cmd.Err.Error()
		}
		return w.WriteString("acknowledge,"+msg+".")
	case HistoryList:
		total, err:= w.WriteString("historylist,\ndate,size\n")
		for _, e:= range cmd.Events {
			if err!=nil {
				break
			}
			var written int
			written, err = w.WriteString(e.Date.Format(time.RFC3339Nano) +","+strconv.Itoa(e.Size)+"\n")
			total += written
		}
		w.WriteString(".")
		return total, err

	default:
		return 0, UnknownCommand
	}
}

/*
register,n1,/tmp/test,9001,8.
acknowledge,.

unregister,n1,/tmp/test.
acknowledge,.
history,n1,/tmp/test.

register n1 /etc1 001 8
register n2 /etc2 002 5
register n3 /etc3 003 6
register n4 /etc4 004 4
unregister n1 /etc1

*/

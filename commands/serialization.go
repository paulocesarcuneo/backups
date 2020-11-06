package commands

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

var RegisterIntervalInvalid = errors.New("interval invalid format")
var RegisterInvalidFormat = errors.New("invalid format for register command")
var UnregisterInvalidFormat = errors.New("invalid format for unregister command")
var HistoryInvalidFormat = errors.New("invalid format for history command")
var UnknownCommand = errors.New("Unknown command")
var TransferInvalidSize = errors.New("invalid size value")
var UnableToHandle = errors.New("Unable to handle command")

func tokenize(reader *bufio.Reader) ([]string, error) {
	data, err := reader.ReadString('.')
	if err != nil {
		return nil, err
	}
	tokens := strings.Split(strings.TrimSuffix(data, "."), ",")
	return tokens, nil
}

func ReadCommand(reader *bufio.Reader) (Command, error) {
	tokens, err := tokenize(reader)
	if err != nil {
		return nil, err
	}
	cmdName := strings.TrimSpace(strings.ToUpper(tokens[0]))
	tokens = tokens[1:]
	switch cmdName {
	case "REGISTER":
		interval, err := strconv.Atoi(tokens[3])
		if err != nil || interval <= 0 {
			return nil, RegisterIntervalInvalid
		}
		return Register{
			Name:     tokens[0],
			Path:     tokens[1],
			Port:     tokens[2],
			Interval: interval}, nil
	case "UNREGISTER":
		cause := ""
		if len(tokens) == 3 {
			cause = tokens[2]
		}
		return UnRegister{
			Name:  tokens[0],
			Path:  tokens[1],
			Cause: cause}, nil
	case "HISTORY":
		return History{
			Name: tokens[0],
			Path: tokens[1]}, nil
	case "ARCHIVE":
		return Archive{Path: tokens[0]}, nil
	case "TRANSFER":
		size, err := strconv.ParseInt(tokens[0], 10, 64)
		if err != nil {
			return nil, TransferInvalidSize
		}
		return Transfer{Size: size, Reader: reader}, nil
	case "ACKNOWLEDGE":
		return Acknowledge{}, nil
	case "ERROR":
		return Error{Error: tokens[0]}, nil
	default:
		return nil, errors.New("Unknown Command:" + cmdName)
	}
}

func WriteCommand(w *bufio.Writer, acmd Command) (int, error) {
	defer w.Flush()

	switch cmd := acmd.(type) {
	case Register:
		return w.WriteString("register," + cmd.Name + "," + cmd.Path + "," + cmd.Port + "," + strconv.Itoa(cmd.Interval) + ".")
	case UnRegister:
		return w.WriteString("unregister," + cmd.Name + "," + cmd.Path + "," + cmd.Cause + ".")
	case History:
		return w.WriteString("history," + cmd.Name + "," + cmd.Path + ".")
	case Archive:
		return w.WriteString("archive," + cmd.Path + ".")
	case Transfer:
		written, err := w.WriteString("transfer," + strconv.FormatInt(cmd.Size, 10) + ".")
		if err != nil {
			return written, err
		}
		writtenFile, err := cmd.WriteTo(w)
		switch file := cmd.Reader.(type) {
		case *os.File:
			file.Close()
		default:
		}
		return written + int(writtenFile), err
	case Acknowledge:
		return w.WriteString("acknowledge.")
	case Error:
		return w.WriteString("error," + cmd.Error + ".")
	case HistoryList:
		total, err := w.WriteString("historylist,\ndate,size\n")
		for _, e := range cmd.Events {
			if err != nil {
				break
			}
			var written int
			written, err = w.WriteString(e.Date.Format(time.RFC3339Nano) + "," + strconv.FormatInt(e.Size, 10) + "\n")
			total += written
		}
		w.WriteString(".")
		return total, err
	default:
		return 0, UnknownCommand
	}
}

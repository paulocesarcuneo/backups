package synchronizer

import (
	. "backups/commands"
	. "backups/storeregistry"
	"backups/quit"
	"bufio"
	"io"
	"log"
	"net"
	"strings"
	"time"
	"strconv"
)

const SyncFolder = "/tmp/synchronizer"

type Synchronizer struct {
	archiverConn net.Conn
	attemptRound int
	register Register
}

func NewSynchronizer(syncRequest Register) (*Synchronizer, error) {
	conn, err := net.Dial("tcp", syncRequest.Port)
	if err != nil {
		return nil, err
	}
	return &Synchronizer{archiverConn: conn, attemptRound: 0, register: syncRequest}, nil
}

func (me *Synchronizer) sync() (*StoreEvent, error) {
	reader := bufio.NewReader(me.archiverConn)
	writer := bufio.NewWriter(me.archiverConn)
	_, err := WriteCommand(writer, Archive{Path: me.register.Path})
	if err != nil {
		return nil, err
	}
	writer.Flush()
	err, cmd := ReadCommand(reader)
	if err != nil {
		return nil, err
	}
	switch transfer :=cmd.(type) {
	case Transfer:
		timeOfTransfer := time.Now()
		localPath := syncFileName(me.attemptRound, me.register.Path)
		status, err := transfer.WriteFile(localPath)
		me.attemptRound++
		return &StoreEvent{Size:transfer.Size,
			Date:timeOfTransfer,
			Path:me.register.Path,
			Node:me.register.Name,
			Status: status,
			LocalPath: localPath}, err
	default:
		return nil, UnableToHandle
	}

}

func syncFileName(attemptRound int, path string) string {
	return SyncFolder +"/"+ strconv.Itoa(attemptRound) + strings.ReplaceAll(path, "/", "_")
}

func (me *Synchronizer) Close() {
	me.archiverConn.Close()
}


func SynchronizeTask(req Register, events chan StoreEvent, directorChan chan Command) chan Command {
	in:= make(chan Command)
	go func(){
		q := quit.Sub()
		log.Println("Synchronizer start:", req)
		var synchroMan *Synchronizer
		var err error
		loop: for {
			select {
			case <-q:
				break loop
			case cmd := <- in:
				switch v:=cmd.(type) {
				case Register:
					log.Println("Synchronizer: updated ", req , " to ", v)
					synchroMan.register = req
				case UnRegister:
					log.Println("Synchronizer: unregister ")
					quit.UnSub(q)
					break loop
				default:
					log.Println("Synchronizer: unhandled command ", cmd)
				}
			case <-time.After(time.Duration(req.Interval) * time.Second):
				log.Println("Synchronizer: syncing")
				if synchroMan == nil {
					synchroMan, err = NewSynchronizer(req)
					if err != nil {
						log.Println("Synchronizer: can't connect ", req)
						continue loop
					}
				}
				event, err := synchroMan.sync()
				log.Println("Synchronizer: event err ", event, err)
				switch{
				case err == io.EOF:
					log.Println("Syncronizer: failed ", err)
					// directorChan <- Unreachable{Name: req.Name, Path:req.Path, Cause: err.Error()}
					continue loop
				case event == nil && err != nil :
					log.Println("Synchronizer: skip sync ", req)
					continue loop
				default:
					log.Println("Synchronizer: synced", event)
					events <- *event
				}
			}
		}
		if synchroMan != nil {
			synchroMan.Close()
		}
	}()
	return in
}

package reception


import (
	. "backups/commands"
	. "backups/storeregistry"
	"backups/quit"
	"log"
	"errors"
)

type Request struct {
	Input Command
	Output chan Command
}

var UnhandledCommand = errors.New("Unable to handle request")
func Receptionist(director chan Command, store *StoreRegistry) chan Request {
	in:= make(chan Request)
	go func(){
		quitC := quit.Sub()
		loop: for {
			select{
				case <-quitC:
				break loop
				case req := <-in:
				switch cmd:= req.Input.(type) {
				case Register:
					log.Println("reception: handle register")
					director <- cmd
					req.Output <- Acknowledge{}
				case UnRegister:
					log.Println("reception: handle unregister")
					director<- cmd
					req.Output <- Acknowledge{}
				case History:
					log.Println("reception: handle history")
					data := store.FetchHistory(cmd.Name, cmd.Path)
					req.Output <- formatHistory(data)
				default:
					log.Println("Reception: unable to handle request", cmd)
					req.Output<- Acknowledge{Err: UnhandledCommand}
				}
			}
		}
	}()
	return in
}

func formatHistory(history []StoreEvent) HistoryList {
	var events []HistoryEntry
	for _,se := range history {
		events = append(events, HistoryEntry{Size:se.Size, Date: se.Date})
	}
	return HistoryList{Events:events}
}

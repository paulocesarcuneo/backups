package synchronizer


import (
	. "backups/commands"
	. "backups/storeregistry"
	"backups/quit"
	"log"
	"time"
)

func Synchronizer(req Register, storekeeper chan StoreEvent) chan Command {
	in:= make(chan Command)
	go func(){
		q := quit.Sub()
		loop: for {
			select {
			case <-q:
				break loop
			case cmd := <- in:
				switch v:=cmd.(type) {
					case Register:
						log.Println("Synchronizer: updated ", req , " to ", v)
						req = v
					case UnRegister:
						log.Println("Synchronizer: unregister")
						quit.UnSub(q)
						break loop
					default:
					    log.Println("Synchronizer: unhandled command", cmd)
				}
			case <-time.After(time.Duration(req.Interval) * time.Second):
				log.Println("Synchronizer: syncing", req)

				storekeeper <- StoreEvent{Size:0,
					Date:time.Now(),
					Path:req.Path,
					Node:req.Name,
					LocalPath: req.Path}
			}
		}
	}()
	return in
}

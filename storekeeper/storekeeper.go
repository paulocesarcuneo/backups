package storekeeper

import (
	. "backups/storeregistry"
	"backups/quit"
	"log"
)

func StoreKeeper(store *StoreRegistry) chan StoreEvent {
	in:= make(chan StoreEvent)
	go func(){
		quit:= quit.Sub()
		loop: for{
			select {
			case <-quit:
				break loop
			case req:=<-in:
				log.Println("StoreKeeper: register event", req)
				history := store.FetchHistory(req.Node, req.Path)
				if len(history) > 10 {
					log.Println("StoreKeeper: house keeping")
				}
				store.Register(req)
			}
		}
	}()
	return in
}

package storekeeper

import (
	"backups/quit"
	"backups/config"
	. "backups/storeregistry"
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
				store.Register(req)
				history := store.FetchHistory(req.Node, req.Path)
				if len(history) > config.Config.HistorySize {
					log.Println("StoreKeeper: house keeping")
					store.HouseKeep(req)
				}
			}
		}
	}()
	return in
}

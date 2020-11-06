package synchronizer

import (
	"backups/commands"
	"backups/noderegistry"
	"backups/quit"
	"log"
	"sync"
	"time"
)

type StoreHandler func(commands.Command) commands.Command

type Synchronizer struct {
	NodeRegistry *noderegistry.Registry
	Store        StoreHandler
	Control      *quit.Control
	Request      *chan commands.Command
	lock         sync.Mutex
}

func (synchronizer *Synchronizer) sync(register commands.Register) {
	client := synchronizer.NodeRegistry.Lease(register)

	if client == nil {
		return
	}

	response, err := client.Send(commands.Archive{Path: register.Path})
	if err != nil {
		log.Println("Synchronizer: failed to archive ", err)
	} else {
		switch transfer := response.(type) {
		case commands.Transfer:
			storeResponse := synchronizer.Store(commands.Store{
				Date:     time.Now(),
				Name:     register.Name,
				Path:     register.Path,
				Transfer: transfer})
			log.Println("Synchronizer: Store Response", storeResponse)
		default:
			log.Println("Synchronizer: Archive Response ", transfer)
		}
	}

	if !synchronizer.NodeRegistry.Release(register) {
		time.AfterFunc(time.Duration(register.Interval)*time.Second, func() {
			*synchronizer.Request <- register
		})
	}
}

func (synchronizer *Synchronizer) Handle(command interface{}) {
	switch register := command.(type) {
	case commands.Register:
		log.Println("Synchronizer: register ", register)
		synchronizer.sync(register)
	default:
		log.Println("Synchronizer: unhandled command")
	}
}

func (synchronizer *Synchronizer) Launch() error {
	quit, err := synchronizer.Control.Sub()
	if err != nil {
		return err
	}
	go func() {
	loop:
		for {
			select {
			case <-quit:
				break loop
			case cmd := <-*synchronizer.Request:
				synchronizer.Handle(cmd)
			}
		}
		quit <- true
	}()
	return nil
}

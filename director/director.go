package director

import (
	"backups/commands"
	"backups/noderegistry"
	"backups/quit"
	"log"
)

type Director struct {
	SynchronizersIn *chan commands.Command
	DirectorIn      chan commands.Command
	NodeRegistry    *noderegistry.Registry
	Control         *quit.Control
}

func (director *Director) Launch() {
	go func() {
		quit, err := director.Control.Sub()
		if err != nil {
			log.Println("Director: Launch failed ", err)
			return
		}
		log.Println("Director: stated")
	loop:
		for {
			select {
			case <-quit:
				break loop
			case req := <-director.DirectorIn:
				log.Println("Director: req received")
				switch cmd := req.(type) {
				case commands.Register:
					log.Println("Director: register ", cmd)
					err = director.NodeRegistry.Register(cmd)
					if err != nil {
						log.Println("Director: register failed", err)
						continue
					}
					*director.SynchronizersIn <- cmd
				case commands.UnRegister:
					err = director.NodeRegistry.UnRegister(cmd)
					if err != nil {
						log.Println("Director: unregister failed", err)
					}
				default:
					log.Println("Director: unable to handle command", cmd)
				}
			}
		}
		quit <- true
	}()
}

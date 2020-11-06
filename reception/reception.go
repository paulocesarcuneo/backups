package reception

import (
	"backups/commands"
	"backups/storeregistry"
	"log"
)

type Reception struct {
	Director chan commands.Command
	Store    *storeregistry.StoreRegistry
}

const UnhandledCommand = "Reception: Unable to handle request"

func (reception *Reception) Handle(request commands.Command) commands.Command {
	switch cmd := request.(type) {
	case commands.Register:
		log.Println("Reception: handle register")
		reception.Director <- cmd
		return commands.Acknowledge{}
	case commands.UnRegister:
		log.Println("Reception: handle unregister")
		reception.Director <- cmd
		return commands.Acknowledge{}
	case commands.History:
		log.Println("Reception: handle history")
		return reception.Store.FetchHistory(cmd.Name, cmd.Path)
	default:
		log.Println("Reception: unable to handle request", cmd)
		return commands.Error{Error: UnhandledCommand}
	}
}

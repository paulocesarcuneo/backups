package director


import (
	. "backups/commands"
	. "backups/storeregistry"
	. "backups/noderegistry"
	"backups/quit"
	"log"
)

func Director(storekeeper chan StoreEvent,
	sync func(Register, chan StoreEvent) chan Command) chan Command {
	in:=make(chan Command)
	go func() {
		quit:= quit.Sub()
		registry := NodeRegistry{Nodes: make(map[string]*NodeEntry)}
		loop: for{
			select {
			case <-quit:
				break loop
			case req :=<- in:
				switch cmd := req.(type) {
				case Register:
					log.Println("Director: register ", cmd)
					var syncIn chan Command
					node := registry.Get(cmd.Name)
					if node != nil {
						node.In <- cmd
						syncIn = node.In
					} else {
						syncIn = sync(cmd, storekeeper)
					}
					registry.Add(cmd, syncIn)
				case UnRegister:
					node := registry.Get(cmd.Name)
					node.In <- cmd
					log.Println("Director: unregister", cmd, node)
					registry.Rem(cmd)
				default:
					log.Println("Director: unable to handle command", cmd)
				}
			}
		}
	}()
	return in
}

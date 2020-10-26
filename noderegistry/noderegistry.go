package noderegistry

import (
	. "backups/commands"
)

type NodeRegistry struct {
	Nodes map[string]*NodeEntry
}

type NodeEntry struct {
    In chan Command
	Register Register
}

func (n NodeRegistry) Get(name string) *NodeEntry {
	return n.Nodes[name]
}

func (n *NodeRegistry) Add(cmd Register, syncChan chan Command) {
	n.Nodes[cmd.Name] = &NodeEntry{syncChan, cmd}
}

func (n *NodeRegistry) Rem(cmd UnRegister) {
	delete(n.Nodes, cmd.Name)
}

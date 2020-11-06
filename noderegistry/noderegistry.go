package noderegistry

import (
	"backups/commands"
	"errors"
	"sync"
)

type Registry struct {
	Nodes map[Entry]*Client
	lock  sync.Mutex
}

func NewRegistry() Registry {
	return Registry{
		Nodes: make(map[Entry]*Client),
		lock:  sync.Mutex{},
	}
}

type Client struct {
	Delete bool
	Active bool
	Client *commands.Client
}

type Entry struct {
	Name string
	Path string
}

func (registry *Registry) Register(register commands.Register) error {
	registry.lock.Lock()
	defer registry.lock.Unlock()

	entry := Entry{
		Name: register.Name,
		Path: register.Path,
	}

	client := registry.Nodes[entry]
	url := register.URL()

	if client == nil {
		registry.Nodes[entry] = &Client{Active: false, Client: &commands.Client{URL: register.URL()}}
		return nil
	}
	if url != client.Client.URL {
		return errors.New("Registered to a different PORT")
	}
	return nil
}

func (registry *Registry) Lease(register commands.Register) *commands.Client {
	registry.lock.Lock()
	defer registry.lock.Unlock()

	entry := Entry{
		Name: register.Name,
		Path: register.Path,
	}

	client := registry.Nodes[entry]

	if client == nil {
		return nil
	}

	client.Active = true

	if client.Delete {
		client.Client.Close()
		delete(registry.Nodes, entry)
		return nil
	}

	return client.Client
}

func (registry *Registry) Release(register commands.Register) bool {
	registry.lock.Lock()
	defer registry.lock.Unlock()

	entry := Entry{
		Name: register.Name,
		Path: register.Path,
	}

	client := registry.Nodes[entry]
	client.Active = false
	if client.Delete {
		client.Client.Close()
		delete(registry.Nodes, entry)
		return true
	}
	return false
}

func (registry *Registry) UnRegister(unregister commands.UnRegister) error {
	registry.lock.Lock()
	defer registry.lock.Unlock()

	entry := Entry{
		Name: unregister.Name,
		Path: unregister.Path,
	}

	client := registry.Nodes[entry]
	if client == nil {
		return errors.New("Register doesnt exists.")
	}
	client.Delete = true
	if !client.Active {
		client.Client.Close()
		delete(registry.Nodes, entry)
	}
	return nil
}

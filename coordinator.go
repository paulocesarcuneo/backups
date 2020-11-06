package main

import (
	"backups/commands"
	"backups/config"
	"backups/director"
	"backups/noderegistry"
	"backups/quit"
	"backups/reception"
	"backups/socketserver"
	"backups/storeregistry"
	"backups/synchronizer"
)

func main() {
	control := quit.NewControl()

	storage := storeregistry.NewStoreRegistry(config.Config.HistorySize)

	synchronizerChan := make(chan commands.Command)
	registry := noderegistry.NewRegistry()
	synchronizer := synchronizer.Synchronizer{
		Control:      &control,
		NodeRegistry: &registry,
		Request:      &synchronizerChan,
		Store:        storage.Handle,
	}
	for i := 0; i < config.Config.Threads; i++ {
		synchronizer.Launch()
	}

	directorChan := make(chan commands.Command)
	director := director.Director{
		SynchronizersIn: &synchronizerChan,
		Control:         &control,
		DirectorIn:      directorChan,
		NodeRegistry:    &registry,
	}

	reception := reception.Reception{
		Director: directorChan,
		Store:    &storage}

	server := socketserver.Server{
		Control:   &control,
		ListenURL: config.Config.ServerUrl,
		Workers:   config.Config.Threads,
		Handle:    commands.TCPCommandAdapter(reception.Handle),
	}

	director.Launch()
	server.Serve()
	control.WaitForTermSignal()
}

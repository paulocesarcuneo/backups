package main

import (
	"backups/archiver"
	"backups/commands"
	"backups/config"
	"backups/quit"
	"backups/socketserver"
)

func main() {
	control := quit.NewControl()
	archiver := archiver.NewArchiver()
	server := socketserver.Server{
		ListenURL: config.Config.ClientUrl,
		Workers:   config.Config.Threads,
		Control:   &control,
		Handle:    commands.TCPCommandAdapter(archiver.Handle),
	}
	server.Serve()
	control.WaitForTermSignal()
}

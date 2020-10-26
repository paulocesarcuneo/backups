package config

import (
	"os"
	"strconv"
)

type Configuration struct {
	HistorySize int
	ClientUrl string
	ServerUrl string
	Threads int
}

// TODO handle config parsing errors
func loadEnv() Configuration {
	historySize := 10
	historySizeStr := os.Getenv("HISTORY_SIZE")
	if historySizeStr != "" {
		historySize , _ = strconv.Atoi(historySizeStr)
	}
	threads := 10
	threadsStr := os.Getenv("THREADS")
	if threadsStr != "" {
		threads, _  = strconv.Atoi(threadsStr)
	}
	serverUrl := os.Getenv("SERVER_URL")
	if serverUrl == "" {
		serverUrl = ":9000"
	}
	clientUrl := os.Getenv("CLIENT_URL")
	if clientUrl == "" {
		clientUrl = ":9001"
	}
	return Configuration{
		HistorySize: historySize,
		ClientUrl: clientUrl,
		ServerUrl: serverUrl,
		Threads: threads}
}

var Config Configuration = loadEnv()

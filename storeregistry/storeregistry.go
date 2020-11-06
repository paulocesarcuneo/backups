package storeregistry

import (
	"backups/commands"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const StoreFolder = "/tmp/storage/"

type StoreEvent struct {
	Size      int64
	Date      time.Time
	Path      string
	Name      string
	Status    string
	LocalPath string
}

type StoreRegistry struct {
	HistorySizeLimit int
	Data             map[string]map[string][]StoreEvent
	lock             sync.Mutex
}

func NewStoreRegistry(historySize int) StoreRegistry {
	return StoreRegistry{
		HistorySizeLimit: historySize,
		Data:             make(map[string]map[string][]StoreEvent),
	}
}

func (s *StoreRegistry) FetchHistory(name, path string) commands.Command {
	nodeEvents := s.Data[name]
	if nodeEvents != nil {
		history := nodeEvents[path]
		var events []commands.HistoryEntry
		for _, se := range history {
			events = append(events, commands.HistoryEntry{Size: se.Size, Date: se.Date})
		}
		return commands.HistoryList{Events: events}
	}
	return commands.Acknowledge{}
}

func (store *StoreRegistry) HouseKeep(e StoreEvent) (err error) {
	nodeHistory := store.Data[e.Name]
	if nodeHistory == nil {
		store.Data[e.Name] = make(map[string][]StoreEvent)
	}
	nodePathHistory := store.Data[e.Name][e.Path]
	nodePathHistory = append(nodePathHistory, e)
	if len(nodePathHistory) > store.HistorySizeLimit {
		err = os.Remove(nodePathHistory[0].LocalPath)
		nodePathHistory = nodePathHistory[1:]
	}
	store.Data[e.Name][e.Path] = nodePathHistory
	return
}

func (store *StoreRegistry) write(request commands.Store) (string, error) {
	currentSize := len(store.Data[request.Name][request.Path])
	localPath := StoreFolder + request.Name + strings.ReplaceAll(request.Path, "/", "_") + strconv.Itoa(currentSize)
	written, err := request.Transfer.WriteFile(localPath)
	log.Println("Storage: File written", written, err)
	if err != nil {
		return "", err
	}
	return localPath, nil
}

func (store *StoreRegistry) Store(request commands.Store) error {
	store.lock.Lock()
	defer store.lock.Unlock()

	localPath, err := store.write(request)
	if err != nil {
		return err
	}
	storeEvent := StoreEvent{
		Name:      request.Name,
		Path:      request.Path,
		Date:      request.Date,
		Size:      request.Transfer.Size,
		LocalPath: localPath,
	}
	store.HouseKeep(storeEvent)
	return nil
}

func (store *StoreRegistry) Handle(request commands.Command) commands.Command {
	switch cmd := request.(type) {
	case commands.Store:
		if err := store.Store(cmd); err != nil {
			return commands.Error{Error: err.Error()}
		}
		return commands.Acknowledge{}
	case commands.History:
		return store.FetchHistory(cmd.Name, cmd.Name)
	default:
		return commands.Error{Error: commands.UnableToHandle.Error()}
	}
}

func init() {
	if _, err := os.Stat(StoreFolder); os.IsNotExist(err) {
		os.Mkdir(StoreFolder, os.ModePerm)
	}
}

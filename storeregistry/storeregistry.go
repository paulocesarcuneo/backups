package storeregistry

import (
	"time"
	"os"
)

type StoreEvent struct {
	Size int
	Date time.Time
	Path string
	Node string
	Status string
	LocalPath string
}

type StoreRegistry struct {
	Data map[string]map[string][]StoreEvent
}

func (s StoreRegistry) FetchHistory(name, path string) []StoreEvent {
	nodeEvents:= s.Data[name]
	if nodeEvents != nil {
		return nodeEvents[path]
	}
	return nil
}

func (s * StoreRegistry) HouseKeep(e StoreEvent) {
	history := s.Data[e.Node][e.Path]
	oldest := history[0]
	os.Remove(oldest.LocalPath)
	s.Data[e.Node][e.Path] = history[1:]
}

func (s* StoreRegistry) Register(e StoreEvent) {
	nodeData := s.Data[e.Node]
	if nodeData == nil {
		nodeData = make(map[string][]StoreEvent)
		s.Data[e.Node] = nodeData
	}
	history := nodeData[e.Path]
	s.Data[e.Node][e.Path] = append(history, e)
}

func NewStoreRegistry() StoreRegistry {
	return StoreRegistry{Data: make(map[string]map[string][]StoreEvent)}
}

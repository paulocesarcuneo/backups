package storeregistry

import (
	"time"
)

type StoreEvent struct {
	Size int
	Date time.Time
	Path string
	Node string
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

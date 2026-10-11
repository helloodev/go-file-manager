package fm

import (
	"strconv"
	"sync"
)

type IndexStore struct {
	sync.RWMutex
	paths []string
	index map[string]int
}

var Store = &IndexStore{
	paths: []string{""},
	index: make(map[string]int),
}

func (s *IndexStore) GetID(path string) string {
	if path == "" || path == "." {
		return "0"
	}

	s.Lock()
	defer s.Unlock()

	if idx, exists := s.index[path]; exists {
		return strconv.Itoa(idx)
	}

	newIdx := len(s.paths)
	s.paths = append(s.paths, path)
	s.index[path] = newIdx

	return strconv.Itoa(newIdx)
}

func (s *IndexStore) GetPath(IdStr string) string {
	idx, err := strconv.Atoi(IdStr)
	if err != nil || idx < 0 {
		return ""
	}

	s.RLock()
	defer s.RUnlock()

	if idx < len(s.paths) {
		return s.paths[idx]
	}
	return ""
}
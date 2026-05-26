package core_domain

import "sync"

type StopList struct {
	mu   sync.RWMutex
	List map[string]struct{} `json:"list"`
}

func NewStopList(l []string) *StopList {
	stopList := &StopList{
		List: make(map[string]struct{}, len(l)),
	}

	for _, item := range l {
		stopList.List[item] = struct{}{}
	}

	return stopList
}

func (sl *StopList) Add(item string) {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	sl.List[item] = struct{}{}
}

func (sl *StopList) Remove(item string) {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	delete(sl.List, item)
}

func (sl *StopList) IsBlocked(item string) bool {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	_, ok := sl.List[item]
	return ok
}

func (sl *StopList) Items() []string {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	result := make([]string, 0, len(sl.List))

	for item := range sl.List {
		result = append(result, item)
	}

	return result
}

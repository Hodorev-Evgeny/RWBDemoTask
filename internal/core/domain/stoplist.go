package core_domain

import "sync"

type StopList struct {
	mu   sync.RWMutex
	List []string `json:"list"`
}

func NewStopList(l []string) *StopList {
	return &StopList{
		List: l,
	}
}

func (sl *StopList) Add(item string) {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	for _, v := range sl.List {
		if v == item {
			return
		}
	}

	sl.List = append(sl.List, item)
}

func (sl *StopList) Remove(item string) {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	for i, v := range sl.List {
		if v == item {
			sl.List = append(sl.List[:i], sl.List[i+1:]...)
			return
		}
	}
}

func (sl *StopList) IsBlocked(item string) bool {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	for _, v := range sl.List {
		if v == item {
			return true
		}
	}

	return false
}

func (sl *StopList) Items() []string {
	sl.mu.RLock()
	defer sl.mu.RUnlock()

	result := make([]string, len(sl.List))
	copy(result, sl.List)

	return result
}

package state

import (
	"sync"

	"github.com/looplab/fsm"
)

type State struct {
	mu sync.Mutex
	m  map[int64]*fsm.FSM
}

func NewState() *State {
	return &State{
		m: make(map[int64]*fsm.FSM),
	}
}

func (fsmMap *State) Add(key int64, fsm *fsm.FSM) {
	fsmMap.mu.Lock()
	defer fsmMap.mu.Unlock()
	fsmMap.m[key] = fsm
}

func (fsmMap *State) Get(key int64) (*fsm.FSM, bool) {
	fsmMap.mu.Lock()
	defer fsmMap.mu.Unlock()
	val, ok := fsmMap.m[key]
	return val, ok
}

func (fsmMap *State) Delete(key int64) {
	fsmMap.mu.Lock()
	defer fsmMap.mu.Unlock()
	delete(fsmMap.m, key)
}

package music

import "sync"

type VoteManager struct {
	mu    sync.RWMutex
	votes map[string]*Vote
}

func newVoterManager() *VoteManager {
	return &VoteManager{
		votes: make(map[string]*Vote),
	}
}

func (vm *VoteManager) Set(typ string, vote *Vote) {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	vm.votes[typ] = vote
}

func (vm *VoteManager) Get(typ string) *Vote {
	vm.mu.RLock()
	defer vm.mu.RUnlock()

	return vm.votes[typ]
}

func (vm *VoteManager) Delete(typ string) {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	delete(vm.votes, typ)
}

func (vm *VoteManager) Clear() {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	vm.votes = make(map[string]*Vote)
}

func (vm *VoteManager) Update(people int) {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	for key, vote := range vm.votes {
		if vote.Total()*4 < people*3 ||
			people <= 3 {
			vote.Callback()
			delete(vm.votes, key)
		}
	}
}

func (vm *VoteManager) RemoveLocals() {
	vm.mu.Lock()
	defer vm.mu.Unlock()

	for key, vote := range vm.votes {
		if !vote.Global {
			delete(vm.votes, key)
		}
	}
}

type Vote struct {
	mu       sync.RWMutex
	Global   bool
	Callback func()
	userSet  map[string]struct{}
}

func newVote() *Vote {
	return &Vote{
		userSet: make(map[string]struct{}),
	}
}

func (v *Vote) Approve(userID string) {
	v.mu.Lock()
	defer v.mu.Unlock()

	v.userSet[userID] = struct{}{}
}

func (v *Vote) Total() int {
	v.mu.RLock()
	defer v.mu.RUnlock()

	return len(v.userSet)
}

package cache

import "sync"

type syncSign map[string]bool

var (
	SyncSignMap syncSign
	locker      sync.Mutex
)

func init() {
	SyncSignMap = make(syncSign, 0)
}

func (m syncSign) Read(contractAddr string) (bool, bool) {
	locker.Lock()
	b1, b2 := m[contractAddr]
	locker.Unlock()
	return b1, b2
}

func (m syncSign) Insert(contractAddr string) {
	locker.Lock()
	m[contractAddr] = true
	locker.Unlock()
}

func (m syncSign) Delete(contractAddr string) {
	locker.Lock()
	delete(m, contractAddr)
	locker.Unlock()
}

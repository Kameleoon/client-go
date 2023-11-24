package storage

import (
	"time"

	cmap "github.com/orcaman/concurrent-map/v2"
)

type VisitorManager interface {
	GetVisitor(visitorCode string) Visitor
	GetOrCreateVisitor(visitorCode string) Visitor

	Enumerate(f func(string, Visitor) bool)
	Len() int

	Close()
}

type VisitorManagerImpl struct {
	visitors         cmap.ConcurrentMap[string, *VisitorImpl]
	expirationPeriod time.Duration
	purgeTicker      *time.Ticker
	stopChan         chan struct{}
}

func NewVisitorManagerImpl(expirationPeriod time.Duration) *VisitorManagerImpl {
	vm := &VisitorManagerImpl{
		visitors:         cmap.New[*VisitorImpl](),
		expirationPeriod: expirationPeriod,
		purgeTicker:      time.NewTicker(expirationPeriod),
		stopChan:         make(chan struct{}),
	}
	go func() {
		for {
			select {
			case <-vm.purgeTicker.C:
				vm.purge()
			case <-vm.stopChan:
				return
			}
		}
	}()
	return vm
}

func (vm *VisitorManagerImpl) ExpirationPeriod() time.Duration {
	return vm.expirationPeriod
}

func (vm *VisitorManagerImpl) Close() {
	vm.stop()
}
func (vm *VisitorManagerImpl) stop() {
	vm.purgeTicker.Stop()
	vm.stopChan <- struct{}{}
}

func (vm *VisitorManagerImpl) GetVisitor(visitorCode string) Visitor {
	// It is essential to update a visitor's last activity time before the visitor can be removed.
	// However, the used map type `cmap.ConcurrentMap` does not provide a "tryGet" method
	// with callback support. That is the reason why `RemoveCb` method is used as a "tryGet".
	var visitor Visitor
	vm.visitors.RemoveCb(visitorCode, func(vc string, v *VisitorImpl, exists bool) bool {
		if v != nil {
			v.UpdateLastActivityTime()
			visitor = v
		}
		return false
	})
	return visitor
}
func (vm *VisitorManagerImpl) GetOrCreateVisitor(visitorCode string) Visitor {
	return vm.visitors.Upsert(visitorCode, nil, func(exist bool, former, _ *VisitorImpl) *VisitorImpl {
		if former != nil {
			former.UpdateLastActivityTime()
			return former
		}
		return NewVisitorImpl()
	})
}

func (vm *VisitorManagerImpl) Enumerate(f func(string, Visitor) bool) {
	for kv := range vm.visitors.IterBuffered() {
		if !f(kv.Key, kv.Val) {
			return
		}
	}
}
func (vm *VisitorManagerImpl) Len() int {
	return vm.visitors.Count()
}

func (vm *VisitorManagerImpl) purge() {
	expiredDT := time.Now().Add(-vm.expirationPeriod)
	var vrs []struct {
		vc string
		v  *VisitorImpl
	}
	vm.visitors.IterCb(func(vc string, v *VisitorImpl) {
		if v.LastActivityTime().Before(expiredDT) {
			vrs = append(vrs, struct {
				vc string
				v  *VisitorImpl
			}{vc: vc, v: v})
		}
	})
	for _, vr := range vrs {
		vm.visitors.RemoveCb(vr.vc, func(key string, v *VisitorImpl, exists bool) bool {
			return v.LastActivityTime().Before(expiredDT)
		})
	}
}

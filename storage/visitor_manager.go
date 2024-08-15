package storage

import (
	"time"

	"github.com/Kameleoon/client-go/v3/logging"
	"github.com/Kameleoon/client-go/v3/managers/data"
	"github.com/Kameleoon/client-go/v3/types"
	cmap "github.com/orcaman/concurrent-map/v2"
)

type VisitorManager interface {
	GetVisitor(visitorCode string) Visitor
	GetOrCreateVisitor(visitorCode string) Visitor

	AddData(visitorCode string, data ...types.Data) Visitor

	Enumerate(f func(string, Visitor) bool)
	Len() int

	Clear()

	Close()
}

type VisitorManagerImpl struct {
	dataManager      data.DataManager
	visitors         cmap.ConcurrentMap[string, *VisitorImpl]
	expirationPeriod time.Duration
	purgeTicker      *time.Ticker
	stopChan         chan struct{}
}

func NewVisitorManagerImpl(
	dataManager data.DataManager, expirationPeriod time.Duration,
) *VisitorManagerImpl {
	logging.Debug("CALL: NewVisitorManagerImpl(expirationPeriod: %s)", expirationPeriod)
	vm := &VisitorManagerImpl{
		dataManager:      dataManager,
		visitors:         cmap.New[*VisitorImpl](),
		expirationPeriod: expirationPeriod,
		purgeTicker:      time.NewTicker(expirationPeriod),
		stopChan:         make(chan struct{}, 8),
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
	logging.Debug("RETURN: NewVisitorManagerImpl(expirationPeriod: %s) -> (VisitorManagerImpl)",
		expirationPeriod)
	return vm
}

func (vm *VisitorManagerImpl) ExpirationPeriod() time.Duration {
	return vm.expirationPeriod
}

func (vm *VisitorManagerImpl) Close() {
	vm.stop()
}
func (vm *VisitorManagerImpl) stop() {
	logging.Debug("CALL: VisitorManagerImpl.stop()")
	vm.purgeTicker.Stop()
	if len(vm.stopChan) == 0 {
		vm.stopChan <- struct{}{}
	}
	logging.Debug("RETURN: VisitorManagerImpl.stop()")
}

func (vm *VisitorManagerImpl) GetVisitor(visitorCode string) Visitor {
	// It is essential to update a visitor's last activity time before the visitor can be removed.
	// However, the used map type `cmap.ConcurrentMap` does not provide a "tryGet" method
	// with callback support. That is the reason why `RemoveCb` method is used as a "tryGet".
	logging.Debug("CALL: VisitorManagerImpl.GetVisitor(visitorCode: %s)", visitorCode)
	var visitor Visitor
	vm.visitors.RemoveCb(visitorCode, func(vc string, v *VisitorImpl, exists bool) bool {
		if v != nil {
			v.UpdateLastActivityTime()
			visitor = v
		}
		return false
	})
	logging.Debug("RETURN: VisitorManagerImpl.GetVisitor(visitorCode: %s) -> (visitor: %s)",
		visitorCode, visitor)
	return visitor
}

func (vm *VisitorManagerImpl) GetOrCreateVisitor(visitorCode string) Visitor {
	return vm.getOrCreateVisitor(visitorCode)
}
func (vm *VisitorManagerImpl) getOrCreateVisitor(visitorCode string) *VisitorImpl {
	logging.Debug("CALL: VisitorManagerImpl.getOrCreateVisitor(visitorCode: %s)", visitorCode)
	visitor := vm.visitors.Upsert(visitorCode, nil, func(exist bool, former, _ *VisitorImpl) *VisitorImpl {
		if former != nil {
			former.UpdateLastActivityTime()
			return former
		}
		return NewVisitorImpl()
	})
	logging.Debug("RETURN: VisitorManagerImpl.getOrCreateVisitor(visitorCode: %s) -> (visitor)", visitorCode)
	return visitor
}

func (vm *VisitorManagerImpl) AddData(visitorCode string, data ...types.Data) Visitor {
	logging.Debug("CALL: VisitorManagerImpl.AddData(visitorCode: %s, data: %s)", visitorCode, data)
	visitor := vm.getOrCreateVisitor(visitorCode)
	cdi := vm.dataManager.DataFile().CustomDataInfo()
	if cdi != nil {
		for i, d := range data {
			if cd, ok := d.(*types.CustomData); ok {
				data[i] = vm.handleCustomData(visitorCode, visitor, cdi, cd)
			}
		}
	}
	visitor.AddData(data...)
	logging.Debug("RETURN: VisitorManagerImpl.AddData(visitorCode: %s, data: %s) -> (visitor)", visitorCode, data)
	return visitor
}

func (vm *VisitorManagerImpl) handleCustomData(
	visitorCode string,
	visitor *VisitorImpl,
	cdi *types.CustomDataInfo,
	cd *types.CustomData,
) types.Data {
	// We shouldn't send custom data with local only type
	if (cdi != nil) && cdi.IsLocalOnly(cd.ID()) {
		cd.MarkAsSent()
	}
	// If mappingIdentifier is passed, we should link anonymous visitor with real unique userId.
	// After authorization, customer must be able to continue work with userId, but hash for variation
	// should be calculated based on anonymous visitor code, that's why set MappingIdentifier to visitor.
	if isMappingIdentifier(cdi, cd) {
		visitor.SetMappingIdentifier(&visitorCode)
		userId := cd.Values()[0]
		if visitorCode != userId {
			vm.visitors.Set(userId, cloneVisitorImpl(visitor))
			logging.Info("Linked anonymous visitor '%s' with user '%s'", visitorCode, userId)
		}
		return types.NewMappingIdentifier(cd)
	}
	return cd
}

func isMappingIdentifier(cdi *types.CustomDataInfo, cd types.ICustomData) bool {
	return (cdi != nil) && (cd.ID() == cdi.MappingIdentifierIndex()) && (len(cd.Values()) > 0) && (cd.Values()[0] != "")
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
	logging.Debug("CALL: VisitorManagerImpl.purge()")
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
	logging.Debug("RETURN: VisitorManagerImpl.purge()")
}

func (vm *VisitorManagerImpl) Clear() {
	logging.Debug("CALL: VisitorManagerImpl.Clear()")
	vm.visitors.Clear()
	logging.Debug("RETURN: VisitorManagerImpl.Clear()")
}

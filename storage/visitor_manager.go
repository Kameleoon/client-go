package storage

import (
	"time"

	"github.com/Kameleoon/client-go/v3/logging"
	"github.com/Kameleoon/client-go/v3/types"
	cmap "github.com/orcaman/concurrent-map/v2"
)

type VisitorManager interface {
	CustomDataInfo() *types.CustomDataInfo
	SetCustomDataInfo(value *types.CustomDataInfo)

	GetVisitor(visitorCode string) Visitor
	GetOrCreateVisitor(visitorCode string) Visitor

	AddData(visitorCode string, data ...types.Data) Visitor

	Enumerate(f func(string, Visitor) bool)
	Len() int

	Clear()

	Close()
}

type VisitorManagerImpl struct {
	customDataInfo   *types.CustomDataInfo
	visitors         cmap.ConcurrentMap[string, *VisitorImpl]
	expirationPeriod time.Duration
	purgeTicker      *time.Ticker
	stopChan         chan struct{}
	logger           logging.Logger
}

func NewVisitorManagerImpl(expirationPeriod time.Duration, logger logging.Logger) *VisitorManagerImpl {
	vm := &VisitorManagerImpl{
		visitors:         cmap.New[*VisitorImpl](),
		expirationPeriod: expirationPeriod,
		purgeTicker:      time.NewTicker(expirationPeriod),
		stopChan:         make(chan struct{}),
		logger:           logger,
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

func (vm *VisitorManagerImpl) CustomDataInfo() *types.CustomDataInfo {
	return vm.customDataInfo
}
func (vm *VisitorManagerImpl) SetCustomDataInfo(value *types.CustomDataInfo) {
	vm.customDataInfo = value
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
	return vm.getOrCreateVisitor(visitorCode)
}
func (vm *VisitorManagerImpl) getOrCreateVisitor(visitorCode string) *VisitorImpl {
	return vm.visitors.Upsert(visitorCode, nil, func(exist bool, former, _ *VisitorImpl) *VisitorImpl {
		if former != nil {
			former.UpdateLastActivityTime()
			return former
		}
		return NewVisitorImpl()
	})
}

func (vm *VisitorManagerImpl) AddData(visitorCode string, data ...types.Data) Visitor {
	visitor := vm.getOrCreateVisitor(visitorCode)
	cdi := vm.customDataInfo
	if cdi != nil {
		for _, d := range data {
			if cd, ok := d.(*types.CustomData); ok {
				vm.handleCustomData(visitorCode, visitor, cdi, cd)
			}
		}
	}
	visitor.AddData(vm.logger, data...)
	return visitor
}
func (vm *VisitorManagerImpl) handleCustomData(
	visitorCode string,
	visitor *VisitorImpl,
	cdi *types.CustomDataInfo,
	cd *types.CustomData,
) {
	// We shouldn't send custom data with local only type
	if cdi.IsLocalOnly(cd.ID()) {
		cd.MarkAsSent()
	}
	// If mappingIdentifier is passed, we should link anonymous visitor with real unique userId.
	// After authorization, customer must be able to continue work with userId, but hash for variation
	// should be calculated based on anonymous visitor code, that's why set MappingIdentifier to visitor.
	if cdi.IsMappingIdentifier(cd.ID()) && (len(cd.Values()) > 0) {
		targetVisitorCode := cd.Values()[0]
		if targetVisitorCode != "" {
			cd.SetIsMappingIdentifier(true)
			visitor.SetMappingIdentifier(&visitorCode)
			if visitorCode != targetVisitorCode {
				vm.visitors.Set(targetVisitorCode, visitor)
			}
		}
	}
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

func (vm *VisitorManagerImpl) Clear() {
	vm.visitors.Clear()
}

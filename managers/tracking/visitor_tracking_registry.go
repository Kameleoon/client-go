package tracking

import (
	"sync"

	"github.com/Kameleoon/client-go/v3/storage"
	cmap "github.com/orcaman/concurrent-map/v2"
)

type VisitorTrackingRegistry interface {
	Add(visitorCode string)
	AddAll(visitorCodes []string)
	Extract() VisitorCodeCollection
}

// Rwmx VisitorTrackingRegistry with ConcurrentMap

const (
	DefaultStorageLimit                           = 1_0000_00
	DefaultExtractionLimit                        = 20_000
	LimitedExtractionThresholdCoefficient         = 2
	RemovalFactor                         float64 = 0.8
)

type RwmxCMapVisitorTrackingRegistry struct {
	visitorManager  storage.VisitorManager
	storageLimit    int
	extractionLimit int
	mutex           sync.RWMutex
	visitors        *cmap.ConcurrentMap[string, struct{}]
}

func NewRwmxCMapVisitorTrackingRegistry(
	visitorManager storage.VisitorManager, storageLimit int, extractionLimit int,
) *RwmxCMapVisitorTrackingRegistry {
	visitors := cmap.New[struct{}]()
	return &RwmxCMapVisitorTrackingRegistry{
		visitorManager:  visitorManager,
		storageLimit:    storageLimit,
		extractionLimit: extractionLimit,
		visitors:        &visitors,
	}
}

func (vtr *RwmxCMapVisitorTrackingRegistry) Add(visitorCode string) {
	vtr.mutex.RLock()
	defer vtr.mutex.RUnlock()
	vtr.visitors.Set(visitorCode, struct{}{})
}

func (vtr *RwmxCMapVisitorTrackingRegistry) AddAll(visitorCodes []string) {
	vtr.mutex.RLock()
	for _, visitorCode := range visitorCodes {
		vtr.visitors.Set(visitorCode, struct{}{})
	}
	vtr.mutex.RUnlock()
	if vtr.visitors.Count() > vtr.storageLimit {
		vtr.mutex.Lock()
		defer vtr.mutex.Unlock()
		vtr.eraseNonexistentVisitors()
		vtr.eraseToStorageLimit()
	}
}

// Not thread-safe
func (vtr *RwmxCMapVisitorTrackingRegistry) eraseNonexistentVisitors() {
	var visitorsToRemove []string
	vtr.visitors.IterCb(func(vc string, v struct{}) {
		if vtr.visitorManager.GetVisitor(vc) == nil {
			visitorsToRemove = append(visitorsToRemove, vc)
		}
	})
	for _, vc := range visitorsToRemove {
		vtr.visitors.Remove(vc)
	}
}

// Not thread-safe
func (vtr *RwmxCMapVisitorTrackingRegistry) eraseToStorageLimit() {
	visitorsToRemoveCount := vtr.visitors.Count() - int(float64(vtr.storageLimit)*RemovalFactor)
	if visitorsToRemoveCount <= 0 {
		return
	}
	visitorsToRemove := vtr.visitors.Keys()[:visitorsToRemoveCount]
	for _, vc := range visitorsToRemove {
		vtr.visitors.Remove(vc)
	}
}

func (vtr *RwmxCMapVisitorTrackingRegistry) Extract() VisitorCodeCollection {
	if vtr.shouldExtractAllBeUsed() {
		return vtr.extractAll(true)
	}
	return vtr.extractLimited()
}

func (vtr *RwmxCMapVisitorTrackingRegistry) shouldExtractAllBeUsed() bool {
	return vtr.visitors.Count() < vtr.extractionLimit*LimitedExtractionThresholdCoefficient
}

func (vtr *RwmxCMapVisitorTrackingRegistry) extractAll(lock bool) VisitorCodeCollection {
	newVisitors := cmap.New[struct{}]()
	if lock {
		vtr.mutex.Lock()
	}
	oldVisitors := vtr.visitors
	vtr.visitors = &newVisitors
	if lock {
		vtr.mutex.Unlock()
	}
	return CMapVisitorCodeCollection{visitorCodes: oldVisitors}
}

func (vtr *RwmxCMapVisitorTrackingRegistry) extractLimited() VisitorCodeCollection {
	vtr.mutex.Lock()
	defer vtr.mutex.Unlock()
	if vtr.shouldExtractAllBeUsed() {
		return vtr.extractAll(false)
	}
	extracted := vtr.visitors.Keys()[:vtr.extractionLimit]
	for _, vc := range extracted {
		vtr.visitors.Remove(vc)
	}
	return SliceVisitorCodeCollection{visitorCodes: extracted}
}

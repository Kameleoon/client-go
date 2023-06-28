package hybrid

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Kameleoon/client-go/v2/logging"
	"github.com/Kameleoon/client-go/v2/storage"
)

const (
	tcInit                  = "window.kameleoonQueue=window.kameleoonQueue||[];"
	tcAssignVariationFormat = "window.kameleoonQueue.push(['Experiments.assignVariation',%d,%d]);"
	tcTriggerFormat         = "window.kameleoonQueue.push(['Experiments.trigger',%d,true]);"
)

// Represents a tool type that is supposed to manage Hybrid integration.
type HybridManager interface {
	// Assigns a variation for an experiment of a visitor.
	AddVariation(visitorCode string, experimentId int, variationId int)
	// Generates an Engine Tracking Code based on assigned variations of the specified visitor.
	GetEngineTrackingCode(visitor string) string
}

type HybridManagerImpl struct {
	sync.Mutex
	expirationTime time.Duration
	cacheFactory   storage.CacheFactory
	cache          storage.Cache
	logger         logging.Logger
}

func NewHybridManagerImpl(
	expirationTime time.Duration,
	cacheFactory storage.CacheFactory,
	logger logging.Logger) (*HybridManagerImpl, error) {

	if mainCache, err := cacheFactory.Create(expirationTime, true); err == nil {
		manager := &HybridManagerImpl{
			expirationTime: expirationTime,
			cacheFactory:   cacheFactory,
			cache:          mainCache,
			logger:         logger,
		}
		manager.log("hybridManager was successfully initialized")
		return manager, nil
	} else {
		if logger != nil {
			logger.Printf("hybridManager can't be initialized due error: %v", err)
		}
		return nil, err
	}
}

func (hm *HybridManagerImpl) AddVariation(visitorCode string, experimentId int, variationId int) {
	hm.Lock()
	defer hm.Unlock()
	// try to find cache and create if not found
	cache, exist := hm.cache.Get(visitorCode)
	visitorCache, ok := cache.(storage.Cache)
	var cacheError error
	if !(exist && ok) {
		if vc, err := hm.cacheFactory.Create(hm.expirationTime, false); err == nil {
			visitorCache = vc
		} else {
			cacheError = err
		}
	}
	// if cache is found or created, add experiment and variation ids
	if visitorCache != nil {
		visitorCache.Set(experimentId, variationId)
		hm.cache.Set(visitorCode, visitorCache)
		hm.log("hybridManager succesfully added variation for visitorCode: %s, experiment: %d, variation: %d",
			visitorCode, experimentId, variationId)
	} else {
		hm.log(
			"hybridManager failed to add variation for visitorCode: %s, experiment: %d, variation: %d; error: %s",
			visitorCode, experimentId, variationId, cacheError)
	}
}

func (hm *HybridManagerImpl) GetEngineTrackingCode(visitorCode string) string {
	var trackingCode strings.Builder
	trackingCode.WriteString(tcInit)
	if value, ok := hm.cache.Get(visitorCode); ok {
		if cache, ok := value.(storage.Cache); ok {
			for k, v := range cache.ActualValues() {
				trackingCode.WriteString(fmt.Sprintf(tcAssignVariationFormat, k, v))
				trackingCode.WriteString(fmt.Sprintf(tcTriggerFormat, k))
			}
		}
	}
	return trackingCode.String()
}

func (hm *HybridManagerImpl) log(format string, args ...interface{}) {
	if hm.logger != nil {
		hm.logger.Printf(format, args...)
	}
}

package hybrid

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Kameleoon/client-go/v3/logging"

	"github.com/Kameleoon/client-go/v3/storage"
	"github.com/Kameleoon/client-go/v3/types"
)

const (
	tcInit                  = "window.kameleoonQueue=window.kameleoonQueue||[];"
	tcAssignVariationFormat = "window.kameleoonQueue.push(['Experiments.assignVariation',%d,%d]);"
	tcTriggerFormat         = "window.kameleoonQueue.push(['Experiments.trigger',%d,true]);"
)

// Represents a tool type that is supposed to manage Hybrid integration.
type HybridManager interface {
	// Generates an Engine Tracking Code based on assigned variations of the specified visitor.
	GetEngineTrackingCode(variations storage.DataMapStorage[int, *types.AssignedVariation]) string
}

type HybridManagerImpl struct {
	expirationTime time.Duration
}

func NewHybridManagerImpl(expirationTime time.Duration) (*HybridManagerImpl, error) {
	logging.Debug("CALL: NewHybridManagerImpl(expirationTime: %s)", expirationTime)
	var err error
	var hybridManagerImpl *HybridManagerImpl
	if expirationTime <= 0 {
		err = errors.New("'expirationTime' must be a postitive value")
		logging.Error("HybridManager isn't initialized properly, "+
			"GetEngineTrackingCode method isn't available for call. error %s", err)
	} else {
		hybridManagerImpl = &HybridManagerImpl{expirationTime: expirationTime}
	}
	logging.Debug("RETURN: NewHybridManagerImpl(expirationTime: %s) -> (hybridManagerImpl, error: %s)",
		expirationTime, err)
	return hybridManagerImpl, err
}

func (hm *HybridManagerImpl) GetEngineTrackingCode(
	variations storage.DataMapStorage[int, *types.AssignedVariation],
) string {
	logging.Debug("CALL: HybridManagerImpl.GetEngineTrackingCode(variations: %s)", variations)
	var trackingCode strings.Builder
	trackingCode.WriteString(tcInit)
	if variations != nil {
		expiredTime := time.Now().Add(-hm.expirationTime)
		variations.Enumerate(func(av *types.AssignedVariation) bool {
			if av.AssignmentTime().After(expiredTime) {
				trackingCode.WriteString(fmt.Sprintf(tcAssignVariationFormat, av.ExperimentId(), av.VariationId()))
				trackingCode.WriteString(fmt.Sprintf(tcTriggerFormat, av.ExperimentId()))
			}
			return true
		})
	}
	trackingCodeString := trackingCode.String()
	logging.Debug("RETURN: HybridManagerImpl.GetEngineTrackingCode(variations: %s) -> (trackingCode: %s)",
		variations, trackingCodeString)
	return trackingCodeString
}

package hybrid

import (
	"errors"
	"fmt"
	"strings"
	"time"

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
	if expirationTime <= 0 {
		return nil, errors.New("'expirationTime' must be a postitive value")
	}
	return &HybridManagerImpl{
		expirationTime: expirationTime,
	}, nil
}

func (hm *HybridManagerImpl) GetEngineTrackingCode(
	variations storage.DataMapStorage[int, *types.AssignedVariation],
) string {
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
	return trackingCode.String()
}

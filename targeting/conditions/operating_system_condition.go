package conditions

import (
	"github.com/Kameleoon/client-go/v3/logging"
	"fmt"

	"github.com/Kameleoon/client-go/v3/errs"
	"github.com/Kameleoon/client-go/v3/types"
	"github.com/Kameleoon/client-go/v3/utils"
	"github.com/segmentio/encoding/json"
)

func NewOperatingSystemCondition(c types.TargetingCondition) *OperatingSystemCondition {
	osType, ok := types.ParseOperatingSystemType(c.OS)
	if !ok {
		logging.Error("Undefined OS for 'OS' condition: %s", c.OS)
	}
	return &OperatingSystemCondition{
		TargetingConditionBase: types.TargetingConditionBase{
			Type:    c.Type,
			Include: c.Include,
		},
		OsType: osType,
	}
}

type OperatingSystemCondition struct {
	types.TargetingConditionBase
	OsType types.OperatingSystemType
}

func (c *OperatingSystemCondition) UnmarshalJSON(data []byte) error {
	var cm = struct {
		OS string `json:"os,omitempty"`
	}{}
	if err := json.Unmarshal(data, &cm); err != nil {
		return err
	}
	osType, ok := types.ParseOperatingSystemType(cm.OS)
	c.OsType = osType
	if !ok {
		return errs.NewInternalError(fmt.Sprintf("undefined OS '%s' for OS condition", cm.OS))
	}
	return nil
}

func (c OperatingSystemCondition) CheckTargeting(targetData interface{}) bool {
	os, ok := targetData.(*types.OperatingSystem)
	return ok && (os != nil) && (os.Type() == c.OsType)
}

func (c *OperatingSystemCondition) String() string {
	return utils.JsonToString(c)
}

package conditions

import (
	"github.com/Kameleoon/client-go/v3/types"
	"github.com/Kameleoon/client-go/v3/utils"
)

func NewDeviceCondition(c types.TargetingCondition) *DeviceCondition {
	return &DeviceCondition{
		TargetingConditionBase: types.TargetingConditionBase{
			Type:    c.Type,
			Include: c.Include,
		},
		Device: c.Device,
	}
}

type DeviceCondition struct {
	types.TargetingConditionBase
	Device types.DeviceType `json:"device"`
}

func (c *DeviceCondition) CheckTargeting(targetData interface{}) bool {
	device, ok := targetData.(*types.Device)
	return ok && (device != nil) && (device.Type() == c.Device)
}

func (c *DeviceCondition) String() string {
	return utils.JsonToString(c)
}

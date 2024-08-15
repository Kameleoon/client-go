package configuration

import (
	"github.com/Kameleoon/client-go/v3/logging"
	"github.com/Kameleoon/client-go/v3/types"
	"fmt"
)

type Configuration struct {
	CustomDataInfo *types.CustomDataInfo `json:"customData"`
	Settings       Settings              `json:"configuration"`
	FeatureFlags   []FeatureFlag         `json:"featureFlags"`
}

func (c Configuration) String() string {
	return fmt.Sprintf("Configuration{CustomDataInfo:%v,Settings:%v,FeatureFlags:'%v'}",
		c.CustomDataInfo, c.Settings, logging.ObjectToString(c.FeatureFlags))
}

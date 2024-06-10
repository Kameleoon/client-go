package configuration

import "github.com/Kameleoon/client-go/v3/types"

type Configuration struct {
	CustomDataInfo *types.CustomDataInfo `json:"customData"`
	Settings       Settings              `json:"configuration"`
	FeatureFlags   []FeatureFlag         `json:"featureFlags"`
}

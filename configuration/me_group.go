package configuration

import (
	"sort"

	"github.com/Kameleoon/client-go/v3/types"
)

type MEGroup struct {
	featureFlags []types.IFeatureFlag
}

func NewMEGroup(featureFlags []types.IFeatureFlag) *MEGroup {
	sort.Slice(featureFlags, func(i, j int) bool {
		return featureFlags[i].GetId() < featureFlags[j].GetId()
	})
	return &MEGroup{featureFlags: featureFlags}
}

func (meg *MEGroup) FeatureFlags() []types.IFeatureFlag {
	return meg.featureFlags
}

func (meg *MEGroup) GetFeatureFlagByHash(hash float64) types.IFeatureFlag {
	idx := int(hash * float64(len(meg.featureFlags)))
	if idx >= len(meg.featureFlags) {
		idx = len(meg.featureFlags) - 1
	}
	return meg.featureFlags[idx]
}

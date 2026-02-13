package conditions

import (
	"github.com/Kameleoon/client-go/v3/logging"
	"github.com/Kameleoon/client-go/v3/types"
	"github.com/Kameleoon/client-go/v3/utils"
)

func NewVersionCondition(c types.TargetingCondition) *VersionCondition {
	return &VersionCondition{
		TargetingConditionBase: types.TargetingConditionBase{
			Type:    c.Type,
			Include: c.Include,
		},
		Version:          c.Version,
		VersionMatchType: c.VersionMatchType,
	}
}

type VersionCondition struct {
	types.TargetingConditionBase
	Version          string             `json:"version"`
	VersionMatchType types.OperatorType `json:"versionMatchType,omitempty"`
}

func (c *VersionCondition) CheckTargeting(targetData interface{}) bool {
	av, ok := targetData.(*types.ApplicationVersion)
	return ok && av != nil && c.CompareWithVersion(av.Version)
}

func (c *VersionCondition) CompareWithVersion(targetVersion string) bool {
	condition, err := utils.NewVersionFromString(c.Version)
	if err != nil {
		logging.Error("Failed to parse version '%s' for '%s' condition", c.Version, c.Type)
		return false
	}
	target, err := utils.NewVersionFromString(targetVersion)
	if err != nil {
		logging.Error("Failed to parse version '%s' for target in '%s' condition", targetVersion, c.Type)
		return false
	}

	cmp := target.CompareTo(condition)
	switch c.VersionMatchType {
	case types.OperatorEqual:
		return cmp == 0
	case types.OperatorGreater:
		return cmp > 0
	case types.OperatorLower:
		return cmp < 0
	default:
		logging.Error("Unexpected comparing operation for '%s' condition: %s", c.Type, c.VersionMatchType)
		return false
	}
}

func (c VersionCondition) String() string {
	return utils.JsonToString(c)
}

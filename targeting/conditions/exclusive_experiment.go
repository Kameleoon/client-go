package conditions

import (
	"github.com/Kameleoon/client-go/v3/logging"
	"github.com/Kameleoon/client-go/v3/storage"
	"github.com/Kameleoon/client-go/v3/types"
	"github.com/Kameleoon/client-go/v3/utils"
)

const (
	CampaignTypeExperiment      = "EXPERIMENT"
	CampaignTypePersonalization = "PERSONALIZATION"
	CampaignTypeAny             = "ANY"
)

type ExclusiveExperimentCondition struct {
	types.TargetingConditionBase
	campaignType string
}

func NewExclusiveExperimentCondition(c types.TargetingCondition) *ExclusiveExperimentCondition {
	return &ExclusiveExperimentCondition{
		TargetingConditionBase: types.TargetingConditionBase{
			Type:    c.Type,
			Include: true,
		},
		campaignType: c.CampaignType,
	}
}

func (c *ExclusiveExperimentCondition) CheckTargeting(targetData interface{}) bool {
	if targetingData, ok := targetData.(TargetingDataExclusiveExperiment); ok {
		switch c.campaignType {
		case CampaignTypeExperiment:
			return c.checkExperiment(targetingData)
		case CampaignTypePersonalization:
			return c.checkPersonalization(targetingData)
		case CampaignTypeAny:
			return c.checkPersonalization(targetingData) && c.checkExperiment(targetingData)
		}
		logging.Error("Unexpected campaign type for %s condition: %s", c.Type, c.campaignType)
	}
	return false
}

func (*ExclusiveExperimentCondition) checkExperiment(td TargetingDataExclusiveExperiment) bool {
	if td.Variations == nil {
		return true
	}
	variationCount := td.Variations.Len()
	return (variationCount == 0) || ((variationCount == 1) && (td.Variations.Get(td.CurrentExperimentId) != nil))
}

func (*ExclusiveExperimentCondition) checkPersonalization(td TargetingDataExclusiveExperiment) bool {
	return (td.Personalizations == nil) || (td.Personalizations.Len() == 0)
}

func (c ExclusiveExperimentCondition) String() string {
	return utils.JsonToString(c)
}

type TargetingDataExclusiveExperiment struct {
	CurrentExperimentId int
	Variations          storage.DataMapStorage[int, *types.AssignedVariation]
	Personalizations    storage.DataMapStorage[int, *types.Personalization]
}

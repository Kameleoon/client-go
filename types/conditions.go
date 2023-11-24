package types

import (
	"strings"

	"github.com/segmentio/encoding/json"
)

type Condition interface {
	GetType() TargetingType
	GetInclude() bool
	CheckTargeting(interface{}) bool
	String() string
}

type ConditionsData struct {
	FirstLevelOrOperators []bool                 `json:"firstLevelOrOperators"`
	FirstLevel            []ConditionsFirstLevel `json:"firstLevel"`
}

type ConditionsFirstLevel struct {
	OrOperators []bool               `json:"orOperators"`
	Conditions  []TargetingCondition `json:"conditions"`
}

// const (
// 	targetingConditionStaticFieldId                 = "id"
// 	targetingConditionStaticFieldValue              = "value"
// 	targetingConditionStaticFieldType               = "targetingType"
// 	targetingConditionStaticFieldOperator           = "valueMatchType"
// 	targetingConditionStaticFieldWeight             = "weight"
// 	targetingConditionStaticFieldIndex              = "customDataIndex"
// 	targetingConditionStaticFieldIsInclude          = "isInclude"
// 	targetingConditionStaticFieldExperiment         = "experiment"
// 	targetingConditionStaticFieldVariation          = "variation"
// 	targetingConditionStaticFieldVariationMatchType = "variationMatchType"
// 	targetingConditionStaticFieldBrowser            = "browser"
// )

type BrowserConditionType string

const (
	BrowserConditionTypeChrome  BrowserConditionType = "CHROME"
	BrowserConditionTypeIE      BrowserConditionType = "IE"
	BrowserConditionTypeFirefox BrowserConditionType = "FIREFOX"
	BrowserConditionTypeSafari  BrowserConditionType = "SAFARI"
	BrowserConditionTypeOpera   BrowserConditionType = "OPERA"
	BrowserConditionTypeOther   BrowserConditionType = "OTHER"
)

// var targetingConditionStaticFields = [...]string{targetingConditionStaticFieldValue,
// 	targetingConditionStaticFieldType, targetingConditionStaticFieldOperator,
// 	targetingConditionStaticFieldIndex, targetingConditionStaticFieldIsInclude, targetingConditionStaticFieldExperiment,
// 	targetingConditionStaticFieldVariation, targetingConditionStaticFieldVariationMatchType,
// 	targetingConditionStaticFieldBrowser, targetingConditionStaticFieldVersionMatchType}

type TargetingConditionBase struct {
	Type    TargetingType `json:"targetingType"`
	Include bool          `json:"isInclude,omitempty"`
}

type TargetingCondition struct {
	TargetingConditionBase
	Value              interface{}          `json:"value,omitempty"`
	Operator           OperatorType         `json:"valueMatchType,omitempty"`
	Index              string               `json:"customDataIndex,omitempty"`
	Experiment         int                  `json:"experiment,omitempty"`
	Variation          int                  `json:"variation,omitempty"`
	VariationMatchType OperatorType         `json:"variationMatchType,omitempty"`
	VersionMatchType   OperatorType         `json:"versionMatchType,omitempty"`
	Browser            BrowserConditionType `json:"browser,omitempty"`
	Version            string               `json:"version,omitempty"`
	Device             DeviceType           `json:"device,omitempty"`
	VisitorCode        string               `json:"visitorCode,omitempty"`
	MatchType          OperatorType         `json:"matchType,omitempty"`
	SdkLanguage        string               `json:"sdkLanguage,omitempty"`
	Title              string               `json:"title,omitempty"`
	Url                string               `json:"url,omitempty"`
	GoalId             int                  `json:"goalId,omitempty"`
}

// func (c *TargetingCondition) UnmarshalJSON(b []byte) error {
// 	rest := make(map[string]json.RawMessage)
// 	if err := json.Unmarshal(b, &rest); err != nil {
// 		return err
// 	}
// 	var err error
// 	for _, field := range targetingConditionStaticFields {
// 		value, exist := rest[field]
// 		if !exist {
// 			continue
// 		}
// 		switch field {
// 		case targetingConditionStaticFieldType:
// 			err = json.Unmarshal(value, &c.TargetingConditionBase.Type)
// 		case targetingConditionStaticFieldValue:
// 			err = json.Unmarshal(value, &c.Value)
// 		case targetingConditionStaticFieldOperator:
// 			err = json.Unmarshal(value, &c.Operator)
// 		case targetingConditionStaticFieldIndex:
// 			err = json.Unmarshal(value, &c.Index)
// 		case targetingConditionStaticFieldIsInclude:
// 			err = json.Unmarshal(value, &c.TargetingConditionBase.Include)
// 		case targetingConditionStaticFieldExperiment:
// 			err = json.Unmarshal(value, &c.Experiment)
// 		case targetingConditionStaticFieldVariation:
// 			err = json.Unmarshal(value, &c.Variation)
// 		case targetingConditionStaticFieldVariationMatchType:
// 			err = json.Unmarshal(value, &c.VariationMatchType)
// 		case targetingConditionStaticFieldBrowser:
// 			err = json.Unmarshal(value, &c.Browser)
// 		}
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	return nil
// }

func (c *TargetingConditionBase) String() string {
	var s strings.Builder
	b, _ := json.Marshal(c)
	s.Write(b)
	return s.String()
}

func (c TargetingConditionBase) GetType() TargetingType {
	return c.Type
}

func (c TargetingConditionBase) GetInclude() bool {
	return c.Include
}

func (c TargetingConditionBase) CheckTargeting(_ []Data) bool {
	return true
}

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

type BrowserConditionType string

const (
	BrowserConditionTypeChrome  BrowserConditionType = "CHROME"
	BrowserConditionTypeIE      BrowserConditionType = "IE"
	BrowserConditionTypeFirefox BrowserConditionType = "FIREFOX"
	BrowserConditionTypeSafari  BrowserConditionType = "SAFARI"
	BrowserConditionTypeOpera   BrowserConditionType = "OPERA"
	BrowserConditionTypeOther   BrowserConditionType = "OTHER"
)

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

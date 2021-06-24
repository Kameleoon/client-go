package types

import (
	"strings"

	"github.com/segmentio/encoding/json"
)

type Condition interface {
	GetType() TargetingType
	SetType(TargetingType)
	GetInclude() bool
	SetInclude(bool)
	CheckTargeting([]TargetingData) bool
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

const (
	targetingConditionStaticFieldId       = "id"
	targetingConditionStaticFieldValue    = "value"
	targetingConditionStaticFieldType     = "targetingType"
	targetingConditionStaticFieldOperator = "valueMatchType"
	targetingConditionStaticFieldWeight   = "weight"
	targetingConditionStaticFieldIndex    = "customDataIndex"
	targetingConditionStaticFieldInclude  = "include"
)

var targetingConditionStaticFields = [...]string{targetingConditionStaticFieldId, targetingConditionStaticFieldValue,
	targetingConditionStaticFieldType, targetingConditionStaticFieldOperator, targetingConditionStaticFieldWeight,
	targetingConditionStaticFieldIndex, targetingConditionStaticFieldInclude}

type TargetingCondition struct {
	Rest     map[string]json.RawMessage `json:"-"`
	Value    interface{}                `json:"value,omitempty"`
	Type     TargetingType              `json:"targetingType"`
	Operator OperatorType               `json:"valueMatchType,omitempty"`
	Index    string                     `json:"customDataIndex,omitempty"`
	ID       int                        `json:"id"`
	Weight   int                        `json:"weight,omitempty"`
	Include  *bool                      `json:"include,omitempty"`
}

func (c *TargetingCondition) UnmarshalJSON(b []byte) error {
	c.Rest = make(map[string]json.RawMessage)
	if err := json.Unmarshal(b, &c.Rest); err != nil {
		return err
	}
	var value json.RawMessage
	var exist bool
	var err error
	for _, field := range targetingConditionStaticFields {
		value, exist = c.Rest[field]
		if !exist {
			continue
		}
		delete(c.Rest, field)
		switch field {
		case targetingConditionStaticFieldType:
			err = json.Unmarshal(value, &c.Type)
		case targetingConditionStaticFieldId:
			err = json.Unmarshal(value, &c.ID)
		case targetingConditionStaticFieldValue:
			err = json.Unmarshal(value, &c.Value)
		case targetingConditionStaticFieldOperator:
			err = json.Unmarshal(value, &c.Operator)
		case targetingConditionStaticFieldWeight:
			err = json.Unmarshal(value, &c.Weight)
		case targetingConditionStaticFieldIndex:
			err = json.Unmarshal(value, &c.Index)
		case targetingConditionStaticFieldInclude:
			err = json.Unmarshal(value, &c.Include)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *TargetingCondition) String() string {
	var s strings.Builder
	b, _ := json.Marshal(c)
	s.Write(b)
	return s.String()
}

func (c TargetingCondition) GetType() TargetingType {
	return c.Type
}

func (c *TargetingCondition) SetType(tt TargetingType) {
	c.Type = tt
}

func (c TargetingCondition) GetInclude() bool {
	if c.Include == nil {
		return true
	}
	return *c.Include
}

func (c *TargetingCondition) SetInclude(i bool) {
	c.Include = &i
}

func (c TargetingCondition) CheckTargeting(_ []TargetingData) bool {
	return true
}

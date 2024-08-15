package conditions

import (
	"github.com/Kameleoon/client-go/v3/logging"
	"regexp"
	"strings"

	"github.com/Kameleoon/client-go/v3/types"
	"github.com/Kameleoon/client-go/v3/utils"
)

func NewCookieCondition(c types.TargetingCondition) *CookieCondition {
	value, valueCast := c.Value.(string)
	if !valueCast {
		value = ""
	}
	return &CookieCondition{
		TargetingConditionBase: types.TargetingConditionBase{
			Type:    c.Type,
			Include: c.Include,
		},
		ConditionName:  c.Name,
		NameMatchType:  c.NameMatchType,
		ConditionValue: value,
		ValueMatchType: c.Operator,
	}
}

type CookieCondition struct {
	types.TargetingConditionBase
	ConditionName  string             `json:"name,omitempty"`
	NameMatchType  types.OperatorType `json:"nameMatchType,omitempty"`
	ConditionValue string             `json:"value,omitempty"`
	ValueMatchType types.OperatorType `json:"valueMatchType,omitempty"`
}

func (c *CookieCondition) CheckTargeting(targetData interface{}) bool {
	cookie, ok := targetData.(*types.Cookie)
	if ok && (cookie != nil) {
		return c.checkValues(c.selectValues(cookie))
	}
	return false
}

func (c *CookieCondition) selectValues(cookie *types.Cookie) []string {
	var values []string
	switch c.NameMatchType {
	case types.OperatorExact:
		if value, ok := cookie.Cookies()[c.ConditionName]; ok {
			values = append(values, value)
		}
	case types.OperatorContains:
		for name, value := range cookie.Cookies() {
			if strings.Contains(name, c.ConditionName) {
				values = append(values, value)
			}
		}
	case types.OperatorRegExp:
		for name, value := range cookie.Cookies() {
			matched, err := regexp.MatchString(c.ConditionName, name)
			if (err == nil) && matched {
				values = append(values, value)
			}
		}
	default:
		logging.Error("Unexpected comparing operation for 'Cookie' condition (name): %s", c.NameMatchType)
	}
	return values
}

func (c *CookieCondition) checkValues(values []string) bool {
	switch c.ValueMatchType {
	case types.OperatorExact:
		for _, v := range values {
			if v == c.ConditionValue {
				return true
			}
		}
	case types.OperatorContains:
		for _, v := range values {
			if strings.Contains(v, c.ConditionValue) {
				return true
			}
		}
	case types.OperatorRegExp:
		for _, v := range values {
			matched, err := regexp.MatchString(c.ConditionValue, v)
			if (err == nil) && matched {
				return true
			}
		}
	default:
		logging.Error("Unexpected comparing operation for 'Cookie' condition (value): %s", c.ValueMatchType)
	}
	return false
}

func (c CookieCondition) String() string {
	return utils.JsonToString(c)
}

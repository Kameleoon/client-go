package conditions

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/Kameleoon/client-go/v3/types"
	"github.com/Kameleoon/client-go/v3/utils"
)

type StringValueCondition struct {
	types.TargetingConditionBase
	Value         string             `json:"value"`
	MatchType     types.OperatorType `json:"matchType"`
	ConditionType string             `json:"conditionType"`
}

func (c *StringValueCondition) CheckTargeting(targetData interface{}) bool {
	value, ok := targetData.(string)
	return ok && c.checkTargeting(value)
}

func (c *StringValueCondition) checkTargeting(value string) bool {
	switch c.MatchType {
	case types.OperatorExact:
		return value == c.Value
	case types.OperatorContains:
		return strings.Contains(value, c.Value)
	case types.OperatorRegExp:
		matched, err := regexp.MatchString(c.Value, value)
		return err == nil && matched
	default:
		fmt.Printf("unexpected comparing operation for %v condition: %v\n", c.TargetingConditionBase.Type, c.MatchType)
		return false
	}
}

func (c *StringValueCondition) String() string {
	return utils.JsonToString(c)
}

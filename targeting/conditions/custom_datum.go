package conditions

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Kameleoon/client-go/v2/types"
	"github.com/Kameleoon/client-go/v2/utils"
	"github.com/segmentio/encoding/json"
)

func NewCustomDatum(c types.TargetingCondition) *CustomDatum {
	if c.Value == nil {
		c.Value = ""
	}
	return &CustomDatum{
		TargetingConditionBase: c.TargetingConditionBase,
		Operator:               c.Operator,
		Value:                  c.Value,
		Index:                  c.Index,
		Include:                c.Include,
	}
}

type CustomDatum struct {
	types.TargetingConditionBase
	Value    interface{}         `json:"value"`
	Type     types.TargetingType `json:"targetingType"`
	Operator types.OperatorType  `json:"valueMatchType"`
	Index    string              `json:"customDataIndex"`
	Include  bool                `json:"include"`
}

func (c *CustomDatum) CheckTargeting(targetData interface{}) bool {
	if _, ok := targetData.([]types.TargetingData); !ok {
		return false
	}
	var arrayCustomData []*types.CustomData
	for _, td := range targetData.([]types.TargetingData) {
		if td.Data.DataType() != types.DataTypeCustom {
			continue
		}
		custom, ok := td.Data.(*types.CustomData)
		if ok && custom.ID == c.Index {
			arrayCustomData = append(arrayCustomData, custom)
		}
	}
	if len(arrayCustomData) == 0 {
		return c.Operator == types.OperatorUndefined
	}
	customData := arrayCustomData[len(arrayCustomData)-1]
	return c.checkTargeting(customData.GetValues())
}

func (c *CustomDatum) checkTargeting(customDataValues []string) bool {
	switch c.Operator {
	case types.OperatorContains:
		value, ok := c.Value.(string)
		if !ok {
			return false
		}
		return c.contains(customDataValues, func(customDataValue string) bool {
			return strings.Contains(customDataValue, value)
		})
	case types.OperatorExact:
		return c.contains(customDataValues, func(customDataValue string) bool {
			return customDataValue == c.Value
		})
	case types.OperatorRegExp:
		pattern, ok := c.Value.(string)
		if !ok {
			return false
		}
		return c.contains(customDataValues, func(customDataValue string) bool {
			matched, err := regexp.MatchString(pattern, customDataValue)
			return err == nil && matched
		})
	case types.OperatorLower, types.OperatorGreater, types.OperatorEqual:
		var number float64
		switch v := c.Value.(type) {
		case string:
			if val, err := strconv.ParseFloat(v, 64); err == nil {
				number = val
			} else {
				return false
			}
		case int:
			number = float64(v)
		case float64:
			number = v
		default:
			return false
		}
		return c.contains(customDataValues, func(customDataValue string) bool {
			if value, err := strconv.ParseFloat(customDataValue, 64); err == nil {
				switch c.Operator {
				case types.OperatorLower:
					return value < number
				case types.OperatorEqual:
					return value == number
				case types.OperatorGreater:
					return value > number
				}
			}
			return false
		})
	case types.OperatorIsTrue:
		return c.contains(customDataValues, func(customDataValue string) bool {
			val, err := strconv.ParseBool(customDataValue)
			return err == nil && val
		})
	case types.OperatorIsFalse:
		return c.contains(customDataValues, func(customDataValue string) bool {
			val, err := strconv.ParseBool(customDataValue)
			return err == nil && !val
		})
	case types.OperatorIsAmongValues:
		var values []interface{}
		if err := json.Unmarshal([]byte(c.Value.(string)), &values); err == nil {
			mapValues := utils.MapToStringDict(values, func(dict *map[string]interface{}, element interface{}) {
				elementString := fmt.Sprintf("%v", element)
				(*dict)[elementString] = true
			})
			return c.contains(customDataValues, func(customDataValue string) bool {
				_, exist := mapValues[customDataValue]
				return exist
			})
		} else {
			fmt.Println(err)
		}
	}
	return false
}

func (c *CustomDatum) contains(customDataValues []string, callback func(string) bool) bool {
	for _, val := range customDataValues {
		if callback(val) {
			return true
		}
	}
	return false
}

func (c *CustomDatum) String() string {
	return utils.JsonToString(c)
}

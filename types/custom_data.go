package types

import (
	"fmt"
	"strings"

	"github.com/Kameleoon/client-go/v2/network"
	"github.com/Kameleoon/client-go/v2/utils"
)

const customDataEventType = "customData"

// TODO: remove Value in next major version and make Values public
// It's need to have backward compatibility
type CustomData struct {
	ID     string
	Value  interface{}
	values []string
}

func NewCustomData(id string, values ...string) *CustomData {
	return &CustomData{
		ID:     id,
		Value:  nil,
		values: values,
	}
}

func (c CustomData) QueryEncode() string {
	if (c.Value == nil) && (len(c.values) == 0) {
		return ""
	}
	qb := network.NewQueryBuilder()
	qb.Append(network.QPEventType, customDataEventType)
	qb.Append(network.QPIndex, c.ID)
	qb.Append(network.QPValuesCountMap, c.encodeValues())
	qb.Append(network.QPOverwrite, "true")
	qb.Append(network.QPNonce, network.GetNonce())
	return qb.String()
}

func (c CustomData) encodeValues() string {
	sb := strings.Builder{}
	sb.WriteString("{\"")
	if c.Value != nil {
		s := utils.EscapeJsonStringControlSymbols(fmt.Sprint(c.Value))
		sb.WriteString(s)
		sb.WriteString("\":1")
	} else {
		for i, value := range c.values {
			if i > 0 {
				sb.WriteString(",\"")
			}
			s := strings.ReplaceAll(value, "\\", "\\\\")
			s = strings.ReplaceAll(s, "\"", "\\\"")
			sb.WriteString(s)
			sb.WriteString("\":1")
		}
	}
	sb.WriteRune('}')
	return sb.String()
}

func (c CustomData) DataType() DataType {
	return DataTypeCustom
}

func (c CustomData) GetValues() []string {
	if c.Value != nil {
		values := [1]string{fmt.Sprint(c.Value)}
		return values[:]
	} else {
		return c.values
	}
}

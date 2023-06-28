package types

import (
	"fmt"
	"strings"
)

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
	if c.Value == nil && len(c.values) == 0 {
		return ""
	}
	var val strings.Builder
	c.addStringValues(&val)
	valueToCount := EncodeURIComponent("valueToCount", val.String())
	var b strings.Builder
	b.WriteString("eventType=customData&index=")
	b.WriteString(c.ID)
	b.WriteString("&")
	b.WriteString(valueToCount)
	b.WriteString("&overwrite=true&nonce=")
	b.WriteString(GetNonce())
	return b.String()
}

func (c CustomData) DataType() DataType {
	return DataTypeCustom
}

func (c CustomData) addStringValues(val *strings.Builder) {
	val.WriteString(`[`)
	if c.Value != nil {
		val.WriteString(fmt.Sprintf(`["%s",1]`, c.Value))
	} else {
		for i, value := range c.values {
			val.WriteString(fmt.Sprintf(`["%s",1]`, value))
			if i < len(c.values)-1 {
				val.WriteString(`,`)
			}
		}
	}
	val.WriteString(`]`)
}

func (c CustomData) GetValues() []string {
	if c.Value != nil {
		values := [1]string{fmt.Sprint(c.Value)}
		return values[:]
	} else {
		return c.values
	}
}

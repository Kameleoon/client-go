package types

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Kameleoon/client-go/v2/utils"
	"github.com/segmentio/encoding/json"

	"net/url"
)

type Data interface {
	QueryEncode() string
	DataType() DataType
}

type TargetingData struct {
	Data
	LastActivityTime time.Time
}

type DataCell struct {
	Data  []TargetingData
	Index map[int]struct{}
}

func (d *DataCell) MarshalJSON() ([]byte, error) {
	return json.Marshal(&d.Data)
}

func (d *DataCell) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &d.Data)
}

func (d *DataCell) String() string {
	b, _ := d.MarshalJSON()
	var s strings.Builder
	s.Write(b)
	return s.String()
}

const NonceLength = 16

type DataType string

const (
	DataTypeCustom     DataType = "CUSTOM"
	DataTypeBrowser    DataType = "BROWSER"
	DataTypeConversion DataType = "CONVERSION"
	DataTypeDevice     DataType = "DEVICE"
	DataTypePageView   DataType = "PAGE_VIEW"
	DataTypeUserAgent  DataType = "USER_AGENT"
)

func GetNonce() string {
	return utils.GetRandomString(NonceLength)
}

type EventData struct {
	Type  DataType
	Value map[string]json.RawMessage
}

func (c *EventData) UnmarshalJSON(b []byte) error {
	c.Value = make(map[string]json.RawMessage)
	err := json.Unmarshal(b, c.Value)
	if t, exist := c.Value["type"]; exist {
		delete(c.Value, "type")
		c.Type = DataType(t)
	}
	return err
}

func (c EventData) QueryEncode() string {
	var b strings.Builder
	b.WriteString("eventType=")
	b.WriteString(string(c.Type))
	b.WriteString("&nonce=")
	b.WriteString(GetNonce())
	if len(c.Value) == 0 {
		return b.String()
	}
	for k, v := range c.Value {
		b.WriteByte('&')
		b.WriteString(k)
		b.WriteByte('=')
		b.Write(v)
	}
	return b.String()
}

func (c EventData) DataType() DataType {
	return c.Type
}

//TODO: remove Value in next major version and make Values public
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

type BrowserType int

const (
	BrowserTypeChrome BrowserType = iota
	BrowserTypeIE
	BrowserTypeFirefox
	BrowserTypeSafari
	BrowserTypeOpera
	BrowserTypeOther
)

type Browser struct {
	Type BrowserType
}

func (b Browser) QueryEncode() string {
	var sb strings.Builder
	sb.WriteString("eventType=staticData&browserIndex=")
	sb.WriteString(utils.WritePositiveInt(int(b.Type)))
	sb.WriteString("&nonce=")
	sb.WriteString(GetNonce())
	return sb.String()
}

func (b Browser) DataType() DataType {
	return DataTypeBrowser
}

type PageView struct {
	URL       string
	Title     string
	Referrers []int
}

func (v PageView) QueryEncode() string {
	var b strings.Builder
	b.WriteString("eventType=page&")
	b.WriteString(EncodeURIComponent("href", v.URL))
	b.WriteString("&title=")
	b.WriteString(v.Title)
	if len(v.Referrers) > 0 {
		b.WriteString("&referrersIndices=[")
		b.WriteString(utils.ArrayToString(v.Referrers, ","))
		b.WriteByte(']')
	}
	b.WriteString("&nonce=")
	b.WriteString(GetNonce())

	return b.String()
}

func (v PageView) DataType() DataType {
	return DataTypePageView
}

type DeviceType string

const (
	DeviceTypeDesktop DeviceType = "DESKTOP"
	DeviceTypePhone   DeviceType = "PHONE"
	DeviceTypeTablet  DeviceType = "TABLET"
)

type Device struct {
	Type DeviceType
}

func (device Device) QueryEncode() string {
	var b strings.Builder
	b.WriteString("eventType=staticData&deviceType=")
	b.WriteString(string(device.Type))
	b.WriteString("&nonce=")
	b.WriteString(GetNonce())
	return b.String()
}

func (device Device) DataType() DataType {
	return DataTypeDevice
}

type Conversion struct {
	GoalID   int
	Revenue  float64
	Negative bool
}

func (c Conversion) QueryEncode() string {
	var b strings.Builder
	b.WriteString("eventType=conversion&goalId=")
	b.WriteString(utils.WritePositiveInt(c.GoalID))
	b.WriteString("&revenue=")
	b.WriteString(strconv.FormatFloat(c.Revenue, 'f', -1, 64))
	b.WriteString("&negative=")
	b.WriteString(strconv.FormatBool(c.Negative))
	b.WriteString("&nonce=")
	b.WriteString(GetNonce())
	return b.String()
}

func (c Conversion) DataType() DataType {
	return DataTypeConversion
}

type UserAgent struct {
	Value string
}

func (ua UserAgent) QueryEncode() string {
	return ""
}

func (ua UserAgent) DataType() DataType {
	return DataTypeUserAgent
}

func EncodeURIComponent(key string, value string) string {
	parameters := url.Values{}
	parameters.Add(key, value)
	encoded := parameters.Encode()

	encoded = strings.ReplaceAll(encoded, "+", "%20")
	encoded = strings.ReplaceAll(encoded, "%21", "!")
	encoded = strings.ReplaceAll(encoded, "%27", "'")
	encoded = strings.ReplaceAll(encoded, "%28", "(")
	encoded = strings.ReplaceAll(encoded, "%29", ")")
	encoded = strings.ReplaceAll(encoded, "%2A", "*")

	return encoded
}

func Remove(seq []Data, index int) []Data {
	seq[index] = seq[len(seq)-1]
	return seq[:len(seq)-1]
}

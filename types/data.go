package types

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/encoding/json"

	"github.com/Kameleoon/client-go/utils"
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
	DataTypeInterest   DataType = "INTEREST"
	DataTypePageView   DataType = "PAGE_VIEW"
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

type CustomData struct {
	ID    string
	Value interface{}
}

func (c CustomData) QueryEncode() string {
	var val strings.Builder
	val.WriteString(`[["`)
	val.WriteString(fmt.Sprint(c.Value))
	val.WriteString(`",1]]`)
	var b strings.Builder
	b.WriteString("eventType=customData&index=")
	b.WriteString(c.ID)
	b.WriteString("&valueToCount=")
	b.WriteString(val.String())
	b.WriteString("&overwrite=true&nonce=")
	b.WriteString(GetNonce())
	return b.String()
}

func (c CustomData) DataType() DataType {
	return DataTypeCustom
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
	sb.WriteString("eventType=staticData&browser=")
	sb.WriteString(utils.WriteUint(int(b.Type)))
	sb.WriteString("&nonce=")
	sb.WriteString(GetNonce())
	return sb.String()
}

func (b Browser) DataType() DataType {
	return DataTypeBrowser
}

type PageView struct {
	URL      string
	Title    string
	Referrer int
}

func (v PageView) QueryEncode() string {
	var b strings.Builder
	b.WriteString("eventType=page&href=")
	b.WriteString(v.URL)
	b.WriteString("&title=")
	b.WriteString(v.Title)
	b.WriteString("&keyPages=[]")
	if v.Referrer == 0 {
		b.WriteString("&referrers=[")
		b.WriteString(strconv.Itoa(v.Referrer))
		b.WriteByte(']')
	}
	b.WriteString("&nonce=")
	b.WriteString(GetNonce())

	return b.String()
}

func (v PageView) DataType() DataType {
	return DataTypePageView
}

type Interest struct {
	Index int
}

func (i Interest) QueryEncode() string {
	var b strings.Builder
	b.WriteString("eventType=interests&indexes=[")
	b.WriteString(strconv.Itoa(i.Index))
	b.WriteString("]&fresh=true&nonce=")
	b.WriteString(GetNonce())
	return b.String()
}

func (i Interest) DataType() DataType {
	return DataTypeInterest
}

type Conversion struct {
	GoalID   int
	Revenue  float64
	Negative bool
}

func (c Conversion) QueryEncode() string {
	var b strings.Builder
	b.WriteString("eventType=conversion&goalId=")
	b.WriteString(utils.WriteUint(c.GoalID))
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

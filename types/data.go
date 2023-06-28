package types

import (
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

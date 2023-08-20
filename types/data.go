package types

import (
	"strings"
	"time"

	"github.com/segmentio/encoding/json"
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

type DataType string

const (
	DataTypeCustom     DataType = "CUSTOM"
	DataTypeBrowser    DataType = "BROWSER"
	DataTypeConversion DataType = "CONVERSION"
	DataTypeDevice     DataType = "DEVICE"
	DataTypePageView   DataType = "PAGE_VIEW"
	DataTypeUserAgent  DataType = "USER_AGENT"
)

func Remove(seq []Data, index int) []Data {
	seq[index] = seq[len(seq)-1]
	return seq[:len(seq)-1]
}

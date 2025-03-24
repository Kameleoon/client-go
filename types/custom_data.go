package types

import (
	"fmt"
	"strings"

	"github.com/Kameleoon/client-go/v3/logging"

	"github.com/Kameleoon/client-go/v3/utils"
)

type ICustomData interface {
	Data
	Sendable
	ID() int
	Values() []string
}

const customDataEventType = "customData"

type CustomData struct {
	duplicationUnsafeSendableBase
	id     int
	values []string
}

func NewCustomData(id int, values ...string) *CustomData {
	return &CustomData{
		id:     id,
		values: values,
	}
}

func (cd CustomData) String() string {
	return fmt.Sprintf("CustomData{id:%d,values:%s}", cd.id, logging.ObjectToString(cd.values))
}

func (cd *CustomData) dataRestriction() {
	// This method is required to separate external type `Data` from `BaseData` types
}

func (cd *CustomData) ID() int {
	return cd.id
}

func (cd *CustomData) Values() []string {
	return cd.values
}

func (cd *CustomData) QueryEncode() string {
	nonce := cd.Nonce()
	if len(nonce) == 0 {
		return ""
	}
	qb := utils.NewQueryBuilder()
	qb.Append(utils.QPEventType, customDataEventType)
	qb.Append(utils.QPIndex, fmt.Sprint(cd.id))
	qb.Append(utils.QPValuesCountMap, cd.encodeValues())
	qb.Append(utils.QPOverwrite, "true")
	qb.Append(utils.QPNonce, nonce)
	return qb.String()
}

func (cd *CustomData) encodeValues() string {
	sb := strings.Builder{}
	sb.WriteString("{")
	for i, value := range cd.values {
		if i > 0 {
			sb.WriteString(",\"")
		} else {
			sb.WriteString("\"")
		}
		s := strings.ReplaceAll(value, "\\", "\\\\")
		s = strings.ReplaceAll(s, "\"", "\\\"")
		sb.WriteString(s)
		sb.WriteString("\":1")
	}
	sb.WriteRune('}')
	return sb.String()
}

func (cd *CustomData) DataType() DataType {
	return DataTypeCustom
}

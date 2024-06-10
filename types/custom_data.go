package types

import (
	"fmt"
	"strings"

	"github.com/Kameleoon/client-go/v3/utils"
)

const customDataEventType = "customData"

type CustomData struct {
	duplicationUnsafeSendableBase
	id                  int
	values              []string
	isMappingIdentifier bool
}

func NewCustomData(id int, values ...string) *CustomData {
	return &CustomData{
		id:     id,
		values: values,
	}
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

func (cd *CustomData) IsMappingIdentifier() bool {
	return cd.isMappingIdentifier
}
func (cd *CustomData) SetIsMappingIdentifier(value bool) {
	cd.isMappingIdentifier = value
}

func (cd *CustomData) QueryEncode() string {
	if len(cd.values) == 0 {
		return ""
	}
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
	if cd.isMappingIdentifier {
		qb.Append(utils.QPMappingIdentifier, "true")
	}
	return qb.String()
}
func (cd *CustomData) encodeValues() string {
	sb := strings.Builder{}
	sb.WriteString("{\"")
	for i, value := range cd.values {
		if i > 0 {
			sb.WriteString(",\"")
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

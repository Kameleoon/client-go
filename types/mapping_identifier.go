package types

import (
	"github.com/Kameleoon/client-go/v3/utils"
)

type MappingIdentifier struct {
	CustomData
}

func NewMappingIdentifier(cd *CustomData) *MappingIdentifier {
	return &MappingIdentifier{CustomData: *cd}
}

func (mi *MappingIdentifier) Unsent() bool {
	return true
}
func (mi *MappingIdentifier) Transmitting() bool {
	return false
}
func (mi *MappingIdentifier) Sent() bool {
	return false
}

func (mi *MappingIdentifier) QueryEncode() string {
	mip := utils.NewQueryBuilder().Append(utils.QPMappingIdentifier, "true").String()
	return mi.CustomData.QueryEncode() + "&" + mip
}

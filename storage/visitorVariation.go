package storage

import "time"

type VisitorVariation struct {
	// uses uint32 for saving memory space
	// Need to be fixed if variation id or unix time > 4294967295 (Feb 07 2106)
	VariationId    uint32
	AssignmentDate uint32
}

func NewVisitorVariation(variationId uint32) *VisitorVariation {
	return &VisitorVariation{VariationId: variationId, AssignmentDate: uint32(time.Now().Unix())}
}

func (variation VisitorVariation) isValid(respoolTime int) bool {
	return variation.AssignmentDate > uint32(respoolTime)
}

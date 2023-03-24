package storage

import "time"

type VisitorVariation struct {
	VariationId    uint32
	AssignmentDate uint32
}

func NewVisitorVariation(variationId uint32) *VisitorVariation {
	return &VisitorVariation{VariationId: variationId, AssignmentDate: uint32(time.Now().Unix())}
}

func (variation VisitorVariation) isValid(respoolTime uint32) bool {
	return variation.AssignmentDate > respoolTime
}

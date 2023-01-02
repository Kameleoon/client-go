package storage

import "time"

type VisitorVariation struct {
	VariationId    int
	AssignmentDate int64
}

func NewVisitorVariation(variationId int) *VisitorVariation {
	return &VisitorVariation{VariationId: variationId, AssignmentDate: time.Now().Unix()}
}

func (variation VisitorVariation) isValid(respoolTime int64) bool {
	return variation.AssignmentDate > respoolTime
}

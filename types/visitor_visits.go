package types

import "fmt"

type VisitorVisits struct {
	previousVisitTimestamps []int64
}

func NewVisitorVisits(previousVisitTimestamps []int64) *VisitorVisits {
	return &VisitorVisits{
		previousVisitTimestamps: previousVisitTimestamps,
	}
}

func (vv VisitorVisits) String() string {
	return fmt.Sprintf(
		"VisitorVisits{previousVisitTimestamps:%v}",
		vv.previousVisitTimestamps,
	)
}

func (vv *VisitorVisits) PreviousVisitTimestamps() []int64 {
	if vv != nil {
		return vv.previousVisitTimestamps
	}
	return []int64{}
}

func (vv *VisitorVisits) DataType() DataType {
	return DataTypeVisitorVisits
}

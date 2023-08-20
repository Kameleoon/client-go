package network

import "fmt"

const experimentEventType = "experiment"

type ExperimentEvent struct {
	ExperimentId int
	VariationId  int
}

func (e ExperimentEvent) QueryEncode() string {
	qb := NewQueryBuilder()
	qb.Append(QPEventType, experimentEventType)
	qb.Append(QPExperimentId, fmt.Sprint(e.ExperimentId))
	qb.Append(QPVariationId, fmt.Sprint(e.VariationId))
	qb.Append(QPNonce, GetNonce())
	return qb.String()
}

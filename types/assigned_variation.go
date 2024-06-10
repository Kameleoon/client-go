package types

import (
	"fmt"
	"time"

	"github.com/Kameleoon/client-go/v3/utils"
)

const assignedVariationType = "experiment"

type AssignedVariation struct {
	duplicationUnsafeSendableBase
	experimentId   int
	variationId    int
	ruleType       RuleType
	assignmentTime time.Time
}

func NewAssignedVariation(experimentId int, variationId int, ruleType RuleType) *AssignedVariation {
	return &AssignedVariation{
		experimentId:   experimentId,
		variationId:    variationId,
		ruleType:       ruleType,
		assignmentTime: time.Now(),
	}
}

func NewAssignedVariationWithTime(
	experimentId int, variationId int, ruleType RuleType, assignmentTime time.Time,
) *AssignedVariation {
	return &AssignedVariation{
		experimentId:   experimentId,
		variationId:    variationId,
		ruleType:       ruleType,
		assignmentTime: assignmentTime,
	}
}

func (av *AssignedVariation) ExperimentId() int {
	return av.experimentId
}

func (av *AssignedVariation) VariationId() int {
	return av.variationId
}

func (av *AssignedVariation) RuleType() RuleType {
	return av.ruleType
}

func (av *AssignedVariation) AssignmentTime() time.Time {
	return av.assignmentTime
}

func (av *AssignedVariation) IsValid(respoolTime int) bool {
	return (respoolTime == 0) || (av.assignmentTime.Unix() >= int64(respoolTime))
}

func (av *AssignedVariation) QueryEncode() string {
	nonce := av.Nonce()
	if len(nonce) == 0 {
		return ""
	}
	qb := utils.NewQueryBuilder()
	qb.Append(utils.QPEventType, assignedVariationType)
	qb.Append(utils.QPExperimentId, fmt.Sprint(av.experimentId))
	qb.Append(utils.QPVariationId, fmt.Sprint(av.variationId))
	qb.Append(utils.QPNonce, nonce)
	return qb.String()
}

func (av *AssignedVariation) DataType() DataType {
	return DataTypeAssignedVariation
}

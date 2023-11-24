package types

import "github.com/Kameleoon/client-go/v3/utils"

const activityEventType = "activity"

type ActivityEvent struct {
	duplicationUnsafeSendableBase
}

func NewActivityEvent() *ActivityEvent {
	return new(ActivityEvent)
}

func (ae *ActivityEvent) QueryEncode() string {
	nonce := ae.Nonce()
	if len(nonce) == 0 {
		return ""
	}
	qb := utils.NewQueryBuilder()
	qb.Append(utils.QPEventType, activityEventType)
	qb.Append(utils.QPNonce, nonce)
	return qb.String()
}

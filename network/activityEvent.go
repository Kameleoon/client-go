package network

const activityEventType = "activity"

type ActivityEvent struct{}

func (ActivityEvent) QueryEncode() string {
	qb := NewQueryBuilder()
	qb.Append(QPEventType, activityEventType)
	qb.Append(QPNonce, GetNonce())
	return qb.String()
}

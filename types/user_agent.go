package types

type UserAgent struct {
	value string
}

func NewUserAgent(value string) UserAgent {
	return UserAgent{value: value}
}

func (ua UserAgent) dataRestriction() {}

func (ua UserAgent) Value() string {
	return ua.value
}

func (ua UserAgent) DataType() DataType {
	return DataTypeUserAgent
}

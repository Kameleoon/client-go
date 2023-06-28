package types

type UserAgent struct {
	Value string
}

func (ua UserAgent) QueryEncode() string {
	return ""
}

func (ua UserAgent) DataType() DataType {
	return DataTypeUserAgent
}

package types

type QueryEncodable interface {
	QueryEncode() string
}

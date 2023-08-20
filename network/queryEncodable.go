package network

type QueryEncodable interface {
	QueryEncode() string
}

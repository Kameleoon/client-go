package types

type Settings interface {
	RealTimeUpdate() bool
	IsConsentRequired() bool
	DataApiDomain() string
}

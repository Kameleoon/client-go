package types

type Settings interface {
	RealTimeUpdate() bool
	IsConsentRequired() bool
	BlockingBehaviourIfConsentNotGiven() ConsentBlockingBehaviour
	DataApiDomain() string
}

package types

import "github.com/Kameleoon/client-go/v3/logging"

type ConsentBlockingBehaviour = int8

const (
	PartiallyBlockedByConsent  ConsentBlockingBehaviour = 0
	CompletelyBlockedByConsent ConsentBlockingBehaviour = 1
)

const (
	PartiallyBlockedByConsentStr  = "PARTIALLY_BLOCK"
	CompletelyBlockedByConsentStr = "FULLY_BLOCK"
)

func ConsentBlockingBehaviourFromStr(str string) ConsentBlockingBehaviour {
	switch str {
	case PartiallyBlockedByConsentStr:
		return PartiallyBlockedByConsent
	case CompletelyBlockedByConsentStr:
		return CompletelyBlockedByConsent
	}
	logging.Warning("Unexpected consent blocking type %s", str)
	return PartiallyBlockedByConsent
}

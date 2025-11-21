package types

type LegalConsent = int8

const (
	LegalConsentUnknown  LegalConsent = 0
	LegalConsentGiven    LegalConsent = 1
	LegalConsentNotGiven LegalConsent = 2
)

package types

import (
	"time"
)

type WhenTimeoutType string

const (
	WhenTimeoutRun             WhenTimeoutType = "RUN"
	WhenTimeoutDisableForPage  WhenTimeoutType = "DISABLE_FOR_PAGE"
	WhenTimeoutDisableForVisit WhenTimeoutType = "DISABLE_FOR_VISIT"
)

type DataStorageType string

const (
	DataStorageStandardCookie DataStorageType = "STANDARD_COOKIE"
	DataStorageLocalStorage   DataStorageType = "LOCAL_STORAGE"
	DataStorageCustomCookie   DataStorageType = "CUSTOM_COOKIE"
)

type IndicatorType string

const (
	IndicatorsRetentionRate     IndicatorType = "RETENTION_RATE"
	IndicatorsNumberOfPagesSeen IndicatorType = "NUMBER_OF_PAGES_SEEN"
	IndicatorsDwellTime         IndicatorType = "DWELL_TIME"
)

type EventMethodType string

const (
	EventMethodClick     EventMethodType = "CLICK"
	EventMethodMousedown EventMethodType = "MOUSEDOWN"
	EventMethodMouseup   EventMethodType = "MOUSEUP"
)

type SiteResponse struct {
	ID                  int                   `json:"id"`
	URL                 string                `json:"url"`
	Description         string                `json:"description"`
	Code                string                `json:"code"`
	BehaviorWhenTimeout WhenTimeoutType       `json:"behaviorWhenTimeout"`
	DataStorage         DataStorageType       `json:"dataStorage"`
	TrackingScript      string                `json:"trackingScript"`
	DomainNames         []string              `json:"domainNames"`
	Indicators          []IndicatorType       `json:"indicators"`
	DateCreated         TimeNoTZ              `json:"dateCreated"`
	IsScriptActive      bool                  `json:"isScriptActive"`
	CaptureEventMethod  EventMethodType       `json:"captureEventMethod"`
	IsAudienceUsed      bool                  `json:"isAudienceUsed"`
	IsKameleoonEnabled  bool                  `json:"isKameleoonEnabled"`
	Experiment          ExperimentConfig      `json:"experimentConfig"`
	Personalization     PersonalizationConfig `json:"personalizationConfig"`
	Audience            AudienceConfig        `json:"audienceConfig"`
}

type ConsentType string

const (
	ConsentOff         ConsentType = "OFF"
	ConsentRequired    ConsentType = "REQUIRED"
	ConsentInteractive ConsentType = "INTERACTIVE"
	ConsentIABTCF      ConsentType = "IABTCF"
)

type ConsentOptoutType string

const (
	ConsentOptoutRun   ConsentOptoutType = "RUN"
	ConsentOptoutBlock ConsentOptoutType = "BLOCK"
)

type BeforeConsentType string

const (
	BeforeConsentNone      BeforeConsentType = "NONE"
	BeforeConsentPartially BeforeConsentType = "PARTIALLY"
	BeforeConsentAll       BeforeConsentType = "ALL"
)

type PersonalizationConfig struct {
	PersonalizationsDeviation        float64           `json:"personalizationsDeviation"`
	IsSameTypePersonalizationEnabled bool              `json:"isSameTypePersonalizationEnabled"`
	IsSameJqueryInjectionAllowed     bool              `json:"isSameJqueryInjectionAllowed"`
	PersonalizationConsent           ConsentType       `json:"personalizationConsent"`
	PersonalizationConsentOptout     ConsentOptoutType `json:"personalizationConsentOptout"`
	BeforePersonalizationConsent     BeforeConsentType `json:"beforePersonalizationConsent"`
}

type TimeNoTZ time.Time

const TimeNoTZLayout = `"2006-01-02T15:04:05"`

// UnmarshalJSON Parses the json string in the custom format
func (t *TimeNoTZ) UnmarshalJSON(date []byte) error {
	nt, err := time.Parse(TimeNoTZLayout, string(date))
	*t = TimeNoTZ(nt)
	return err
}

// MarshalJSON writes a quoted string in the custom format
func (t TimeNoTZ) MarshalJSON() ([]byte, error) {
	return time.Time(t).AppendFormat(nil, TimeNoTZLayout), nil
}

// String returns the time in the custom format
func (t TimeNoTZ) String() string {
	return time.Time(t).Format(TimeNoTZLayout)
}

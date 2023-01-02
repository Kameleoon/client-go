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

type TimeNoTZ time.Time
type TimeTZ struct {
	time.Time
}

const TimeNoTZLayout = `"2006-01-02T15:04:05"`
const TimeTZLayout = `"2006-01-02T15:04:05-0700"`

// UnmarshalJSON Parses the json string in the custom format
func (t *TimeNoTZ) UnmarshalJSON(date []byte) error {
	nt, err := time.Parse(TimeNoTZLayout, string(date))
	*t = TimeNoTZ(nt)
	return err
}

// UnmarshalJSON Parses the json string in the custom format
func (t *TimeTZ) UnmarshalJSON(date []byte) error {
	if string(date) != "null" {
		nt, err := time.Parse(TimeTZLayout, string(date))
		t.Time = nt
		return err
	}
	return nil
}

// MarshalJSON writes a quoted string in the custom format
func (t TimeNoTZ) MarshalJSON() ([]byte, error) {
	return time.Time(t).AppendFormat(nil, TimeNoTZLayout), nil
}

// String returns the time in the custom format
func (t TimeNoTZ) String() string {
	return time.Time(t).Format(TimeNoTZLayout)
}

// UnmarshalJSON string to int
// func (value *string) UnmarshalJSON(date []byte) error {
// 	nt, err := time.Parse(TimeNoTZLayout, string(date))
// 	*t = TimeNoTZ(nt)
// 	return err
// }

package types

import "time"

const (
	STATUS_ACTIVE              = "ACTIVE"
	FEATURE_STATUS_DEACTIVATED = "DEACTIVATED"
)

type FeatureFlag struct {
	ID                 int             `json:"id"`
	Name               string          `json:"name"`
	IdentificationKey  string          `json:"identificationKey"`
	Description        string          `json:"description"`
	Tags               []string        `json:"tags"`
	SiteID             int             `json:"siteId"`
	ExpositionRate     float64         `json:"expositionRate"`
	TargetingSegmentID int             `json:"targetingSegmentId"`
	TargetingSegment   interface{}     `json:"targetingSegment,omitempty"`
	VariationsID       []int           `json:"variations,omitempty"`
	Variations         []Variation     `json:"-"`
	Goals              []int           `json:"goals"`
	SDKLanguageType    SDKLanguageType `json:"sdkLanguageType"`
	Status             string          `json:"status"`
	DateCreated        TimeNoTZ        `json:"dateCreated"`
	DateModified       TimeNoTZ        `json:"dateModified"`
	RespoolTime        RespoolTime     `json:"respoolTime"`
	FeatureStatus      string          `json:"featureStatus"`
	Schedules          []Schedule      `json:"schedules"`
}

type Schedule struct {
	DateStart *TimeTZ `json:"dateStart,omitempty"`
	DateEnd   *TimeTZ `json:"dateEnd,omitempty"`
}

func (ff *FeatureFlag) IsScheduleActive() bool {
	/// if featureStatus == `DEACTIVATED` or no schedules then need to return current status
	currentStatus := ff.Status == STATUS_ACTIVE
	if ff.FeatureStatus == FEATURE_STATUS_DEACTIVATED || len(ff.Schedules) == 0 {
		return currentStatus
	}
	currentTime := time.Now()
	/// need to find if currentTime is in any period -> active or not -> not activate
	for _, schedule := range ff.Schedules {
		if (schedule.DateStart == nil || schedule.DateStart.Before(currentTime)) && (schedule.DateEnd == nil || schedule.DateEnd.After(currentTime)) {
			return true
		}
	}
	return false
}

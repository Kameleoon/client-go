package configuration

import (
	"time"

	"github.com/Kameleoon/client-go/targeting"
	"github.com/Kameleoon/client-go/types"
	"github.com/segmentio/encoding/json"
)

const (
	STATUS_ACTIVE              = "ACTIVE"
	FEATURE_STATUS_DEACTIVATED = "DEACTIVATED"
)

type FeatureFlag struct {
	types.FeatureFlag
	TargetingSegment *targeting.Segment `json:"-"`
}

func (ff *FeatureFlag) UnmarshalJSON(data []byte) error {
	type FeatureFlagHidden FeatureFlag
	if err := json.Unmarshal(data, (*FeatureFlagHidden)(ff)); err != nil {
		return err
	}
	if ff.Segment.ID != 0 {
		ff.TargetingSegment = targeting.NewSegment(&ff.Segment)
	}
	return nil
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

func (ff FeatureFlag) SiteCodeEnabled() bool {
	return ff.SiteEnabled
}

func (ff FeatureFlag) GetTargetingSegment() *targeting.Segment {
	return ff.TargetingSegment
}

func (ff FeatureFlag) GetId() int {
	return ff.ID
}

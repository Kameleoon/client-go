package types

type Segment interface {
	CheckTargeting(data TargetingDataGetter) bool
}

package configuration

import "github.com/Kameleoon/client-go/v2/targeting"

type TargetingObjectInterface interface {
	GetTargetingSegment() *targeting.Segment
}

type SiteCodeEnabledInterface interface {
	SiteCodeEnabled() bool
}

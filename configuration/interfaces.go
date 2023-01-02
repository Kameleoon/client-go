package configuration

import "github.com/Kameleoon/client-go/targeting"

type TargetingObjectInterface interface {
	GetTargetingSegment() *targeting.Segment
}

type SiteCodeEnabledInterface interface {
	SiteCodeEnabled() bool
}

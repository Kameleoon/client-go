package configuration

import "github.com/Kameleoon/client-go/v3/targeting"

type TargetingObjectInterface interface {
	GetTargetingSegment() *targeting.Segment
}

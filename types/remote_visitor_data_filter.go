package types

import (
	"fmt"

	"github.com/Kameleoon/client-go/v3/logging"
	"github.com/Kameleoon/client-go/v3/utils"
)

type RemoteVisitorDataFilter struct {
	PreviousVisitAmount int // The value must be between 1 and 25.
	CurrentVisit        bool
	CustomData          bool
	PageViews           bool
	Geolocation         bool
	Device              bool
	Browser             bool
	OperatingSystem     bool
	Conversion          bool
	Experiments         bool
	Kcs                 bool
	VisitorCode         interface{} // Type: `bool`. Defaults to `true` if not set.
}

func DefaultRemoteVisitorDataFilter() RemoteVisitorDataFilter {
	return RemoteVisitorDataFilter{
		PreviousVisitAmount: 1,
		CurrentVisit:        true,
		CustomData:          true,
		VisitorCode:         true,
	}
}

func (r *RemoteVisitorDataFilter) ApplyDefaultValues() {
	var ok bool
	if r.VisitorCode, ok = utils.Deopt[bool](r.VisitorCode, true); !ok {
		logging.Warning(
			"Failed to deopt RemoteVisitorDataFilter.VisitorCode. Expected bool or nil value, got '%s'", r.VisitorCode,
		)
	}
}

func (r RemoteVisitorDataFilter) String() string {
	return fmt.Sprintf(
		"RemoteVisitorDataFilter{PreviousVisitAmount:%d,CurrentVisit:%t,CustomData:%t,PageViews:%t,Geolocation:%t,"+
			"Device:%t,Browser:%t,OperatingSystem:%t,Conversions:%t,Experiments:%t,Kcs:%t,VisitorCode:%t}",
		r.PreviousVisitAmount,
		r.CurrentVisit,
		r.CustomData,
		r.PageViews,
		r.Geolocation,
		r.Device,
		r.Browser,
		r.OperatingSystem,
		r.Conversion,
		r.Experiments,
		r.Kcs,
		r.VisitorCode,
	)
}

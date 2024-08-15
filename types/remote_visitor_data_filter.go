package types

import "fmt"

type RemoteVisitorDataFilter struct {
	PreviousVisitAmount int
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
}

func DefaultRemoteVisitorDataFilter() RemoteVisitorDataFilter {
	return RemoteVisitorDataFilter{
		PreviousVisitAmount: 1,
		CurrentVisit:        true,
		CustomData:          true,
	}
}

func (r RemoteVisitorDataFilter) String() string {
	return fmt.Sprintf(
		"RemoteVisitorDataFilter{PreviousVisitAmount:%d,CurrentVisit:%t,CustomData:%t,PageViews:%t,"+
			"Geolocation:%t,Device:%t,Browser:%t,OperatingSystem:%t,Conversions:%t,Experiments:%t,Kcs:%t}",
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
	)
}

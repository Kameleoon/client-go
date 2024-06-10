package types

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

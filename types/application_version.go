package types

type ApplicationVersion struct {
	Version string `json:"version"`
}

func NewApplicationVersion(version string) *ApplicationVersion {
	return &ApplicationVersion{Version: version}
}

func (av *ApplicationVersion) dataRestriction() {}

func (av *ApplicationVersion) DataType() DataType {
	return DataTypeApplicationVersion
}

func (av *ApplicationVersion) String() string {
	return "ApplicationVersion{version:" + av.Version + "}"
}

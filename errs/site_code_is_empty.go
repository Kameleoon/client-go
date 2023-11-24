package errs

type SiteCodeIsEmpty struct {
	KameleoonError
}

func NewSiteCodeIsEmpty(msg string) *SiteCodeIsEmpty {
	return &SiteCodeIsEmpty{NewKameleoonError("Site code is empty")}
}

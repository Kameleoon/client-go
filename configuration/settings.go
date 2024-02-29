package configuration

import "encoding/json"

const consentTypeRequired = "REQUIRED"

type Settings struct {
	realTimeUpdate    bool
	isConsentRequired bool
	dataApiDomain     string
}

func (s *Settings) UnmarshalJSON(data []byte) error {
	var sm = struct {
		RealTimeUpdate bool   `json:"realTimeUpdate"`
		ConsentType    string `json:"consentType"`
		DataApiDomain  string `json:"dataApiDomain,omitempty"`
	}{}
	if err := json.Unmarshal(data, &sm); err != nil {
		return err
	}
	s.realTimeUpdate = sm.RealTimeUpdate
	s.isConsentRequired = sm.ConsentType == consentTypeRequired
	s.dataApiDomain = sm.DataApiDomain
	return nil
}

func (s Settings) RealTimeUpdate() bool {
	return s.realTimeUpdate
}

func (s Settings) IsConsentRequired() bool {
	return s.isConsentRequired
}

func (s Settings) DataApiDomain() string {
	return s.dataApiDomain
}

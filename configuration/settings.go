package configuration

import "encoding/json"

const consentTypeRequired = "REQUIRED"

type Settings struct {
	RealTimeUpdate    bool
	IsConsentRequired bool
}

func (s *Settings) UnmarshalJSON(data []byte) error {
	var sm = struct {
		RealTimeUpdate bool   `json:"realTimeUpdate"`
		ConsentType    string `json:"consentType"`
	}{}
	if err := json.Unmarshal(data, &sm); err != nil {
		return err
	}
	s.RealTimeUpdate = sm.RealTimeUpdate
	s.IsConsentRequired = sm.ConsentType == consentTypeRequired
	return nil
}

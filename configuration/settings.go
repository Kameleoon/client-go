package configuration

import (
	"encoding/json"
	"fmt"

	"github.com/Kameleoon/client-go/v3/types"
)

const consentTypeRequired = "REQUIRED"

type Settings struct {
	realTimeUpdate                     bool
	isConsentRequired                  bool
	blockingBehaviourIfConsentNotGiven types.ConsentBlockingBehaviour
	dataApiDomain                      string
}

func (s Settings) String() string {
	return fmt.Sprintf("Settings{realTimeUpdate:%v,isConsentRequired:%v,dataApiDomain:'%v'}",
		s.realTimeUpdate, s.isConsentRequired, s.dataApiDomain)
}

func (s *Settings) UnmarshalJSON(data []byte) error {
	var sm = struct {
		RealTimeUpdate        bool   `json:"realTimeUpdate"`
		ConsentType           string `json:"consentType"`
		ConsentOptOutBehavior string `json:"consentOptOutBehavior"`
		DataApiDomain         string `json:"dataApiDomain,omitempty"`
	}{}
	if err := json.Unmarshal(data, &sm); err != nil {
		return err
	}
	s.realTimeUpdate = sm.RealTimeUpdate
	s.isConsentRequired = sm.ConsentType == consentTypeRequired
	s.blockingBehaviourIfConsentNotGiven = types.ConsentBlockingBehaviourFromStr(sm.ConsentOptOutBehavior)
	s.dataApiDomain = sm.DataApiDomain
	return nil
}

func (s Settings) RealTimeUpdate() bool {
	return s.realTimeUpdate
}

func (s Settings) IsConsentRequired() bool {
	return s.isConsentRequired
}

func (s Settings) BlockingBehaviourIfConsentNotGiven() types.ConsentBlockingBehaviour {
	return s.blockingBehaviourIfConsentNotGiven
}

func (s Settings) DataApiDomain() string {
	return s.dataApiDomain
}

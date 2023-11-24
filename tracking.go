package kameleoon

import (
	"github.com/Kameleoon/client-go/v3/storage"
	"github.com/Kameleoon/client-go/v3/types"
)

func (c *kameleoonClient) isConsentProvided(visitor storage.Visitor) bool {
	return !c.dataFile.Settings().IsConsentRequired || ((visitor != nil) && visitor.LegalConsent())
}

func (c *kameleoonClient) sendTrackingRequest(visitorCode string, visitor storage.Visitor, forceRequest bool) {
	if visitor == nil {
		visitor = c.visitorManager.GetVisitor(visitorCode)
		if (visitor == nil) && c.dataFile.Settings().IsConsentRequired {
			return
		}
	}
	consent := c.isConsentProvided(visitor)
	var unsent []types.Sendable
	var userAgent string
	if visitor != nil {
		userAgent = visitor.UserAgent()
		unsent = c.selectUnsentData(visitor, consent)
	}
	if len(unsent) == 0 {
		if forceRequest && consent {
			unsent = append(unsent, types.NewActivityEvent())
		} else {
			return
		}
	}
	go func() {
		if c.makeTrackingRequest(visitorCode, userAgent, unsent) {
			for _, qe := range unsent {
				qe.MarkAsSent()
			}
		}
	}()
}

func (c *kameleoonClient) makeTrackingRequest(visitorCode string, userAgent string, data []types.Sendable) (sent bool) {
	if len(data) == 0 {
		return false
	}
	token := c.token
	c.log("Start post to tracking")
	out, err := c.networkManager.SendTrackingData(visitorCode, data, userAgent, token, -1)
	if err != nil {
		c.log("Failed to post tracking data, error: %v", err)
		return false
	}
	c.log("Post to tracking done")
	return out
}

func (c *kameleoonClient) selectUnsentData(visitor storage.Visitor, consent bool) []types.Sendable {
	var unsent []types.Sendable
	if consent {
		visitor.EnumerateSendableData(func(s types.Sendable) bool {
			if !s.Sent() {
				unsent = append(unsent, s)
			}
			return true
		})
	} else {
		visitor.Conversions().Enumerate(func(c *types.Conversion) bool {
			if !c.Sent() {
				unsent = append(unsent, c)
			}
			return true
		})
		visitor.Variations().Enumerate(func(av *types.AssignedVariation) bool {
			if !av.Sent() && (av.RuleType() == types.RuleTypeTargetedDelivery) {
				unsent = append(unsent, av)
			}
			return true
		})
	}
	return unsent
}

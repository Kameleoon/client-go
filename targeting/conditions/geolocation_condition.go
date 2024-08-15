package conditions

import (
	"strings"

	"github.com/Kameleoon/client-go/v3/types"
	"github.com/Kameleoon/client-go/v3/utils"
)

func NewGeolocationCondition(c types.TargetingCondition) *GeolocationCondition {
	return &GeolocationCondition{
		TargetingConditionBase: types.TargetingConditionBase{
			Type:    c.Type,
			Include: c.Include,
		},
		Country: c.Country,
		Region:  c.Region,
		City:    c.City,
	}
}

type GeolocationCondition struct {
	types.TargetingConditionBase
	Country string `json:"country,omitempty"`
	Region  string `json:"region,omitempty"`
	City    string `json:"city,omitempty"`
}

func (c *GeolocationCondition) CheckTargeting(targetData interface{}) bool {
	geolocation, ok := targetData.(*types.Geolocation)
	return ok && (geolocation != nil) &&
		(c.Country != "") && strings.EqualFold(geolocation.Country(), c.Country) &&
		((c.Region == "") || strings.EqualFold(geolocation.Region(), c.Region)) &&
		((c.City == "") || strings.EqualFold(geolocation.City(), c.City))
}

func (c GeolocationCondition) String() string {
	return utils.JsonToString(c)
}

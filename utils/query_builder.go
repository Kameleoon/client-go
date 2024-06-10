package utils

import (
	"net/url"
	"strings"
)

// QueryBuilder

type QueryBuilder struct {
	params url.Values
}

func NewQueryBuilder() *QueryBuilder {
	var qb QueryBuilder
	qb.params = url.Values{}
	return &qb
}

func (qb *QueryBuilder) Append(name string, value string) {
	if (len(name) > 0) && (len(value) > 0) {
		qb.params.Set(name, value)
	}
}

func (qb *QueryBuilder) String() string {
	replacer := strings.NewReplacer("+", "%20")
	return replacer.Replace(qb.params.Encode())
}

const (
	QPBrowserIndex            = "browserIndex"
	QPBrowserVersion          = "browserVersion"
	QPCity                    = "city"
	QPConversion              = "conversion"
	QPCountry                 = "country"
	QPClientId                = "client_id"
	QPClientSecret            = "client_secret"
	QPCurrentVisit            = "currentVisit"
	QPCustomData              = "customData"
	QPDeviceType              = "deviceType"
	QPEnvironment             = "environment"
	QPEventType               = "eventType"
	QPExperiment              = "experiment"
	QPExperimentId            = "id"
	QPGeolocation             = "geolocation"
	QPGoalId                  = "goalId"
	QPGrantType               = "grant_type"
	QPHref                    = "href"
	QPIndex                   = "index"
	QPKcs                     = "kcs"
	QPKey                     = "key"
	QPLatitude                = "latitude"
	QPLongitude               = "longitude"
	QPMappingIdentifier       = "mappingIdentifier"
	QPMappingValue            = "mappingValue"
	QPMaxNumberPreviousVisits = "maxNumberPreviousVisits"
	QPNegative                = "negative"
	QPNonce                   = "nonce"
	QPOs                      = "os"
	QPOsIndex                 = "osIndex"
	QPOverwrite               = "overwrite"
	QPPage                    = "page"
	QPPostalCode              = "postalCode"
	QPReferrersIndices        = "referrersIndices"
	QPRegion                  = "region"
	QPRevenue                 = "revenue"
	QPSdkName                 = "sdkName"
	QPSdkVersion              = "sdkVersion"
	QPSiteCode                = "siteCode"
	QPStaticData              = "staticData"
	QPTitle                   = "title"
	QPTimestamp               = "ts"
	QPValuesCountMap          = "valuesCountMap"
	QPVariationId             = "variationId"
	QPVersion                 = "version"
	QPVisitorCode             = "visitorCode"
)

package network

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
	replacer := strings.NewReplacer(
		"+", "%20",
		"%21", "!",
		"%2A", "*",
		"%27", "'",
		"%28", "(",
		"%29", ")",
	)
	return replacer.Replace(qb.params.Encode())
}

const (
	QPBrowserIndex            = "browserIndex"
	QPBrowserVersion          = "browserVersion"
	QPClientId                = "client_id"
	QPClientSecret            = "client_secret"
	QPCurrentVisit            = "currentVisit"
	QPCustomData              = "customData"
	QPDeviceType              = "deviceType"
	QPEnvironment             = "environment"
	QPEventType               = "eventType"
	QPExperimentId            = "id"
	QPGoalId                  = "goalId"
	QPGrantType               = "grant_type"
	QPHref                    = "href"
	QPIndex                   = "index"
	QPKey                     = "key"
	QPMaxNumberPreviousVisits = "maxNumberPreviousVisits"
	QPNegative                = "negative"
	QPNonce                   = "nonce"
	QPOverwrite               = "overwrite"
	QPReferrersIndices        = "referrersIndices"
	QPRevenue                 = "revenue"
	QPSdkName                 = "sdkName"
	QPSdkVersion              = "sdkVersion"
	QPSiteCode                = "siteCode"
	QPTitle                   = "title"
	QPTimestamp               = "ts"
	QPValuesCountMap          = "valuesCountMap"
	QPVariationId             = "variationId"
	QPVersion                 = "version"
	QPVisitorCode             = "visitorCode"
)

package network

import (
	"fmt"
	"strings"
)

const (
	trackingPath                  = "/visit/events"
	visitorDataPath               = "/visit/visitor"
	experimentsConfigurationsPath = "/visit/experimentsConfigurations"
	getDataPath                   = "/map/map"
	postDataPath                  = "/map/maps"
	configurationApiUrl           = "https://client-config.kameleoon.com/mobile"
	realTimeConfigurationUrl      = "https://events.kameleoon.com:8110/sse"
	oauthTokenUrl                 = "https://api.kameleoon.com/oauth/token"
)

const (
	DefaultDataApiUrl = "https://data.kameleoon.io"
	TestDataApiUrl    = "https://data.kameleoon.net"
)

type UrlProvider struct {
	SiteCode   string
	DataApiUrl string
	SdkName    string
	SdkVersion string
}

func (up *UrlProvider) MakeTrackingUrl(visitorCode string) string {
	qb := NewQueryBuilder()
	qb.Append(QPSdkName, up.SdkName)
	qb.Append(QPSdkVersion, up.SdkVersion)
	qb.Append(QPSiteCode, up.SiteCode)
	qb.Append(QPVisitorCode, visitorCode)
	return up.DataApiUrl + trackingPath + "?" + qb.String()
}

func (up *UrlProvider) MakeVisitorDataGetUrl(visitorCode string) string {
	qb := NewQueryBuilder()
	qb.Append(QPSiteCode, up.SiteCode)
	qb.Append(QPVisitorCode, visitorCode)
	qb.Append(QPCurrentVisit, "true")
	qb.Append(QPMaxNumberPreviousVisits, "1")
	qb.Append(QPCustomData, "true")
	qb.Append(QPVersion, "0")
	return up.DataApiUrl + visitorDataPath + "?" + qb.String()
}

func (up *UrlProvider) MakeApiDataGetRequestUrl(key string) string {
	qb := NewQueryBuilder()
	qb.Append(QPSiteCode, up.SiteCode)
	qb.Append(QPKey, key)
	return up.DataApiUrl + getDataPath + "?" + qb.String()
}

func (up *UrlProvider) MakeApiDataPostRequestUrl(key string) string {
	panic("`MakeApiDataPostRequestUrl` is not implemented!")
}

func (up *UrlProvider) MakeConfigurationUrl(environment string, ts int64) string {
	type param struct {
		name, value string
	}
	params := make([]param, 0, 3)
	params = append(params, param{name: QPSiteCode, value: up.SiteCode})
	if len(environment) > 0 {
		params = append(params, param{name: QPEnvironment, value: environment})
	}
	if ts != -1 {
		params = append(params, param{name: QPTimestamp, value: fmt.Sprint(ts)})
	}
	sb := strings.Builder{}
	sb.WriteString(configurationApiUrl)
	sb.WriteRune('?')
	for i := 0; i < len(params); i++ {
		if i > 0 {
			sb.WriteRune('&')
		}
		sb.WriteString(params[i].name)
		sb.WriteRune('=')
		sb.WriteString(params[i].value)
	}
	return sb.String()
}

func (up *UrlProvider) MakeRealTimeUrl() string {
	qb := NewQueryBuilder()
	qb.Append(QPSiteCode, up.SiteCode)
	return realTimeConfigurationUrl + "?" + qb.String()
}

func (up *UrlProvider) MakeBearerTokenUrl() string {
	return oauthTokenUrl
}

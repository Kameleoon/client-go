package network

import (
	"fmt"
	"strings"

	"github.com/Kameleoon/client-go/v3/utils"
)

const (
	trackingPath                  = "/visit/events"
	visitorDataPath               = "/visit/visitor"
	experimentsConfigurationsPath = "/visit/experimentsConfigurations"
	getDataPath                   = "/map/map"
	postDataPath                  = "/map/maps"
	configurationApiUrlFormat     = "https://%s.kameleoon.eu/sdk-config"
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
	qb := utils.NewQueryBuilder()
	qb.Append(utils.QPSdkName, up.SdkName)
	qb.Append(utils.QPSdkVersion, up.SdkVersion)
	qb.Append(utils.QPSiteCode, up.SiteCode)
	qb.Append(utils.QPVisitorCode, visitorCode)
	return up.DataApiUrl + trackingPath + "?" + qb.String()
}

func (up *UrlProvider) MakeVisitorDataGetUrl(visitorCode string) string {
	qb := utils.NewQueryBuilder()
	qb.Append(utils.QPSiteCode, up.SiteCode)
	qb.Append(utils.QPVisitorCode, visitorCode)
	qb.Append(utils.QPCurrentVisit, "true")
	qb.Append(utils.QPMaxNumberPreviousVisits, "1")
	qb.Append(utils.QPCustomData, "true")
	qb.Append(utils.QPVersion, "0")
	return up.DataApiUrl + visitorDataPath + "?" + qb.String()
}

func (up *UrlProvider) MakeApiDataGetRequestUrl(key string) string {
	qb := utils.NewQueryBuilder()
	qb.Append(utils.QPSiteCode, up.SiteCode)
	qb.Append(utils.QPKey, key)
	return up.DataApiUrl + getDataPath + "?" + qb.String()
}

func (up *UrlProvider) MakeApiDataPostRequestUrl(key string) string {
	panic("`MakeApiDataPostRequestUrl` is not implemented!")
}

func (up *UrlProvider) MakeConfigurationUrl(environment string, ts int64) string {
	type param struct {
		name, value string
	}
	params := make([]param, 0, 2)
	if len(environment) > 0 {
		params = append(params, param{name: utils.QPEnvironment, value: environment})
	}
	if ts != -1 {
		params = append(params, param{name: utils.QPTimestamp, value: fmt.Sprint(ts)})
	}
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf(configurationApiUrlFormat, up.SiteCode))
	for i := 0; i < len(params); i++ {
		if i == 0 {
			sb.WriteRune('?')
		} else {
			sb.WriteRune('&')
		}
		sb.WriteString(params[i].name)
		sb.WriteRune('=')
		sb.WriteString(params[i].value)
	}
	return sb.String()
}

func (up *UrlProvider) MakeRealTimeUrl() string {
	qb := utils.NewQueryBuilder()
	qb.Append(utils.QPSiteCode, up.SiteCode)
	return realTimeConfigurationUrl + "?" + qb.String()
}

func (up *UrlProvider) MakeBearerTokenUrl() string {
	return oauthTokenUrl
}

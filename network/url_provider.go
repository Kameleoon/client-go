package network

import (
	"fmt"
	"strings"

	"github.com/Kameleoon/client-go/v3/types"
	"github.com/Kameleoon/client-go/v3/utils"
)

const (
	DefaultDataApiDomain = "data.kameleoon.io"
	TestDataApiDomain    = "data.kameleoon.net"
	trackingPath         = "/visit/events"
	visitorDataPath      = "/visit/visitor"
	getDataPath          = "/map/map"

	configurationApiUrlFormat = "https://sdk-config.kameleoon.eu/%s"

	realTimeConfigurationUrl = "https://events.kameleoon.com:8110/sse"

	oauthTokenUrl = "https://api.kameleoon.com/oauth/token"
)

type UrlProvider interface {
	MakeTrackingUrl(visitorCode string, isUniqueIdentifier bool) string
	MakeVisitorDataGetUrl(visitorCode string, filter types.RemoteVisitorDataFilter, isUniqueIdentifier bool) string
	MakeApiDataGetRequestUrl(key string) string
	MakeConfigurationUrl(environment string, ts int64) string
	MakeRealTimeUrl() string
	MakeAccessTokenUrl() string

	ApplyDataApiDomain(dataApiDomain string)

	SiteCode() string
	DataApiDomain() string
	SdkName() string
	SdkVersion() string
}

type UrlProviderImpl struct {
	siteCode      string
	dataApiDomain string
	sdkName       string
	sdkVersion    string
}

func NewUrlProviderImpl(siteCode string, dataApiDomain string, sdkName string, sdkVersion string) *UrlProviderImpl {
	return &UrlProviderImpl{
		siteCode:      siteCode,
		dataApiDomain: dataApiDomain,
		sdkName:       sdkName,
		sdkVersion:    sdkVersion,
	}
}

func (up *UrlProviderImpl) SiteCode() string {
	return up.siteCode
}

func (up *UrlProviderImpl) DataApiDomain() string {
	return up.dataApiDomain
}

func (up *UrlProviderImpl) SdkName() string {
	return up.sdkName
}

func (up *UrlProviderImpl) SdkVersion() string {
	return up.sdkVersion
}

func getUserIdQP(isUniqueIdentifier bool) string {
	if isUniqueIdentifier {
		return utils.QPMappingValue
	}
	return utils.QPVisitorCode
}

func (up *UrlProviderImpl) MakeTrackingUrl(visitorCode string, isUniqueIdentifier bool) string {
	qb := utils.NewQueryBuilder()
	qb.Append(utils.QPSdkName, up.sdkName)
	qb.Append(utils.QPSdkVersion, up.sdkVersion)
	qb.Append(utils.QPSiteCode, up.siteCode)
	qb.Append(getUserIdQP(isUniqueIdentifier), visitorCode)
	return fmt.Sprintf("https://%s%s?%s", up.dataApiDomain, trackingPath, qb.String())
}

func (up *UrlProviderImpl) MakeVisitorDataGetUrl(
	visitorCode string,
	filter types.RemoteVisitorDataFilter,
	isUniqueIdentifier bool,
) string {
	qb := utils.NewQueryBuilder()
	qb.Append(utils.QPSiteCode, up.siteCode)
	qb.Append(getUserIdQP(isUniqueIdentifier), visitorCode)
	qb.Append(utils.QPMaxNumberPreviousVisits, fmt.Sprintf("%d", filter.PreviousVisitAmount))
	qb.Append(utils.QPVersion, "0")
	addFlagParamIfRequired(qb, utils.QPKcs, filter.Kcs)
	addFlagParamIfRequired(qb, utils.QPCurrentVisit, filter.CurrentVisit)
	addFlagParamIfRequired(qb, utils.QPCustomData, filter.CustomData)
	addFlagParamIfRequired(qb, utils.QPConversion, filter.Conversion)
	addFlagParamIfRequired(qb, utils.QPGeolocation, filter.Geolocation)
	addFlagParamIfRequired(qb, utils.QPExperiment, filter.Experiments)
	addFlagParamIfRequired(qb, utils.QPPage, filter.PageViews)
	addFlagParamIfRequired(qb, utils.QPStaticData, filter.Device || filter.Browser || filter.OperatingSystem)
	return fmt.Sprintf("https://%s%s?%s", up.dataApiDomain, visitorDataPath, qb.String())
}

func addFlagParamIfRequired(qb *utils.QueryBuilder, paramName string, state bool) {
	if state {
		qb.Append(paramName, "true")
	}
}

func (up *UrlProviderImpl) MakeApiDataGetRequestUrl(key string) string {
	qb := utils.NewQueryBuilder()
	qb.Append(utils.QPSiteCode, up.siteCode)
	qb.Append(utils.QPKey, key)
	return fmt.Sprintf("https://%s%s?%s", up.dataApiDomain, getDataPath, qb.String())
}

func (up *UrlProviderImpl) MakeConfigurationUrl(environment string, ts int64) string {
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
	sb.WriteString(fmt.Sprintf(configurationApiUrlFormat, up.siteCode))
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

func (up *UrlProviderImpl) MakeRealTimeUrl() string {
	qb := utils.NewQueryBuilder()
	qb.Append(utils.QPSiteCode, up.siteCode)
	return realTimeConfigurationUrl + "?" + qb.String()
}

func (up *UrlProviderImpl) MakeAccessTokenUrl() string {
	return oauthTokenUrl
}

func (up *UrlProviderImpl) ApplyDataApiDomain(dataApiDomain string) {
	if dataApiDomain != "" {
		up.dataApiDomain = dataApiDomain
	}
}

package kameleoon

import (
	"strings"
	"sync"
	"time"

	"github.com/Kameleoon/client-go/v3/configuration"
	"github.com/Kameleoon/client-go/v3/errs"
	"github.com/Kameleoon/client-go/v3/hybrid"
	"github.com/Kameleoon/client-go/v3/network"
	"github.com/Kameleoon/client-go/v3/network/cookie"
	"github.com/Kameleoon/client-go/v3/storage"
	"github.com/Kameleoon/client-go/v3/targeting/conditions"
	"github.com/Kameleoon/client-go/v3/types"
	"github.com/Kameleoon/client-go/v3/utils"
	"github.com/segmentio/encoding/json"
	"github.com/valyala/fasthttp"
)

type KameleoonClient interface {
	WaitInit() error

	// GetVisitorCode should be called to get the Kameleoon visitorCode for the current visitor.
	//
	// This is especially important when using Kameleoon in a mixed front-end and back-end environment,
	// where user identification consistency must be guaranteed.
	//
	// The implementation logic is described here:
	// First we check if a kameleoonVisitorCode cookie or query parameter associated with the current HTTP request can be
	// found. If so, we will use this as the visitor identifier. If no cookie / parameter is found in the current
	// request, we either randomly generate a new identifier, or use the defaultVisitorCode argument as identifier if it
	// is passed. This allows our customers to use their own identifiers as visitor codes, should they wish to.
	// This can have the added benefit of matching Kameleoon visitors with their own users without any additional
	// look-ups in a matching table.
	GetVisitorCode(request *fasthttp.Request, response *fasthttp.Response, defaultVisitorCode ...string) (string, error)

	SetLegalConsent(visitorCode string, consent bool, response ...*fasthttp.Response) error

	// AddData associate various Data to a visitor.
	//
	// Note that this method doesn't return any value and doesn't interact with the
	// Kameleoon back-end servers by itself. Instead, the declared data is saved for future sending via the flush method.
	// This reduces the number of server calls made, as data is usually grouped into a single server call triggered by
	// the execution of the flush method.
	AddData(visitorCode string, allData ...types.Data) error

	// TrackConversion on a particular goal
	//
	// This method requires visitorCode and goalID to track conversion on this particular goal.
	// In addition, this method also accepts revenue as a third optional argument to track revenue.
	// This method is non-blocking as the server call is made asynchronously.
	TrackConversion(visitorCode string, goalID int) error

	TrackConversionRevenue(visitorCode string, goalID int, revenue float64) error

	// FlushVisitor the associated data.
	//
	// The data added with the method AddData, is not directly sent to the kameleoon servers.
	// It's stored and accumulated until it is sent automatically by the TriggerExperiment or TrackConversion methods.
	// With this method you can manually send it.
	FlushVisitor(visitorCode string) error

	FlushAll()

	// GetFeatureVariationKey returns a variation key for visitor code
	//
	// This method takes a visitorCode and featureKey as mandatory arguments and
	// returns a variation assigned for a given visitor
	// If such a user has never been associated with any feature flag rules, the SDK returns a default variation key
	// You have to make sure that proper error handling is set up in your code as shown in the example to the right to catch potential exceptions.
	//
	// returns FeatureNotFound error
	// returns VisitorCodeNotValid error
	// returns FeatureEnvironmentDisabled error
	GetFeatureVariationKey(visitorCode string, featureKey string) (string, error)

	// GetFeatureVariable retrieves a feature variable value from assigned for visitor variation
	//
	// A feature variable can be changed easily via our web application.
	//
	// returns FeatureNotFound error
	// returns VisitorCodeNotValid error
	// returns FeatureEnvironmentDisabled error
	// returns FeatureVariableNotFound error
	// returns VariationNotFound error
	GetFeatureVariable(visitorCode string, featureKey string, variableKey string) (interface{}, error)

	// IsFeatureActive checks if feature is active for a visitor or not
	// (returns true / false instead of variation key)
	// This method takes a visitorCode and featureKey as mandatory arguments to check
	// if the specified feature will be active for a given user.
	// If such a user has never been associated with this feature flag, the SDK returns a boolean value randomly
	// (true if the user should have this feature or false if not).
	// You have to make sure that proper error handling is set up in your code as shown in the example to the right to catch potential exceptions.
	//
	// returns FeatureNotFound error
	// returns VisitorCodeNotValid
	IsFeatureActive(visitorCode string, featureKey string) (bool, error)

	// GetFeatureVariationVariables retrieves all feature variable values for a given variation
	//
	// This method takes a featureKey and variationKey as mandatory arguments and
	// returns a list of variables for a given variation key
	// A feature variable can be changed easily via our web application.
	//
	// returns FeatureNotFound error
	// returns FeatureEnvironmentDisabled error
	// returns VariationNotFound error
	GetFeatureVariationVariables(featureKey string, variationKey string) (map[string]interface{}, error)

	// The GetRemoteData method allows you to retrieve data (according to a key passed as
	// argument)stored on a remote Kameleoon server. Usually data will be stored on our remote servers
	// via the use of our Data API. This method, along with the availability of our highly scalable servers
	// for this purpose, provides a convenient way to quickly store massive amounts of data that
	// can be later retrieved for each of your visitors / users.
	//
	// returns Network timeout error
	GetRemoteData(key string, timeout ...time.Duration) ([]byte, error)

	GetRemoteVisitorData(visitorCode string, addData bool, timeout ...time.Duration) ([]types.Data, error)

	OnUpdateConfiguration(handler func())

	// GetFeatureList returns a list of all feature flag keys
	GetFeatureList() []string

	// GetActiveFeatureListForVisitor returns a list of active feature flag keys for a visitor
	//
	// returns VisitorCodeNotValid error when visitor code is not valid
	GetActiveFeatureListForVisitor(visitorCode string) ([]string, error)

	GetEngineTrackingCode(visitorCode string) string
}

type kameleoonClient struct {
	siteCode       string
	cfg            *KameleoonClientConfig
	visitorManager storage.VisitorManager
	networkManager network.NetworkManager
	cookieManager  cookie.CookieManager
	hybridManager  hybrid.HybridManager

	m                          sync.Mutex
	readiness                  *kameleoonClientReadiness
	token                      string
	dataFile                   *configuration.DataFile
	configurationUpdateService configuration.ConfigurationUpdateService
	updateConfigurationHandler func()
	closed                     bool
}

func newClient(siteCode string, cfg *KameleoonClientConfig) (*kameleoonClient, error) {
	if len(siteCode) == 0 {
		return nil, errs.NewSiteCodeIsEmpty("Provided siteCode is empty")
	}
	if err := cfg.defaults(); err != nil {
		return nil, err
	}
	hybridManager, err := hybrid.NewHybridManagerImpl(5 * time.Second)
	if err != nil {
		cfg.Logger.Printf("HybridManager isn't initialized properly, " +
			"GetEngineTrackingCode method isn't available for call")
	}
	np := network.NewNetProviderImpl(cfg.Network.ReadTimeout, cfg.Network.WriteTimeout,
		cfg.Network.MaxConnsPerHost, cfg.Network.ProxyURL)
	up := &network.UrlProvider{SiteCode: siteCode, DataApiUrl: cfg.dataApiUrl,
		SdkName: utils.SdkName, SdkVersion: utils.SdkVersion}
	nm := network.NewNetworkManagerImpl(cfg.Environment, cfg.DefaultTimeout, np, up, cfg.Logger)
	client := newClientInternal(siteCode, cfg, nm, hybridManager)
	return client, nil
}
func newClientInternal(siteCode string, cfg *KameleoonClientConfig, networkManager network.NetworkManager,
	hybridManager hybrid.HybridManager) *kameleoonClient {
	client := &kameleoonClient{
		siteCode:       siteCode,
		cfg:            cfg,
		readiness:      newKameleoonClientReadiness(),
		visitorManager: newVisitorManager(cfg),
		hybridManager:  hybridManager,
		networkManager: networkManager,
		cookieManager:  cookie.NewCookieManagerImpl(cfg.TopLevelDomain),
	}
	go client.updateConfigInitially()
	return client
}
func newVisitorManager(cfg *KameleoonClientConfig) storage.VisitorManager {
	return storage.NewVisitorManagerImpl(cfg.SessionDuration)
}

func (c *kameleoonClient) WaitInit() error {
	return c.readiness.Wait()
}

func (c *kameleoonClient) close() {
	if !c.closed {
		c.m.Lock()
		if c.closed {
			c.m.Unlock()
		} else {
			c.closed = true
			c.m.Unlock()
			c.visitorManager.Close()
		}
	}
}

func (c *kameleoonClient) GetVisitorCode(request *fasthttp.Request, response *fasthttp.Response,
	defaultVisitorCode ...string) (string, error) {
	return c.cookieManager.GetOrAdd(request, response, defaultVisitorCode...)
}

func (c *kameleoonClient) SetLegalConsent(visitorCode string, consent bool, response ...*fasthttp.Response) error {
	if err := utils.ValidateVisitorCode(visitorCode); err != nil {
		return err
	}
	v := c.visitorManager.GetOrCreateVisitor(visitorCode)
	v.SetLegalConsent(consent)
	if len(response) > 0 {
		c.cookieManager.Update(visitorCode, consent, response[0])
	}
	return nil
}

func (c *kameleoonClient) AddData(visitorCode string, allData ...types.Data) error {
	//var stats runtime.MemStats
	//runtime.ReadMemStats(&stats))
	if err := utils.ValidateVisitorCode(visitorCode); err != nil {
		return err
	}
	v := c.visitorManager.GetOrCreateVisitor(visitorCode)
	v.AddData(c.cfg.Logger, allData...)
	return nil
}

func (c *kameleoonClient) TrackConversion(visitorCode string, goalID int) error {
	return c.trackConversion(visitorCode, goalID)
}

func (c *kameleoonClient) TrackConversionRevenue(visitorCode string, goalID int, revenue float64) error {
	return c.trackConversion(visitorCode, goalID, revenue)
}

func (c *kameleoonClient) trackConversion(visitorCode string, goalID int, revenue ...float64) error {
	if err := utils.ValidateVisitorCode(visitorCode); err != nil {
		return err
	}
	var conv *types.Conversion
	if len(revenue) > 0 {
		conv = types.NewConversionWithRevenue(goalID, revenue[0])
	} else {
		conv = types.NewConversion(goalID)
	}
	c.AddData(visitorCode, conv)
	c.FlushVisitor(visitorCode)
	return nil
}

func (c *kameleoonClient) FlushVisitor(visitorCode string) error {
	if err := utils.ValidateVisitorCode(visitorCode); err != nil {
		return err
	}
	c.sendTrackingRequest(visitorCode, nil, true)
	return nil
}

func (c *kameleoonClient) FlushAll() {
	c.visitorManager.Enumerate(func(vc string, v storage.Visitor) bool {
		c.sendTrackingRequest(vc, v, false)
		return true
	})
}

func (c *kameleoonClient) GetFeatureVariationKey(visitorCode string, featureKey string) (string, error) {
	_, variationKey, err := c.getFeatureVariationKey(visitorCode, featureKey)
	return variationKey, err
}

// getFeatureVariationKey is a helper method for getting variation key for feature flag
func (c *kameleoonClient) getFeatureVariationKey(visitorCode string, featureKey string) (*configuration.FeatureFlag, string, error) {
	// validate that visitor code is acceptable else throw VisitorCodeNotValid exception
	if err := utils.ValidateVisitorCode(visitorCode); err != nil {
		return nil, string(types.VARIATION_OFF), err
	}
	// find enabled feature flag else return an error
	featureFlag, err := c.dataFile.GetFeatureFlag(featureKey)
	if err != nil {
		return nil, string(types.VARIATION_OFF), err
	}
	variation, rule := c.calculateVariationRuleForFeature(visitorCode, &featureFlag)
	// get variation key from feature flag
	variationKey := c.calculateVariationKey(variation, rule, &featureFlag.DefaultVariationKey)

	visitor := c.assignFeatureVariation(visitorCode, rule, variation)
	c.sendTrackingRequest(visitorCode, visitor, true)

	return &featureFlag, variationKey, nil
}
func (c *kameleoonClient) assignFeatureVariation(visitorCode string,
	rule *configuration.Rule, variation *types.VariationByExposition) storage.Visitor {
	var visitor storage.Visitor
	if rule != nil {
		if (variation != nil) && (variation.VariationID != nil) {
			visitor = c.visitorManager.GetOrCreateVisitor(visitorCode)
			asVariation := types.NewAssignedVariation(rule.ExperimentId, *variation.VariationID, rule.Type)
			visitor.AssignVariation(asVariation)
		}
	}
	return visitor
}

func (c *kameleoonClient) calculateVariationKey(varByExp *types.VariationByExposition,
	rule *configuration.Rule, defaultVariationKey *string) string {
	if varByExp != nil {
		return varByExp.VariationKey
	} else if rule != nil && rule.IsExperimentType() {
		return string(types.VARIATION_OFF)
	} else {
		return *defaultVariationKey
	}
}

// getVariationRuleForFeature is a helper method for calculate variation key for feature flag
func (c *kameleoonClient) calculateVariationRuleForFeature(
	visitorCode string, featureFlag *configuration.FeatureFlag) (*types.VariationByExposition, *configuration.Rule) {
	// no rules -> return DefaultVariationKey
	for _, rule := range featureFlag.Rules {
		//check if visitor is targeted for rule, else next rule
		if c.checkTargeting(visitorCode, featureFlag.Id, &rule) {

			// Disable searching in variation storage (uncommented if you need use variation storage)
			// check for saved variation for rule if it's experimentation rule
			// if savedVariation, found := c.getSavedVariationForRule(visitorCode, &rule); found {
			// 	return savedVariation, &rule, false
			// }

			//uses for rule exposition
			hashRule := utils.GetHashDoubleRule(visitorCode, rule.Id, rule.RespoolTime)
			//check main expostion for rule with hashRule
			if hashRule <= rule.Exposition {
				if rule.IsTargetDeliveryType() {
					var variation *types.VariationByExposition
					if len(rule.VariationByExposition) > 0 {
						variation = &rule.VariationByExposition[0]
					}
					return variation, &rule
				}
				//uses for variation's expositions
				hashVariation := utils.GetHashDoubleRule(visitorCode, rule.ExperimentId, rule.RespoolTime)
				// get variation with new hashVariation
				variation := rule.GetVariationByHash(hashVariation)
				if variation != nil {
					return variation, &rule
				}
			}
			if rule.IsTargetDeliveryType() {
				break
			}
		}
	}
	return nil, nil
}

func (c *kameleoonClient) getSavedVariationForRule(visitorCode string, rule *configuration.Rule) (*types.VariationByExposition, bool) {
	if (rule != nil) && rule.IsExperimentType() && (rule.ExperimentId != 0) {
		v := c.visitorManager.GetVisitor(visitorCode)
		if v != nil {
			if variation := v.Variations().Get(rule.ExperimentId); variation != nil {
				return rule.GetVariation(variation.VariationId()), true
			}
		}
	}
	return nil, false
}

func (c *kameleoonClient) GetFeatureVariable(visitorCode string, featureKey string, variableKey string) (interface{}, error) {
	featureFlag, variationKey, err := c.getFeatureVariationKey(visitorCode, featureKey)
	if err != nil {
		return nil, err
	}
	variation, exist := featureFlag.GetVariationByKey(variationKey)
	if !exist {
		return nil, errs.NewFeatureVariationNotFound(featureKey, variationKey)
	}
	variable, exist := variation.GetVariableByKey(variableKey)
	if !exist {
		return nil, errs.NewFeatureVariableNotFound(featureKey, variationKey, variableKey)
	}
	return parseFeatureVariableV2(variable), nil
}

func (c *kameleoonClient) IsFeatureActive(visitorCode string, featureKey string) (bool, error) {
	variationKey, err := c.GetFeatureVariationKey(visitorCode, featureKey)
	switch err.(type) {
	case *errs.FeatureEnvironmentDisabled:
		return false, nil
	default:
		return variationKey != string(types.VARIATION_OFF), err
	}
}

func (c *kameleoonClient) GetFeatureVariationVariables(featureKey string,
	variationKey string) (map[string]interface{}, error) {
	featureFlag, err := c.dataFile.GetFeatureFlag(featureKey)
	if err != nil {
		return nil, err
	}
	variation, exist := featureFlag.GetVariationByKey(variationKey)
	if !exist {
		return nil, errs.NewFeatureVariationNotFound(featureKey, variationKey)
	}
	mapVariableValues := make(map[string]interface{})
	for _, variable := range variation.Variables {
		mapVariableValues[variable.Key] = parseFeatureVariableV2(&variable)
	}
	return mapVariableValues, nil
}

func parseFeatureVariableV2(variable *types.Variable) interface{} {
	var value interface{}
	switch variable.Type {
	case "JSON":
		if valueString, ok := variable.Value.(string); ok {
			if err := json.Unmarshal([]byte(valueString), &value); err != nil {
				value = nil
			}
		}
	default:
		value = variable.Value
	}
	return value
}

func (c *kameleoonClient) GetRemoteData(key string, timeout ...time.Duration) ([]byte, error) {
	timeoutValue := time.Duration(-1)
	if len(timeout) > 0 {
		timeoutValue = timeout[0]
	}
	c.log("Retrieve data from remote source (key '%s')", key)
	out, err := c.networkManager.GetRemoteData(key, timeoutValue)
	if err != nil {
		c.log("Failed retrieve data from remote source: %v", err)
		return nil, err
	}
	return out, nil
}

func (c *kameleoonClient) GetRemoteVisitorData(visitorCode string, addData bool,
	timeout ...time.Duration) ([]types.Data, error) {
	timeoutValue := time.Duration(-1)
	if len(timeout) > 0 {
		timeoutValue = timeout[0]
	}
	out, err := c.networkManager.GetRemoteVisitorData(visitorCode, timeoutValue)
	if err != nil {
		return nil, err
	}
	var dataList []types.Data
	if dataList, err = parseCustomDataList(out); err != nil {
		return nil, err
	}
	if addData {
		err = c.AddData(visitorCode, dataList...)
	}
	return dataList, err
}

func parseCustomDataList(raw json.RawMessage) ([]types.Data, error) {
	list := remoteVisitorDataList{}
	if err := json.Unmarshal(raw, &list); err != nil {
		return nil, err
	}
	latestRecord := list.latestRecord()
	if latestRecord == nil {
		return []types.Data{}, nil
	}
	customDataList := make([]types.Data, 0, len(latestRecord.CustomDataEvents))
	for _, event := range latestRecord.CustomDataEvents {
		if event.Data != nil {
			customDataList = append(customDataList, event.Data.toCustomData())
		}
	}
	return customDataList, nil
}

type remoteVisitorCustomData struct {
	Id     int            `json:"index"`
	Values map[string]int `json:"valuesCountMap"`
}

func (rvcd *remoteVisitorCustomData) toCustomData() *types.CustomData {
	values := make([]string, 0, len(rvcd.Values))
	for v := range rvcd.Values {
		values = append(values, v)
	}
	return types.NewCustomData(rvcd.Id, values...)
}

type remoteVisitorEvent struct {
	Data *remoteVisitorCustomData `json:"data"`
}

type remoteVisitorDataVisit struct {
	CustomDataEvents []remoteVisitorEvent `json:"customDataEvents"`
}

type remoteVisitorDataList struct {
	CurrentVisit   *remoteVisitorDataVisit  `json:"currentVisit"`
	PreviousVisits []remoteVisitorDataVisit `json:"previousVisits"`
}

func (list remoteVisitorDataList) latestRecord() *remoteVisitorDataVisit {
	if list.CurrentVisit != nil {
		return list.CurrentVisit
	}
	if (list.PreviousVisits != nil) && (len(list.PreviousVisits) > 0) {
		return &list.PreviousVisits[0]
	}
	return nil
}

func (c *kameleoonClient) log(format string, args ...interface{}) {
	c.cfg.Logger.Printf(format, args...)
}

type oauthResp struct {
	Token string `json:"access_token"`
}

func (c *kameleoonClient) fetchToken() error {
	c.log("Fetching bearer token")
	out, err := c.networkManager.FetchBearerToken(c.cfg.ClientID, c.cfg.ClientSecret, -1)
	resp := oauthResp{}
	if err == nil {
		err = json.Unmarshal(out, &resp)
	}
	if err != nil {
		c.log("Failed to fetch bearer token: %v", err)
		return err
	}
	c.log("Bearer Token is fetched: %s", resp.Token)
	var b strings.Builder
	b.WriteString("Bearer ")
	b.WriteString(resp.Token)
	token := b.String()
	c.m.Lock()
	c.token = token
	c.m.Unlock()
	return nil
}

func (c *kameleoonClient) updateConfigInitially() {
	url := c.networkManager.GetUrlProvider().MakeRealTimeUrl()
	err := c.configurationUpdateService.Start(c.cfg.RefreshInterval, url, c.fetchConfig, nil, c.cfg.Logger)
	c.readiness.set(err)
}

func (c *kameleoonClient) OnUpdateConfiguration(handler func()) {
	c.updateConfigurationHandler = handler
}

func (c *kameleoonClient) fetchConfig(ts int64) error {
	// if err := c.fetchToken(); err != nil {
	// 	return err
	// }

	if clientConfig, err := c.requestClientConfig(c.siteCode, ts); err == nil {
		c.m.Lock()
		c.dataFile = configuration.NewDataFile(clientConfig, c.cfg.Environment)
		c.configurationUpdateService.UpdateSettings(clientConfig.Settings)
		c.updateCookieManagerConsentRequired()
		c.m.Unlock()
		if ts != -1 {
			c.updateConfigurationHandler()
		}
		return nil
	} else {
		c.log("Failed to fetch: %v", err)
		return err
	}
}

func (c *kameleoonClient) updateCookieManagerConsentRequired() {
	consentRequired := c.dataFile.Settings().IsConsentRequired && !c.dataFile.HasAnyTargetedDeliveryRule()
	c.cookieManager.SetConsentRequired(consentRequired)
}

func (c *kameleoonClient) requestClientConfig(siteCode string, ts int64) (configuration.Configuration, error) {
	if ts == -1 {
		c.log("Fetching configuration")
	} else {
		c.log("Fetching configuration for TS:%v", ts)
	}
	var campaigns configuration.Configuration
	out, err := c.networkManager.FetchConfiguration(ts, -1)
	if err == nil {
		if len(out) == 0 {
			err = errs.NewInternalError("Response is empty")
		} else {
			err = json.Unmarshal(out, &campaigns)
		}
	}
	if err == nil {
		c.log("Configuraiton fetched: %v", campaigns)
	} else {
		c.log("Failed to fetch client-config: %v", err)
	}
	return campaigns, err
}

/*
func (c *kameleoonClient) getValidSavedVariation(visitorCode string, experiment *configuration.Experiment) (int, bool) {
	//get saved variation
	if savedVariationId, exist := c.variationStorage.GetVariationId(visitorCode, experiment.ID); exist {
		// get actual respoolTime value for saved variation
		respoolTimeValue := 0
		for _, respoolTime := range experiment.RespoolTime {
			if respoolTime.VariationId == savedVariationId {
				respoolTimeValue = int(respoolTime.Value)
			}
		}
		// checking variation for validity along with actual respoolTime
		return c.variationStorage.IsVariationValid(visitorCode, experiment.ID, respoolTimeValue)
	}
	return 0, false
}
//*/

func (c *kameleoonClient) checkTargeting(visitorCode string, campaignId int, expOrFForRule configuration.TargetingObjectInterface) bool {
	segment := expOrFForRule.GetTargetingSegment()
	visitor := c.visitorManager.GetVisitor(visitorCode)
	return segment == nil || segment.CheckTargeting(func(targetingType types.TargetingType) interface{} {
		return c.getConditionData(targetingType, visitor, visitorCode, campaignId)
	})
}

func (c *kameleoonClient) getConditionData(targetingType types.TargetingType, visitor storage.Visitor,
	visitorCode string, campaignId int) interface{} {
	var conditionData interface{}
	switch targetingType {
	case types.TargetingCustomDatum:
		if visitor != nil {
			conditionData = visitor.CustomData()
		}
	case types.TargetingBrowser:
		if visitor != nil {
			conditionData = visitor.Browser()
		}
	case types.TargetingDeviceType:
		if visitor != nil {
			conditionData = visitor.Device()
		}
	case types.TargetingPageTitle:
		fallthrough
	case types.TargetingPageUrl:
		if visitor != nil {
			conditionData = visitor.PageViewVisits()
		}
	case types.TargetingConversions:
		if visitor != nil {
			conditionData = visitor.Conversions()
		}
	case types.TargetingVisitorCode:
		conditionData = visitorCode
	case types.TargetingSDKLanguage:
		conditionData = &types.TargetedDataSdk{Language: utils.SdkName, Version: utils.SdkVersion}
	case types.TargetingTargetExperiment:
		if visitor != nil {
			conditionData = visitor.Variations()
		}
	case types.TargetingExclusiveExperiment:
		if visitor != nil {
			conditionData = conditions.TargetedDataExclusiveExperiment{
				ExperimentId:     campaignId,
				VariationStorage: visitor.Variations(),
			}
		}
	}
	return conditionData
}

func (c *kameleoonClient) GetFeatureList() []string {
	c.m.Lock()
	defer c.m.Unlock()
	arrayKeys := make([]string, 0, len(c.dataFile.FeatureFlags()))
	for _, ff := range c.dataFile.FeatureFlags() {
		arrayKeys = append(arrayKeys, ff.FeatureKey)
	}
	return arrayKeys
}

func (c *kameleoonClient) GetActiveFeatureListForVisitor(visitorCode string) ([]string, error) {
	if err := utils.ValidateVisitorCode(visitorCode); err != nil {
		return []string{}, err
	}
	c.m.Lock()
	defer c.m.Unlock()
	arrayIds := make([]string, 0, len(c.dataFile.FeatureFlags()))
	for _, ff := range c.dataFile.FeatureFlags() {
		variation, rule := c.calculateVariationRuleForFeature(visitorCode, &ff)
		if ff.GetVariationKey(variation, rule) != string(types.VARIATION_OFF) {
			arrayIds = append(arrayIds, ff.FeatureKey)
		}
	}
	return arrayIds, nil
}

func (c *kameleoonClient) GetEngineTrackingCode(visitorCode string) string {
	if c.hybridManager == nil {
		c.log("HybridManager wasn't initialized properly. GetEngineTrackingCode method isn't avaiable")
		return ""
	}
	visitor := c.visitorManager.GetVisitor(visitorCode)
	var variations storage.DataMapStorage[int, *types.AssignedVariation]
	if visitor != nil {
		variations = visitor.Variations()
	}
	return c.hybridManager.GetEngineTrackingCode(variations)
}

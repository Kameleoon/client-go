package kameleoon

import (
	"sync"
	"time"

	"github.com/Kameleoon/client-go/v3/configuration"
	"github.com/Kameleoon/client-go/v3/errs"
	"github.com/Kameleoon/client-go/v3/managers/hybrid"
	"github.com/Kameleoon/client-go/v3/managers/remotedata"
	"github.com/Kameleoon/client-go/v3/managers/warehouse"
	"github.com/Kameleoon/client-go/v3/network"
	"github.com/Kameleoon/client-go/v3/network/cookie"
	"github.com/Kameleoon/client-go/v3/storage"
	"github.com/Kameleoon/client-go/v3/targeting"
	"github.com/Kameleoon/client-go/v3/types"
	"github.com/Kameleoon/client-go/v3/utils"
	"github.com/segmentio/encoding/json"
	"github.com/valyala/fasthttp"
)

// Deprecated: Please use `VisitorWarehouseAudienceOptParams` instead
type VisitorWarehouseAudienceParams struct {
	VisitorCode     string
	CustomDataIndex int
	WarehouseKey    string        // optional
	Timeout         time.Duration // optional
}

type VisitorWarehouseAudienceOptParams struct {
	WarehouseKey string
	Timeout      time.Duration
}

type RemoteVisitorDataOptParams struct {
	IsUniqueIdentifier bool
	Timeout            time.Duration
}

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
	TrackConversion(visitorCode string, goalID int, isUniqueIdentifier ...bool) error

	TrackConversionRevenue(visitorCode string, goalID int, revenue float64, isUniqueIdentifier ...bool) error

	// FlushVisitor the associated data.
	//
	// The data added with the method AddData, is not directly sent to the kameleoon servers.
	// It's stored and accumulated until it is sent automatically by the TriggerExperiment or TrackConversion methods.
	// With this method you can manually send it.
	FlushVisitor(visitorCode string, isUniqueIdentifier ...bool) error

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
	GetFeatureVariationKey(visitorCode string, featureKey string, isUniqueIdentifier ...bool) (string, error)

	// GetFeatureVariable retrieves a feature variable value from assigned for visitor variation
	//
	// A feature variable can be changed easily via our web application.
	//
	// returns FeatureNotFound error
	// returns VisitorCodeNotValid error
	// returns FeatureEnvironmentDisabled error
	// returns FeatureVariableNotFound error
	// returns VariationNotFound error
	GetFeatureVariable(
		visitorCode string, featureKey string, variableKey string, isUniqueIdentifier ...bool,
	) (interface{}, error)

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
	IsFeatureActive(visitorCode string, featureKey string, isUniqueIdentifier ...bool) (bool, error)

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

	// GetVisitorWarehouseAudience retrieves data associated with a visitor's warehouse audiences and adds
	// it to the visitor. Retrieves all audience data associated with the visitor in your data warehouse using the
	// specified `visitorCode` and `warehouseKey`. The `warehouseKey` is typically your internal user
	// ID. The `customDataIndex` parameter corresponds to the Kameleoon custom data that Kameleoon uses
	// to target your visitors. You can refer to the
	// warehouse targeting documentation (https://help.kameleoon.com/warehouse-audience-targeting/)
	// for additional details. The method passes the result to the returned `CustomData` object,
	// confirming that the data has been added to the visitor and is available for targeting purposes.
	//
	// Parameters:
	// - VisitorCode: The unique identifier of the visitor for whom you want to retrieve and add the
	//   data.
	// - WarehouseKey: The key to identify the warehouse data, typically your internal user ID. The values is optional.
	// - CustomDataIndex: An integer representing the index of the custom data you want to use to
	//   target your BigQuery Audiences.
	// - Timeout: Time to wait for the response
	//
	// Returns:
	// - A `CustomData` instance confirming that the data has been added to the visitor.
	// - An error if the visitor code is empty or longer than 255 characters.
	GetVisitorWarehouseAudience(params VisitorWarehouseAudienceParams) (*types.CustomData, error)

	// GetVisitorWarehouseAudienceWithOptParams retrieves data associated with a visitor's warehouse audiences and adds
	// it to the visitor. Retrieves all audience data associated with the visitor in your data warehouse using the
	// specified `visitorCode` and `warehouseKey`. The `warehouseKey` is typically your internal user
	// ID. The `customDataIndex` parameter corresponds to the Kameleoon custom data that Kameleoon uses
	// to target your visitors. You can refer to the
	// warehouse targeting documentation (https://help.kameleoon.com/warehouse-audience-targeting/)
	// for additional details. The method passes the result to the returned `CustomData` object,
	// confirming that the data has been added to the visitor and is available for targeting purposes.
	//
	// Parameters:
	// - visitorCode: The unique identifier of the visitor for whom you want to retrieve and add the
	//   data.
	// - customDataIndex: An integer representing the index of the custom data you want to use to
	//   target your BigQuery Audiences.
	// - params: An object with optional parameters: WarehouseKey and Timeout
	//
	// Returns:
	// - A `CustomData` instance confirming that the data has been added to the visitor.
	// - An error if the visitor code is empty or longer than 255 characters.
	GetVisitorWarehouseAudienceWithOptParams(
		visitorCode string, customDataIndex int, params ...VisitorWarehouseAudienceOptParams,
	) (*types.CustomData, error)

	GetRemoteVisitorData(visitorCode string, addData bool, timeout ...time.Duration) ([]types.Data, error)

	GetRemoteVisitorDataWithOptParams(
		visitorCode string, addData bool, filter types.RemoteVisitorDataFilter, params ...RemoteVisitorDataOptParams,
	) ([]types.Data, error)

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
	cfg               *KameleoonClientConfig
	visitorManager    storage.VisitorManager
	networkManager    network.NetworkManager
	cookieManager     cookie.CookieManager
	hybridManager     hybrid.HybridManager
	warehouseManager  warehouse.WarehouseManager
	targetingManager  targeting.TargetingManager
	remoteDataManager remotedata.RemoteDataManager

	m                          sync.Mutex
	readiness                  *kameleoonClientReadiness
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
	np := network.NewNetProviderImpl(cfg.Network.ReadTimeout, cfg.Network.WriteTimeout,
		cfg.Network.MaxConnsPerHost, cfg.Network.ProxyURL)
	up := network.NewUrlProviderImpl(siteCode, network.DefaultDataApiDomain, utils.SdkName, utils.SdkVersion)
	atsf := &network.AccessTokenSourceFactoryImpl{ClientId: cfg.ClientID, ClientSecret: cfg.ClientSecret}
	nm := network.NewNetworkManagerImpl(cfg.Environment, cfg.DefaultTimeout, np, up, atsf, cfg.Logger)
	vm := newVisitorManager(cfg)
	hm, err := hybrid.NewHybridManagerImpl(5 * time.Second)
	if err != nil {
		cfg.Logger.Printf("HybridManager isn't initialized properly, " +
			"GetEngineTrackingCode method isn't available for call")
	}
	tm := targeting.NewTargetingManager(vm)
	rdm := remotedata.NewRemoteDataManager(nm, vm, cfg.Logger)
	client := newClientInternal(cfg, nm, vm, hm, tm, rdm)
	return client, nil
}

func newClientInternal(
	cfg *KameleoonClientConfig,
	networkManager network.NetworkManager,
	visitorManager storage.VisitorManager,
	hybridManager hybrid.HybridManager,
	targetingManager targeting.TargetingManager,
	remoteDataManager remotedata.RemoteDataManager,
) *kameleoonClient {
	client := &kameleoonClient{
		cfg:               cfg,
		readiness:         newKameleoonClientReadiness(),
		dataFile:          configuration.NewDataFile(configuration.Configuration{}, cfg.Environment),
		visitorManager:    visitorManager,
		hybridManager:     hybridManager,
		networkManager:    networkManager,
		cookieManager:     cookie.NewCookieManagerImpl(cfg.TopLevelDomain),
		warehouseManager:  warehouse.NewWarehouseManagerImpl(networkManager, visitorManager, cfg.Logger),
		targetingManager:  targetingManager,
		remoteDataManager: remoteDataManager,
	}
	client.targetingManager.SetDataFile(client.dataFile)
	go client.updateConfigInitially()
	return client
}

func newVisitorManager(cfg *KameleoonClientConfig) storage.VisitorManager {
	return storage.NewVisitorManagerImpl(cfg.SessionDuration, nil)
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
	c.visitorManager.AddData(visitorCode, allData...)
	return nil
}

func (c *kameleoonClient) TrackConversion(visitorCode string, goalID int, isUniqueIdentifier ...bool) error {
	var isUniqueIdentifierValue bool
	if len(isUniqueIdentifier) > 0 {
		isUniqueIdentifierValue = isUniqueIdentifier[0]
	}
	return c.trackConversion(visitorCode, isUniqueIdentifierValue, goalID)
}

func (c *kameleoonClient) TrackConversionRevenue(
	visitorCode string, goalID int, revenue float64, isUniqueIdentifier ...bool,
) error {
	var isUniqueIdentifierValue bool
	if len(isUniqueIdentifier) > 0 {
		isUniqueIdentifierValue = isUniqueIdentifier[0]
	}
	return c.trackConversion(visitorCode, isUniqueIdentifierValue, goalID, revenue)
}

func (c *kameleoonClient) trackConversion(
	visitorCode string, isUniqueIdentifier bool, goalID int, revenue ...float64,
) error {
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
	c.FlushVisitor(visitorCode, isUniqueIdentifier)
	return nil
}

func (c *kameleoonClient) FlushVisitor(visitorCode string, isUniqueIdentifier ...bool) error {
	if err := utils.ValidateVisitorCode(visitorCode); err != nil {
		return err
	}
	var isUniqueIdentifierValue bool
	if len(isUniqueIdentifier) > 0 {
		isUniqueIdentifierValue = isUniqueIdentifier[0]
	}
	c.sendTrackingRequest(visitorCode, nil, true, isUniqueIdentifierValue)
	return nil
}

func (c *kameleoonClient) FlushAll() {
	c.visitorManager.Enumerate(func(vc string, v storage.Visitor) bool {
		c.sendTrackingRequest(vc, v, false, false)
		return true
	})
}

func (c *kameleoonClient) GetFeatureVariationKey(
	visitorCode string, featureKey string, isUniqueIdentifier ...bool,
) (string, error) {
	var isUniqueIdentifierValue bool
	if len(isUniqueIdentifier) > 0 {
		isUniqueIdentifierValue = isUniqueIdentifier[0]
	}
	_, variationKey, err := c.getFeatureVariationKey(visitorCode, isUniqueIdentifierValue, featureKey)
	return variationKey, err
}

// getFeatureVariationKey is a helper method for getting variation key for feature flag
func (c *kameleoonClient) getFeatureVariationKey(
	visitorCode string, isUniqueIdentifier bool, featureKey string,
) (*configuration.FeatureFlag, string, error) {
	// validate that visitor code is acceptable else throw VisitorCodeNotValid exception
	if err := utils.ValidateVisitorCode(visitorCode); err != nil {
		return nil, string(types.VariationOff), err
	}
	// find enabled feature flag else return an error
	featureFlag, err := c.dataFile.GetFeatureFlag(featureKey)
	if err != nil {
		return nil, string(types.VariationOff), err
	}
	variation, rule := c.calculateVariationRuleForFeature(visitorCode, featureFlag)
	// get variation key from feature flag
	variationKey := c.calculateVariationKey(variation, rule, &featureFlag.DefaultVariationKey)

	visitor := c.assignFeatureVariation(visitorCode, rule, variation)
	c.sendTrackingRequest(visitorCode, visitor, true, isUniqueIdentifier)

	return featureFlag, variationKey, nil
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
		return string(types.VariationOff)
	} else {
		return *defaultVariationKey
	}
}

// getVariationRuleForFeature is a helper method for calculate variation key for feature flag
func (c *kameleoonClient) calculateVariationRuleForFeature(
	visitorCode string, featureFlag *configuration.FeatureFlag,
) (*types.VariationByExposition, *configuration.Rule) {
	codeForHash := visitorCode
	if visitor := c.visitorManager.GetVisitor(visitorCode); (visitor != nil) && (visitor.MappingIdentifier() != nil) {
		codeForHash = *visitor.MappingIdentifier()
	}
	// no rules -> return DefaultVariationKey
	for _, rule := range featureFlag.Rules {
		//check if visitor is targeted for rule, else next rule
		if c.targetingManager.CheckTargeting(visitorCode, rule.ExperimentId, &rule) {

			// Disable searching in variation storage (uncommented if you need use variation storage)
			// check for saved variation for rule if it's experimentation rule
			// if savedVariation, found := c.getSavedVariationForRule(visitorCode, &rule); found {
			// 	return savedVariation, &rule, false
			// }

			//uses for rule exposition
			hashRule := utils.GetHashDoubleRule(codeForHash, rule.Id, rule.RespoolTime)
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
				hashVariation := utils.GetHashDoubleRule(codeForHash, rule.ExperimentId, rule.RespoolTime)
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

// func (c *kameleoonClient) getSavedVariationForRule(visitorCode string, rule *configuration.Rule) (*types.VariationByExposition, bool) {
// 	if (rule != nil) && rule.IsExperimentType() && (rule.ExperimentId != 0) {
// 		v := c.visitorManager.GetVisitor(visitorCode)
// 		if v != nil {
// 			if variation := v.Variations().Get(rule.ExperimentId); variation != nil {
// 				return rule.GetVariation(variation.VariationId()), true
// 			}
// 		}
// 	}
// 	return nil, false
// }

func (c *kameleoonClient) GetFeatureVariable(
	visitorCode string, featureKey string, variableKey string, isUniqueIdentifier ...bool,
) (interface{}, error) {
	var isUniqueIdentifierValue bool
	if len(isUniqueIdentifier) > 0 {
		isUniqueIdentifierValue = isUniqueIdentifier[0]
	}
	featureFlag, variationKey, err := c.getFeatureVariationKey(visitorCode, isUniqueIdentifierValue, featureKey)
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

func (c *kameleoonClient) IsFeatureActive(
	visitorCode string, featureKey string, isUniqueIdentifier ...bool,
) (bool, error) {
	variationKey, err := c.GetFeatureVariationKey(visitorCode, featureKey, isUniqueIdentifier...)
	switch err.(type) {
	case *errs.FeatureEnvironmentDisabled:
		return false, nil
	default:
		return variationKey != string(types.VariationOff), err
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
	return c.remoteDataManager.GetData(key, timeout...)
}

func (c *kameleoonClient) GetRemoteVisitorData(
	visitorCode string,
	addData bool,
	timeout ...time.Duration,
) ([]types.Data, error) {
	filter := types.DefaultRemoteVisitorDataFilter()
	return c.remoteDataManager.GetVisitorData(visitorCode, filter, addData, false, timeout...)
}

func (c *kameleoonClient) GetRemoteVisitorDataWithOptParams(
	visitorCode string, addData bool, filter types.RemoteVisitorDataFilter, params ...RemoteVisitorDataOptParams,
) ([]types.Data, error) {
	var p RemoteVisitorDataOptParams
	if len(params) > 0 {
		p = params[0]
	}
	var timeout []time.Duration
	if p.Timeout > 0 {
		timeout = []time.Duration{p.Timeout}
	}
	return c.remoteDataManager.GetVisitorData(
		visitorCode, filter, addData, p.IsUniqueIdentifier, timeout...,
	)
}

func (c *kameleoonClient) log(format string, args ...interface{}) {
	c.cfg.Logger.Printf(format, args...)
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
	if clientConfig, err := c.requestClientConfig(ts); err == nil {
		c.updateDataFile(configuration.NewDataFile(clientConfig, c.cfg.Environment))
		if ts != -1 {
			c.updateConfigurationHandler()
		}
		return nil
	} else {
		c.log("Failed to fetch: %v", err)
		return err
	}
}

func (c *kameleoonClient) updateDataFile(df *configuration.DataFile) {
	c.m.Lock()
	defer c.m.Unlock()
	c.dataFile = df
	c.configurationUpdateService.UpdateSettings(df.Settings())
	c.networkManager.GetUrlProvider().ApplyDataApiDomain(df.Settings().DataApiDomain())
	c.updateCookieManagerConsentRequired()
	c.visitorManager.SetCustomDataInfo(df.CustomDataInfo())
	c.targetingManager.SetDataFile(df)
}

func (c *kameleoonClient) updateCookieManagerConsentRequired() {
	consentRequired := c.dataFile.Settings().IsConsentRequired() && !c.dataFile.HasAnyTargetedDeliveryRule()
	c.cookieManager.SetConsentRequired(consentRequired)
}

func (c *kameleoonClient) requestClientConfig(ts int64) (configuration.Configuration, error) {
	if ts == -1 {
		c.log("Fetching configuration")
	} else {
		c.log("Fetching configuration for TS:%v", ts)
	}
	var campaigns configuration.Configuration
	out, err := c.networkManager.FetchConfiguration(ts)
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
		variation, rule := c.calculateVariationRuleForFeature(visitorCode, ff)
		if ff.GetVariationKey(variation, rule) != string(types.VariationOff) {
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

func (c *kameleoonClient) GetVisitorWarehouseAudience(params VisitorWarehouseAudienceParams) (*types.CustomData, error) {
	return c.warehouseManager.GetVisitorWarehouseAudience(
		params.VisitorCode, params.WarehouseKey, params.CustomDataIndex, params.Timeout)
}

func (c *kameleoonClient) GetVisitorWarehouseAudienceWithOptParams(
	visitorCode string, customDataIndex int, params ...VisitorWarehouseAudienceOptParams,
) (*types.CustomData, error) {
	var p VisitorWarehouseAudienceOptParams
	if len(params) > 0 {
		p = params[0]
	}
	return c.warehouseManager.GetVisitorWarehouseAudience(
		visitorCode, p.WarehouseKey, customDataIndex, p.Timeout,
	)
}

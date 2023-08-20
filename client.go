package kameleoon

import (
	"strings"
	"sync"
	"time"

	"github.com/Kameleoon/client-go/v2/configuration"
	"github.com/Kameleoon/client-go/v2/hybrid"
	"github.com/Kameleoon/client-go/v2/network"
	"github.com/Kameleoon/client-go/v2/storage"
	"github.com/Kameleoon/client-go/v2/types"
	"github.com/Kameleoon/client-go/v2/utils"
	"github.com/cornelk/hashmap"
	"github.com/segmentio/encoding/json"
)

const (
	SdkLanguage = "GO"
	SdkVersion  = "2.2.0" // IMPORTANT!!! SCRIPTS USES THIS VALUE, DO NOT RENAME/FORMAT - ONLY CHANGE VALUE.
)

const (
	REFERENCE                     = 0
	KAMELEOON_VISITOR_CODE_LENGTH = 255
)

type Client struct {
	Data             *hashmap.HashMap
	Cfg              *Config
	networkManager   network.NetworkManager
	variationStorage storage.VariationStorage
	hybridManager    hybrid.HybridManager

	m, mUA                     sync.Mutex
	init                       bool
	initError                  error
	token                      string
	experiments                []configuration.Experiment
	featureFlagsV2             []configuration.FeatureFlagV2
	userAgents                 map[string]types.UserAgent
	configurationUpdateService configuration.ConfigurationUpdateService
	updateConfigurationHandler func()
}

func NewClient(cfg *Config) *Client {
	cfg.defaults()
	hybridManager, errHybrid := hybrid.NewHybridManagerImpl(5*time.Second,
		&storage.CacheFactoryImpl{}, cfg.Logger)
	if errHybrid != nil {
		cfg.Logger.Printf("HybridManager isn't initialized properly, " +
			"GetEngineTrackingCode method isn't available for call")
	}
	np := network.NewNetProviderImpl(cfg.Network.KameleoonClient, cfg.Network.ReadTimeout, cfg.Network.WriteTimeout,
		cfg.Network.MaxConnsPerHost, cfg.Network.ProxyURL)
	up := &network.UrlProvider{SiteCode: cfg.SiteCode, DataApiUrl: cfg.dataApiUrl,
		SdkName: SdkLanguage, SdkVersion: SdkVersion}
	nm := network.NewNetworkManagerImpl(cfg.Environment, cfg.Timeout, np, up, cfg.Logger)
	return newClient(cfg, nm, hybridManager)
}

func newClient(cfg *Config, networkManager network.NetworkManager, hybridManager hybrid.HybridManager) *Client {
	c := &Client{
		Cfg:              cfg,
		networkManager:   networkManager,
		variationStorage: storage.NewVariationStorage(),
		Data:             new(hashmap.HashMap),
		userAgents:       make(map[string]types.UserAgent),
		hybridManager:    hybridManager,
	}
	go c.updateConfig()
	return c
}

func (c *Client) RunWhenReady(cb func(c *Client, err error)) {
	c.m.Lock()
	if c.init || c.initError != nil {
		c.m.Unlock()
		cb(c, c.initError)
		return
	}
	c.m.Unlock()

	t := time.NewTicker(time.Second)
	defer t.Stop()
	for range t.C {
		c.m.Lock()
		if c.init || c.initError != nil {
			c.m.Unlock()
			cb(c, c.initError)
			return
		}
		c.m.Unlock()
	}
}

// TriggerExperiment trigger an experiment.
//
// If such a visitorCode has never been associated with any variation, the SDK returns a randomly selected variation.
// If a user with a given visitor_code is already registered with a variation, it will detect the previously
// registered variation and return the variation_id.
// You have to make sure that proper error handling is set up in your code as shown in the example to the right to
// catch potential exceptions.
//
// returns ErrFeatureConfigNotFound error when experiment configuration is not found
// returns ErrNotAllocated error when visitor triggered the experiment, but did not activate it.
// returns ErrVisitorCodeNotValid error when visitor code is not valid
// Usually, this happens because the user has been associated with excluded traffic
// returns NotTargeted error when visitor is not targeted by the experiment, as the associated targeting segment conditions were not fulfilled.
// He should see the reference variation

func (c *Client) TriggerExperiment(visitorCode string, experimentId int) (int, error) {
	return c.triggerExperiment(visitorCode, experimentId)
}

func (c *Client) triggerExperiment(visitorCode string, experimentId int) (int, error) {
	if _, err := c.validateVisitorCode(visitorCode); err != nil {
		return -1, err
	}
	ex, err := c.getExperiment(experimentId)
	if err != nil {
		return -1, err
	}
	if err := c.checkSiteCodeEnable(&ex); err != nil {
		return -1, err
	}

	var variationId *int
	var noneVariation bool
	targeted := c.checkTargeting(visitorCode, ex.ID, &ex)
	if targeted {
		variationId = c.calculateVariationForExperiment(visitorCode, &ex)
		noneVariation = variationId == nil
		c.saveVariation(visitorCode, &ex.ID, variationId)
	}
	c.sendTrackingRequest(visitorCode, &ex.ID, variationId)
	if !targeted {
		return -1, newErrNotTargeted(visitorCode)
	}
	if noneVariation {
		return -1, newErrNotAllocated(visitorCode)
	}
	return *variationId, nil
}

// return parameters: first - variationId , second - should to save to variation storage
func (c *Client) calculateVariationForExperiment(visitorCode string, exp *configuration.Experiment) *int {

	// Disable searching in variation storage (uncommented if you need use variation storage)
	// if savedVariationId, exist := c.getValidSavedVariation(visitorCode, exp); exist {
	// 	return &savedVariationId, false
	// }

	threshold := getHashDouble(visitorCode, exp.ID, exp.RespoolTime)
	for _, deviation := range exp.Deviations {
		threshold -= deviation.Value
		if threshold < 0 {

			// Disable saving in variation storage (uncommented if you need use variation storage)
			// if return true as second argument the variation will be saved
			// return &deviation.VariationId, true

			return &deviation.VariationId
		}
	}
	return nil
}

// AddData associate various Data to a visitor.
//
// Note that this method doesn't return any value and doesn't interact with the
// Kameleoon back-end servers by itself. Instead, the declared data is saved for future sending via the flush method.
// This reduces the number of server calls made, as data is usually grouped into a single server call triggered by
// the execution of the flush method.
func (c *Client) AddData(visitorCode string, allData ...types.Data) error {
	// TODO think about memory size and c.Cfg.VisitorDataMaxSize
	//var stats runtime.MemStats
	//runtime.ReadMemStats(&stats))
	if _, err := c.validateVisitorCode(visitorCode); err != nil {
		return err
	}
	data := make([]types.Data, 0, len(allData))
	for _, element := range allData {
		if ua, ok := element.(*types.UserAgent); ok {
			c.addUserAgent(visitorCode, ua)
		} else {
			data = append(data, element)
		}
	}
	t := time.Now()
	td := make([]types.TargetingData, len(data))
	for i := 0; i < len(data); i++ {
		td[i] = types.TargetingData{
			LastActivityTime: t,
			Data:             data[i],
		}
	}
	actual, exist := c.Data.Get(visitorCode)
	if !exist {
		c.Data.Set(visitorCode, &types.DataCell{
			Data:  td,
			Index: make(map[int]struct{}),
		})
		return nil
	}
	cell, ok := actual.(*types.DataCell)
	if !ok {
		c.Data.Set(visitorCode, &types.DataCell{
			Data:  td,
			Index: make(map[int]struct{}),
		})
		return nil
	}
	cell.Data = append(cell.Data, td...)
	c.Data.Set(visitorCode, cell)
	return nil
}

func (c *Client) getDataCell(visitorCode string) *types.DataCell {
	val, exist := c.Data.Get(visitorCode)
	if !exist {
		return nil
	}
	cell, ok := val.(*types.DataCell)
	if !ok {
		return nil
	}
	return cell
}

// TrackConversion on a particular goal
//
// This method requires visitorCode and goalID to track conversion on this particular goal.
// In addition, this method also accepts revenue as a third optional argument to track revenue.
// The visitorCode usually is identical to the one that was used when triggering the experiment.
// This method is non-blocking as the server call is made asynchronously.
func (c *Client) TrackConversion(visitorCode string, goalID int) error {
	return c.trackConversion(visitorCode, goalID)
}

func (c *Client) TrackConversionRevenue(visitorCode string, goalID int, revenue float64) error {
	return c.trackConversion(visitorCode, goalID, revenue)
}

func (c *Client) trackConversion(visitorCode string, goalID int, revenue ...float64) error {
	if _, err := c.validateVisitorCode(visitorCode); err != nil {
		return err
	}
	conv := types.Conversion{GoalId: goalID}
	if len(revenue) > 0 {
		conv.Revenue = revenue[0]
	}
	c.AddData(visitorCode, &conv)
	c.FlushVisitor(visitorCode)
	return nil
}

// FlushVisitor the associated data.
//
// The data added with the method AddData, is not directly sent to the kameleoon servers.
// It's stored and accumulated until it is sent automatically by the TriggerExperiment or TrackConversion methods.
// With this method you can manually send it.
func (c *Client) FlushVisitor(visitorCode string) error {
	if _, err := c.validateVisitorCode(visitorCode); err != nil {
		return err
	}
	c.sendTrackingRequest(visitorCode, nil, nil)
	return nil
}

func (c *Client) FlushAll() {
	for kv := range c.Data.Iter() {
		if cell, valueOk := kv.Value.(*types.DataCell); valueOk {
			if len(cell.Data) == len(cell.Index) {
				continue
			}
			if visitorCode, keyOk := kv.Key.(string); keyOk {
				unsent, lim := c.selectUnsentData(cell)
				go func() {
					sent := c.makeTrackingRequest(visitorCode, unsent, nil, nil)
					if sent {
						c.markDataAsSent(cell, lim)
					}
				}()
			}
		}
	}
}

// GetVariationAssociatedData returns JSON Data associated with a variation.
//
// The JSON data usually represents some metadata of the variation, and can be configured on our web application
// interface or via our Automation API.
// This method takes the variationID as a parameter and will return the data as a json string.
// It will return an error if the variation ID is wrong or corresponds to an experiment that is not yet online.
//
// returns VariationNotFound error if the variation is not found.
func (c *Client) GetVariationAssociatedData(variationID int) ([]byte, error) {
	c.m.Lock()
	for _, ex := range c.experiments {
		for _, v := range ex.Variations {
			if v.ID == variationID {
				c.m.Unlock()
				return []byte(v.CustomJson), nil
			}
		}
	}
	c.m.Unlock()
	return nil, newErrVariationNotFound(utils.WriteUint(uint(variationID)))
}

// GetFeatureVariationKey returns a variation key for visitor code
//
// This method takes a visitorCode and featureKey as mandatory arguments and
// returns a variation assigned for a given visitor
// If such a user has never been associated with any feature flag rules, the SDK returns a default variation key
// You have to make sure that proper error handling is set up in your code as shown in the example to the right to catch potential exceptions.
//
// returns ErrFeatureConfigNotFound error
// returns ErrVisitorCodeNotValid
func (c *Client) GetFeatureVariationKey(visitorCode string, featureKey string) (string, error) {
	_, variationKey, err := c.getFeatureVariationKey(visitorCode, featureKey)
	return variationKey, err
}

// getFeatureVariationKey is a helper method for getting variation key for feature flag
func (c *Client) getFeatureVariationKey(visitorCode string, featureKey string) (*configuration.FeatureFlagV2, string, error) {
	// validate that visitor code is acceptable else throw VisitorCodeNotValid exception
	if _, err := c.validateVisitorCode(visitorCode); err != nil {
		return nil, string(types.VARIATION_OFF), err
	}
	//find feature flag else throw ErrFeatureConfigNotFound error
	featureFlag, err := c.getFeatureFlag(featureKey)
	if err != nil {
		return nil, string(types.VARIATION_OFF), err
	}
	variation, rule := c.calculateVariationRuleForFeature(visitorCode, &featureFlag)
	// get variation key from feature flag
	variationKey := c.calculateVariationKey(variation, rule, &featureFlag.DefaultVariationKey)

	var experimentId *int
	var variationId *int
	if rule != nil {
		experimentId = &rule.ExperimentId
		if variation != nil {
			variationId = variation.VariationID
		}
	}
	c.sendTrackingRequest(visitorCode, experimentId, variationId)
	// save variationId to variation storage if it wasn't saved before
	c.saveVariation(visitorCode, experimentId, variationId)
	return &featureFlag, variationKey, nil
}

func (c *Client) calculateVariationKey(varByExp *types.VariationByExposition,
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
func (c *Client) calculateVariationRuleForFeature(
	visitorCode string, featureFlag *configuration.FeatureFlagV2) (*types.VariationByExposition, *configuration.Rule) {
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
			hashRule := getHashDoubleRule(visitorCode, rule.Id, rule.RespoolTime)
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
				hashVariation := getHashDoubleRule(visitorCode, rule.ExperimentId, rule.RespoolTime)
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

func (c *Client) getSavedVariationForRule(visitorCode string, rule *configuration.Rule) (*types.VariationByExposition, bool) {
	if (rule != nil) && rule.IsExperimentType() && (rule.ExperimentId != 0) {
		if savedVariationId, exist := c.variationStorage.GetVariationId(visitorCode, rule.ExperimentId); exist {
			return rule.GetVariation(savedVariationId), true
		}
	}
	return nil, false
}

// GetFeatureVariable retrieves a feature variable value from assigned for visitor variation
//
// A feature variable can be changed easily via our web application.
//
// returns ErrFeatureConfigNotFound error
// returns ErrVisitorCodeNotValid
// returns ErrFeatureVariableNotFound error
// returns ErrVariationNotFound error
func (c *Client) GetFeatureVariable(visitorCode string, featureKey string, variableKey string) (interface{}, error) {
	featureFlag, variationKey, err := c.getFeatureVariationKey(visitorCode, featureKey)
	if err != nil {
		return nil, err
	}
	variation, exist := featureFlag.GetVariationByKey(variationKey)
	if !exist {
		return nil, newErrVariationNotFound(featureKey)
	}
	variable, exist := variation.GetVariableByKey(variableKey)
	if !exist {
		return nil, newErrFeatureVariableNotFound(variableKey)
	}
	return parseFeatureVariableV2(variable), nil
}

// IsFeatureActive checks if feature is active for a visitor or not
// (returns true / false instead of variation key)
// This method takes a visitorCode and featureKey as mandatory arguments to check
// if the specified feature will be active for a given user.
// If such a user has never been associated with this feature flag, the SDK returns a boolean value randomly
// (true if the user should have this feature or false if not).
// You have to make sure that proper error handling is set up in your code as shown in the example to the right to catch potential exceptions.
//
// returns ErrFeatureConfigNotFound error
// returns ErrVisitorCodeNotValid
func (c *Client) IsFeatureActive(visitorCode string, featureKey string) (bool, error) {
	variationKey, err := c.GetFeatureVariationKey(visitorCode, featureKey)
	return variationKey != string(types.VARIATION_OFF), err
}

// GetFeatureAllVariables retrieves all feature variable values for a given variation
//
// This method takes a featureKey and variationKey as mandatory arguments and
// returns a list of variables for a given variation key
// A feature variable can be changed easily via our web application.
//
// returns ErrFeatureConfigNotFound error
// returns ErrVariationNotFound error
func (c *Client) GetFeatureAllVariables(featureKey string, variationKey string) (map[string]interface{}, error) {
	featureFlag, err := c.getFeatureFlag(featureKey) //find feature flag else throw ErrFeatureConfigNotFound error
	if err != nil {
		return nil, newErrFeatureConfigNotFound(featureKey)
	}
	variation, exist := featureFlag.GetVariationByKey(variationKey)
	if !exist {
		return nil, newErrVariationNotFound(featureKey)
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

// ObtainFeatureVariable retrieves a feature variable.
//
// A feature variable can be changed easily via our web application.
//
// returns FeatureConfigurationNotFound error
// returns FeatureVariableNotFound error
// func (c *Client) ObtainFeatureVariable(featureKey interface{}, variableKey string) (interface{}, error) {
// 	ff, err := c.findFeatureFlag(featureKey)
// 	if err != nil {
// 		return nil, err
// 	}
// 	var customJson interface{}
// 	for _, v := range ff.Variations {
// 		cj := make(map[string]interface{})

// 		stringData := string(v.CustomJson[:])
// 		stringData = strings.ReplaceAll(stringData, "\\\\\\", "KameleoonTmpSymbol")
// 		stringData = strings.ReplaceAll(stringData, "\\", "")
// 		stringData = strings.ReplaceAll(stringData, "KameleoonTmpSymbol", "\\")
// 		stringData = stringData[1 : len(stringData)-1]

// 		if err = json.Unmarshal([]byte(stringData), &cj); err != nil {
// 			continue
// 		}
// 		if val, exist := cj[variableKey]; exist {
// 			customJson = val
// 		}
// 	}
// 	if customJson == nil {
// 		return nil, newErrFeatureVariableNotFound("Feature variable not found")
// 	}
// 	return parseFeatureVariable(customJson), nil
// }

// func parseFeatureVariable(customJson interface{}) interface{} {
// 	var value interface{}
// 	if mapInterface, ok := customJson.(map[string]interface{}); ok {
// 		switch mapInterface["type"] {
// 		case "JSON":
// 			if valueString, ok := mapInterface["value"].(string); ok {
// 				if err := json.Unmarshal([]byte(valueString), &value); err != nil {
// 					value = nil
// 				}
// 			}
// 		default:
// 			value = mapInterface["value"]
// 		}
// 	}
// 	return value
// }

// The GetRemoteData method allows you to retrieve data (according to a key passed as
// argument)stored on a remote Kameleoon server. Usually data will be stored on our remote servers
// via the use of our Data API. This method, along with the availability of our highly scalable servers
// for this purpose, provides a convenient way to quickly store massive amounts of data that
// can be later retrieved for each of your visitors / users.
//
// returns Network timeout error
func (c *Client) GetRemoteData(key string, timeout ...time.Duration) ([]byte, error) {
	timeoutValue := time.Duration(-1)
	if len(timeout) > 0 {
		timeoutValue = timeout[0]
	}
	outChan := make(chan json.RawMessage)
	errChan := make(chan error)
	c.log("Retrieve data from remote source (key '%s')", key)
	c.networkManager.GetRemoteData(key, timeoutValue, outChan, errChan)
	select {
	case out := <-outChan:
		return out, nil
	case err := <-errChan:
		if innerErr, isInnerErr := err.(network.ErrUnexpectedResponseStatus); isInnerErr {
			err = newErrUnexpectedStatusCode(innerErr.Code)
		}
		c.log("Failed retrieve data from remote source: %v", err)
		return nil, err
	}
}

// Deprecated: Please use `GetRemoteData`
func (c *Client) RetrieveDataFromRemoteSource(key string, timeout ...time.Duration) ([]byte, error) {
	return c.GetRemoteData(key, timeout...)
}

func (c *Client) getExperiment(id int) (configuration.Experiment, error) {
	c.m.Lock()
	defer c.m.Unlock()
	for _, ex := range c.experiments {
		if ex.ID == id {
			return ex, nil
		}
	}
	return configuration.Experiment{}, newErrExperimentConfigNotFound(utils.WritePositiveInt(id))
}

func (c *Client) getFeatureFlag(featureKey string) (configuration.FeatureFlagV2, error) {
	c.m.Lock()
	defer c.m.Unlock()
	for _, featureFlag := range c.featureFlagsV2 {
		if featureFlag.FeatureKey == featureKey {
			return featureFlag, nil
		}
	}
	return configuration.FeatureFlagV2{}, newErrFeatureConfigNotFound(featureKey)
}

func (c *Client) log(format string, args ...interface{}) {
	c.Cfg.Logger.Printf(format, args...)
}

type oauthResp struct {
	Token string `json:"access_token"`
}

func (c *Client) fetchToken() error {
	c.log("Fetching bearer token")
	outChan := make(chan json.RawMessage)
	errChan := make(chan error)
	c.networkManager.FetchBearerToken(c.Cfg.ClientID, c.Cfg.ClientSecret, time.Duration(-1), outChan, errChan)
	resp := oauthResp{}
	var out json.RawMessage
	var err error
	select {
	case out = <-outChan:
		err = json.Unmarshal(out, &resp)
	case err = <-errChan:
	}
	if err != nil {
		c.log("Failed to fetch bearer token: %v", err)
		return err
	}
	c.log("Bearer Token is fetched: %s", resp.Token)
	var b strings.Builder
	b.WriteString("Bearer ")
	b.WriteString(resp.Token)
	c.m.Lock()
	c.token = b.String()
	c.m.Unlock()
	return nil
}

func (c *Client) updateConfig() {
	url := c.networkManager.GetUrlProvider().MakeRealTimeUrl()
	err := c.configurationUpdateService.Start(c.Cfg.ConfigUpdateInterval, url, c.fetchConfig, nil, c.Cfg.Logger)
	c.m.Lock()
	c.init = true
	c.initError = err
	c.m.Unlock()
}

func (c *Client) OnUpdateConfiguration(handler func()) {
	c.updateConfigurationHandler = handler
}

func (c *Client) fetchConfig(ts int64) error {
	// if err := c.fetchToken(); err != nil {
	// 	return err
	// }

	if clientConfig, err := c.requestClientConfig(c.Cfg.SiteCode, ts); err == nil {
		c.m.Lock()
		c.configurationUpdateService.UpdateSettings(clientConfig.Settings)
		c.experiments = clientConfig.Experiments
		// c.featureFlags = clientConfig.FeatureFlags
		c.featureFlagsV2 = clientConfig.FeatureFlagsV2
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

func (c *Client) requestClientConfig(siteCode string, ts int64) (configuration.Configuration, error) {
	if ts == -1 {
		c.log("Fetching configuration")
	} else {
		c.log("Fetching configuration for TS:%v", ts)
	}
	var campaigns configuration.Configuration
	outChan := make(chan json.RawMessage)
	errChan := make(chan error)
	c.networkManager.FetchConfiguration(ts, time.Duration(-1), outChan, errChan)
	var out json.RawMessage
	var err error
	select {
	case out = <-outChan:
		if len(out) == 0 {
			err = ErrEmptyResponse
		} else {
			err = json.Unmarshal(out, &campaigns)
		}
	case err = <-errChan:
	}
	if err == nil {
		c.log("Configuraiton fetched: %v", campaigns)
	} else {
		c.log("Failed to fetch client-config: %v", err)
	}
	return campaigns, err
}

func (c *Client) checkSiteCodeEnable(campaign configuration.SiteCodeEnabledInterface) error {
	if !campaign.SiteCodeEnabled() {
		return newSiteCodeDisabled(c.Cfg.SiteCode)
	}
	return nil
}

func (c *Client) addUserAgent(visitorCode string, ua *types.UserAgent) {
	c.mUA.Lock()
	defer c.mUA.Unlock()
	if len(c.userAgents) > c.Cfg.UserAgentMaxSize {
		c.userAgents = make(map[string]types.UserAgent)
	}
	c.userAgents[visitorCode] = *ua
}

func (c *Client) getUserAgent(visitorCode string) string {
	c.mUA.Lock()
	defer c.mUA.Unlock()
	if ua, ok := c.userAgents[visitorCode]; ok {
		return ua.Value
	}
	return ""
}

func (c *Client) getValidSavedVariation(visitorCode string, experiment *configuration.Experiment) (int, bool) {
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

func (c *Client) checkTargeting(visitorCode string, campaignId int, expOrFForRule configuration.TargetingObjectInterface) bool {
	segment := expOrFForRule.GetTargetingSegment()
	return segment == nil || segment.CheckTargeting(func(targetingType types.TargetingType) interface{} {
		return c.getConditionData(targetingType, visitorCode, campaignId)
	})
}

func (c *Client) getConditionData(targetingType types.TargetingType, visitorCode string, campaignId int) interface{} {
	var conditionData interface{}
	switch targetingType {
	case types.TargetingCustomDatum:
		fallthrough
	case types.TargetingBrowser:
		fallthrough
	case types.TargetingDeviceType:
		fallthrough
	case types.TargetingPageTitle:
		fallthrough
	case types.TargetingConversions:
		fallthrough
	case types.TargetingPageUrl:
		if cell := c.getDataCell(visitorCode); cell != nil {
			conditionData = cell.Data
		} else {
			conditionData = []types.TargetingData{}
		}
	case types.TargetingVisitorCode:
		conditionData = &visitorCode
	case types.TargetingSDKLanguage:
		conditionData = &types.TargetedDataSdk{Language: SdkLanguage, Version: SdkVersion}
	case types.TargetingTargetExperiment:
		conditionData = c.variationStorage.GetMapSavedVariationId(visitorCode)
	case types.TargetingExclusiveExperiment:
		conditionData = &types.TargetedDataExclusiveExperiment{ExperimentId: campaignId,
			VisitorVariationStorage: c.variationStorage.GetMapSavedVariationId(visitorCode)}
	}
	return conditionData
}

// GetExperimentList returns a list of all experiment ids
func (c *Client) GetExperimentList() []int {
	c.m.Lock()
	defer c.m.Unlock()
	arrayIds := make([]int, 0, len(c.experiments))
	for _, exp := range c.experiments {
		arrayIds = append(arrayIds, int(exp.ID))
	}
	return arrayIds
}

// GetExperimentListForVisitor returns a list of all experiment ids targeted for a visitor
// if onlyAllocated is `true` returns a list of allocated experiments for a visitor
//
// returns ErrVisitorCodeNotValid error when visitor code is not valid
func (c *Client) GetExperimentListForVisitor(visitorCode string, onlyAllocated bool) ([]int, error) {
	if _, err := c.validateVisitorCode(visitorCode); err != nil {
		return []int{}, err
	}
	c.m.Lock()
	defer c.m.Unlock()
	arrayIds := make([]int, 0, len(c.experiments))
	for _, exp := range c.experiments {
		// experiment should be only targeted if onlyAllocated == false
		// experiment should be targeted & has variation if onlyAllocated == true
		if c.checkTargeting(visitorCode, exp.ID, &exp) {
			if onlyAllocated {
				if variationId := c.calculateVariationForExperiment(visitorCode, &exp); variationId == nil {
					continue
				}
			}
			arrayIds = append(arrayIds, exp.ID)
		}
	}
	return arrayIds, nil
}

// GetFeatureList returns a list of all feature flag keys
func (c *Client) GetFeatureList() []string {
	c.m.Lock()
	defer c.m.Unlock()
	arrayKeys := make([]string, 0, len(c.featureFlagsV2))
	for _, ff := range c.featureFlagsV2 {
		arrayKeys = append(arrayKeys, ff.FeatureKey)
	}
	return arrayKeys
}

// GetActiveFeatureListForVisitor returns a list of active feature flag keys for a visitor
//
// returns ErrVisitorCodeNotValid error when visitor code is not valid
func (c *Client) GetActiveFeatureListForVisitor(visitorCode string) ([]string, error) {
	if _, err := c.validateVisitorCode(visitorCode); err != nil {
		return []string{}, err
	}
	c.m.Lock()
	defer c.m.Unlock()
	arrayIds := make([]string, 0, len(c.featureFlagsV2))
	for _, ff := range c.featureFlagsV2 {
		variation, rule := c.calculateVariationRuleForFeature(visitorCode, &ff)
		if ff.GetVariationKey(variation, rule) != string(types.VARIATION_OFF) {
			arrayIds = append(arrayIds, ff.FeatureKey)
		}
	}
	return arrayIds, nil
}

func (c *Client) GetEngineTrackingCode(visitorCode string) string {
	if c.hybridManager == nil {
		c.log("HybridManager wasn't initialized properly. GetEngineTrackingCode method isn't avaiable")
		return ""
	}
	return c.hybridManager.GetEngineTrackingCode(visitorCode)
}

func (c *Client) saveVariation(visitorCode string, experimentId *int, variationId *int) {
	if experimentId != nil && variationId != nil {
		if c.hybridManager != nil {
			c.hybridManager.AddVariation(visitorCode, *experimentId, *variationId)
		}
		c.variationStorage.UpdateVariation(visitorCode, *experimentId, *variationId)
	}
}

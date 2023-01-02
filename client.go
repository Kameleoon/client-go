package kameleoon

import (
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/Kameleoon/client-go/v2/configuration"
	"github.com/Kameleoon/client-go/v2/storage"
	"github.com/Kameleoon/client-go/v2/types"
	"github.com/Kameleoon/client-go/v2/utils"
	"github.com/cornelk/hashmap"
	"github.com/segmentio/encoding/json"
	"github.com/valyala/fasthttp"
)

const SDKVersion = "2.0.2"

const (
	API_URL                       = "https://api.kameleoon.com"
	API_OAUTH                     = "https://api.kameleoon.com/oauth/token"
	API_SSX_URL                   = "https://api-ssx.kameleoon.com"
	API_DATA_URL                  = "https://api-data.kameleoon.com"
	API_CLIENT_CONFIG_URL         = "https://client-config.kameleoon.com"
	REFERENCE                     = 0
	KAMELEOON_VISITOR_CODE_LENGTH = 255
	USER_AGENT_MAX_COUNT          = 50000
)

type Client struct {
	Data             *hashmap.HashMap
	Cfg              *Config
	network          networkClient
	variationStorage *storage.VariationStorage

	m           sync.Mutex
	init        bool
	initError   error
	token       string
	experiments []configuration.Experiment
	// featureFlags   []configuration.FeatureFlag
	featureFlagsV2 []configuration.FeatureFlagV2
	userAgents     map[string]types.UserAgent
}

func NewClient(cfg *Config) *Client {
	c := &Client{
		Cfg:              cfg,
		network:          newNetworkClient(&cfg.Network),
		variationStorage: storage.NewVariationStorage(),
		Data:             new(hashmap.HashMap),
		userAgents:       make(map[string]types.UserAgent),
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
func (c *Client) TriggerExperiment(visitorCode string, experimentID int) (int, error) {
	return c.triggerExperiment(visitorCode, experimentID)
}

func (c *Client) triggerExperiment(visitorCode string, experimentID int) (int, error) {
	if _, err := c.validateVisitorCode(visitorCode); err != nil {
		return -1, err
	}

	var ex configuration.Experiment
	c.m.Lock()
	for i, e := range c.experiments {
		if e.ID == experimentID {
			ex = e
			break
		}
		if i == len(c.experiments)-1 {
			c.m.Unlock()
			return -1, newErrExperimentConfigNotFound(utils.WriteUint(experimentID))
		}
	}
	c.m.Unlock()
	req := trackingRequest{
		Type:         TrackingRequestExperiment,
		VisitorCode:  visitorCode,
		ExperimentID: ex.ID,
		UserAgent:    c.getUserAgent(visitorCode),
	}

	if err := c.checkSiteCodeEnable(&ex); err != nil {
		return -1, err
	}
	if !c.checkTargeting(visitorCode, ex.ID, &ex) {
		return -1, newErrNotTargeted(visitorCode)
	}

	variationId := REFERENCE
	noneVariation := true
	if savedVariationId, exist := c.getValidSavedVariation(visitorCode, &ex); exist {
		variationId = savedVariationId
		noneVariation = false
	} else {
		calculatedVariationId := c.getVariationForExperiment(ex, visitorCode)
		if calculatedVariationId != nil {
			variationId = *calculatedVariationId
			noneVariation = false
			c.variationStorage.UpdateVariation(visitorCode, experimentID, variationId)
		}
	}
	req.VariationID = variationId
	req.NoneVariation = noneVariation
	go c.postTrackingAsync(req)
	if noneVariation {
		return -1, newErrNotAllocated(visitorCode)
	}
	return variationId, nil
}

func (c *Client) getVariationForExperiment(exp configuration.Experiment, visitorCode string) *int {
	threshold := getHashDouble(visitorCode, exp.ID, exp.RespoolTime)
	for _, deviation := range exp.Deviations {
		threshold -= deviation.Value
		if threshold < 0 {
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
	conv := types.Conversion{GoalID: goalID}
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
	go c.postTrackingAsync(trackingRequest{
		Type:        TrackingRequestData,
		VisitorCode: visitorCode,
		UserAgent:   c.getUserAgent(visitorCode),
	})
	return nil
}

func (c *Client) FlushAll() {
	for kv := range c.Data.Iter() {
		if visitorCode, ok := kv.Key.(string); ok {
			c.FlushVisitor(visitorCode)
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
				return v.CustomJson, nil
			}
		}
	}
	c.m.Unlock()
	return nil, newErrVariationNotFound(utils.WriteUint(variationID))
}

// REMOVE LATER
// ActivateFeature activates a feature toggle.
//
// This method takes a visitorCode and feature_key (or featureID) as mandatory arguments to check
// if the specified feature will be active for a given user.
// If such a user has never been associated with this feature flag, the SDK returns a boolean value randomly
// (true if the user should have this feature or false if not).
// If a user with a given visitorCode is already registered with this feature flag, it will detect the previous featureFlag value.
// You have to make sure that proper error handling is set up in your code as shown in the example to the right to catch potential exceptions.
//
// returns ErrFeatureConfigNotFound error
// returns ErrNotTargeted error
// returns ErrVisitorCodeNotValid error
// func (c *Client) ActivateFeature(visitorCode string, featureKey interface{}) (bool, error) {
// 	return c.activateFeature(visitorCode, featureKey)
// }

// func (c *Client) activateFeature(visitorCode string, featureKey interface{}) (bool, error) {
// 	if _, err := c.validateVisitorCode(visitorCode); err != nil {
// 		return false, err
// 	}
// 	ff, err := c.findFeatureFlag(featureKey)
// 	if err != nil {
// 		return false, err
// 	}
// 	req := trackingRequest{
// 		Type:         TrackingRequestExperiment,
// 		VisitorCode:  visitorCode,
// 		ExperimentID: ff.ID,
// 		UserAgent:    c.getUserAgent(visitorCode),
// 	}
// 	if err := c.checkSiteCodeEnable(&ff); err != nil {
// 		return false, err
// 	}

// 	if !c.checkTargeting(visitorCode, ff.ID, ff) {
// 		return false, newErrNotTargeted(visitorCode)
// 	}

// 	if !ff.IsScheduleActive() {
// 		return false, nil
// 	}

// 	threshold := getHashDouble(visitorCode, ff.ID, nil)
// 	if threshold >= 1-ff.ExpositionRate {
// 		if len(ff.Variations) > 0 {
// 			req.VariationID = ff.Variations[0].ID
// 		}
// 		go c.postTrackingAsync(req)
// 		return true, nil
// 	}
// 	req.VariationID = REFERENCE
// 	req.NoneVariation = true
// 	go c.postTrackingAsync(req)
// 	return false, nil
// }

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
	if _, err := c.validateVisitorCode(visitorCode); err != nil { // validate that visitor code is acceptable else throw VisitorCodeNotValid exception
		return nil, string(types.VARIATION_OFF), err
	}
	featureFlag, err := c.findFeatureFlagV2(featureKey) //find feature flag else throw ErrFeatureConfigNotFound error
	if err != nil {
		return nil, string(types.VARIATION_OFF), err
	}
	variationKey, rule := c.calculateVariationKeyForFeature(&featureFlag, visitorCode)
	if rule != nil && rule.Type == string(types.RuleTypeExperimentation) { //send tracking request if rule is type of EXPERIMENTATION
		c.sendTrackingRequest(visitorCode, variationKey, rule)
	}
	return &featureFlag, variationKey, nil
}

// getVariationKeyForFeatureFlag is a helper method for calculate variation key for feature flag
func (c *Client) calculateVariationKeyForFeature(featureFlag *configuration.FeatureFlagV2, visitorCode string) (string, *configuration.Rule) {
	if len(featureFlag.Rules) > 0 { // no rules -> return DefaultVariationKey
		hashRule := getHashDoubleV2(visitorCode, featureFlag.ID, "")               //uses for rule exposition
		hashVariation := getHashDoubleV2(visitorCode, featureFlag.ID, "variation") //uses for variation's expositions
		for _, rule := range featureFlag.Rules {
			if c.checkTargeting(visitorCode, featureFlag.ID, &rule) { //check if visitor is targeted for rule, else next rule
				if hashRule < rule.Exposition { //check main expostion for rule with hashRule
					variationKey := rule.GetVariationKey(hashVariation) // get variation key with new hashVariation
					if variationKey != nil {                            // variationKey can be nil for experiment rules only, for targeted rule will be always true
						return *variationKey, &rule
					}
				} else if rule.Type == string(types.RuleTypeTargetedDelivery) { //if visitor is targeted for targeted rule then break cycle -> return default
					break
				}
			}
		}
	}
	return featureFlag.DefaultVariationKey, nil
}

// sendTrackingRequest is a helper method for sending tracking requests for new FF v2
func (c *Client) sendTrackingRequest(visitorCode string, variationKey string, rule *configuration.Rule) {
	if rule.ExperimentID != nil {
		req := trackingRequest{
			Type:         TrackingRequestExperiment,
			VisitorCode:  visitorCode,
			ExperimentID: *rule.ExperimentID,
			UserAgent:    c.getUserAgent(visitorCode),
			VariationID:  *rule.GetVariationIdByKey(variationKey),
		}
		go c.postTrackingAsync(req)
	} else {
		c.log("An attempt to send a request with null experimentId was blocked")
	}
}

// GetFeatureVariable retrieves a feature variable value from assigned for visitor variation
//
// A feature variable can be changed easily via our web application.
//
// returns ErrFeatureConfigNotFound error
// returns ErrVisitorCodeNotValid
// returns ErrFeatureVariableNotFound error
// returns ErrVariationNotFound error
func (c *Client) GetFeatureVariable(visitorCode string, featureKey string, variableName string) (interface{}, error) {
	featureFlag, variationKey, err := c.getFeatureVariationKey(visitorCode, featureKey)
	if err != nil {
		return nil, err
	}
	variation, exist := featureFlag.GetVariationByKey(variationKey)
	if !exist {
		return nil, newErrVariationNotFound(featureKey)
	}
	variable, exist := variation.GetVariableByKey(variableName)
	if !exist {
		return nil, newErrFeatureVariableNotFound(variableName)
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

// GetFeatureAllVariables retrieves all feature variable values
//
// This method takes a featureKey and variationKey as mandatory arguments and
// returns a variation assigned for a given visitor
// A feature variable can be changed easily via our web application.
//
// returns ErrFeatureConfigNotFound error
// returns ErrVariationNotFound error
func (c *Client) GetFeatureAllVariables(featureKey string, variationKey string) (map[string]interface{}, error) {
	featureFlag, err := c.findFeatureFlagV2(featureKey) //find feature flag else throw ErrFeatureConfigNotFound error
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

// The RetrieveDataFromRemoteSource method allows you to retrieve data (according to a key passed as
// argument)stored on a remote Kameleoon server. Usually data will be stored on our remote servers
// via the use of our Data API. This method, along with the availability of our highly scalable servers
// for this purpose, provides a convenient way to quickly store massive amounts of data that
// can be later retrieved for each of your visitors / users.
//
// returns Network timeout error
func (c *Client) RetrieveDataFromRemoteSource(key string, timeout ...time.Duration) ([]byte, error) {
	r := request{
		URL:          c.buildAPIDataPath(key),
		Method:       MethodGet,
		ContentType:  HeaderContentTypeJson,
		ClientHeader: c.Cfg.Network.KameleoonClient,
	}
	if len(timeout) > 0 {
		r.Timeout = timeout[0]
	} else {
		r.Timeout = DefaultDoTimeout
	}
	var data []byte
	cb := func(resp *fasthttp.Response, err error) error {
		if err != nil {
			return err
		}
		if resp.StatusCode() >= fasthttp.StatusBadRequest {
			return ErrBadStatus
		}
		data = resp.Body()
		return nil
	}
	c.log("Retrieve data from remote source: %v", r)
	if err := c.network.Do(r, cb); err != nil {
		c.log("Failed retrieve data from remote source: %v", err)
		return nil, err
	}

	return data, nil
}

func (c *Client) buildAPIDataPath(key string) string {
	var url strings.Builder
	url.WriteString(API_DATA_URL)
	url.WriteString("/data?siteCode=")
	url.WriteString(c.Cfg.SiteCode)
	url.WriteString("&")
	url.WriteString(types.EncodeURIComponent("key", key))
	return url.String()
}

// func (c *Client) findFeatureFlag(featureKey interface{}) (configuration.FeatureFlag, error) {
// 	var flag configuration.FeatureFlag

// 	c.m.Lock()
// 	switch key := featureKey.(type) {
// 	case string:
// 		for i, featureFlag := range c.featureFlags {
// 			if featureFlag.IdentificationKey == key {
// 				flag = featureFlag
// 				break
// 			}
// 			if i == len(c.featureFlags)-1 {
// 				c.m.Unlock()
// 				return flag, newErrFeatureConfigNotFound(key)
// 			}
// 		}
// 	case int:
// 		for i, featureFlag := range c.featureFlags {
// 			if featureFlag.ID == key {
// 				flag = featureFlag
// 				break
// 			}
// 			if i == len(c.featureFlags)-1 {
// 				c.m.Unlock()
// 				return flag, newErrFeatureConfigNotFound(strconv.Itoa(key))
// 			}
// 		}
// 	default:
// 		c.m.Unlock()
// 		return flag, ErrInvalidFeatureKeyType
// 	}
// 	c.m.Unlock()

// 	return flag, nil
// }

func (c *Client) findFeatureFlagV2(featureKey string) (configuration.FeatureFlagV2, error) {
	var flag configuration.FeatureFlagV2
	c.m.Lock()
	defer c.m.Unlock()
	for _, featureFlag := range c.featureFlagsV2 {
		if featureFlag.FeatureKey == featureKey {
			return featureFlag, nil
		}
	}
	return flag, newErrFeatureConfigNotFound(featureKey)
}

// func (c *Client) GetExperiment(id int) *configuration.Experiment {
// 	c.m.Lock()
// 	for i, ex := range c.experiments {
// 		if ex.ID == id {
// 			c.m.Unlock()
// 			return &c.experiments[i]
// 		}
// 	}
// 	c.m.Unlock()
// 	return nil
// }

// func (c *Client) GetFeatureFlag(id int) *configuration.FeatureFlag {
// 	c.m.Lock()
// 	for i, ff := range c.featureFlags {
// 		if ff.ID == id {
// 			c.m.Unlock()
// 			return &c.featureFlags[i]
// 		}
// 	}
// 	c.m.Unlock()
// 	return nil
// }

func (c *Client) log(format string, args ...interface{}) {
	if !c.Cfg.VerboseMode {
		return
	}
	if len(args) == 0 {
		c.Cfg.Logger.Printf(format)
		return
	}
	c.Cfg.Logger.Printf(format, args...)
}

type oauthResp struct {
	Token string `json:"access_token"`
}

func (c *Client) fetchToken() error {
	c.log("Fetching bearer token")
	form := url.Values{
		"grant_type":    []string{"client_credentials"},
		"client_id":     []string{c.Cfg.ClientID},
		"client_secret": []string{c.Cfg.ClientSecret},
	}
	resp := oauthResp{}
	r := request{
		Method:      MethodPost,
		URL:         API_OAUTH,
		ContentType: HeaderContentTypeForm,
		BodyString:  form.Encode(),
	}

	err := c.network.Do(r, respCallbackJson(&resp))
	if err != nil {
		c.log("Failed to fetch bearer token: %v", err)
		return err
	} else {
		c.log("Bearer Token is fetched: %s", resp.Token)
	}
	var b strings.Builder
	b.WriteString("Bearer ")
	b.WriteString(resp.Token)
	c.m.Lock()
	c.token = b.String()
	c.m.Unlock()
	return nil
}

func (c *Client) updateConfig() {
	c.log("Start-up, fetching is starting")
	err := c.fetchConfig()
	c.m.Lock()
	c.init = true
	c.initError = err
	c.m.Unlock()
	ticker := time.NewTicker(c.Cfg.ConfigUpdateInterval)
	c.log("Scheduled job to fetch configuration is starting")
	for range ticker.C {
		c.fetchConfig()
	}
}

func (c *Client) fetchConfig() error {
	// if err := c.fetchToken(); err != nil {
	// 	return err
	// }

	if clientConfig, err := c.requestClientConfig(c.Cfg.SiteCode); err == nil {
		c.m.Lock()
		c.experiments = clientConfig.Experiments
		// c.featureFlags = clientConfig.FeatureFlags
		c.featureFlagsV2 = clientConfig.FeatureFlagsV2
		c.m.Unlock()
		return nil
	} else {
		c.log("Failed to fetch: %v", err)
		return err
	}
}

func (c *Client) requestClientConfig(siteCode string) (configuration.Configuration, error) {
	c.log("Fetching Configuration")
	var campaigns configuration.Configuration
	uri, err := c.buildFetchPathClientConfig(API_CLIENT_CONFIG_URL, c.Cfg.SiteCode, c.Cfg.Environment)
	if err != nil {
		return campaigns, err
	}
	req := request{
		Method:      MethodGet,
		URL:         uri,
		ContentType: HeaderContentTypeJson,
	}
	cb := func(resp *fasthttp.Response, err error) error {
		if err != nil {
			return err
		}
		b := resp.Body()
		if len(b) == 0 {
			return ErrEmptyResponse
		}
		var res configuration.Configuration
		err = json.Unmarshal(b, &res)
		if err != nil {
			return err
		}
		campaigns = res
		return nil
	}
	err = c.network.Do(req, cb)
	if err != nil {
		c.log("Failed to fetch: %v, request: %v", err, req)
	}

	c.log("Configuraiton fetched: %v", campaigns)
	return campaigns, err
}

func (c *Client) buildFetchPathClientConfig(base string, siteCode string, environment string) (string, error) {
	var buf strings.Builder
	buf.WriteString(base)
	buf.WriteString("/mobile")
	isFirst := true
	addValue := func(name string, value string) {
		if !isFirst {
			buf.WriteByte('&')
		} else {
			buf.WriteByte('?')
			isFirst = false
		}
		buf.WriteString(name)
		buf.WriteByte('=')
		buf.WriteString(value)
	}
	addValue("siteCode", siteCode)
	if len(environment) > 0 {
		addValue("environment", environment)
	}
	return buf.String(), nil
}

func (c *Client) checkSiteCodeEnable(campaign configuration.SiteCodeEnabledInterface) error {
	if !campaign.SiteCodeEnabled() {
		return newSiteCodeDisabled(c.Cfg.SiteCode)
	}
	return nil
}

func (c *Client) addUserAgent(visitorCode string, ua *types.UserAgent) {
	if len(c.userAgents) > USER_AGENT_MAX_COUNT {
		for k := range c.userAgents {
			delete(c.userAgents, k)
		}
	}
	c.userAgents[visitorCode] = *ua
}

func (c *Client) getUserAgent(visitorCode string) string {
	if ua, ok := c.userAgents[visitorCode]; ok {
		return ua.Value
	}
	return ""
}

func (c *Client) getValidSavedVariation(visitorCode string, experiment *configuration.Experiment) (int, bool) {
	if savedVariationId, exist := c.variationStorage.GetVariationId(visitorCode, experiment.ID); exist {
		respoolTimeValue := 0
		for _, respoolTime := range experiment.RespoolTime {
			if respoolTime.VariationId == savedVariationId {
				respoolTimeValue = int(respoolTime.Value)
			}
		}
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
		if cell := c.getDataCell(visitorCode); cell != nil {
			conditionData = cell.Data
		} else {
			conditionData = []types.TargetingData{}
		}
	case types.TargetingTargetExperiment:
		conditionData = c.variationStorage.GetMapSavedVariationId(visitorCode)
	case types.TargetingExclusiveExperiment:
		conditionData = &types.ExclusiveExperimentTargetedData{ExperimentId: campaignId,
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
		arrayIds = append(arrayIds, exp.ID)
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
		if c.checkTargeting(visitorCode, exp.ID, &exp) &&
			(!onlyAllocated || c.getVariationForExperiment(exp, visitorCode) != nil) {
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
		variationKey, _ := c.calculateVariationKeyForFeature(&ff, visitorCode)
		if variationKey != string(types.VARIATION_OFF) {
			arrayIds = append(arrayIds, ff.FeatureKey)
		}
	}
	return arrayIds, nil
}

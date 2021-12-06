package kameleoon

import (
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cornelk/hashmap"
	"github.com/segmentio/encoding/json"
	"github.com/valyala/fasthttp"

	"github.com/Kameleoon/client-go/targeting"
	"github.com/Kameleoon/client-go/types"
	"github.com/Kameleoon/client-go/utils"
)

const sdkVersion = "1.0.3"

const (
	API_URL     = "https://api.kameleoon.com"
	API_OAUTH   = "https://api.kameleoon.com/oauth/token"
	API_SSX_URL = "https://api-ssx.kameleoon.com"
	REFERENCE   = "0"
)

type Client struct {
	Data *hashmap.HashMap
	Cfg  *Config
	rest restClient

	m            sync.Mutex
	init         bool
	initError    error
	token        string
	experiments  []types.Experiment
	featureFlags []types.FeatureFlag
}

func NewClient(cfg *Config) *Client {
	c := &Client{
		Cfg:  cfg,
		rest: newRESTClient(&cfg.REST),
		Data: new(hashmap.HashMap),
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
// returns ExperimentConfigurationNotFound error when experiment configuration is not found
// returns NotActivated error when visitor triggered the experiment, but did not activate it.
// Usually, this happens because the user has been associated with excluded traffic
// returns NotTargeted error when visitor is not targeted by the experiment, as the associated targeting segment conditions were not fulfilled.
// He should see the reference variation
func (c *Client) TriggerExperiment(visitorCode string, experimentID int) (int, error) {
	return c.triggerExperiment(visitorCode, experimentID)
}

func (c *Client) TriggerExperimentTimeout(visitorCode string, experimentID int, timeout time.Duration) (int, error) {
	return c.triggerExperiment(visitorCode, experimentID, timeout)
}

func (c *Client) triggerExperiment(visitorCode string, experimentID int, timeout ...time.Duration) (int, error) {
	var ex types.Experiment
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
	}
	if !c.Cfg.BlockingMode {
		var data []types.TargetingData
		if cell := c.getDataCell(visitorCode); cell != nil {
			data = cell.Data
		}
		segment, ok := ex.TargetingSegment.(*targeting.Segment)
		if ok && !segment.CheckTargeting(data) {
			return -1, newErrNotTargeted(visitorCode)
		}

		threshold := getHashDouble(ex.ID, visitorCode, ex.RespoolTime)
		keys := make([]int, 0, len(ex.Deviations))
		for k := range ex.Deviations {
			keyInt, _ := strconv.Atoi(k)
			keys = append(keys, keyInt)
		}
		sort.Ints(keys)
		for _, kInt := range keys {
			k := strconv.Itoa(kInt)
			v := ex.Deviations[k]
			threshold -= v
			if threshold >= 0 {
				continue
			}
			req.VariationID = k
			go c.postTrackingAsync(req)
			return utils.ParseUint(k)
		}

		req.VariationID = REFERENCE
		req.NoneVariation = true
		go c.postTrackingAsync(req)
		return -1, newErrNotActivated(visitorCode)
	}

	data := c.selectSendData()
	var sb strings.Builder
	for _, dataCell := range data {
		for i := 0; i < len(dataCell.Data); i++ {
			if _, exist := dataCell.Index[i]; exist {
				continue
			}
			sb.WriteString(dataCell.Data[i].QueryEncode())
			sb.WriteByte('\n')
		}
	}

	r := request{
		URL:          c.buildTrackingPath(c.Cfg.TrackingURL, req),
		Method:       MethodPost,
		ContentType:  HeaderContentTypeText,
		ClientHeader: c.Cfg.TrackingVersion,
	}
	if len(timeout) > 0 {
		r.Timeout = timeout[0]
	}
	var id string
	cb := func(resp *fasthttp.Response, err error) error {
		if err != nil {
			return err
		}
		if resp.StatusCode() >= fasthttp.StatusBadRequest {
			return ErrBadStatus
		}
		id = string(resp.Body())
		return err
	}
	c.log("Trigger experiment request: %v", r)
	if err := c.rest.Do(r, cb); err != nil {
		c.log("Failed to trigger experiment: %v", err)
		return -1, err
	}
	switch id {
	case "", "null":
		return -1, newErrNotTargeted(visitorCode)
	case "0":
		return -1, newErrNotActivated(visitorCode)
	}

	return utils.ParseUint(id)
}

// AddData associate various Data to a visitor.
//
// Note that this method doesn't return any value and doesn't interact with the
// Kameleoon back-end servers by itself. Instead, the declared data is saved for future sending via the flush method.
// This reduces the number of server calls made, as data is usually grouped into a single server call triggered by
// the execution of the flush method.
func (c *Client) AddData(visitorCode string, data ...types.Data) {
	// TODO think about memory size and c.Cfg.VisitorDataMaxSize
	//var stats runtime.MemStats
	//runtime.ReadMemStats(&stats)
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
		return
	}
	cell, ok := actual.(*types.DataCell)
	if !ok {
		c.Data.Set(visitorCode, &types.DataCell{
			Data:  td,
			Index: make(map[int]struct{}),
		})
		return
	}
	cell.Data = append(cell.Data, td...)
	c.Data.Set(visitorCode, cell)
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
func (c *Client) TrackConversion(visitorCode string, goalID int) {
	c.trackConversion(visitorCode, goalID)
}

func (c *Client) TrackConversionRevenue(visitorCode string, goalID int, revenue float64) {
	c.trackConversion(visitorCode, goalID, revenue)
}

func (c *Client) trackConversion(visitorCode string, goalID int, revenue ...float64) {
	conv := types.Conversion{GoalID: goalID}
	if len(revenue) > 0 {
		conv.Revenue = revenue[0]
	}
	c.AddData(visitorCode, &conv)
	c.FlushVisitor(visitorCode)
}

// FlushVisitor the associated data.
//
// The data added with the method AddData, is not directly sent to the kameleoon servers.
// It's stored and accumulated until it is sent automatically by the TriggerExperiment or TrackConversion methods.
// With this method you can manually send it.
func (c *Client) FlushVisitor(visitorCode string) {
	go c.postTrackingAsync(trackingRequest{
		Type:        TrackingRequestData,
		VisitorCode: visitorCode,
	})
}

func (c *Client) FlushAll() {
	go c.postTrackingAsync(trackingRequest{
		Type: TrackingRequestData,
	})
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

// ActivateFeature activates a feature toggle.
//
// This method takes a visitorCode and feature_key (or featureID) as mandatory arguments to check
// if the specified feature will be active for a given user.
// If such a user has never been associated with this feature flag, the SDK returns a boolean value randomly
// (true if the user should have this feature or false if not).
// If a user with a given visitorCode is already registered with this feature flag, it will detect the previous featureFlag value.
// You have to make sure that proper error handling is set up in your code as shown in the example to the right to catch potential exceptions.
//
// returns FeatureConfigurationNotFound error
// returns NotTargeted error
func (c *Client) ActivateFeature(visitorCode string, featureKey interface{}) (bool, error) {
	return c.activateFeature(visitorCode, featureKey)
}

func (c *Client) ActivateFeatureTimeout(visitorCode string, featureKey interface{}, timeout time.Duration) (bool, error) {
	return c.activateFeature(visitorCode, featureKey, timeout)
}

func (c *Client) activateFeature(visitorCode string, featureKey interface{}, timeout ...time.Duration) (bool, error) {
	ff, err := c.getFeatureFlag(featureKey)
	if err != nil {
		return false, err
	}
	req := trackingRequest{
		Type:         TrackingRequestExperiment,
		VisitorCode:  visitorCode,
		ExperimentID: ff.ID,
	}
	if !c.Cfg.BlockingMode {
		var data []types.TargetingData
		if cell := c.getDataCell(visitorCode); cell != nil {
			data = cell.Data
		}

		segment, ok := ff.TargetingSegment.(*targeting.Segment)
		if ok && !segment.CheckTargeting(data) {
			return false, newErrNotTargeted(visitorCode)
		}

		threshold := getHashDouble(ff.ID, visitorCode, nil)
		if threshold <= ff.ExpositionRate {
			if len(ff.VariationsID) > 0 {
				req.VariationID = utils.WriteUint(ff.VariationsID[0])
			}
			go c.postTrackingAsync(req)
			return true, nil
		}
		req.VariationID = REFERENCE
		req.NoneVariation = true
		go c.postTrackingAsync(req)
		return false, nil
	}

	data := c.selectSendData()
	var sb strings.Builder
	for _, dataCell := range data {
		for i := 0; i < len(dataCell.Data); i++ {
			if _, exist := dataCell.Index[i]; exist {
				continue
			}
			sb.WriteString(dataCell.Data[i].QueryEncode())
			sb.WriteByte('\n')
		}
	}
	r := request{
		URL:          c.buildTrackingPath(c.Cfg.TrackingURL, req),
		Method:       MethodPost,
		ContentType:  HeaderContentTypeText,
		ClientHeader: c.Cfg.TrackingVersion,
	}
	if len(timeout) > 0 {
		r.Timeout = timeout[0]
	}
	var result string
	cb := func(resp *fasthttp.Response, err error) error {
		if err != nil {
			return err
		}
		if resp.StatusCode() >= fasthttp.StatusBadRequest {
			return ErrBadStatus
		}
		result = string(resp.Body())
		return err
	}
	c.log("Activate feature request: %v", r)
	if err = c.rest.Do(r, cb); err != nil {
		c.log("Failed to get activation: %v", err)
		return false, err
	}
	switch result {
	case "", "null":
		return false, newErrFeatureConfigNotFound(visitorCode)
	}
	return true, nil
}

// GetFeatureVariable retrieves a feature variable.
//
// A feature variable can be changed easily via our web application.
//
// returns FeatureConfigurationNotFound error
// returns FeatureVariableNotFound error
func (c *Client) GetFeatureVariable(featureKey interface{}, variableKey string) (interface{}, error) {
	ff, err := c.getFeatureFlag(featureKey)
	if err != nil {
		return nil, err
	}
	var customJson interface{}
	for _, v := range ff.Variations {
		cj := make(map[string]interface{})

		stringData := string(v.CustomJson[:])
		stringData = strings.ReplaceAll(stringData, "\\\\\\", "KameleoonTmpSymbol")
		stringData = strings.ReplaceAll(stringData, "\\", "")
		stringData = strings.ReplaceAll(stringData, "KameleoonTmpSymbol", "\\")
		stringData = stringData[1 : len(stringData)-1]

		if err = json.Unmarshal([]byte(stringData), &cj); err != nil {
			continue
		}
		if val, exist := cj[variableKey]; exist {
			customJson = val
		}
	}
	if customJson == nil {
		return nil, newErrFeatureVariableNotFound("Feature variable not found")
	}
	return c.parseFeatureVariable(customJson), nil
}

func (c *Client) parseFeatureVariable(customJson interface{}) interface{} {
	var value interface{}
	if mapInterface, ok := customJson.(map[string]interface{}); ok {
		switch mapInterface["type"] {
		case "Boolean":
			value, _ = strconv.ParseBool(mapInterface["value"].(string))
		case "Number":
			value, _ = strconv.Atoi(mapInterface["value"].(string))
		case "String":
			value = mapInterface["value"]
		case "JSON":
			if valueString, ok := mapInterface["value"].(string); ok {
				if err := json.Unmarshal([]byte(valueString), &value); err != nil {
					value = nil
				}
			}
		default:
			value = nil
		}
	}
	return value
}

func (c *Client) getFeatureFlag(featureKey interface{}) (types.FeatureFlag, error) {
	var flag types.FeatureFlag

	c.m.Lock()
	switch key := featureKey.(type) {
	case string:
		for i, featureFlag := range c.featureFlags {
			if featureFlag.IdentificationKey == key {
				flag = featureFlag
				break
			}
			if i == len(c.featureFlags)-1 {
				c.m.Unlock()
				return flag, newErrFeatureConfigNotFound(key)
			}
		}
	case int:
		for i, featureFlag := range c.featureFlags {
			if featureFlag.ID == key {
				flag = featureFlag
				break
			}
			if i == len(c.featureFlags)-1 {
				c.m.Unlock()
				return flag, newErrFeatureConfigNotFound(strconv.Itoa(key))
			}
		}
	default:
		c.m.Unlock()
		return flag, ErrInvalidFeatureKeyType
	}
	c.m.Unlock()

	return flag, nil
}

func (c *Client) GetExperiment(id int) *types.Experiment {
	c.m.Lock()
	for i, ex := range c.experiments {
		if ex.ID == id {
			c.m.Unlock()
			return &c.experiments[i]
		}
	}
	c.m.Unlock()
	return nil
}

func (c *Client) GetFeatureFlag(id int) *types.FeatureFlag {
	c.m.Lock()
	for i, ff := range c.featureFlags {
		if ff.ID == id {
			c.m.Unlock()
			return &c.featureFlags[i]
		}
	}
	c.m.Unlock()
	return nil
}

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

	err := c.rest.Do(r, respCallbackJson(&resp))
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
	if err != nil {
		c.log("Failed to fetch: %v", err)
		return
	}
	ticker := time.NewTicker(c.Cfg.ConfigUpdateInterval)
	c.log("Scheduled job to fetch configuration is starting")
	for range ticker.C {
		err = c.fetchConfig()
		if err != nil {
			c.log("Failed to fetch: %v", err)
			return
		}
	}
}

func (c *Client) fetchConfig() error {
	if err := c.fetchToken(); err != nil {
		return err
	}
	// Rest API
	// siteID, err := c.fetchSiteID()
	// if err != nil {
	// 	return err
	// }
	// experiments, err := c.fetchExperiments(siteID)
	// if err != nil {
	// 	return err
	// }
	// featureFlags, err := c.fetchFeatureFlags(siteID)

	//GraphQL
	experiments, err := c.fetchExperimentsGraphQL(c.Cfg.SiteCode)
	fmt.Println(experiments)
	if err != nil {
		return err
	}
	featureFlags, err := c.fetchFeatureFlagsGraphQL(c.Cfg.SiteCode)
	if err != nil {
		return err
	}

	c.m.Lock()
	c.experiments = append(c.experiments, experiments...)
	c.featureFlags = append(c.featureFlags, featureFlags...)
	c.m.Unlock()
	return nil
}

func (c *Client) fetchSite() (*types.SiteResponse, error) {
	c.log("Fetching site")
	filter := []fetchFilter{{
		Field:      "code",
		Operator:   "EQUAL",
		Parameters: []string{c.Cfg.SiteCode},
	}}
	res := []types.SiteResponse{{}}
	cb := func(resp *fasthttp.Response, err error) error {
		if err != nil {
			return err
		}
		b := resp.Body()
		if len(b) == 0 {
			return ErrEmptyResponse
		}
		if b[0] == '[' {
			return json.Unmarshal(b, &res)
		}
		return json.Unmarshal(b, &res[0])
	}
	err := c.fetchOne("/sites", fetchQuery{PerPage: 1}, filter, cb)
	if err != nil {
		return nil, err
	}
	if len(res) == 0 {
		return nil, ErrEmptyResponse
	}
	c.log("Sites are fetched: %v", res)
	return &res[0], err
}

type siteResponseID struct {
	ID int `json:"id"`
}

func (c *Client) fetchSiteID() (int, error) {
	c.log("Fetching site id")
	filter := []fetchFilter{{
		Field:      "code",
		Operator:   "EQUAL",
		Parameters: []string{c.Cfg.SiteCode},
	}}
	res := []siteResponseID{{}}
	cb := func(resp *fasthttp.Response, err error) error {
		if err != nil {
			return err
		}
		b := resp.Body()
		if len(b) == 0 {
			return ErrEmptyResponse
		}
		if b[0] == '[' {
			return json.Unmarshal(b, &res)
		}
		return json.Unmarshal(b, &res[0])
	}
	err := c.fetchOne("/sites", fetchQuery{PerPage: 1}, filter, cb)
	if len(res) == 0 {
		return -1, ErrEmptyResponse
	}
	c.log("Sites are fetched: %v", res)
	return res[0].ID, err
}

func (c *Client) fetchExperimentsGraphQL(siteCode string, perPage ...int) ([]types.Experiment, error) {
	c.log("Fetching experiments")
	pp := -1
	if len(perPage) > 0 {
		pp = perPage[0]
	}
	var ex []types.Experiment
	cb := func(resp *fasthttp.Response, err error) error {
		if err != nil {
			return err
		}
		b := resp.Body()
		if len(b) == 0 {
			return ErrEmptyResponse
		}
		var res ExperimentDataGraphQL
		fmt.Println(string(b))
		err = json.Unmarshal(b, &res)
		if err != nil {
			fmt.Println(err)
			return err
		}
		for _, expQL := range res.Data.Experiments.Edge {
			ex = append(ex, expQL.Transform())
		}
		return nil
	}
	err := c.fetchAllGraphQL(GetExperimentsGraphQL(siteCode), fetchQuery{PerPage: pp}, cb)
	c.log("Experiment are fetched: %v", ex)
	return ex, err
}

func (c *Client) fetchExperiments(siteID int, perPage ...int) ([]types.Experiment, error) {
	c.log("Fetching experiments")
	pp := -1
	if len(perPage) > 0 {
		pp = perPage[0]
	}
	var ex []types.Experiment
	filters := []fetchFilter{
		{
			Field:      "siteId",
			Operator:   "EQUAL",
			Parameters: []string{utils.WriteUint(siteID)},
		},
		{
			Field:      "status",
			Operator:   "IN",
			Parameters: []string{"ACTIVE", "DEVIATED"},
		},
		{
			Field:      "type",
			Operator:   "IN",
			Parameters: []string{string(types.ExperimentTypeServerSide), string(types.ExperimentTypeHybrid)},
		},
	}
	cb := func(resp *fasthttp.Response, err error) error {
		if err != nil {
			return err
		}
		b := resp.Body()
		if len(b) == 0 {
			return ErrEmptyResponse
		}
		res := []types.Experiment{{}}
		if b[0] == '[' {
			err = json.Unmarshal(b, &res)
		} else {
			err = json.Unmarshal(b, &res[0])
		}
		if err != nil {
			return err
		}
		ex = append(ex, res...)
		return nil
	}
	err := c.fetchAll("/experiments", fetchQuery{PerPage: pp}, filters, cb)
	for i := 0; i < len(ex); i++ {
		err = c.completeExperiment(&ex[i])
		if err != nil {
			return nil, err
		}
	}
	c.log("Experiment are fetched: %v", ex)
	return ex, err
}

func (c *Client) completeExperiment(e *types.Experiment) error {
	for _, id := range e.VariationsID {
		variation, err := c.fetchVariation(id)
		if err != nil {
			continue
		}
		e.Variations = append(e.Variations, variation)
	}
	if e.TargetingSegmentID > 0 {
		segment, err := c.fetchSegment(e.TargetingSegmentID)
		if err != nil {
			return err
		}
		if segment.ID == 0 {
			return newErrNotFound("segment id")
		}
		if segment.ConditionsData == nil {
			return newErrNotFound("segment condition data")
		}
		e.TargetingSegment = targeting.NewSegment(segment)
	}
	return nil
}

func (c *Client) fetchVariation(id int) (types.Variation, error) {
	v := types.Variation{}
	var path strings.Builder
	path.WriteString("/variations/")
	path.WriteString(utils.WriteUint(id))
	err := c.fetchOne(path.String(), fetchQuery{}, nil, respCallbackJson(&v))
	return v, err
}

func (c *Client) fetchSegment(id int) (*types.Segment, error) {
	s := &types.Segment{}

	var path strings.Builder
	path.WriteString("/segments/")
	path.WriteString(utils.WriteUint(id))
	err := c.fetchOne(path.String(), fetchQuery{}, nil, respCallbackJson(s))
	return s, err
}

func (c *Client) fetchFeatureFlagsGraphQL(siteCode string, perPage ...int) ([]types.FeatureFlag, error) {
	c.log("Fetching feature flags")
	pp := -1
	if len(perPage) > 0 {
		pp = perPage[0]
	}
	var ff []types.FeatureFlag
	cb := func(resp *fasthttp.Response, err error) error {
		if err != nil {
			return err
		}
		b := resp.Body()
		if len(b) == 0 {
			return ErrEmptyResponse
		}
		fmt.Println(string(b))
		var res FeatureFlagDataGraphQL
		err = json.Unmarshal(b, &res)
		if err != nil {
			fmt.Println(err)
			return err
		}
		for _, ffQL := range res.Data.FeatureFlags.Edge {
			ff = append(ff, ffQL.Transform())
		}
		return nil
	}

	err := c.fetchAllGraphQL(GetFeatureFlagsGraphQL(siteCode), fetchQuery{PerPage: pp}, cb)
	c.log("Feature flags are fetched: %v", ff)
	return ff, err
}

func (c *Client) fetchFeatureFlags(siteID int, perPage ...int) ([]types.FeatureFlag, error) {
	c.log("Fetching feature flags")
	pp := -1
	if len(perPage) > 0 {
		pp = perPage[0]
	}
	var ff []types.FeatureFlag
	filters := []fetchFilter{
		{
			Field:      "siteId",
			Operator:   "EQUAL",
			Parameters: []string{utils.WriteUint(siteID)},
		},
		{
			Field:      "status",
			Operator:   "EQUAL",
			Parameters: []string{"ACTIVE"},
		},
	}
	cb := func(resp *fasthttp.Response, err error) error {
		if err != nil {
			return err
		}
		b := resp.Body()
		if len(b) == 0 {
			return ErrEmptyResponse
		}
		fmt.Println(string(b))
		res := []types.FeatureFlag{{}}
		if b[0] == '[' {
			err = json.Unmarshal(b, &res)
		} else {
			err = json.Unmarshal(b, &res[0])
		}
		if err != nil {
			return err
		}
		ff = append(ff, res...)
		return nil
	}

	err := c.fetchAll("/feature-flags", fetchQuery{PerPage: pp}, filters, cb)
	for i := 0; i < len(ff); i++ {
		err = c.completeFeatureFlag(&ff[i])
		if err != nil {
			return nil, err
		}
	}
	c.log("Feature flags are fetched: %v", ff)
	return ff, err
}

func (c *Client) completeFeatureFlag(ff *types.FeatureFlag) error {
	for _, id := range ff.VariationsID {
		variation, err := c.fetchVariation(id)
		if err != nil {
			continue
		}
		ff.Variations = append(ff.Variations, variation)
	}
	if ff.TargetingSegmentID > 0 {
		segment, err := c.fetchSegment(ff.TargetingSegmentID)
		if err != nil {
			return err
		}
		if segment.ID == 0 {
			return newErrNotFound("segment id")
		}
		if segment.ConditionsData == nil {
			return newErrNotFound("segment condition Data")
		}
		ff.TargetingSegment = targeting.NewSegment(segment)
	}
	return nil
}

type fetchQuery struct {
	PerPage int `url:"perPage,omitempty"`
	Page    int `url:"page,omitempty"`
}

type fetchFilter struct {
	Field      string      `json:"field"`
	Operator   string      `json:"operator"`
	Parameters interface{} `json:"parameters"`
}

func (c *Client) fetchAllGraphQL(queryQL string, q fetchQuery, cb respCallback) error {
	currentPage := 1
	lastPage := -1
	iterator := func(resp *fasthttp.Response, err error) error {
		if resp.StatusCode() >= fasthttp.StatusBadRequest {
			return ErrBadStatus
		}
		var cbErr error
		if cb != nil {
			cbErr = cb(resp, err)
		}
		if lastPage < 0 {
			count := resp.Header.Peek(HeaderPaginationCount)
			lastPage, err = fasthttp.ParseUint(count)
			return err
		}
		return cbErr
	}
	for {
		q.Page = currentPage
		if lastPage >= 0 && currentPage > lastPage {
			break
		}
		err := c.fetchOneGraphQL(queryQL, q, iterator)
		if err != nil {
			break
		}
		currentPage++
	}
	return nil
}

func (c *Client) fetchAll(path string, q fetchQuery, filters []fetchFilter, cb respCallback) error {
	currentPage := 1
	lastPage := -1
	iterator := func(resp *fasthttp.Response, err error) error {
		if resp.StatusCode() >= fasthttp.StatusBadRequest {
			return ErrBadStatus
		}
		var cbErr error
		if cb != nil {
			cbErr = cb(resp, err)
		}
		if lastPage < 0 {
			count := resp.Header.Peek(HeaderPaginationCount)
			lastPage, err = fasthttp.ParseUint(count)
			return err
		}
		return cbErr
	}
	for {
		q.Page = currentPage
		if lastPage >= 0 && currentPage > lastPage {
			break
		}
		err := c.fetchOne(path, q, filters, iterator)
		if err != nil {
			break
		}
		currentPage++
	}
	return nil
}

func (c *Client) fetchOneGraphQL(queryQL string, q fetchQuery, cb respCallback) error {
	uri, err := buildFetchPathGraphQL(API_URL+"/v1/graphql", q)
	if err != nil {
		return err
	}
	req := request{
		Method:      MethodPost,
		URL:         uri,
		ContentType: HeaderContentTypeJson,
		BodyString:  queryQL,
	}
	c.m.Lock()
	req.AuthToken = c.token
	c.m.Unlock()
	if len(req.AuthToken) == 0 {
		return newErrCredentialsNotFound(req.String())
	}
	err = c.rest.Do(req, cb)
	if err != nil {
		c.log("Failed to fetch: %v, request: %v", err, req)
	}
	return err
}

func (c *Client) fetchOne(path string, q fetchQuery, filters []fetchFilter, cb respCallback) error {
	uri, err := buildFetchPath(API_URL, path, q, filters)
	if err != nil {
		return err
	}
	req := request{
		Method:      MethodGet,
		URL:         uri,
		ContentType: HeaderContentTypeJson,
	}
	c.m.Lock()
	req.AuthToken = c.token
	c.m.Unlock()
	if len(req.AuthToken) == 0 {
		return newErrCredentialsNotFound(req.String())
	}
	err = c.rest.Do(req, cb)
	if err != nil {
		c.log("Failed to fetch: %v, request: %v", err, req)
	}
	return err
}

func buildFetchPath(base, path string, q fetchQuery, filters []fetchFilter) (string, error) {
	var buf strings.Builder
	buf.WriteString(base)
	buf.WriteString(path)
	isFirst := true
	writeDelim := func() {
		if !isFirst {
			buf.WriteByte('&')
		} else {
			buf.WriteByte('?')
			isFirst = false
		}
	}
	if q.PerPage > 0 {
		writeDelim()
		buf.WriteString("perPage=")
		buf.WriteString(strconv.Itoa(q.PerPage))
	}
	if q.Page > 0 {
		writeDelim()
		buf.WriteString("page=")
		buf.WriteString(strconv.Itoa(q.Page))
	}
	if len(filters) > 0 {
		writeDelim()
		buf.WriteString("filter=")
		fbuf, err := json.Marshal(filters)
		if err != nil {
			return "", err
		}
		buf.WriteString(url.QueryEscape(string(fbuf)))
	}
	return buf.String(), nil
}

func buildFetchPathGraphQL(base string, q fetchQuery) (string, error) {
	var buf strings.Builder
	buf.WriteString(base)
	isFirst := true
	writeDelim := func() {
		if !isFirst {
			buf.WriteByte('&')
		} else {
			buf.WriteByte('?')
			isFirst = false
		}
	}
	if q.PerPage > 0 {
		writeDelim()
		buf.WriteString("perPage=")
		buf.WriteString(strconv.Itoa(q.PerPage))
	}
	if q.Page > 0 {
		writeDelim()
		buf.WriteString("page=")
		buf.WriteString(strconv.Itoa(q.Page))
	}
	return buf.String(), nil
}

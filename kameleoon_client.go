package kameleoon

import (
	"sync"
	"time"

	"github.com/Kameleoon/client-go/v3/logging"

	"github.com/Kameleoon/client-go/v3/configuration"
	"github.com/Kameleoon/client-go/v3/errs"
	"github.com/Kameleoon/client-go/v3/managers/data"
	"github.com/Kameleoon/client-go/v3/managers/hybrid"
	"github.com/Kameleoon/client-go/v3/managers/remotedata"
	"github.com/Kameleoon/client-go/v3/managers/tracking"
	"github.com/Kameleoon/client-go/v3/managers/warehouse"
	"github.com/Kameleoon/client-go/v3/network"
	"github.com/Kameleoon/client-go/v3/network/cookie"
	"github.com/Kameleoon/client-go/v3/realtime"
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
	IsUniqueIdentifier bool // Deprecated: Please use `UniqueIdentifier` data instead
	Timeout            time.Duration
}

type GetVariationOptParams struct {
	track bool
}

func NewGetVariationOptParams() GetVariationOptParams {
	return GetVariationOptParams{track: true}
}
func (p GetVariationOptParams) Track(value bool) GetVariationOptParams {
	p.track = value
	return p
}

type GetVariationsOptParams struct {
	onlyActive bool
	track      bool
}

func NewGetVariationsOptParams() GetVariationsOptParams {
	return GetVariationsOptParams{onlyActive: false, track: true}
}
func (p GetVariationsOptParams) OnlyActive(value bool) GetVariationsOptParams {
	p.onlyActive = value
	return p
}
func (p GetVariationsOptParams) Track(value bool) GetVariationsOptParams {
	p.track = value
	return p
}

type SetForcedVariationOptParams struct {
	forceTargeting bool
}

func NewSetForcedVariationOptParams() SetForcedVariationOptParams {
	return SetForcedVariationOptParams{forceTargeting: true}
}
func (p SetForcedVariationOptParams) ForceTargeting(value bool) SetForcedVariationOptParams {
	p.forceTargeting = value
	return p
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

	FlushVisitorInstantly(visitorCode string) error

	FlushAll(instant ...bool)

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
	//
	// Deprecated: Please use `GetVariation(visitorCode, featureKey, true)` instead
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
	//
	// Deprecated: Please use `GetVariation(visitorCode, featureKey, true)` instead
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
	// returns VisitorCodeNotValid error
	IsFeatureActive(visitorCode string, featureKey string, isUniqueIdentifier ...bool) (bool, error)

	IsFeatureActiveWithTracking(visitorCode string, featureKey string, track bool) (bool, error)

	GetVariation(visitorCode string, featureKey string, params ...GetVariationOptParams) (types.Variation, error)

	GetVariations(visitorCode string, params ...GetVariationsOptParams) (map[string]types.Variation, error)

	// GetFeatureVariationVariables retrieves all feature variable values for a given variation
	//
	// This method takes a featureKey and variationKey as mandatory arguments and
	// returns a list of variables for a given variation key
	// A feature variable can be changed easily via our web application.
	//
	// returns FeatureNotFound error
	// returns FeatureEnvironmentDisabled error
	// returns VariationNotFound error
	//
	// Deprecated: Please use `GetVariation(visitorCode, featureKey, false)` instead
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

	// Deprecated: Please use `GetRemoteVisitorDataWithFilter`
	GetRemoteVisitorDataWithOptParams(
		visitorCode string, addData bool, filter types.RemoteVisitorDataFilter, params ...RemoteVisitorDataOptParams,
	) ([]types.Data, error)

	GetRemoteVisitorDataWithFilter(
		visitorCode string, addData bool, filter types.RemoteVisitorDataFilter, params ...RemoteVisitorDataOptParams,
	) ([]types.Data, error)

	OnUpdateConfiguration(handler func())

	// GetFeatureList returns a list of all feature flag keys
	GetFeatureList() []string

	// GetActiveFeatureListForVisitor returns a list of active feature flag keys for a visitor
	//
	// returns VisitorCodeNotValid error when visitor code is not valid
	//
	// Deprecated: Please use `GetActiveFeatures`
	GetActiveFeatureListForVisitor(visitorCode string) ([]string, error)

	// GetActiveFeatures returns a map that contains the assigned variations of the active features
	// using the keys of the corresponding active features.
	//
	// returns VisitorCodeNotValid error when visitor code is not valid
	//
	// Deprecated: Please use `GetVariations(visitorCode, true, false)` instead
	GetActiveFeatures(visitorCode string) (map[string]types.Variation, error)

	GetEngineTrackingCode(visitorCode string) string

	// Sets or resets a forced variation for a visitor in a specific experiment,
	// so the experiment will be evaluated to the variation for the visitor.
	//
	// In order to reset the forced variation set the `variationKey` parameter to empty string.
	// If the forced variation you want to reset does not exist, the method will have no effect.
	//
	// Parameters:
	// - visitorCode: The unique visitor code identifying the visitor.
	// - experimentId: The identifier of the experiment you want to set/reset the forced variation for.
	// - variationKey: The identifier of the variation you want the experiment to be evaluated to.
	//   Set to empty string to reset the forced variation.
	// - ForceTargeting: If `true`, the visitor will be targeted to the experiment regardless its conditions.
	//   Otherwise, the normal targeting logic will be preserved. Optional (defaults to `true`).
	//
	// May return one of the following errors:
	// - VisitorCodeInvalid:
	//   The provided visitor code is invalid.
	// - FeatureExperimentNotFound:
	//   The experiment id does not exist in the feature flag.
	// - VariationNotFound:
	//   The variation key does not belong to the experiment.
	SetForcedVariation(
		visitorCode string, experimentId int, variationKey string, params ...SetForcedVariationOptParams,
	) error
}

type kameleoonClient struct {
	cfg                  *KameleoonClientConfig
	visitorManager       storage.VisitorManager
	networkManager       network.NetworkManager
	cookieManager        cookie.CookieManager
	hybridManager        hybrid.HybridManager
	warehouseManager     warehouse.WarehouseManager
	targetingManager     targeting.TargetingManager
	remoteDataManager    remotedata.RemoteDataManager
	trackingManager      tracking.TrackingManager
	configurationManager configuration.ConfigurationManager

	m           sync.Mutex
	readiness   *kameleoonClientReadiness
	dataManager data.DataManager
	closed      bool
}

func newClient(siteCode string, cfg *KameleoonClientConfig) (*kameleoonClient, error) {
	logging.Info("CALL: newClient(siteCode: %s, config: %s)", siteCode, cfg)
	if len(siteCode) == 0 {
		err := errs.NewSiteCodeIsEmpty("Provided siteCode is empty")
		logging.Info("RETURN: newClient(siteCode: %s, config: %s) -> (client, error: %s)",
			siteCode, cfg, err)
		return nil, err
	}
	if err := cfg.defaults(); err != nil {
		logging.Info("RETURN: newClient(siteCode: %s, config: %s) -> (client, error: %s)",
			siteCode, cfg, err)
		return nil, err
	}

	if cfg.Logger != nil {
		logging.SetOldLogger(cfg.Logger)
	}

	if cfg.VerboseMode && logging.GetLogLevel() == logging.WARNING {
		logging.SetLogLevel(logging.INFO)
	}

	df := configuration.NewDataFile(configuration.Configuration{}, cfg.Environment)
	dm := data.NewDataManagerImpl(df)
	np := network.NewNetProviderImpl(cfg.Network.ReadTimeout, cfg.Network.WriteTimeout,
		cfg.Network.MaxConnsPerHost, cfg.Network.ProxyURL)
	up := network.NewUrlProviderImpl(siteCode, network.DefaultDataApiDomain, utils.SdkName, utils.SdkVersion)
	atsf := &network.AccessTokenSourceFactoryImpl{ClientId: cfg.ClientID, ClientSecret: cfg.ClientSecret}
	nm := network.NewNetworkManagerImpl(cfg.Environment, cfg.DefaultTimeout, np, up, atsf)
	vm := newVisitorManager(dm, cfg)
	hm, _ := hybrid.NewHybridManagerImpl(5*time.Second, dm)
	tarM := targeting.NewTargetingManager(dm, vm)
	rdm := remotedata.NewRemoteDataManager(dm, nm, vm)
	trM := tracking.NewTrackingManagerImpl(dm, nm, vm, cfg.TrackingInterval)
	cm := configuration.NewConfigurationManager(dm, nm, &realtime.NetSseClient{}, cfg.RefreshInterval, cfg.Environment)
	client := newClientInternal(cfg, dm, nm, vm, hm, tarM, rdm, trM, cm)
	logging.Info("RETURN: newClient(siteCode: %s, config: %s) -> (client, error: <nil>)",
		siteCode, cfg)
	return client, nil
}

func newClientInternal(
	cfg *KameleoonClientConfig,
	dataManager data.DataManager,
	networkManager network.NetworkManager,
	visitorManager storage.VisitorManager,
	hybridManager hybrid.HybridManager,
	targetingManager targeting.TargetingManager,
	remoteDataManager remotedata.RemoteDataManager,
	trackingManager tracking.TrackingManager,
	configurationManager configuration.ConfigurationManager,
) *kameleoonClient {
	client := &kameleoonClient{
		cfg:                  cfg,
		readiness:            newKameleoonClientReadiness(),
		dataManager:          dataManager,
		visitorManager:       visitorManager,
		hybridManager:        hybridManager,
		networkManager:       networkManager,
		cookieManager:        cookie.NewCookieManagerImpl(dataManager, visitorManager, cfg.TopLevelDomain),
		warehouseManager:     warehouse.NewWarehouseManagerImpl(networkManager, visitorManager),
		targetingManager:     targetingManager,
		remoteDataManager:    remoteDataManager,
		trackingManager:      trackingManager,
		configurationManager: configurationManager,
	}
	go client.updateConfigInitially()
	return client
}

func newVisitorManager(dm data.DataManager, cfg *KameleoonClientConfig) storage.VisitorManager {
	return storage.NewVisitorManagerImpl(dm, cfg.SessionDuration)
}

func (c *kameleoonClient) WaitInit() error {
	logging.Info("CALL: kameleoonClient.WaitInit()")
	err := c.readiness.Wait()
	logging.Info("RETURN: kameleoonClient.WaitInit() -> (error: %s)", err)
	return err
}

func (c *kameleoonClient) close() {
	logging.Debug("CALL: kameleoonClient.close()")
	if !c.closed {
		c.m.Lock()
		if c.closed {
			c.m.Unlock()
		} else {
			c.closed = true
			c.m.Unlock()
			c.visitorManager.Close()
			c.trackingManager.Close()
		}
	}
	logging.Debug("RETURN: kameleoonClient.close()")
}

func (c *kameleoonClient) GetVisitorCode(request *fasthttp.Request, response *fasthttp.Response,
	defaultVisitorCode ...string) (string, error) {
	logging.Info("CALL: kameleoonClient.GetVisitorCode(request, response, defaultVisitorCode: %s)",
		defaultVisitorCode)
	visitorCode, err := c.cookieManager.GetOrAdd(request, response, defaultVisitorCode...)
	logging.Info(
		"RETURN: kameleoonClient.GetVisitorCode(request, response, defaultVisitorCode: %s) -> "+
			"(visitorCode: %s, error: %s)", defaultVisitorCode, visitorCode, err)
	return visitorCode, err
}

func (c *kameleoonClient) SetLegalConsent(visitorCode string, consent bool, response ...*fasthttp.Response) error {
	logging.Info("CALL: kameleoonClient.SetLegalConsent(visitorCode: %s, consent: %s, response)",
		visitorCode, consent)
	err := utils.ValidateVisitorCode(visitorCode)
	if err == nil {
		v := c.visitorManager.GetOrCreateVisitor(visitorCode)
		v.SetLegalConsent(consent)
		if len(response) > 0 {
			c.cookieManager.Update(visitorCode, consent, response[0])
		}
	}
	logging.Info("RETURN: kameleoonClient.SetLegalConsent(visitorCode: %s, consent: %s, response) -> (error: %s)",
		visitorCode, consent, err)
	return err
}

func (c *kameleoonClient) AddData(visitorCode string, allData ...types.Data) error {
	//var stats runtime.MemStats
	//runtime.ReadMemStats(&stats))
	logging.Info("CALL: kameleoonClient.AddData(visitorCode: %s, allData: %s)", visitorCode, allData)
	err := utils.ValidateVisitorCode(visitorCode)
	if err == nil {
		c.visitorManager.AddData(visitorCode, allData...)
	}
	logging.Info("RETURN: kameleoonClient.AddData(visitorCode: %s, allData: %s) -> (error: %s)",
		visitorCode, allData, err)
	return err
}

func (c *kameleoonClient) TrackConversion(visitorCode string, goalID int, isUniqueIdentifier ...bool) error {
	logging.Info(
		"CALL: kameleoonClient.TrackConversion(visitorCode: %s, goalID: %s, isUniqueIdentifier: %s)",
		visitorCode, goalID, isUniqueIdentifier)
	if len(isUniqueIdentifier) > 0 {
		c.setUniqueIdentifier(visitorCode, isUniqueIdentifier[0])
	}
	err := c.trackConversion(visitorCode, goalID)
	logging.Info(
		"RETURN: kameleoonClient.TrackConversion(visitorCode: %s, goalID: %s, isUniqueIdentifier: %s) -> (error: %s)",
		visitorCode, goalID, isUniqueIdentifier, err)
	return err
}

func (c *kameleoonClient) TrackConversionRevenue(
	visitorCode string, goalID int, revenue float64, isUniqueIdentifier ...bool,
) error {
	logging.Info(
		"CALL: kameleoonClient.TrackConversionRevenue(visitorCode: %s, goalID: %s, revenue: %s,"+
			" isUniqueIdentifier: %s)", visitorCode, goalID, revenue, isUniqueIdentifier)
	if len(isUniqueIdentifier) > 0 {
		c.setUniqueIdentifier(visitorCode, isUniqueIdentifier[0])
	}
	err := c.trackConversion(visitorCode, goalID, revenue)
	logging.Info(
		"RETURN: kameleoonClient.TrackConversionRevenue(visitorCode: %s, goalID: %s, revenue: %s,"+
			" isUniqueIdentifier: %s) -> (error: %s)", visitorCode, goalID, revenue, isUniqueIdentifier, err)
	return err
}

func (c *kameleoonClient) trackConversion(
	visitorCode string, goalID int, revenue ...float64,
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
	c.trackingManager.AddVisitorCode(visitorCode)
	return nil
}

func (c *kameleoonClient) FlushVisitor(visitorCode string, isUniqueIdentifier ...bool) error {
	logging.Info("CALL: kameleoonClient.FlushVisitor(visitorCode: %s, isUniqueIdentifier: %s)",
		visitorCode, isUniqueIdentifier)
	err := utils.ValidateVisitorCode(visitorCode)
	if err == nil {
		if len(isUniqueIdentifier) > 0 {
			c.setUniqueIdentifier(visitorCode, isUniqueIdentifier[0])
		}
		c.trackingManager.AddVisitorCode(visitorCode)
	}
	logging.Info("RETURN: kameleoonClient.FlushVisitor(visitorCode: %s, isUniqueIdentifier: %s) -> (error: %s)",
		visitorCode, isUniqueIdentifier, err)
	return err
}

func (c *kameleoonClient) FlushVisitorInstantly(visitorCode string) error {
	logging.Info("CALL: kameleoonClient.FlushVisitorInstantly(visitorCode: %s)", visitorCode)
	err := utils.ValidateVisitorCode(visitorCode)
	if err == nil {
		c.trackingManager.TrackVisitor(visitorCode)
	}
	logging.Info("RETURN: kameleoonClient.FlushVisitorInstantly(visitorCode: %s) -> (error: %s)", visitorCode, err)
	return err
}

func (c *kameleoonClient) FlushAll(instant ...bool) {
	logging.Info("CALL: kameleoonClient.FlushAll(instant: %s)", instant)
	c.visitorManager.Enumerate(func(vc string, v storage.Visitor) bool {
		notEmpty := false
		v.EnumerateSendableData(func(s types.Sendable) bool {
			notEmpty = true
			return false
		})
		if notEmpty {
			c.trackingManager.AddVisitorCode(vc)
		}
		return true
	})
	if (len(instant) > 0) && instant[0] {
		c.trackingManager.TrackAll()
	}
	logging.Info("RETURN: kameleoonClient.FlushAll(instant: %s)", instant)
}

func (c *kameleoonClient) GetFeatureVariationKey(
	visitorCode string, featureKey string, isUniqueIdentifier ...bool,
) (string, error) {
	logging.Info(
		"CALL: kameleoonClient.GetFeatureVariationKey(visitorCode: %s, featureKey: %s, isUniqueIdentifier: %s)",
		visitorCode, featureKey, isUniqueIdentifier)
	if len(isUniqueIdentifier) > 0 {
		c.setUniqueIdentifier(visitorCode, isUniqueIdentifier[0])
	}
	_, variationKey, err := c.getFeatureVariationKey(visitorCode, featureKey)
	logging.Info(
		"RETURN: kameleoonClient.GetFeatureVariationKey(visitorCode: %s, featureKey: %s, isUniqueIdentifier: %s) "+
			"-> (variationKey: %s, error: %s)", visitorCode, featureKey, isUniqueIdentifier, variationKey, err)
	return variationKey, err
}

// getFeatureVariationKey is a helper method for getting variation key for feature flag
func (c *kameleoonClient) getFeatureVariationKey(
	visitorCode string, featureKey string,
) (featureFlag types.FeatureFlag, variationKey string, err error) {
	logging.Debug(
		"CALL: kameleoonClient.getFeatureVariationKey(visitorCode: %s, featureKey: %s)", visitorCode, featureKey,
	)
	defer func() {
		logging.Debug(
			"RETURN: kameleoonClient.getFeatureVariationKey(visitorCode: %s, featureKey: %s) -> (featureFlag: %s, "+
				"variationKey: %s, error: %s)", visitorCode, featureKey, featureFlag, variationKey, err,
		)
	}()
	if err = utils.ValidateVisitorCode(visitorCode); err != nil {
		variationKey = string(types.VariationOff)
		return
	}
	// find enabled feature flag else return an error
	if featureFlag, err = c.dataManager.DataFile().GetFeatureFlag(featureKey); err != nil {
		variationKey = string(types.VariationOff)
		return
	}
	visitor := c.visitorManager.GetVisitor(visitorCode)
	evalExp := c.evaluate(visitor, visitorCode, featureFlag, true, true)
	// get variation key from feature flag
	defaultVariationKey := featureFlag.GetDefaultVariationKey()
	variationKey = c.calculateVariationKey(evalExp, defaultVariationKey)
	c.trackingManager.AddVisitorCode(visitorCode)
	return
}

func (c *kameleoonClient) saveVariation(
	visitorCode string, evalExp *evaluatedExperiment, track bool,
) {
	if (evalExp == nil) || (evalExp.experiment.ExperimentId == 0) || (evalExp.varByExp.VariationID == nil) {
		return
	}
	logging.Debug(
		"CALL: kameleoonClient.saveVariation(visitorCode: %s, evalExp: %s, track: %s)",
		visitorCode, evalExp, track,
	)
	visitor := c.visitorManager.GetOrCreateVisitor(visitorCode)
	asVariation := types.NewAssignedVariation(
		evalExp.experiment.ExperimentId, *evalExp.varByExp.VariationID, evalExp.ruleType,
	)
	if !track {
		asVariation.MarkAsSent()
	}
	visitor.AssignVariation(asVariation)
	logging.Debug(
		"RETURN: kameleoonClient.saveVariation(visitorCode: %s, evalExp: %s, track: %s)",
		visitorCode, evalExp, track,
	)
}

func (c *kameleoonClient) calculateVariationKey(evalExp *evaluatedExperiment, defaultVariationKey string) string {
	logging.Debug(
		"CALL: kameleoonClient.calculateVariationKey(evalExp: %s, defaultVariationKey: %s)",
		evalExp, defaultVariationKey,
	)
	var variationKey string
	if evalExp != nil {
		variationKey = evalExp.varByExp.VariationKey
	} else {
		variationKey = defaultVariationKey
	}
	logging.Debug(
		"RETURN: kameleoonClient.calculateVariationKey(evalExp: %s, defaultVariationKey: %s) -> "+
			"(variationKey: %s)", evalExp, defaultVariationKey, variationKey,
	)
	return variationKey
}

func getCodeForHash(visitor storage.Visitor, visitorCode string) string {
	// use mappingIdentifier instead of visitorCode if it was set up
	if (visitor != nil) && (visitor.MappingIdentifier() != nil) {
		return *visitor.MappingIdentifier()
	}
	return visitorCode
}

// getVariationRuleForFeature is a helper method for calculate variation key for feature flag
func (c *kameleoonClient) calculateVariationRuleForFeature(
	visitorCode string, featureFlag types.FeatureFlag,
) (selected *evaluatedExperiment) {
	logging.Debug(
		"CALL: kameleoonClient.calculateVariationRuleForFeature(visitorCode: %s, featureFlag: %s)",
		visitorCode, featureFlag,
	)
	defer func() {
		logging.Debug(
			"RETURN: kameleoonClient.calculateVariationRuleForFeature(visitorCode: %s, featureFlag: %s)"+
				" -> (evalExp: %s)",
			visitorCode, featureFlag, selected,
		)
	}()
	visitor := c.visitorManager.GetVisitor(visitorCode)
	codeForHash := getCodeForHash(visitor, visitorCode)
	// no rules -> return DefaultVariationKey
	for _, rule := range featureFlag.GetRules() {
		var forcedVariation *types.ForcedExperimentVariation
		if visitor != nil {
			forcedVariation = visitor.GetForcedExperimentVariation(rule.GetRuleBase().ExperimentId)
			if (forcedVariation != nil) && forcedVariation.ForceTargeting() {
				// Forcing experiment variation in force-targeting mode
				return newEvaluatedExperimentFromVarByExpRule(forcedVariation.VarByExp(), rule)
			}
		}

		//check if visitor is targeted for rule, else next rule
		if c.targetingManager.CheckTargeting(visitorCode, rule.GetRuleBase().ExperimentId, rule) {
			if forcedVariation != nil {
				// Forcing experiment variation in targeting-only mode
				return newEvaluatedExperimentFromVarByExpRule(forcedVariation.VarByExp(), rule)
			}

			// Disable searching in variation storage (uncommented if you need use variation storage)
			// check for saved variation for rule if it's experimentation rule
			// if savedVariation, found := c.getSavedVariationForRule(visitorCode, &rule); found {
			// 	return savedVariation, &rule, false
			// }

			//uses for rule exposition
			hashRule := utils.GetHashDoubleRule(codeForHash, rule.GetRuleBase().Id, rule.GetRuleBase().RespoolTime)
			logging.Debug("Calculated rule hash %s for code %s", hashRule, codeForHash)
			//check main expostion for rule with hashRule
			if hashRule <= rule.GetRuleBase().Exposition {
				if rule.IsTargetDeliveryType() {
					var variation *types.VariationByExposition
					if len(rule.GetRuleBase().VariationsByExposition) > 0 {
						variation = &rule.GetRuleBase().VariationsByExposition[0]
					}
					return newEvaluatedExperimentFromVarByExpRule(variation, rule)
				}
				//uses for variation's expositions
				hashVariation := utils.GetHashDoubleRule(
					codeForHash, rule.GetRuleBase().ExperimentId, rule.GetRuleBase().RespoolTime,
				)
				logging.Debug("Calculated variation hash %s for code %s", hashVariation, codeForHash)
				// get variation with new hashVariation
				variation := rule.GetVariationByHash(hashVariation)
				if variation != nil {
					return newEvaluatedExperimentFromVarByExpRule(variation, rule)
				}
			}
			if rule.IsTargetDeliveryType() {
				break
			}
		}
	}
	return nil
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
	logging.Info(
		"CALL: kameleoonClient.GetFeatureVariable(visitorCode: %s, featureKey: %s, variableKey: %s,"+
			" isUniqueIdentifier: %s)", visitorCode, featureKey, variableKey, isUniqueIdentifier)
	if len(isUniqueIdentifier) > 0 {
		c.setUniqueIdentifier(visitorCode, isUniqueIdentifier[0])
	}
	var variableValue interface{}
	featureFlag, variationKey, err := c.getFeatureVariationKey(visitorCode, featureKey)
	if err == nil {
		variation, exist := featureFlag.GetVariationByKey(variationKey)
		if !exist {
			err = errs.NewFeatureVariationNotFound(featureKey, variationKey)
		} else {
			variable, exist := variation.GetVariableByKey(variableKey)
			if !exist {
				err = errs.NewFeatureVariableNotFound(featureKey, variationKey, variableKey)
			} else {
				variableValue = parseFeatureVariable(variable)
			}
		}
	}

	logging.Info(
		"RETURN: kameleoonClient.GetFeatureVariable(visitorCode: %s, featureKey: %s, variableKey: %s, "+
			"isUniqueIdentifier: %s) -> (variable: %s, err: %s)",
		visitorCode, featureKey, variableKey, isUniqueIdentifier, variableValue, err)
	return variableValue, err
}

func (c *kameleoonClient) IsFeatureActive(
	visitorCode string, featureKey string, isUniqueIdentifier ...bool,
) (isFeatureActive bool, err error) {
	logging.Info(
		"CALL: kameleoonClient.IsFeatureActive(visitorCode: %s, featureKey: %s, isUniqueIdentifier: %s)",
		visitorCode, featureKey, isUniqueIdentifier,
	)
	defer func() {
		logging.Info(
			"RETURN: kameleoonClient.IsFeatureActive(visitorCode: %s, featureKey: %s, isUniqueIdentifier: %s) -> "+
				"(isFeatureActive: %s, err: %s)", visitorCode, featureKey, isUniqueIdentifier, isFeatureActive, err,
		)
	}()
	if err = utils.ValidateVisitorCode(visitorCode); err != nil {
		return
	}
	if len(isUniqueIdentifier) > 0 {
		c.setUniqueIdentifier(visitorCode, isUniqueIdentifier[0])
	}
	isFeatureActive, err = c.isFeatureActive(visitorCode, featureKey, true)
	return
}

func (c *kameleoonClient) IsFeatureActiveWithTracking(
	visitorCode string, featureKey string, track bool,
) (isFeatureActive bool, err error) {
	logging.Info(
		"CALL: kameleoonClient.IsFeatureActiveWithTracking(visitorCode: %s, featureKey: %s, track: %s)",
		visitorCode, featureKey, track,
	)
	defer func() {
		logging.Info(
			"RETURN: kameleoonClient.IsFeatureActiveWithTracking(visitorCode: %s, featureKey: %s, track: %s) -> "+
				"(isFeatureActive: %s, err: %s)", visitorCode, featureKey, track, isFeatureActive, err,
		)
	}()
	if err = utils.ValidateVisitorCode(visitorCode); err != nil {
		return
	}
	isFeatureActive, err = c.isFeatureActive(visitorCode, featureKey, track)
	return
}

func (c *kameleoonClient) isFeatureActive(
	visitorCode string, featureKey string, track bool,
) (isFeatureActive bool, err error) {
	logging.Debug(
		"CALL: kameleoonClient.isFeatureActive(visitorCode: %s, featureKey: %s, track: %s)",
		visitorCode, featureKey, track,
	)
	defer func() {
		logging.Debug(
			"RETURN: kameleoonClient.isFeatureActive(visitorCode: %s, featureKey: %s, track: %s) -> "+
				"(isFeatureActive: %s, err: %s)", visitorCode, featureKey, track, isFeatureActive, err,
		)
	}()
	var featureFlag types.FeatureFlag
	if featureFlag, err = c.dataManager.DataFile().GetFeatureFlag(featureKey); err != nil {
		switch err.(type) {
		case *errs.FeatureEnvironmentDisabled:
			return false, nil
		default:
			return false, err
		}
	}
	variationKey, _ := c.getVariationInfo(visitorCode, featureFlag, track)
	isFeatureActive = variationKey != string(types.VariationOff)
	if track {
		c.trackingManager.AddVisitorCode(visitorCode)
	}
	return
}

func (c *kameleoonClient) GetVariation(
	visitorCode string, featureKey string, params ...GetVariationOptParams,
) (externalVariation types.Variation, err error) {
	logging.Info(
		"CALL: kameleoonClient.GetVariation(visitorCode: %s, featureKey: %s, params: %s)",
		visitorCode, featureKey, params,
	)
	defer func() {
		logging.Info(
			"RETURN: kameleoonClient.GetVariation(visitorCode: %s, featureKey: %s, params: %s) -> "+
				"(variation: %s, err: %s)", visitorCode, featureKey, params, externalVariation, err,
		)
	}()
	var p GetVariationOptParams
	if len(params) > 0 {
		p = params[0]
	} else {
		p = NewGetVariationOptParams()
	}
	if err = utils.ValidateVisitorCode(visitorCode); err != nil {
		return
	}
	var featureFlag types.FeatureFlag
	if featureFlag, err = c.dataManager.DataFile().GetFeatureFlag(featureKey); err != nil {
		return
	}
	variationKey, evalExp := c.getVariationInfo(visitorCode, featureFlag, p.track)
	variation, _ := featureFlag.GetVariationByKey(variationKey)
	externalVariation = createExternalVariation(variation, evalExp)
	if p.track {
		c.trackingManager.AddVisitorCode(visitorCode)
	}
	return
}

func (c *kameleoonClient) GetVariations(
	visitorCode string, params ...GetVariationsOptParams,
) (variations map[string]types.Variation, err error) {
	logging.Info(
		"CALL: kameleoonClient.GetVariations(visitorCode: %s, params: %s)",
		visitorCode, params,
	)
	defer func() {
		logging.Info(
			"RETURN: kameleoonClient.GetVariations(visitorCode: %s, params: %s) -> (variations: %s, err: %s)",
			visitorCode, params, variations, err,
		)
	}()
	var p GetVariationsOptParams
	if len(params) > 0 {
		p = params[0]
	} else {
		p = NewGetVariationsOptParams()
	}
	if err = utils.ValidateVisitorCode(visitorCode); err != nil {
		return
	}
	variations = make(map[string]types.Variation)
	for _, ff := range c.dataManager.DataFile().GetOrderedFeatureFlags() {
		if !ff.GetEnvironmentEnabled() {
			continue
		}
		variationKey, evalExp := c.getVariationInfo(visitorCode, ff, p.track)
		if p.onlyActive && (variationKey == string(types.VariationOff)) {
			continue
		}
		variation, _ := ff.GetVariationByKey(variationKey)
		variations[ff.GetFeatureKey()] = createExternalVariation(variation, evalExp)
	}
	if p.track {
		c.trackingManager.AddVisitorCode(visitorCode)
	}
	return
}

func (c *kameleoonClient) getVariationInfo(
	visitorCode string, featureFlag types.FeatureFlag, track bool,
) (variationKey string, evalExp *evaluatedExperiment) {
	logging.Debug(
		"CALL: kameleoonClient.getVariationInfo(visitorCode: %s, featureFlag: %s, track: %s)",
		visitorCode, featureFlag, track,
	)
	visitor := c.visitorManager.GetVisitor(visitorCode)
	evalExp = c.evaluate(visitor, visitorCode, featureFlag, track, true)
	defaultVariationKey := featureFlag.GetDefaultVariationKey()
	variationKey = c.calculateVariationKey(evalExp, defaultVariationKey)
	logging.Debug(
		"RETURN: kameleoonClient.getVariationInfo(visitorCode: %s, featureFlag: %s, track: %s)"+
			" -> (variationKey: %s, evalExp: %s)", visitorCode, featureFlag, track, variationKey, evalExp,
	)
	return
}

func (c *kameleoonClient) evaluate(
	visitor storage.Visitor, visitorCode string, featureFlag types.FeatureFlag, track, save bool,
) (evalExp *evaluatedExperiment) {
	logging.Debug(
		"CALL: kameleoonClient.evaluate(visitor, visitorCode: %s, featureFlag: %s, track: %s, save: %s)",
		visitorCode, featureFlag, track, save,
	)
	var forcedVariation *types.ForcedFeatureVariation
	if visitor != nil {
		forcedVariation = visitor.GetForcedFeatureVariation(featureFlag.GetFeatureKey())
	}
	if forcedVariation != nil {
		evalExp = newEvaluatedExperimentFromForcedVariation(forcedVariation)
	} else if c.isVisitorNotInHoldout(visitor, visitorCode, track, save) {
		evalExp = c.calculateVariationRuleForFeature(visitorCode, featureFlag)
	}
	if save && ((forcedVariation == nil) || !forcedVariation.Simulated()) {
		c.saveVariation(visitorCode, evalExp, track)
	}
	logging.Debug(
		"RETURN: kameleoonClient.evaluate(visitor, visitorCode: %s, featureFlag: %s, track: %s, save: %s)"+
			" -> (evalExp: %s)", visitorCode, featureFlag, track, save, evalExp,
	)
	return
}

func (c *kameleoonClient) isVisitorNotInHoldout(visitor storage.Visitor, visitorCode string, track, save bool) bool {
	const inHoldoutVariationKey = "in-holdout"
	logging.Debug(
		"CALL: kameleoonClient.isVisitorNotInHoldout(visitor, visitorCode: %s, track: %s, save: %s)",
		visitorCode, track, save,
	)
	isNotInHoldout := true
	if holdout := c.dataManager.DataFile().Holdout(); holdout != nil {
		codeForHash := getCodeForHash(visitor, visitorCode)
		variationHash := utils.GetHashDouble(codeForHash, holdout.ExperimentId)
		logging.Debug("Calculated holdout hash %s for code %s", variationHash, codeForHash)
		if varByExp := holdout.GetVariationByHash(variationHash); varByExp != nil {
			isNotInHoldout = varByExp.VariationKey != inHoldoutVariationKey
			if save {
				evalExp := &evaluatedExperiment{
					varByExp:   varByExp,
					experiment: holdout,
					ruleType:   types.RuleTypeExperimentation,
				}
				c.saveVariation(visitorCode, evalExp, track)
			}
		}
	}
	logging.Debug(
		"RETURN: kameleoonClient.isVisitorNotInHoldout(visitor, visitorCode: %s, track: %s, save: %s)"+
			" -> (isNotInHoldout: %s)", visitorCode, track, save, isNotInHoldout,
	)
	return isNotInHoldout
}

func createExternalVariation(
	internalVariation *types.VariationFeatureFlag, evalExp *evaluatedExperiment,
) (variation types.Variation) {
	logging.Debug("CALL: createExternalVariation(internalVariation: %s, evalExp: %s)", internalVariation, evalExp)
	defer func() {
		logging.Debug(
			"RETURN: createExternalVariation(internalVariation: %s, evalExp: %s) -> (variation: %s)",
			internalVariation, evalExp, variation,
		)
	}()
	variables := make(map[string]types.Variable)
	if internalVariation != nil {
		for _, internalVariable := range internalVariation.Variables {
			variables[internalVariable.Key] = types.Variable{
				Key:   internalVariable.Key,
				Type:  internalVariable.Type,
				Value: parseFeatureVariable(&internalVariable),
			}
		}
	}
	var variationKey string
	if internalVariation != nil {
		variationKey = internalVariation.Key
	}
	var variationId *int
	var experimentId *int
	if evalExp != nil {
		variationId = utils.Reref(evalExp.varByExp.VariationID)
		experimentId = utils.Reref(&evalExp.experiment.ExperimentId)
	}
	variation = types.Variation{
		Key:          variationKey,
		VariationID:  variationId,
		ExperimentID: experimentId,
		Variables:    variables,
	}
	return
}

func (c *kameleoonClient) GetFeatureVariationVariables(
	featureKey string, variationKey string,
) (map[string]interface{}, error) {
	logging.Info(
		"CALL: kameleoonClient.GetFeatureVariationVariables(featureKey: %s, variationKey: %s)",
		featureKey, variationKey)
	var mapVariableValues map[string]interface{}
	featureFlag, err := c.dataManager.DataFile().GetFeatureFlag(featureKey)
	if err == nil {
		variation, exist := featureFlag.GetVariationByKey(variationKey)
		if !exist {
			err = errs.NewFeatureVariationNotFound(featureKey, variationKey)
		} else {
			mapVariableValues = make(map[string]interface{})
			for _, variable := range variation.Variables {
				mapVariableValues[variable.Key] = parseFeatureVariable(&variable)
			}
		}
	}
	logging.Info(
		"RETURN: kameleoonClient.GetFeatureVariationVariables(featureKey: %s, variationKey: %s) -> "+
			"(variables: %s, error: %s)", featureKey, variationKey, mapVariableValues, err)
	return mapVariableValues, err
}

func parseFeatureVariable(variable *types.Variable) interface{} {
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
	logging.Info("CALL: kameleoonClient.GetRemoteData(key: %s, timeout: %s)", key, timeout)
	remoteData, err := c.remoteDataManager.GetData(key, timeout...)
	logging.Info("RETURN: kameleoonClient.GetRemoteData(key: %s, timeout: %s) -> (remoteData: %s, error: %s)",
		remoteData, err)
	return remoteData, err
}

func (c *kameleoonClient) GetRemoteVisitorData(
	visitorCode string,
	addData bool,
	timeout ...time.Duration,
) ([]types.Data, error) {
	logging.Info("CALL: kameleoonClient.GetRemoteVisitorData(visitorCode: %s, addData: %s, timeout: %s)",
		visitorCode, addData, timeout)
	filter := types.DefaultRemoteVisitorDataFilter()
	visitorData, err := c.remoteDataManager.GetVisitorData(visitorCode, filter, addData, timeout...)
	logging.Info(
		"RETURN: kameleoonClient.GetRemoteVisitorData(visitorCode: %s, addData: %s, timeout: %s) -> "+
			"(visitorData: %s, error: %s)", visitorCode, addData, timeout, visitorData, err)
	return visitorData, err
}

func (c *kameleoonClient) GetRemoteVisitorDataWithOptParams(
	visitorCode string, addData bool, filter types.RemoteVisitorDataFilter, params ...RemoteVisitorDataOptParams,
) ([]types.Data, error) {
	logging.Info(
		"CALL: kameleoonClient.GetRemoteVisitorDataWithOptParams(visitorCode: %s, addData: %s, filter: %s, params: %s)",
		visitorCode, addData, filter, params)
	var p RemoteVisitorDataOptParams
	if len(params) > 0 {
		p = params[0]
	}
	var timeout []time.Duration
	if p.Timeout > 0 {
		timeout = []time.Duration{p.Timeout}
	}
	c.setUniqueIdentifier(visitorCode, p.IsUniqueIdentifier)
	visitorData, err := c.remoteDataManager.GetVisitorData(visitorCode, filter, addData, timeout...)
	logging.Info(
		"RETURN: kameleoonClient.GetRemoteVisitorDataWithOptParams(visitorCode: %s, addData: %s, filter: %s, "+
			"params: %s) -> (visitorData: %s, error: %s)", visitorCode, addData, filter, params, visitorData, err)
	return visitorData, err
}

func (c *kameleoonClient) GetRemoteVisitorDataWithFilter(
	visitorCode string, addData bool, filter types.RemoteVisitorDataFilter, params ...RemoteVisitorDataOptParams,
) ([]types.Data, error) {
	logging.Info(
		"CALL: kameleoonClient.GetRemoteVisitorDataWithFilter(visitorCode: %s, addData: %s, filter: %s, params: %s)",
		visitorCode, addData, filter, params)
	var p RemoteVisitorDataOptParams
	if len(params) > 0 {
		p = params[0]
	}
	var timeout []time.Duration
	if p.Timeout > 0 {
		timeout = []time.Duration{p.Timeout}
	}
	remoteVisitorData, err := c.remoteDataManager.GetVisitorData(visitorCode, filter, addData, timeout...)
	logging.Info(
		"RETURN: kameleoonClient.GetRemoteVisitorDataWithFilter(visitorCode: %s, addData: %s, filter: %s,"+
			" params: %s) -> (remoteVisitorData: %s, error: %s)",
		visitorCode, addData, filter, params, remoteVisitorData, err)
	return remoteVisitorData, err
}

func (c *kameleoonClient) updateConfigInitially() {
	logging.Debug("CALL: kameleoonClient.updateConfigInitially()")
	err := c.configurationManager.Start()
	c.readiness.set(err)
	logging.Debug("RETURN: kameleoonClient.updateConfigInitially()")
}

func (c *kameleoonClient) OnUpdateConfiguration(handler func()) {
	c.configurationManager.OnUpdateConfiguration(handler)
	logging.Info("CALL/RETURN: kameleoonClient.OnUpdateConfiguration(handler)")
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
	logging.Info("CALL: kameleoonClient.GetFeatureList()")
	featureFlags := c.dataManager.DataFile().GetFeatureFlags()
	arrayKeys := make([]string, 0, len(featureFlags))
	for _, ff := range featureFlags {
		arrayKeys = append(arrayKeys, ff.GetFeatureKey())
	}
	logging.Info("RETURN: kameleoonClient.GetFeatureList() -> (features: %s)", arrayKeys)
	return arrayKeys
}

func (c *kameleoonClient) GetActiveFeatureListForVisitor(visitorCode string) ([]string, error) {
	logging.Info("CALL: kameleoonClient.GetActiveFeatureListForVisitor(visitorCode: %s)", visitorCode)
	err := utils.ValidateVisitorCode(visitorCode)
	var arrayIds []string
	if err == nil {
		visitor := c.visitorManager.GetVisitor(visitorCode)
		featureFlags := c.dataManager.DataFile().GetOrderedFeatureFlags()
		arrayIds = make([]string, 0, len(featureFlags))
		for _, ff := range featureFlags {
			if !ff.GetEnvironmentEnabled() {
				continue
			}
			evalExp := c.evaluate(visitor, visitorCode, ff, false, false)
			variationKey := c.calculateVariationKey(evalExp, ff.GetDefaultVariationKey())
			if variationKey != string(types.VariationOff) {
				arrayIds = append(arrayIds, ff.GetFeatureKey())
			}
		}
	} else {
		arrayIds = []string{}
	}
	logging.Info(
		"RETURN: kameleoonClient.GetActiveFeatureListForVisitor(visitorCode: %s) -> (activeFeatures: %s, error: %s)",
		visitorCode, arrayIds, err)
	return arrayIds, err
}

func (c *kameleoonClient) GetActiveFeatures(visitorCode string) (map[string]types.Variation, error) {
	logging.Info("CALL: kameleoonClient.GetActiveFeatures(visitorCode: %s)", visitorCode)
	if err := utils.ValidateVisitorCode(visitorCode); err != nil {
		logging.Info(
			"RETURN: kameleoonClient.GetActiveFeatures(visitorCode: %s) -> (activeFeatures: <nil>, error: %s)",
			visitorCode, err)
		return nil, err
	}
	visitor := c.visitorManager.GetVisitor(visitorCode)
	mapActiveFeatures := make(map[string]types.Variation)
	for _, ff := range c.dataManager.DataFile().GetOrderedFeatureFlags() {
		if !ff.GetEnvironmentEnabled() {
			continue
		}
		evalExp := c.evaluate(visitor, visitorCode, ff, false, false)
		variationKey := c.calculateVariationKey(evalExp, ff.GetDefaultVariationKey())
		if variationKey == string(types.VariationOff) {
			continue
		}
		variation, _ := ff.GetVariationByKey(variationKey)
		mapActiveFeatures[ff.GetFeatureKey()] = createExternalVariation(variation, evalExp)
	}
	logging.Info(
		"RETURN: kameleoonClient.GetActiveFeatures(visitorCode: %s) -> (activeFeatures: %s, error: <nil>)",
		visitorCode, mapActiveFeatures)
	return mapActiveFeatures, nil
}

func (c *kameleoonClient) GetEngineTrackingCode(visitorCode string) string {
	logging.Info("CALL: kameleoonClient.GetEngineTrackingCode(visitorCode: %s)", visitorCode)
	var engineTrackingCode string
	if c.hybridManager == nil {
		logging.Error("HybridManager wasn't initialized properly. GetEngineTrackingCode method isn't avaiable")
		engineTrackingCode = ""
	} else {
		visitor := c.visitorManager.GetVisitor(visitorCode)
		var variations storage.DataMapStorage[int, *types.AssignedVariation]
		if visitor != nil {
			variations = visitor.Variations()
		}
		engineTrackingCode = c.hybridManager.GetEngineTrackingCode(variations)
	}
	logging.Info("RETURN: kameleoonClient.GetEngineTrackingCode(visitorCode: %s) -> (engineTrackingCode: %s)",
		visitorCode, engineTrackingCode)
	return engineTrackingCode
}

func (c *kameleoonClient) GetVisitorWarehouseAudience(params VisitorWarehouseAudienceParams) (*types.CustomData, error) {
	logging.Info("CALL: kameleoonClient.GetVisitorWarehouseAudience(params: %s)", params)
	customData, err := c.warehouseManager.GetVisitorWarehouseAudience(
		params.VisitorCode, params.WarehouseKey, params.CustomDataIndex, params.Timeout)
	logging.Info("RETURN: kameleoonClient.GetVisitorWarehouseAudience(params: %s) -> (customData: %s, error: %s)",
		params, customData, err)
	return customData, err
}

func (c *kameleoonClient) GetVisitorWarehouseAudienceWithOptParams(
	visitorCode string, customDataIndex int, params ...VisitorWarehouseAudienceOptParams,
) (*types.CustomData, error) {
	logging.Info(
		"CALL: kameleoonClient.GetVisitorWarehouseAudienceWithOptParams(visitorCode: %s, customDataIndex: %s, "+
			"params: %s)", visitorCode, customDataIndex, params)
	var p VisitorWarehouseAudienceOptParams
	if len(params) > 0 {
		p = params[0]
	}
	customData, err := c.warehouseManager.GetVisitorWarehouseAudience(
		visitorCode, p.WarehouseKey, customDataIndex, p.Timeout,
	)
	logging.Info(
		"RETURN: kameleoonClient.GetVisitorWarehouseAudienceWithOptParams(visitorCode: %s, customDataIndex: %s, "+
			"params: %s) -> (customData: %s, error: %s)", visitorCode, customDataIndex, params, customData, err)
	return customData, err
}

func (c *kameleoonClient) setUniqueIdentifier(visitorCode string, isUniqueIdentifier bool) {
	logging.Warning(
		"The 'isUniqueIdentifier' parameter is deprecated. Please, add 'UniqueIdentifier' to a visitor instead.")
	c.visitorManager.AddData(visitorCode, types.NewUniqueIdentifier(isUniqueIdentifier))
}

func (c *kameleoonClient) SetForcedVariation(
	visitorCode string, experimentId int, variationKey string, params ...SetForcedVariationOptParams,
) (err error) {
	logging.Info(
		"CALL: kameleoonClient.SetForcedVariation(visitorCode: %s, experimentId: %s, variationKey: %s, params: %s)",
		visitorCode, experimentId, variationKey, params,
	)
	defer func() {
		logging.Info(
			"RETURN: kameleoonClient.SetForcedVariation(visitorCode: %s, experimentId: %s, variationKey: %s, "+
				"params: %s) -> (error: %s)", visitorCode, experimentId, variationKey, params, err,
		)
	}()
	if err = utils.ValidateVisitorCode(visitorCode); err != nil {
		return
	}
	visitor := c.visitorManager.GetOrCreateVisitor(visitorCode)
	if variationKey != "" {
		ruleInfo, ruleInfoExists := c.dataManager.DataFile().GetRuleInfoByExpId(experimentId)
		if !ruleInfoExists {
			return errs.NewFeatureExperimentNotFound(experimentId)
		}
		var varByExp *types.VariationByExposition
		if varByExp, err = ruleInfo.Rule.GetVariationByKey(variationKey); err != nil {
			return
		}
		var p SetForcedVariationOptParams
		if len(params) > 0 {
			p = params[0]
		} else {
			p = NewSetForcedVariationOptParams()
		}
		forcedVariation := types.NewForcedExperimentVariation(ruleInfo.Rule, varByExp, p.forceTargeting)
		visitor.AddBaseData(true, forcedVariation)
	} else {
		visitor.ResetForcedVariation(experimentId)
	}
	return
}

type evaluatedExperiment struct {
	varByExp   *types.VariationByExposition
	experiment *types.Experiment
	ruleType   types.RuleType
}

func newEvaluatedExperimentFromVarByExpRule(
	varByExp *types.VariationByExposition, rule types.Rule,
) *evaluatedExperiment {
	return &evaluatedExperiment{
		varByExp:   varByExp,
		experiment: &rule.GetRuleBase().Experiment,
		ruleType:   rule.GetRuleBase().Type,
	}
}

func newEvaluatedExperimentFromForcedVariation(forcedVariation types.ForcedVariation) *evaluatedExperiment {
	if (forcedVariation.VarByExp() != nil) && (forcedVariation.Rule() != nil) {
		return newEvaluatedExperimentFromVarByExpRule(forcedVariation.VarByExp(), forcedVariation.Rule())
	}
	return nil
}

// func newEvaluatedExperimentFromForcedExperimentVariation(
// 	forcedVariation *types.ForcedExperimentVariation,
// ) *evaluatedExperiment {
// 	return newEvaluatedExperimentFromVarByExpRule(forcedVariation.VarByExp(), forcedVariation.Rule())
// }

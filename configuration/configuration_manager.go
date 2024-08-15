package configuration

import (
	"github.com/Kameleoon/client-go/v3/logging"
	"sync"
	"time"

	"github.com/Kameleoon/client-go/v3/managers/data"
	"github.com/Kameleoon/client-go/v3/network"
	"github.com/Kameleoon/client-go/v3/realtime"
	"github.com/segmentio/encoding/json"
)

// Not thread-safe
type ConfigurationManager interface {
	Start() error
	OnUpdateConfiguration(handler func())
}

type configurationManagerImpl struct {
	dataManager    data.DataManager
	networkManager network.NetworkManager
	sseClient      realtime.SseClient

	pollingUpdateInterval time.Duration
	environment           string

	updateConfigurationHandler func()

	mx                           sync.Mutex
	lastTS                       int64
	pollingConfigurationTicker   *time.Ticker
	pollingConfigurationStopChan chan bool
	realTimeConfigurationService *realtime.RealTimeConfigurationService
	realTimeUpdateChan           chan realtime.RealTimeEvent
}

func NewConfigurationManager(dataManager data.DataManager, networkManager network.NetworkManager,
	sseClient realtime.SseClient, pollingUpdateInterval time.Duration, environment string,
) *configurationManagerImpl {
	return &configurationManagerImpl{
		dataManager:           dataManager,
		networkManager:        networkManager,
		sseClient:             sseClient,
		pollingUpdateInterval: pollingUpdateInterval,
		environment:           environment,
	}
}

func (cm *configurationManagerImpl) Start() error {
	logging.Debug("CALL: configurationManagerImpl.Start()")
	ok, err := cm.tryFetch(-1)
	if !ok {
		cm.startPollingConfigurationTickerIfNeeded()
	}
	logging.Debug("RETURN: configurationManagerImpl.Start() -> (error: %s)", err)
	return err
}

func (cm *configurationManagerImpl) OnUpdateConfiguration(handler func()) {
	logging.Debug("CALL: configurationManagerImpl.OnUpdateConfiguration()")
	cm.updateConfigurationHandler = handler
	logging.Debug("RETURN: configurationManagerImpl.OnUpdateConfiguration()")
}

func (cm *configurationManagerImpl) tryFetch(ts int64) (bool, error) {
	logging.Debug("CALL: configurationManagerImpl.tryFetch(ts: %s)", ts)
	if (ts != -1) && (ts < cm.lastTS) {
		logging.Debug("RETURN: configurationManagerImpl.tryFetch(ts: %s) -> (isFetched: false, error: <nil>)", ts)
		return false, nil
	}
	if err := cm.fetchConfig(ts); err != nil {
		logging.Error("Fetch failed: %s", err)
		if cm.dataManager.DataFile().Settings().RealTimeUpdate() {
			cm.manageConfigurationUpdate(false)
			logging.Warning("Switched to polling mode due to failed fetch")
		}
		logging.Debug("RETURN: configurationManagerImpl.tryFetch(ts: %s) -> (isFetched: false, error: %s)", ts, err)
		return false, err
	}
	cm.mx.Lock()
	defer cm.mx.Unlock()
	if ts != -1 {
		if ts < cm.lastTS {
			logging.Debug("RETURN: configurationManagerImpl.tryFetch(ts: %s) -> (isFetched: false, error: <nil>)", ts)
			return false, nil
		}
		cm.lastTS = ts
	}
	cm.manageConfigurationUpdate(cm.dataManager.DataFile().Settings().RealTimeUpdate())
	logging.Debug("RETURN: configurationManagerImpl.tryFetch(ts: %s) -> (isFetched: true, error: <nil>)", ts)
	return true, nil
}

func (cm *configurationManagerImpl) fetchConfig(ts int64) error {
	logging.Debug("CALL: configurationManagerImpl.fetchConfig(ts: %s)", ts)
	clientConfig, err := cm.requestClientConfig(ts)
	if err == nil {
		cm.updateDataFile(NewDataFile(clientConfig, cm.environment))
		if ts != -1 && cm.updateConfigurationHandler != nil {
			cm.updateConfigurationHandler()
		}
	} else {
		logging.Error("Failed to fetch: %s", err)
	}
	logging.Debug("RETURN: configurationManagerImpl.fetchConfig(ts: %s) -> (error: %s)", ts, err)
	return err
}

func (cm *configurationManagerImpl) updateDataFile(df *DataFile) {
	logging.Debug("CALL: configurationManagerImpl.updateDataFile(df: %s)", df)
	cm.mx.Lock()
	defer cm.mx.Unlock()
	cm.dataManager.SetDataFile(df)
	cm.networkManager.GetUrlProvider().ApplyDataApiDomain(df.Settings().DataApiDomain())
	logging.Debug("RETURN: configurationManagerImpl.updateDataFile(df: %s)", df)
}

func (cm *configurationManagerImpl) requestClientConfig(ts int64) (Configuration, error) {
	logging.Debug("CALL: configurationManagerImpl.requestClientConfig(ts: %s)", ts)
	if ts == -1 {
		logging.Info("Fetching configuration")
	} else {
		logging.Info("Fetching configuration for TS:%s", ts)
	}
	var campaigns Configuration

	out, err := cm.networkManager.FetchConfiguration(ts)
	if err == nil {
		err = json.Unmarshal(out, &campaigns)
	}
	if err == nil {
		logging.Info("Configuraiton fetched: %s", campaigns)
	} else {
		logging.Error("Failed to fetch client-config: %s", err)
	}
	logging.Debug("RETURN: configurationManagerImpl.requestClientConfig(ts: %s) -> (campaigns: %s, error: %s)",
		ts, campaigns, err)
	return campaigns, err
}

func (cm *configurationManagerImpl) startPollingConfigurationTickerIfNeeded() {
	if cm.pollingConfigurationTicker != nil {
		return
	}
	logging.Debug("CALL: configurationManagerImpl.startPollingConfigurationTickerIfNeeded()")
	cm.pollingConfigurationStopChan = make(chan bool, 1)
	cm.pollingConfigurationTicker = time.NewTicker(cm.pollingUpdateInterval)
	go func() {
		for halt := false; !halt; {
			select {
			case halt = <-cm.pollingConfigurationStopChan:
			default:
				select {
				case halt = <-cm.pollingConfigurationStopChan:
				case <-cm.pollingConfigurationTicker.C:
					cm.tryFetch(-1)
				}
			}
		}
		cm.mx.Lock()
		cm.pollingConfigurationTicker = nil
		cm.mx.Unlock()
	}()
	logging.Info("Configuration polling is started")
	logging.Debug("RETURN: configurationManagerImpl.startPollingConfigurationTickerIfNeeded()")
}

func (cm *configurationManagerImpl) stopPollingConfigurationTickerIfNeeded() {
	logging.Debug("CALL: configurationManagerImpl.stopPollingConfigurationTickerIfNeeded()")
	if cm.pollingConfigurationTicker != nil {
		cm.pollingConfigurationStopChan <- true
		cm.pollingConfigurationTicker.Stop()
		logging.Info("Configuration polling is stopped")
	}
	logging.Debug("RETURN: configurationManagerImpl.stopPollingConfigurationTickerIfNeeded()")
}

func (cm *configurationManagerImpl) startRealTimeConfigurationServiceIfNeeded() {
	if cm.realTimeConfigurationService != nil {
		return
	}
	logging.Debug("CALL: configurationManagerImpl.startRealTimeConfigurationServiceIfNeeded()")
	cm.realTimeUpdateChan = make(chan realtime.RealTimeEvent, 16)
	cm.realTimeConfigurationService = realtime.NewRealTimeConfigurationService(
		cm.networkManager.GetUrlProvider().MakeRealTimeUrl(), cm.realTimeUpdateChan, cm.sseClient)
	go func() {
		for realTimeEvent := range cm.realTimeUpdateChan {
			cm.tryFetch(realTimeEvent.TimeStamp)
		}
	}()
	logging.Info("Configuration streaming is started")
	logging.Debug("RETURN: configurationManagerImpl.startRealTimeConfigurationServiceIfNeeded()")
}

func (cm *configurationManagerImpl) stopRealTimeConfigurationServiceIfNeeded() {
	logging.Debug("CALL: configurationManagerImpl.stopRealTimeConfigurationServiceIfNeeded()")
	if cm.realTimeConfigurationService != nil {
		cm.realTimeConfigurationService.Close()
		cm.realTimeConfigurationService = nil
		logging.Info("Configuration streaming is stopped")
	}
	logging.Debug("RETURN: configurationManagerImpl.stopRealTimeConfigurationServiceIfNeeded()")
}

func (cm *configurationManagerImpl) manageConfigurationUpdate(realTimeUpdate bool) {
	if realTimeUpdate {
		cm.stopPollingConfigurationTickerIfNeeded()
		cm.startRealTimeConfigurationServiceIfNeeded()
	} else {
		cm.stopRealTimeConfigurationServiceIfNeeded()
		cm.startPollingConfigurationTickerIfNeeded()
	}
}

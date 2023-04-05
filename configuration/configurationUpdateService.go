package configuration

import (
	"sync"
	"time"

	"github.com/Kameleoon/client-go/v2/logging"
	"github.com/Kameleoon/client-go/v2/realtime"
	"github.com/Kameleoon/client-go/v2/types"
)

type ConfigurationUpdateService struct {
	mx                           sync.Mutex
	pollingUpdateInterval        time.Duration
	url                          string
	fetchFunc                    func(int64) error
	sseClientSource              func() realtime.SseClient
	logger                       logging.Logger
	settings                     types.Settings
	lastTS                       int64
	pollingConfigurationTicker   *time.Ticker
	pollingConfigurationStopChan chan bool
	realTimeConfigurationService *realtime.RealTimeConfigurationService
	realTimeUpdateChan           chan realtime.RealTimeEvent
}

func (cus *ConfigurationUpdateService) Start(pollingUpdateInterval time.Duration, siteCode string,
	fetchFunc func(int64) error, sseClientSource func() realtime.SseClient,
	logger logging.Logger) error {
	cus.pollingUpdateInterval = pollingUpdateInterval
	cus.url = "https://events.kameleoon.com:8110/sse?siteCode=" + siteCode
	cus.fetchFunc = fetchFunc
	cus.sseClientSource = sseClientSource
	cus.logger = logger
	cus.log("Start-up, fetching is starting")
	return cus.doInitialFetch()
}

func (cus *ConfigurationUpdateService) UpdateSettings(settings types.Settings) {
	cus.settings = settings
}

func (cus *ConfigurationUpdateService) doInitialFetch() error {
	ok, err := cus.tryFetch(-1)
	if !ok {
		cus.startPollingConfigurationTickerIfNeeded()
	}
	return err
}

func (cus *ConfigurationUpdateService) tryFetch(ts int64) (bool, error) {
	if (ts != -1) && (ts < cus.lastTS) {
		return false, nil
	}
	if err := cus.fetchFunc(ts); err != nil {
		cus.log("Fetch failed: %v", err)
		if cus.settings.RealTimeUpdate {
			cus.settings.RealTimeUpdate = false
			cus.manageConfigurationUpdate()
			cus.log("Switched to polling mode due to failed fetch")
		}
		return false, err
	}
	cus.mx.Lock()
	defer cus.mx.Unlock()
	if ts != -1 {
		if ts < cus.lastTS {
			return false, nil
		}
		cus.lastTS = ts
	}
	cus.manageConfigurationUpdate()
	return true, nil
}

func (cus *ConfigurationUpdateService) startPollingConfigurationTickerIfNeeded() {
	if cus.pollingConfigurationTicker != nil {
		return
	}
	cus.pollingConfigurationStopChan = make(chan bool, 1)
	cus.pollingConfigurationTicker = time.NewTicker(cus.pollingUpdateInterval)
	go func() {
		for halt := false; !halt; {
			select {
			case halt = <-cus.pollingConfigurationStopChan:
			default:
				select {
				case halt = <-cus.pollingConfigurationStopChan:
				case <-cus.pollingConfigurationTicker.C:
					cus.tryFetch(-1)
				}
			}
		}
		cus.mx.Lock()
		cus.pollingConfigurationTicker = nil
		cus.mx.Unlock()
	}()
	cus.log("Configuration polling is started")
}

func (cus *ConfigurationUpdateService) stopPollingConfigurationTickerIfNeeded() {
	if cus.pollingConfigurationTicker != nil {
		cus.pollingConfigurationStopChan <- true
		cus.pollingConfigurationTicker.Stop()
		cus.log("Configuration polling is stopped")
	}
}

func (cus *ConfigurationUpdateService) startRealTimeConfigurationServiceIfNeeded() {
	if cus.realTimeConfigurationService != nil {
		return
	}
	cus.realTimeUpdateChan = make(chan realtime.RealTimeEvent, 16)
	var sse realtime.SseClient
	if cus.sseClientSource == nil {
		sse = &realtime.NetSseClient{}
	} else {
		sse = cus.sseClientSource()
	}
	cus.realTimeConfigurationService = realtime.NewRealTimeConfigurationService(
		cus.url, cus.realTimeUpdateChan, sse, cus.logger)
	go func() {
		for realTimeEvent := range cus.realTimeUpdateChan {
			cus.tryFetch(realTimeEvent.TimeStamp)
		}
	}()
	cus.log("Configuration streaming is started")
}

func (cus *ConfigurationUpdateService) stopRealTimeConfigurationServiceIfNeeded() {
	if cus.realTimeConfigurationService != nil {
		cus.realTimeConfigurationService.Close()
		cus.realTimeConfigurationService = nil
		cus.log("Configuration streaming is stopped")
	}
}

func (cus *ConfigurationUpdateService) manageConfigurationUpdate() {
	if cus.settings.RealTimeUpdate {
		cus.stopPollingConfigurationTickerIfNeeded()
		cus.startRealTimeConfigurationServiceIfNeeded()
	} else {
		cus.stopRealTimeConfigurationServiceIfNeeded()
		cus.startPollingConfigurationTickerIfNeeded()
	}
}

func (cus *ConfigurationUpdateService) log(format string, args ...interface{}) {
	if cus.logger != nil {
		cus.logger.Printf(format, args...)
	}
}

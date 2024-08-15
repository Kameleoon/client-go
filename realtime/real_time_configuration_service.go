package realtime

import (
	"encoding/json"
	"sync"

	"github.com/Kameleoon/client-go/v3/logging"
	net "github.com/subchord/go-sse"
)

const configurationUpdateEvent = "configuration-update-event"

type RealTimeConfigurationService struct {
	closeFlagMx sync.Mutex
	url         string
	updateChan  chan RealTimeEvent
	sse         SseClient
	closeFlag   bool
	closeChan   chan bool
}

func NewRealTimeConfigurationService(url string, updateChan chan RealTimeEvent,
	sse SseClient) *RealTimeConfigurationService {
	rtcs := &RealTimeConfigurationService{
		url:        url,
		updateChan: updateChan,
		sse:        sse,
		closeChan:  make(chan bool, 1),
	}
	go rtcs.run()
	return rtcs
}

func (rtcs *RealTimeConfigurationService) Close() {
	rtcs.closeFlagMx.Lock()
	defer rtcs.closeFlagMx.Unlock()
	if rtcs.closeFlag {
		return
	}
	rtcs.closeFlag = true
	rtcs.closeChan <- true
}

func (rtcs *RealTimeConfigurationService) run() {
	logging.Info("Real-Time Configuration Service started")
	if rtcs.sse == nil {
		logging.Error("SSE Client is not provided, Real-time Configuration Service is not started")
		return
	}
	for !rtcs.closeFlag {
		rtcs.sse.Dispose()
		if err := rtcs.sse.Init(rtcs.url); err != nil {
			logging.Error("Failed to open SSE connection: %s", err)
			continue
		}
		logging.Info("SSE connection open")
		for halt := false; !halt; {
			select {
			case halt = <-rtcs.closeChan:
			case err := <-rtcs.sse.GetErrorChan():
				logging.Error("Error occurred within SSE client: %s", err)
				halt = true
			default:
				select {
				case halt = <-rtcs.closeChan:
				case err := <-rtcs.sse.GetErrorChan():
					halt = true
					logging.Error("Error occurred within SSE client: %s", err)
				case evt := <-rtcs.sse.GetEventChan():
					logging.Info("Got %s SSE event", configurationUpdateEvent)
					if err := rtcs.handleEvent(evt); err != nil {
						logging.Error("Error occurred during SSE event parsing: %s", err)
					}
				}
			}
		}
		logging.Info("SSE connection closed")
	}
	close(rtcs.updateChan)
	rtcs.sse.Dispose()
	logging.Info("Real-Time Configuration Service stopped")
}

func (rtcs *RealTimeConfigurationService) handleEvent(evt net.Event) error {
	b := []byte(evt.GetData())
	var rtEvent RealTimeEvent
	if err := json.Unmarshal(b, &rtEvent); err != nil {
		return err
	}
	rtcs.updateChan <- rtEvent
	return nil
}

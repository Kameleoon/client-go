package realtime

import (
	"encoding/json"
	"sync"

	"github.com/Kameleoon/client-go/v2/logging"
	net "github.com/subchord/go-sse"
)

const configurationUpdateEvent = "configuration-update-event"

type RealTimeConfigurationService struct {
	closeFlagMx sync.Mutex
	url         string
	updateChan  chan RealTimeEvent
	sse         SseClient
	logger      logging.Logger
	closeFlag   bool
	closeChan   chan bool
}

func NewRealTimeConfigurationService(url string, updateChan chan RealTimeEvent, sse SseClient,
	logger logging.Logger) *RealTimeConfigurationService {
	rtcs := &RealTimeConfigurationService{
		url:        url,
		updateChan: updateChan,
		sse:        sse,
		logger:     logger,
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
	rtcs.log("Real-Time Configuration Service started")
	for !rtcs.closeFlag {
		rtcs.sse.Dispose(rtcs.log)
		if err := rtcs.sse.Init(rtcs.url); err != nil {
			rtcs.log("Failed to open SSE connection: %v", err)
			continue
		}
		rtcs.log("SSE connection open")
		for halt := false; !halt; {
			select {
			case halt = <-rtcs.closeChan:
			case err := <-rtcs.sse.GetErrorChan():
				rtcs.log("Error occurred within SSE client: %v", err)
				halt = true
			default:
				select {
				case halt = <-rtcs.closeChan:
				case err := <-rtcs.sse.GetErrorChan():
					halt = true
					rtcs.log("Error occurred within SSE client: %v", err)
				case evt := <-rtcs.sse.GetEventChan():
					rtcs.log("Got '%s' SSE event", configurationUpdateEvent)
					if err := rtcs.handleEvent(evt); err != nil {
						rtcs.log("Error occurred during SSE event parsing: %v", err)
					}
				}
			}
		}
		rtcs.log("SSE connection closed")
	}
	close(rtcs.updateChan)
	rtcs.sse.Dispose(rtcs.log)
	rtcs.log("Real-Time Configuration Service stopped")
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

func (rtcs *RealTimeConfigurationService) log(format string, args ...interface{}) {
	if rtcs.logger != nil {
		rtcs.logger.Printf(format, args...)
	}
}

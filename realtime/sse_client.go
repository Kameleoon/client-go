package realtime

import (
	"github.com/Kameleoon/client-go/v3/logging"
	"net/http"

	net "github.com/subchord/go-sse"
)

type SseClient interface {
	Init(url string) error
	Dispose()
	GetErrorChan() <-chan error
	GetEventChan() <-chan net.Event
}

type NetSseClient struct {
	feed *net.SSEFeed
	sub  *net.Subscription
}

func (sse *NetSseClient) Init(url string) error {
	if err := sse.initFeed(url); err != nil {
		return err
	}
	if err := sse.initSub(); err != nil {
		return err
	}
	return nil
}

func (sse *NetSseClient) initFeed(url string) error {
	headers := map[string][]string{
		http.CanonicalHeaderKey("Accept"):        {"text/event-stream"},
		http.CanonicalHeaderKey("Cache-Control"): {"no-cache"},
		http.CanonicalHeaderKey("Connection"):    {"Keep-Alive"},
	}
	feed, err := net.ConnectWithSSEFeed(url, headers)
	if err != nil {
		return err
	}
	sse.feed = feed
	return nil
}

func (sse *NetSseClient) initSub() error {
	sub, err := sse.feed.Subscribe(configurationUpdateEvent)
	if err != nil {
		return err
	}
	sse.sub = sub
	return nil
}

func (sse *NetSseClient) Dispose() {
	defer func() { // "close of closed channel" panic occurs
		if err := recover(); err != nil {
			logging.Warning("Panic occurred during SSE dispose process: %s", err)
		}
	}()
	if sse.feed != nil {
		sse.feed.Close()
		sse.feed = nil
	}
	if sse.sub != nil {
		sse.sub.Close()
		sse.sub = nil
	}
}

func (sse *NetSseClient) GetErrorChan() <-chan error {
	return sse.sub.ErrFeed()
}

func (sse *NetSseClient) GetEventChan() <-chan net.Event {
	return sse.sub.Feed()
}

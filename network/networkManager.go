package network

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Kameleoon/client-go/v2/logging"
)

// declaration

type NetworkManager interface {
	GetEnvironment() string
	GetDefaultTimeout() time.Duration
	GetNetProvider() NetProvider
	GetUrlProvider() *UrlProvider

	FetchConfiguration(ts int64, timeout time.Duration, out chan<- json.RawMessage, err chan<- error)
	GetRemoteData(key string, timeout time.Duration, out chan<- json.RawMessage, err chan<- error)
	GetVisitorRemoteData(visitorCode string, timeout time.Duration, out chan<- json.RawMessage, err chan<- error)

	SendTrackingData(visitorCode string, lines []QueryEncodable, userAgent string, authToken string,
		timeout time.Duration, out chan<- bool, err chan<- error)
	FetchBearerToken(clientId string, clientSecret string, timeout time.Duration,
		out chan<- json.RawMessage, err chan<- error)
}

// base implementation

type NetworkManagerImpl struct {
	Environment            string
	DefaultTimeout         time.Duration
	NetProvider            NetProvider
	UrlProvider            *UrlProvider
	Logger                 logging.Logger
	TrackingCallRetryDelay time.Duration
}

func NewNetworkManagerImpl(environment string, defaultTimeout time.Duration, netProvider NetProvider,
	urlProvider *UrlProvider, logger logging.Logger) *NetworkManagerImpl {
	return &NetworkManagerImpl{
		Environment:            environment,
		DefaultTimeout:         defaultTimeout,
		NetProvider:            netProvider,
		UrlProvider:            urlProvider,
		Logger:                 logger,
		TrackingCallRetryDelay: DefaultTrackingCallRetryDelay,
	}
}

func (nm *NetworkManagerImpl) GetEnvironment() string {
	return nm.Environment
}
func (nm *NetworkManagerImpl) GetDefaultTimeout() time.Duration {
	return nm.DefaultTimeout
}
func (nm *NetworkManagerImpl) GetNetProvider() NetProvider {
	return nm.NetProvider
}
func (nm *NetworkManagerImpl) GetUrlProvider() *UrlProvider {
	return nm.UrlProvider
}

// API call commons

func (nm *NetworkManagerImpl) ensureTimeout(timeout *time.Duration) {
	if *timeout < 0 {
		*timeout = nm.DefaultTimeout
	}
}

func (nm *NetworkManagerImpl) makeCall(request Request, attemptCount int, retryDelay time.Duration,
	textOut chan<- string, jsonOut chan<- json.RawMessage, statusOut chan<- bool, errChan chan<- error) {
	go func() {
		var err error
		for i := 0; i < attemptCount; i++ {
			if (i > 0) && (retryDelay > 0) {
				time.Sleep(retryDelay)
			}
			response := nm.NetProvider.Call(request)
			if response.Err != nil {
				err = response.Err
				nm.logErrOccurred(request, response.Err)
				continue
			}
			if response.Code/100 != 2 {
				err = ErrUnexpectedResponseStatus{Code: response.Code}
				nm.logUnexpectedCode(request, response.Code)
				continue
			}
			if textOut != nil {
				textOut <- string(response.Body)
			}
			if jsonOut != nil {
				jsonOut <- response.Body
			}
			if statusOut != nil {
				statusOut <- true
			}
			return
		}
		errChan <- err
	}()
}

func (nm *NetworkManagerImpl) logErrOccurred(request Request, err error) {
	nm.Logger.Printf("%s: Error occurred during request: %v", makeErrMsg(request), err)
}
func (nm *NetworkManagerImpl) logUnexpectedCode(request Request, code int) {
	nm.Logger.Printf("%s: Received unexpected status code '%d'", makeErrMsg(request), code)
}

func makeErrMsg(request Request) string {
	if len(request.Data) == 0 {
		return fmt.Sprintf("%s call '%s' failed", request.Method, request.Url)
	}
	return fmt.Sprintf("%s call '%s' (data '%s') failed", request.Method, request.Url, request.Data)
}

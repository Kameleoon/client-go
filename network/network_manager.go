package network

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Kameleoon/client-go/v3/errs"
	"github.com/Kameleoon/client-go/v3/logging"
	"github.com/Kameleoon/client-go/v3/types"
)

const (
	networkCallRetriesNumberCritical    = 2
	NetworkCallAttemptsNumberCritical   = networkCallRetriesNumberCritical + 1 // +1 for initial request
	NetworkCallAttemptsNumberUncritical = 1

	codeUnauthorized = 401
	codeForbidden    = 403
)

// declaration

type NetworkManager interface {
	GetEnvironment() string
	GetDefaultTimeout() time.Duration
	GetNetProvider() NetProvider
	GetUrlProvider() UrlProvider
	GetAccessTokenSource() AccessTokenSource

	// Automation API
	FetchAccessJWToken(clientId string, clientSecret string, timeout time.Duration) (json.RawMessage, error)

	// SDK config API
	FetchConfiguration(ts int64) (json.RawMessage, error)

	// Data API
	GetRemoteData(key string, timeout time.Duration) (json.RawMessage, error)
	GetRemoteVisitorData(visitorCode string, filter types.RemoteVisitorDataFilter, isUniqueIdentifier bool,
		timeout time.Duration) (json.RawMessage, error)
	SendTrackingData(visitorCode string, lines []types.Sendable, userAgent string,
		isUniqueIdentifier bool) (bool, error)
}

// base implementation

type NetworkManagerImpl struct {
	Environment            string
	DefaultTimeout         time.Duration
	NetProvider            NetProvider
	UrlProvider            UrlProvider
	Logger                 logging.Logger
	TrackingCallRetryDelay time.Duration
	accessTokenSource      AccessTokenSource
}

func NewNetworkManagerImpl(
	environment string,
	defaultTimeout time.Duration,
	netProvider NetProvider,
	urlProvider UrlProvider,
	accessTokenSourceFactory AccessTokenSourceFactory,
	logger logging.Logger,
) *NetworkManagerImpl {
	nm := &NetworkManagerImpl{
		Environment:            environment,
		DefaultTimeout:         defaultTimeout,
		NetProvider:            netProvider,
		UrlProvider:            urlProvider,
		Logger:                 logger,
		TrackingCallRetryDelay: DefaultTrackingCallRetryDelay,
	}
	nm.accessTokenSource = accessTokenSourceFactory.create(nm, logger)
	return nm
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

func (nm *NetworkManagerImpl) GetUrlProvider() UrlProvider {
	return nm.UrlProvider
}

func (nm *NetworkManagerImpl) GetAccessTokenSource() AccessTokenSource {
	return nm.accessTokenSource
}

// API call commons

func (nm *NetworkManagerImpl) ensureTimeout(request *Request) {
	if request.Timeout <= 0 {
		request.Timeout = nm.DefaultTimeout
	}
}

func (nm *NetworkManagerImpl) makeCall(request *Request, attemptCount int, retryDelay time.Duration) ([]byte, error) {

	nm.ensureTimeout(request)
	var err error
	var isTokenRejected bool
	for i := 0; i < attemptCount; i++ {
		if (i > 0) && (retryDelay > 0) {
			time.Sleep(retryDelay)
		}
		nm.authorizeIfRequired(request)
		response := nm.NetProvider.Call(request)
		if isTokenRejected, err = nm.processErrors(request, &response); err == nil {
			return response.Body, nil
		}
	}
	if isTokenRejected {
		request.AccessToken = ""
		response := nm.NetProvider.Call(request)
		if _, err = nm.processErrors(request, &response); err == nil {
			return response.Body, nil
		}
	}
	return nil, err
}

func (nm *NetworkManagerImpl) authorizeIfRequired(request *Request) {
	if request.IsAuthRequired {
		request.AccessToken = nm.accessTokenSource.GetToken(request.Timeout)
	}
}

func (nm *NetworkManagerImpl) processErrors(request *Request, response *Response) (bool, error) {
	var err error
	var isTokenRejected bool
	if response.Err != nil {
		err = response.Err
		nm.logErrOccurred(request, response.Err)
	} else if response.Code/100 != 2 {
		err = errs.NewUnexpectedStatusCode(response.Code)
		nm.logUnexpectedCode(request, response.Code)
		if (response.Code == codeUnauthorized || response.Code == codeForbidden) && request.AccessToken != "" {
			nm.accessTokenSource.DiscardToken(request.AccessToken)
			isTokenRejected = true
		}
	}
	return isTokenRejected, err
}

func (nm *NetworkManagerImpl) logErrOccurred(request *Request, err error) {
	nm.Logger.Printf("%s: Error occurred during request: %v", makeErrMsg(request), err)
}
func (nm *NetworkManagerImpl) logUnexpectedCode(request *Request, code int) {
	nm.Logger.Printf("%s: Received unexpected status code '%d'", makeErrMsg(request), code)
}

func makeErrMsg(request *Request) string {
	if len(request.Data) == 0 {
		return fmt.Sprintf("%s call '%s' failed", request.Method, request.Url)
	}
	return fmt.Sprintf("%s call '%s' (data '%s') failed", request.Method, request.Url, request.Data)
}

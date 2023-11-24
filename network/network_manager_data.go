package network

import (
	"encoding/json"
	"time"
)

func (nm *NetworkManagerImpl) GetRemoteData(key string, timeout time.Duration) (json.RawMessage, error) {
	url := nm.UrlProvider.MakeApiDataGetRequestUrl(key)
	nm.ensureTimeout(&timeout)
	request := Request{
		Method:      HttpGet,
		Url:         url,
		ContentType: JsonContentType,
		Timeout:     timeout,
	}
	return nm.makeCall(&request, NetworkCallAttemptsNumberUncritical, -1)
}

func (nm *NetworkManagerImpl) GetRemoteVisitorData(visitorCode string, timeout time.Duration) (json.RawMessage, error) {
	url := nm.UrlProvider.MakeVisitorDataGetUrl(visitorCode)
	nm.ensureTimeout(&timeout)
	request := Request{
		Method:      HttpGet,
		Url:         url,
		ContentType: JsonContentType,
		Timeout:     timeout,
	}
	return nm.makeCall(&request, NetworkCallAttemptsNumberUncritical, -1)
}

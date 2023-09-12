package network

import (
	"encoding/json"
	"time"
)

const fetchConfigurationAttemptNumber = 4

func (nm *NetworkManagerImpl) FetchConfiguration(ts int64, timeout time.Duration) (json.RawMessage, error) {
	url := nm.UrlProvider.MakeConfigurationUrl(nm.Environment, ts)
	nm.ensureTimeout(&timeout)
	request := Request{
		Method:      HttpGet,
		Url:         url,
		ContentType: JsonContentType,
		Timeout:     timeout,
	}
	return nm.makeCall(request, fetchConfigurationAttemptNumber, -1)
}

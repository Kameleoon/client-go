package network

import (
	"encoding/json"
	"time"
)

const fetchConfigurationAttemptNumber = 4

func (nm *NetworkManagerImpl) FetchConfiguration(ts int64, timeout time.Duration,
	out chan<- json.RawMessage, err chan<- error) {
	url := nm.UrlProvider.MakeConfigurationUrl(nm.Environment, ts)
	nm.ensureTimeout(&timeout)
	request := Request{
		Method:      HttpGet,
		Url:         url,
		ContentType: JsonContentType,
		Timeout:     timeout,
	}
	nm.makeCall(request, fetchConfigurationAttemptNumber, time.Duration(-1), nil, out, nil, err)
}

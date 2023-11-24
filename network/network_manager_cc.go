package network

import (
	"encoding/json"
	"time"
)

const (
	HeaderSdkType    = "X-Kameleoon-SDK-Type"
	HeaderSdkVersion = "X-Kameleoon-SDK-Version"
)

func (nm *NetworkManagerImpl) FetchConfiguration(ts int64, timeout time.Duration) (json.RawMessage, error) {
	url := nm.UrlProvider.MakeConfigurationUrl(nm.Environment, ts)
	nm.ensureTimeout(&timeout)
	request := &Request{
		Method:      HttpGet,
		Url:         url,
		ContentType: JsonContentType,
		Timeout:     timeout,
		Headers: map[string]string{
			HeaderSdkType:    nm.UrlProvider.SdkName,
			HeaderSdkVersion: nm.UrlProvider.SdkVersion,
		},
	}
	return nm.makeCall(request, NetworkCallAttemptsNumberCritical, -1)
}

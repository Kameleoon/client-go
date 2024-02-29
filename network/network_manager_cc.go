package network

import (
	"encoding/json"
)

const (
	HeaderSdkType    = "X-Kameleoon-SDK-Type"
	HeaderSdkVersion = "X-Kameleoon-SDK-Version"
)

func (nm *NetworkManagerImpl) FetchConfiguration(ts int64) (json.RawMessage, error) {
	url := nm.UrlProvider.MakeConfigurationUrl(nm.Environment, ts)
	request := &Request{
		Method:      HttpGet,
		Url:         url,
		ContentType: JsonContentType,
		Headers: map[string]string{
			HeaderSdkType:    nm.UrlProvider.SdkName(),
			HeaderSdkVersion: nm.UrlProvider.SdkVersion(),
		},
	}
	return nm.makeCall(request, NetworkCallAttemptsNumberCritical, -1)
}

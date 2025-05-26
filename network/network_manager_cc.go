package network

const (
	HeaderSdkType         = "X-Kameleoon-SDK-Type"
	HeaderSdkVersion      = "X-Kameleoon-SDK-Version"
	HeaderIfModifiedSince = "If-Modified-Since"
	HeaderLastModified    = "Last-Modified"
)

func (nm *NetworkManagerImpl) FetchConfiguration(ts int64, ifModifiedSince string) (FetchedConfiguration, error) {
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
	if ifModifiedSince != "" {
		request.Headers[HeaderIfModifiedSince] = ifModifiedSince
	}
	response, err := nm.makeCall(request, NetworkCallAttemptsNumberCritical, -1, HeaderLastModified)
	if err != nil {
		return FetchedConfiguration{}, err
	}
	lastModified := response.HeadersRead[HeaderLastModified]
	return FetchedConfiguration{Configuration: response.Body, LastModified: lastModified}, nil
}

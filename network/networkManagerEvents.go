package network

import (
	"encoding/json"
	"strings"
	"time"
)

const (
	TrackingCallRetryNumber       = 3
	DefaultTrackingCallRetryDelay = time.Second * 5
)

func (nm *NetworkManagerImpl) GetRemoteData(key string, timeout time.Duration,
	out chan<- json.RawMessage, err chan<- error) {
	url := nm.UrlProvider.MakeApiDataGetRequestUrl(key)
	nm.ensureTimeout(&timeout)
	request := Request{
		Method:      HttpGet,
		Url:         url,
		ContentType: JsonContentType,
		Timeout:     timeout,
	}
	nm.makeCall(request, 1, time.Duration(-1), nil, out, nil, err)
}

func (nm *NetworkManagerImpl) GetVisitorRemoteData(visitorCode string, timeout time.Duration,
	out chan<- json.RawMessage, err chan<- error) {
	url := nm.UrlProvider.MakeVisitorDataGetUrl(visitorCode)
	nm.ensureTimeout(&timeout)
	request := Request{
		Method:      HttpGet,
		Url:         url,
		ContentType: JsonContentType,
		Timeout:     timeout,
	}
	nm.makeCall(request, 1, time.Duration(-1), nil, out, nil, err)
}

func (nm *NetworkManagerImpl) SendTrackingData(visitorCode string, lines []QueryEncodable, userAgent string,
	authToken string, timeout time.Duration, out chan<- bool, err chan<- error) {
	if len(lines) == 0 {
		go func() { out <- false }()
		return
	}
	url := nm.UrlProvider.MakeTrackingUrl(visitorCode)
	nm.ensureTimeout(&timeout)
	data := formTrackingRequestData(lines)
	request := Request{
		Method:      HttpPost,
		Url:         url,
		ContentType: TextContentType,
		AuthToken:   authToken,
		Timeout:     timeout,
		UserAgent:   userAgent,
		Data:        data,
	}
	nm.makeCall(request, TrackingCallRetryNumber+1, nm.TrackingCallRetryDelay, nil, nil, out, err)
}
func formTrackingRequestData(lines []QueryEncodable) string {
	sb := strings.Builder{}
	for _, line := range lines {
		line := line.QueryEncode()
		if len(line) > 0 {
			if sb.Len() > 0 {
				sb.WriteRune('\n')
			}
			sb.WriteString(line)
		}
	}
	return sb.String()
}

package network

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/Kameleoon/client-go/v3/types"
)

const (
	DefaultTrackingCallRetryDelay = time.Second * 5
)

func (nm *NetworkManagerImpl) GetRemoteData(key string, timeout time.Duration) (json.RawMessage, error) {
	url := nm.UrlProvider.MakeApiDataGetRequestUrl(key)
	request := Request{
		Method:         HttpGet,
		Url:            url,
		ContentType:    JsonContentType,
		Timeout:        timeout,
		IsAuthRequired: true,
	}
	return nm.makeCall(&request, NetworkCallAttemptsNumberUncritical, -1)
}

func (nm *NetworkManagerImpl) GetRemoteVisitorData(visitorCode string, timeout time.Duration) (json.RawMessage, error) {
	url := nm.UrlProvider.MakeVisitorDataGetUrl(visitorCode)
	request := Request{
		Method:         HttpGet,
		Url:            url,
		ContentType:    JsonContentType,
		Timeout:        timeout,
		IsAuthRequired: true,
	}
	return nm.makeCall(&request, NetworkCallAttemptsNumberUncritical, -1)
}

func (nm *NetworkManagerImpl) SendTrackingData(visitorCode string, lines []types.Sendable,
	userAgent string) (bool, error) {

	if len(lines) == 0 {
		return false, nil
	}
	url := nm.UrlProvider.MakeTrackingUrl(visitorCode)
	data := formTrackingRequestData(lines)
	request := Request{
		Method:         HttpPost,
		Url:            url,
		ContentType:    TextContentType,
		UserAgent:      userAgent,
		Data:           data,
		Timeout:        nm.DefaultTimeout,
		IsAuthRequired: true,
	}
	_, err := nm.makeCall(&request, NetworkCallAttemptsNumberCritical, nm.TrackingCallRetryDelay)
	if err != nil {
		return false, err
	}
	return true, nil
}

func formTrackingRequestData(qes []types.Sendable) string {
	sb := strings.Builder{}
	for _, qe := range qes {
		line := qe.QueryEncode()
		if len(line) > 0 {
			if sb.Len() > 0 {
				sb.WriteRune('\n')
			}
			sb.WriteString(line)
		}
	}
	return sb.String()
}

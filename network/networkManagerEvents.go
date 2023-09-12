package network

import (
	"strings"
	"time"
)

const (
	TrackingCallRetryNumber       = 3
	DefaultTrackingCallRetryDelay = time.Second * 5
)

func (nm *NetworkManagerImpl) SendTrackingData(visitorCode string, lines []QueryEncodable, userAgent string,
	authToken string, timeout time.Duration) (bool, error) {
	if len(lines) == 0 {
		return false, nil
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
	_, err := nm.makeCall(request, TrackingCallRetryNumber+1, nm.TrackingCallRetryDelay)
	if err != nil {
		return false, err
	}
	return true, nil
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

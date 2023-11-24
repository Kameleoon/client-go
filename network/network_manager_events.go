package network

import (
	"strings"
	"time"

	"github.com/Kameleoon/client-go/v3/types"
)

const (
	DefaultTrackingCallRetryDelay = time.Second * 5
)

func (nm *NetworkManagerImpl) SendTrackingData(visitorCode string, lines []types.Sendable, userAgent string,
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

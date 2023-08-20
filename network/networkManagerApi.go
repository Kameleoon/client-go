package network

import (
	"encoding/json"
	"time"
)

func (nm *NetworkManagerImpl) FetchBearerToken(clientId string, clientSecret string, timeout time.Duration,
	out chan<- json.RawMessage, err chan<- error) {
	url := nm.UrlProvider.MakeBearerTokenUrl()
	nm.ensureTimeout(&timeout)
	data := formFetchBearerTokenData(clientId, clientSecret)
	request := Request{
		Method:      HttpPost,
		Url:         url,
		ContentType: FormContentType,
		Timeout:     timeout,
		Data:        data,
	}
	nm.makeCall(request, 1, time.Duration(-1), nil, out, nil, err)
}
func formFetchBearerTokenData(clientId string, clientSecret string) string {
	qb := NewQueryBuilder()
	qb.Append(QPGrantType, "client_credentials")
	qb.Append(QPClientId, clientId)
	qb.Append(QPClientSecret, clientSecret)
	return qb.String()
}

package network

import (
	"encoding/json"
	"time"
)

func (nm *NetworkManagerImpl) FetchBearerToken(clientId string, clientSecret string,
	timeout time.Duration) (json.RawMessage, error) {
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
	return nm.makeCall(request, 1, -1)
}
func formFetchBearerTokenData(clientId string, clientSecret string) string {
	qb := NewQueryBuilder()
	qb.Append(QPGrantType, "client_credentials")
	qb.Append(QPClientId, clientId)
	qb.Append(QPClientSecret, clientSecret)
	return qb.String()
}

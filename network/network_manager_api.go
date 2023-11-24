package network

import (
	"encoding/json"
	"time"

	"github.com/Kameleoon/client-go/v3/utils"
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
	return nm.makeCall(&request, NetworkCallAttemptsNumberCritical, -1)
}
func formFetchBearerTokenData(clientId string, clientSecret string) string {
	qb := utils.NewQueryBuilder()
	qb.Append(utils.QPGrantType, "client_credentials")
	qb.Append(utils.QPClientId, clientId)
	qb.Append(utils.QPClientSecret, clientSecret)
	return qb.String()
}

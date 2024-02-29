package network

import (
	"encoding/json"
	"time"

	"github.com/Kameleoon/client-go/v3/utils"
)

const (
	GrantType             = "client_credentials"
	HeaderContentTypeName = "Content-Type"
)

func (nm *NetworkManagerImpl) FetchAccessJWToken(clientId string, clientSecret string,
	timeout time.Duration) (json.RawMessage, error) {

	url := nm.UrlProvider.MakeAccessTokenUrl()
	data := nm.formFetchAccessTokenData(clientId, clientSecret)
	request := Request{
		Method:      HttpPost,
		Url:         url,
		ContentType: FormContentType,
		Data:        data,
		Timeout:     timeout,
	}
	return nm.makeCall(&request, NetworkCallAttemptsNumberUncritical, -1)
}

func (nm *NetworkManagerImpl) formFetchAccessTokenData(clientId string, clientSecret string) string {
	qb := utils.NewQueryBuilder()
	qb.Append(utils.QPGrantType, "client_credentials")
	qb.Append(utils.QPClientId, clientId)
	qb.Append(utils.QPClientSecret, clientSecret)
	return qb.String()
}

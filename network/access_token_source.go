package network

import (
	"encoding/json"
	"time"

	"github.com/Kameleoon/client-go/v3/logging"
)

const (
	TokenExpirationGap   = 60   // in seconds
	TokenObsolescenceGap = 1800 // in seconds
)

type AccessTokenSource interface {
	GetToken(timeout time.Duration) string
	DiscardToken(token string)
}

type AccessTokenSourceImpl struct {
	clientId       string
	clientSecret   string
	networkManager NetworkManager
	logger         logging.Logger
	cachedToken    *expiringToken // pointer for thread-safe
	fetching       bool
}

func (ats *AccessTokenSourceImpl) GetToken(timeout time.Duration) string {
	now := time.Now()
	token := ats.cachedToken
	if token != nil && !token.isExpired(now) {
		if !ats.fetching && token.isObsolete(now) {
			ats.fetching = true // set `fetching` here as well to reduce the number of requests until goroutine runned
			go ats.fetchToken(timeout)
		}
		return token.value
	}
	return ats.fetchToken(timeout)

}

func (ats *AccessTokenSourceImpl) DiscardToken(token string) {
	cachedToken := ats.cachedToken
	if cachedToken != nil && cachedToken.value == token {
		ats.cachedToken = nil
	}
}

func (ats *AccessTokenSourceImpl) fetchToken(timeout time.Duration) string {
	ats.fetching = true
	defer func() { ats.fetching = false }()
	ats.logger.Printf("Fetching access token")
	jsonResponse, err := ats.networkManager.FetchAccessJWToken(ats.clientId, ats.clientSecret, timeout)
	if err != nil {
		ats.logger.Printf("Failed to fetch access token: %v", err)
		return ""
	}
	accessTokenResponse := accessTokenResponse{}
	err = json.Unmarshal(jsonResponse, &accessTokenResponse)
	if err != nil {
		ats.logger.Printf("Failed to unmarshal access token: %v", err)
		return ""
	}
	ats.logger.Printf("Access token is fetched: %s", accessTokenResponse.Token)
	ats.handleFetchedToken(accessTokenResponse)
	return accessTokenResponse.Token
}

func (ats *AccessTokenSourceImpl) handleFetchedToken(accessTokenResponse accessTokenResponse) {
	expiresIn := accessTokenResponse.ExpiresIn
	now := time.Now()
	expTime := now.Add(time.Second * time.Duration(expiresIn-TokenExpirationGap))
	var obsTime time.Time
	if expiresIn > TokenObsolescenceGap {
		obsTime = now.Add(time.Second * time.Duration(expiresIn-TokenObsolescenceGap))
	} else {
		obsTime = expTime
		if expiresIn <= TokenExpirationGap {
			ats.logger.Printf(
				"Kameleoon SDK: Access token life time (%ds) is not long enough to cache the token\n", expiresIn)
		} else {
			ats.logger.Printf("Kameleoon SDK: Access token life time (%ds) is not long enough to refresh cached token"+
				"in background\n", expiresIn)
		}
	}
	ats.cachedToken = newExpiringToken(accessTokenResponse.Token, expTime, obsTime)
}

type expiringToken struct {
	value            string
	expirationTime   time.Time
	obsolescenceTime time.Time
}

func newExpiringToken(value string, expirationTime, obsolescenceTime time.Time) *expiringToken {
	return &expiringToken{
		value:            value,
		expirationTime:   expirationTime,
		obsolescenceTime: obsolescenceTime,
	}
}

func (et *expiringToken) isExpired(now time.Time) bool {
	return !now.Before(et.expirationTime)
}

func (et *expiringToken) isObsolete(now time.Time) bool {
	return !now.Before(et.obsolescenceTime)
}

type accessTokenResponse struct {
	Token     string `json:"access_token"`
	ExpiresIn int    `json:"expires_in"`
}

package network

import "github.com/Kameleoon/client-go/v3/logging"

type AccessTokenSourceFactory interface {
	create(networkManager NetworkManager, logger logging.Logger) AccessTokenSource
}

type AccessTokenSourceFactoryImpl struct {
	ClientId     string
	ClientSecret string
}

func (f *AccessTokenSourceFactoryImpl) create(networkManager NetworkManager, logger logging.Logger) AccessTokenSource {
	return &AccessTokenSourceImpl{
		clientId:       f.ClientId,
		clientSecret:   f.ClientSecret,
		networkManager: networkManager,
		logger:         logger,
	}
}

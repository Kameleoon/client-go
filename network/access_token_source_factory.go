package network

type AccessTokenSourceFactory interface {
	create(networkManager NetworkManager) AccessTokenSource
}

type AccessTokenSourceFactoryImpl struct {
	ClientId     string
	ClientSecret string
}

func (f *AccessTokenSourceFactoryImpl) create(networkManager NetworkManager) AccessTokenSource {
	return &AccessTokenSourceImpl{
		clientId:       f.ClientId,
		clientSecret:   f.ClientSecret,
		networkManager: networkManager,
	}
}

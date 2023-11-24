package kameleoon

import (
	cmap "github.com/orcaman/concurrent-map/v2"
)

var KameleoonClientFactory = newKameleoonClientFactory()

type kameleoonClientFactory struct {
	clients cmap.ConcurrentMap[string, *kameleoonClient]
}

func newKameleoonClientFactory() *kameleoonClientFactory {
	return &kameleoonClientFactory{
		clients: cmap.New[*kameleoonClient](),
	}
}

func (cf *kameleoonClientFactory) Create(siteCode string, cfg *KameleoonClientConfig) (KameleoonClient, error) {
	return cf.createWithConfigSource(siteCode, func() (*KameleoonClientConfig, error) {
		return cfg, nil
	})
}

func (cf *kameleoonClientFactory) CreateFromFile(siteCode string, cfgPath string) (KameleoonClient, error) {
	return cf.createWithConfigSource(siteCode, func() (*KameleoonClientConfig, error) {
		return LoadConfig(cfgPath)
	})
}

func (cf *kameleoonClientFactory) createWithConfigSource(siteCode string,
	cfgSrc func() (*KameleoonClientConfig, error)) (KameleoonClient, error) {
	var err error
	return cf.clients.Upsert(siteCode, nil,
		func(exist bool, former, _ *kameleoonClient) *kameleoonClient {
			if former != nil {
				return former
			}
			cfg, cerr := cfgSrc()
			if cerr == nil {
				client, cerr := newClient(siteCode, cfg)
				if cerr == nil {
					return client
				}
			}
			err = cerr
			return nil
		}), err
}

func (cf *kameleoonClientFactory) Forget(siteCode string) {
	cf.clients.RemoveCb(siteCode, func(_ string, client *kameleoonClient, exists bool) bool {
		if client != nil {
			client.close()
		}
		return true
	})
}

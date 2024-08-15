package kameleoon

import (
	"github.com/Kameleoon/client-go/v3/logging"
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
	logging.Info("CALL: KameleoonClientFactory.Create(siteCode: %s, config: %s)", siteCode, cfg)
	client, err := cf.createWithConfigSource(siteCode, func() (*KameleoonClientConfig, error) {
		return cfg, nil
	})
	logging.Info("RETURN: KameleoonClientFactory.Create(siteCode: %s, config: %s) -> (client, error: %s)",
		siteCode, cfg, err)
	return client, err
}

func (cf *kameleoonClientFactory) CreateFromFile(siteCode string, cfgPath string) (KameleoonClient, error) {
	logging.Info("CALL: KameleoonClientFactory.CreateFromFile(siteCode: %s, configPath: %s)", siteCode, cfgPath)
	client, err := cf.createWithConfigSource(siteCode, func() (*KameleoonClientConfig, error) {
		return LoadConfig(cfgPath)
	})
	logging.Info(
		"RETURN: KameleoonClientFactory.CreateFromFile(siteCode: %s, configPath: %s) -> (client, error: %s)",
		siteCode, cfgPath, err)
	return client, err
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
	logging.Info("CALL: KameleoonClientFactory.Forget(siteCode: %s)", siteCode)
	cf.clients.RemoveCb(siteCode, func(_ string, client *kameleoonClient, exists bool) bool {
		if client != nil {
			client.close()
		}
		return true
	})
	logging.Info("RETURN: KameleoonClientFactory.Forget(siteCode: %s)", siteCode)
}

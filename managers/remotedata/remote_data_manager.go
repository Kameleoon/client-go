package remotedata

import (
	"encoding/json"
	"time"

	"github.com/Kameleoon/client-go/v3/logging"
	"github.com/Kameleoon/client-go/v3/network"
	"github.com/Kameleoon/client-go/v3/storage"
	"github.com/Kameleoon/client-go/v3/types"
)

type RemoteDataManager interface {
	GetData(key string, timeout ...time.Duration) ([]byte, error)
	GetVisitorData(
		visitorCode string,
		filter types.RemoteVisitorDataFilter,
		addData bool,
		isUniqueIdentifier bool,
		timeout ...time.Duration,
	) ([]types.Data, error)
}

type remoteDataManagerImpl struct {
	networkManager network.NetworkManager
	visitorManager storage.VisitorManager
	logger         logging.Logger
}

func NewRemoteDataManager(
	networkManager network.NetworkManager,
	visitorManager storage.VisitorManager,
	logger logging.Logger,
) RemoteDataManager {
	return &remoteDataManagerImpl{
		networkManager: networkManager,
		visitorManager: visitorManager,
		logger:         logger,
	}
}

func (rdm *remoteDataManagerImpl) GetData(key string, timeout ...time.Duration) ([]byte, error) {
	timeoutValue := time.Duration(-1)
	if len(timeout) > 0 {
		timeoutValue = timeout[0]
	}
	rdm.logger.Printf("Retrieve data from remote source (key '%s')", key)
	out, err := rdm.networkManager.GetRemoteData(key, timeoutValue)
	if err != nil {
		rdm.logger.Printf("Failed retrieve data from remote source: %v", err)
		return nil, err
	}
	return out, nil
}

func (rdm *remoteDataManagerImpl) GetVisitorData(
	visitorCode string,
	filter types.RemoteVisitorDataFilter,
	addData bool,
	isUniqueIdentifier bool,
	timeout ...time.Duration,
) ([]types.Data, error) {
	// TODO: Uncomment with the next major update
	//if err := utils.ValidateVisitorCode(visitorCode); err != nil {
	//	return nil, err
	//}
	timeoutValue := time.Duration(-1)
	if len(timeout) > 0 {
		timeoutValue = timeout[0]
	}
	out, err := rdm.networkManager.GetRemoteVisitorData(visitorCode, filter, isUniqueIdentifier, timeoutValue)
	if err != nil {
		return nil, err
	}
	var data remoteVisitorData
	if err = json.Unmarshal(out, &data); err != nil {
		return nil, err
	}
	data.MarkVisitorDataAsSent(rdm.visitorManager.CustomDataInfo())
	if addData {
		visitor := rdm.visitorManager.GetOrCreateVisitor(visitorCode)
		visitor.AddBaseData(rdm.logger, false, data.CollectDataToAdd()...)
	}
	return data.CollectVisitorDataToReturn(), nil
}

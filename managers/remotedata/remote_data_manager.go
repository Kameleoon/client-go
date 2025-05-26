package remotedata

import (
	"encoding/json"
	"time"

	"github.com/Kameleoon/client-go/v3/logging"
	"github.com/Kameleoon/client-go/v3/managers/data"
	"github.com/Kameleoon/client-go/v3/network"
	"github.com/Kameleoon/client-go/v3/storage"
	"github.com/Kameleoon/client-go/v3/types"
)

type RemoteDataManager interface {
	GetData(key string, timeout ...time.Duration) ([]byte, error)
	GetVisitorData(
		visitorCode string, filter types.RemoteVisitorDataFilter, addData bool, timeout ...time.Duration,
	) ([]types.Data, error)
}

type remoteDataManagerImpl struct {
	dataManager    data.DataManager
	networkManager network.NetworkManager
	visitorManager storage.VisitorManager
}

func NewRemoteDataManager(
	dataManager data.DataManager,
	networkManager network.NetworkManager,
	visitorManager storage.VisitorManager,
) RemoteDataManager {
	remoteDataManagerImpl := &remoteDataManagerImpl{
		dataManager:    dataManager,
		networkManager: networkManager,
		visitorManager: visitorManager,
	}
	logging.Debug(
		"CALL/RETURN: NewRemoteDataManager(dataManager, networkManager, visitorManager) -> (remoteDataManagerImpl)")
	return remoteDataManagerImpl
}

func (rdm *remoteDataManagerImpl) GetData(key string, timeout ...time.Duration) ([]byte, error) {
	logging.Debug("CALL: remoteDataManagerImpl.GetData(key: %s, timeout: %s)", key, timeout)
	timeoutValue := time.Duration(-1)
	if len(timeout) > 0 {
		timeoutValue = timeout[0]
	}
	out, err := rdm.networkManager.GetRemoteData(key, timeoutValue)
	if err != nil {
		logging.Error("Failed to fetch remote data for %s: %s", key, err)
		out = nil
	}
	logging.Debug(
		"RETURN: remoteDataManagerImpl.GetData(key: %s, timeout: %s) -> (remoteData: %s, error: %s)",
		key, timeout, out, err)
	return out, err
}

func (rdm *remoteDataManagerImpl) GetVisitorData(
	visitorCode string,
	filter types.RemoteVisitorDataFilter,
	addData bool,
	timeout ...time.Duration,
) ([]types.Data, error) {
	logging.Debug(
		"CALL: remoteDataManagerImpl.GetVisitorData(visitorCode: %s, filter: %s, addData: %s, timeout: %s)",
		visitorCode, filter, addData, timeout)
	// TODO: Uncomment with the next major update
	//if err := utils.ValidateVisitorCode(visitorCode); err != nil {
	//	return nil, err
	//}
	timeoutValue := time.Duration(-1)
	if len(timeout) > 0 {
		timeoutValue = timeout[0]
	}
	visitor := rdm.visitorManager.GetVisitor(visitorCode)
	var isUniqueIdentifier bool
	if visitor != nil {
		isUniqueIdentifier = visitor.IsUniqueIdentifier()
	}
	filter.ApplyDefaultValues()
	out, err := rdm.networkManager.GetRemoteVisitorData(visitorCode, filter, isUniqueIdentifier, timeoutValue)
	if err != nil {
		logging.Error("Failed to fetch remote visitor data for %s: %s", visitorCode, err)
		logging.Debug(
			"RETURN: remoteDataManagerImpl.GetVisitorData(visitorCode: %s, filter: %s, addData: %s, "+
				"isUniqueIdentifier: %s, timeout: %s) -> (remoteVisitorData: <nil>, error: %s)",
			visitorCode, filter, addData, isUniqueIdentifier, timeout, err)
		return nil, err
	}
	data := newRemoteVisitorData(filter)
	if err = json.Unmarshal(out, &data); err != nil {
		logging.Debug(
			"RETURN: remoteDataManagerImpl.GetVisitorData(visitorCode: %s, filter: %s, addData: %s, "+
				"isUniqueIdentifier: %s, timeout: %s) -> (remoteVisitorData: <nil>, error: %s)",
			visitorCode, filter, addData, isUniqueIdentifier, timeout, err)
		return nil, err
	}
	data.MarkVisitorDataAsSent(rdm.dataManager.DataFile().CustomDataInfo())
	if addData {
		// Cannot use `visitorManager.AddData` because it could use remote visitor data for mapping
		visitor = rdm.visitorManager.GetOrCreateVisitor(visitorCode)
		visitor.AddBaseData(false, data.CollectDataToAdd()...)
	}
	if (filter.VisitorCode == true) && (data.visitorCode != "") {
		// We apply visitor code from the latest visit fetched from Data API
		visitor = rdm.visitorManager.GetOrCreateVisitor(visitorCode)
		visitor.SetMappingIdentifier(&data.visitorCode)
	}
	visitorData := data.CollectVisitorDataToReturn()
	logging.Debug(
		"RETURN: remoteDataManagerImpl.GetVisitorData(visitorCode: %s, filter: %s, addData: %s, "+
			"isUniqueIdentifier: %s, timeout: %s) -> (remoteVisitorData: %s, error: <nil>)",
		visitorCode, filter, addData, isUniqueIdentifier, timeout, visitorData)
	return visitorData, nil
}

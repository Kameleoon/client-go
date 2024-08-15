package warehouse

import (
	"github.com/Kameleoon/client-go/v3/logging"
	"encoding/json"
	"time"

	"github.com/Kameleoon/client-go/v3/network"
	"github.com/Kameleoon/client-go/v3/storage"
	"github.com/Kameleoon/client-go/v3/types"
	"github.com/Kameleoon/client-go/v3/utils"
)

type WarehouseManager interface {
	GetVisitorWarehouseAudience(
		visitorCode string, warehouseKey string, customDataIndex int, timeout time.Duration) (*types.CustomData, error)
}

type warehouseResponse struct {
	WarehouseAudiences map[string]interface{} `json:"warehouseAudiences,omitempty"`
}

type warehouseManagerImpl struct {
	networkManager network.NetworkManager
	visitorManager storage.VisitorManager
}

func NewWarehouseManagerImpl(networkManager network.NetworkManager,
	visitorManager storage.VisitorManager) *warehouseManagerImpl {
	logging.Debug("CALL: NewWarehouseManagerImpl(networkManager, visitorManager)")
	warehouseManagerImpl := &warehouseManagerImpl{
		networkManager: networkManager,
		visitorManager: visitorManager,
	}
	logging.Debug("RETURN: NewWarehouseManagerImpl(networkManager, visitorManager)")
	return warehouseManagerImpl
}

func (wm *warehouseManagerImpl) GetVisitorWarehouseAudience(
	visitorCode string, warehouseKey string, customDataIndex int, timeout time.Duration) (*types.CustomData, error) {
	logging.Debug(
		"CALL: warehouseManagerImpl.GetVisitorWarehouseAudience(visitorCode: %s, warehouseKey: %s, "+
			"customDataIndex: %s, timeout: %s)", visitorCode, warehouseKey, customDataIndex, timeout)

	if err := utils.ValidateVisitorCode(visitorCode); err != nil {
		logging.Debug(
			"RETURN: warehouseManagerImpl.GetVisitorWarehouseAudience(visitorCode: %s, warehouseKey: %s, "+
				"customDataIndex: %s, timeout: %s) -> (customData: <nil>, error: %s)",
			visitorCode, warehouseKey, customDataIndex, timeout, err)
		return nil, err
	}

	remoteDataKey := remoteDataKey(visitorCode, warehouseKey)

	remoteData, err := wm.networkManager.GetRemoteData(remoteDataKey, timeout)
	if err != nil {
		logging.Debug(
			"RETURN: warehouseManagerImpl.GetVisitorWarehouseAudience(visitorCode: %s, warehouseKey: %s, "+
				"customDataIndex: %s, timeout: %s) -> (customData: <nil>, error: %s)",
			visitorCode, warehouseKey, customDataIndex, timeout, err)
		return nil, err
	}

	var warehouseResponse warehouseResponse
	err = json.Unmarshal(remoteData, &warehouseResponse)
	if err != nil {
		logging.Debug(
			"RETURN: warehouseManagerImpl.GetVisitorWarehouseAudience(visitorCode: %s, warehouseKey: %s, "+
				"customDataIndex: %s, timeout: %s) -> (customData: <nil>, error: %s)",
			visitorCode, warehouseKey, customDataIndex, timeout, err)
		return nil, err
	}

	values := make([]string, 0, len(warehouseResponse.WarehouseAudiences))
	for value := range warehouseResponse.WarehouseAudiences {
		values = append(values, value)
	}
	customData := types.NewCustomData(customDataIndex, values...)

	wm.visitorManager.AddData(visitorCode, customData)

	logging.Debug(
		"RETURN: warehouseManagerImpl.GetVisitorWarehouseAudience(visitorCode: %s, warehouseKey: %s, "+
			"customDataIndex: %s, timeout: %s) -> (customData: %s, error: <nil>)",
		visitorCode, warehouseKey, customDataIndex, timeout, customData)
	return customData, nil
}

func remoteDataKey(visitorCode string, warehouseKey string) string {
	if warehouseKey != "" {
		return warehouseKey
	} else {
		return visitorCode
	}
}

package warehouse

import (
	"encoding/json"
	"time"

	"github.com/Kameleoon/client-go/v3/logging"
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
	logger         logging.Logger
}

func NewWarehouseManagerImpl(networkManager network.NetworkManager,
	visitorManager storage.VisitorManager, logger logging.Logger) *warehouseManagerImpl {

	return &warehouseManagerImpl{
		networkManager: networkManager,
		visitorManager: visitorManager,
		logger:         logger,
	}
}

func (wm *warehouseManagerImpl) GetVisitorWarehouseAudience(
	visitorCode string, warehouseKey string, customDataIndex int, timeout time.Duration) (*types.CustomData, error) {

	if err := utils.ValidateVisitorCode(visitorCode); err != nil {
		return nil, err
	}

	remoteDataKey := remoteDataKey(visitorCode, warehouseKey)

	remoteData, err := wm.networkManager.GetRemoteData(remoteDataKey, timeout)
	if err != nil {
		wm.logger.Printf("Kameleoon SDK: Failed to fetch visitor warehouse audience: %s", err.Error())
		return nil, err
	}

	var warehouseResponse warehouseResponse
	err = json.Unmarshal(remoteData, &warehouseResponse)
	if err != nil {
		wm.logger.Printf("Kameleoon SDK: Failed to handle visitor warehouse audience response: %s", err.Error())
		return nil, err
	}

	values := make([]string, 0, len(warehouseResponse.WarehouseAudiences))
	for value := range warehouseResponse.WarehouseAudiences {
		values = append(values, value)
	}
	customData := types.NewCustomData(customDataIndex, values...)

	wm.visitorManager.GetOrCreateVisitor(visitorCode).AddData(wm.logger, customData)

	return customData, nil
}

func remoteDataKey(visitorCode string, warehouseKey string) string {
	if warehouseKey != "" {
		return warehouseKey
	} else {
		return visitorCode
	}
}

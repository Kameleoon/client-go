package data

import (
	"github.com/Kameleoon/client-go/v3/types"
)

type DataManager interface {
	DataFile() types.DataFile
	IsConsentRequired() bool

	SetDataFile(dataFile types.DataFile)
}

type DataManagerImpl struct {
	container *dataContainer
}

func NewDataManagerImpl(dataFile types.DataFile) *DataManagerImpl {
	return &DataManagerImpl{
		container: newDataContainer(dataFile),
	}
}

func (dm *DataManagerImpl) DataFile() types.DataFile {
	return dm.container.dataFile
}

func (dm *DataManagerImpl) IsConsentRequired() bool {
	return dm.container.isConsentRequired
}

func (dm *DataManagerImpl) SetDataFile(dataFile types.DataFile) {
	dm.container = newDataContainer(dataFile)
}

type dataContainer struct {
	dataFile          types.DataFile
	isConsentRequired bool
}

func newDataContainer(dataFile types.DataFile) *dataContainer {
	return &dataContainer{
		dataFile:          dataFile,
		isConsentRequired: dataFile.Settings().IsConsentRequired() && !dataFile.HasAnyTargetedDeliveryRule(),
	}
}

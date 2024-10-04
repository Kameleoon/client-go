package data

import (
	"github.com/Kameleoon/client-go/v3/types"
)

type DataManager interface {
	DataFile() types.DataFile
	IsVisitorCodeManaged() bool

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

func (dm *DataManagerImpl) IsVisitorCodeManaged() bool {
	return dm.container.isVisitorCodeManaged
}

func (dm *DataManagerImpl) SetDataFile(dataFile types.DataFile) {
	dm.container = newDataContainer(dataFile)
}

type dataContainer struct {
	dataFile             types.DataFile
	isVisitorCodeManaged bool
}

func newDataContainer(dataFile types.DataFile) *dataContainer {
	return &dataContainer{
		dataFile:             dataFile,
		isVisitorCodeManaged: dataFile.Settings().IsConsentRequired() && !dataFile.HasAnyTargetedDeliveryRule(),
	}
}

package data

import (
	"github.com/Kameleoon/client-go/v3/types"
)

type DataManager interface {
	DataFile() types.IDataFile
	IsVisitorCodeManaged() bool

	SetDataFile(dataFile types.IDataFile)
}

type DataManagerImpl struct {
	container *dataContainer
}

func NewDataManagerImpl(dataFile types.IDataFile) *DataManagerImpl {
	return &DataManagerImpl{
		container: newDataContainer(dataFile),
	}
}

func (dm *DataManagerImpl) DataFile() types.IDataFile {
	return dm.container.dataFile
}

func (dm *DataManagerImpl) IsVisitorCodeManaged() bool {
	return dm.container.isVisitorCodeManaged
}

func (dm *DataManagerImpl) SetDataFile(dataFile types.IDataFile) {
	dm.container = newDataContainer(dataFile)
}

type dataContainer struct {
	dataFile             types.IDataFile
	isVisitorCodeManaged bool
}

func newDataContainer(dataFile types.IDataFile) *dataContainer {
	return &dataContainer{
		dataFile:             dataFile,
		isVisitorCodeManaged: dataFile.Settings().IsConsentRequired() && !dataFile.HasAnyTargetedDeliveryRule(),
	}
}

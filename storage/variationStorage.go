package storage

type VariationStorage interface {
	GetVariationId(visitorCode string, experimentId int) (int, bool)
	IsVariationValid(visitorCode string, experimentId int, respoolTime int) (int, bool)
	UpdateVariation(visitorCode string, experimentId int, variationId int)
	GetMapSavedVariationId(visitorCode string) map[int]int
}

type VariationStorageImpl struct {
	storage map[string](map[uint32]*VisitorVariation)
}

func NewVariationStorage() *VariationStorageImpl {
	return &VariationStorageImpl{storage: make(map[string](map[uint32]*VisitorVariation))}
}

func (vs *VariationStorageImpl) GetVariationId(visitorCode string, experimentId int) (int, bool) {
	return vs.IsVariationValid(visitorCode, experimentId, 0)
}

func (vs *VariationStorageImpl) IsVariationValid(visitorCode string, experimentId int, respoolTime int) (int, bool) {
	if storageVisitor, exist := vs.storage[visitorCode]; exist {
		if variation, exist := storageVisitor[uint32(experimentId)]; exist {
			if variation.isValid(respoolTime) {
				return int(variation.VariationId), true
			}
		}
	}
	return 0, false
}

func (vs *VariationStorageImpl) UpdateVariation(visitorCode string, experimentId int, variationId int) {
	_, exist := vs.storage[visitorCode]
	if !exist {
		vs.storage[visitorCode] = make(map[uint32]*VisitorVariation)
	}
	vs.storage[visitorCode][uint32(experimentId)] = NewVisitorVariation(uint32(variationId))
}

func (vs *VariationStorageImpl) GetMapSavedVariationId(visitorCode string) map[int]int {
	if storageVisitor, exist := vs.storage[visitorCode]; exist {
		mapVariations := make(map[int]int)
		for key, value := range storageVisitor {
			mapVariations[int(key)] = int(value.VariationId)
		}
		return mapVariations
	}
	return nil
}

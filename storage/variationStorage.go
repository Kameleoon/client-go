package storage

type VariationStorage struct {
	Storage map[string](map[int]*VisitorVariation)
}

func NewVariationStorage() *VariationStorage {
	return &VariationStorage{Storage: make(map[string](map[int]*VisitorVariation))}
}

func (vs *VariationStorage) GetVariationId(visitorCode string, experimentId int) (int, bool) {
	return vs.IsVariationValid(visitorCode, experimentId, 0)
}

func (vs *VariationStorage) IsVariationValid(visitorCode string, experimentId int, respoolTime int) (int, bool) {
	if storageVisitor, exist := vs.Storage[visitorCode]; exist {
		if variation, exist := storageVisitor[experimentId]; exist {
			if variation.isValid(int64(respoolTime)) {
				return variation.VariationId, true
			}
		}
	}
	return 0, false
}

func (vs *VariationStorage) UpdateVariation(visitorCode string, experimentId int, variationId int) {
	_, exist := vs.Storage[visitorCode]
	if !exist {
		vs.Storage[visitorCode] = make(map[int]*VisitorVariation)
	}
	vs.Storage[visitorCode][experimentId] = NewVisitorVariation(variationId)
}

func (vs *VariationStorage) GetMapSavedVariationId(visitorCode string) map[int]int {
	if storageVisitor, exist := vs.Storage[visitorCode]; exist {
		mapVariations := make(map[int]int)
		for key, value := range storageVisitor {
			mapVariations[key] = value.VariationId
		}
		return mapVariations
	}
	return nil
}

package types

type KcsHeat struct {
	values map[int]map[int]float64
}

func NewKcsHeat(values map[int]map[int]float64) *KcsHeat {
	return &KcsHeat{values: values}
}

func (kh *KcsHeat) Values() map[int]map[int]float64 {
	return kh.values
}

func (kh *KcsHeat) DataType() DataType {
	return DataTypeKcsHeat
}

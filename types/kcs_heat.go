package types

import "fmt"

type KcsHeat struct {
	values map[int]map[int]float64
}

func NewKcsHeat(values map[int]map[int]float64) *KcsHeat {
	return &KcsHeat{values: values}
}

func (kh KcsHeat) String() string {
	return fmt.Sprintf("KcsHeat{values:%v}", kh.values)
}

func (kh *KcsHeat) Values() map[int]map[int]float64 {
	return kh.values
}

func (kh *KcsHeat) DataType() DataType {
	return DataTypeKcsHeat
}

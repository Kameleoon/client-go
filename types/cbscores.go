package types

import (
	"fmt"
	"sort"
)

type CBScores struct {
	// keys = experiment IDs / values = list of variation IDs ordered descending
	// by score (there may be several variation ids with same score)
	values map[int][]VarGroup
}

func NewCBScores(cbsMap map[int][]ScoredVarId) *CBScores {
	values := make(map[int][]VarGroup)
	for cbsKey, cbsValue := range cbsMap {
		values[cbsKey] = extractVarIds(cbsValue)
	}
	return &CBScores{values: values}
}

func extractVarIds(scores []ScoredVarId) []VarGroup {
	grouped := make(map[float64][]int)
	for _, score := range scores {
		grouped[score.Score] = append(grouped[score.Score], score.VariationId)
	}
	keys := make([]float64, 0, len(grouped))
	for key := range grouped {
		keys = append(keys, key)
	}
	sort.Float64s(keys)
	varIds := make([]VarGroup, len(keys))
	for i := 0; i < len(keys); i++ {
		varIds[i] = newVarGroup(grouped[keys[len(keys)-i-1]])
	}
	return varIds
}

func (vg *CBScores) Values() map[int][]VarGroup {
	return vg.values
}

func (*CBScores) DataType() DataType {
	return DataTypeCBScores
}

func (vg CBScores) String() string {
	return fmt.Sprintf("CBScores{values:%v}", vg.values)
}

type ScoredVarId struct {
	VariationId int
	Score       float64
}

type VarGroup struct {
	ids []int
}

func newVarGroup(ids []int) VarGroup {
	sort.Ints(ids)
	return VarGroup{ids: ids}
}

func (vg VarGroup) Ids() []int {
	return vg.ids
}

func (vg VarGroup) String() string {
	return fmt.Sprintf("VarGroup{ids:%v}", vg.ids)
}

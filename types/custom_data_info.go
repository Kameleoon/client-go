package types

import (
	"encoding/json"
	"fmt"
)

const scopeVisitor = "VISITOR"
const undefinedIndex = -1

type CustomDataInfo struct {
	localOnly              map[int]struct{}
	visitorScope           map[int]struct{}
	mappingIdentifierIndex int
}

func NewCustomDataInfo() *CustomDataInfo {
	return &CustomDataInfo{
		localOnly:              make(map[int]struct{}),
		visitorScope:           make(map[int]struct{}),
		mappingIdentifierIndex: undefinedIndex,
	}
}

func (cdi *CustomDataInfo) UnmarshalJSON(data []byte) error {
	var cds = []struct {
		Index               int    `json:"index"`
		LocalOnly           bool   `json:"localOnly"`
		Scope               string `json:"scope"`
		IsMappingIdentifier bool   `json:"isMappingIdentifier"`
	}{}
	if err := json.Unmarshal(data, &cds); err != nil {
		return err
	}
	cdi.localOnly = make(map[int]struct{})
	cdi.visitorScope = make(map[int]struct{})
	cdi.mappingIdentifierIndex = undefinedIndex
	for _, cd := range cds {
		if cd.LocalOnly {
			cdi.localOnly[cd.Index] = struct{}{}
		}
		if cd.Scope == scopeVisitor {
			cdi.visitorScope[cd.Index] = struct{}{}
		}
		if cd.IsMappingIdentifier {
			if cdi.mappingIdentifierIndex != undefinedIndex {
				fmt.Printf("Kameleoon SDK: More than one mapping identifier is set. " +
					"Undefined behavior may occur on cross-device reconciliation.")
			}
			cdi.mappingIdentifierIndex = cd.Index
		}
	}
	return nil
}

func (cdi *CustomDataInfo) MappingIdentifierIndex() int {
	return cdi.mappingIdentifierIndex
}

func (cdi *CustomDataInfo) IsLocalOnly(index int) bool {
	_, ok := cdi.localOnly[index]
	return ok
}

func (cdi *CustomDataInfo) IsVisitorScope(index int) bool {
	_, ok := cdi.visitorScope[index]
	return ok
}

func (cdi *CustomDataInfo) IsMappingIdentifier(index int) bool {
	return index == cdi.mappingIdentifierIndex
}

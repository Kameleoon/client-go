package types

import (
	"encoding/json"

	"github.com/Kameleoon/client-go/v3/logging"
)

const scopeVisitor = "VISITOR"
const undefinedIndex = -1

type CustomDataInfo struct {
	localOnly              map[int]struct{}
	visitorScope           map[int]struct{}
	customDataIndexById    map[int]int
	customDataIndexByName  map[string]int
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
		Id                  int    `json:"id"`
		Index               int    `json:"index"`
		Name                string `json:"name"`
		LocalOnly           bool   `json:"localOnly"`
		Scope               string `json:"scope"`
		IsMappingIdentifier bool   `json:"isMappingIdentifier"`
	}{}
	if err := json.Unmarshal(data, &cds); err != nil {
		return err
	}
	cdi.localOnly = make(map[int]struct{})
	cdi.visitorScope = make(map[int]struct{})
	cdi.customDataIndexById = make(map[int]int)
	cdi.customDataIndexByName = make(map[string]int)
	cdi.mappingIdentifierIndex = undefinedIndex
	for _, cd := range cds {
		if cd.LocalOnly {
			cdi.localOnly[cd.Index] = struct{}{}
		}
		if cd.Scope == scopeVisitor {
			cdi.visitorScope[cd.Index] = struct{}{}
		}
		cdi.customDataIndexById[cd.Id] = cd.Index
		if cd.Name != "" {
			cdi.customDataIndexByName[cd.Name] = cd.Index
		}
		if cd.IsMappingIdentifier {
			if cdi.mappingIdentifierIndex != undefinedIndex {
				logging.Warning("More than one mapping identifier is set. " +
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

func (cdi *CustomDataInfo) GetCustomDataIndexById(customDataId int) (index int, exists bool) {
	index, exists = cdi.customDataIndexById[customDataId]
	return
}

func (cdi *CustomDataInfo) GetCustomDataIndexByName(customDataName string) (index int, exists bool) {
	index, exists = cdi.customDataIndexByName[customDataName]
	return
}

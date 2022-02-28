package kameleoon

import (
	"github.com/Kameleoon/client-go/targeting"
	"github.com/Kameleoon/client-go/types"
)

type SegmentQL struct {
	ID int `json:"id,string"`
	types.Segment
}

type VariationQL struct {
	ID int `json:"id,string"`
	types.Variation
}

func GetExperimentsGraphQL(siteCode string) string {
	return `{
		"operationName": "getExperiments",
		"query": "query getExperiments($first: Int, $after: String, $filter: FilteringExpression, $sort: [SortingParameter!]) { experiments(first: $first, after: $after, filter: $filter, sort: $sort) { edges { node { id name type site { id code isKameleoonEnabled } status variations { id customJson } deviations { variationId value } respoolTime {variationId value } segment { id name conditionsData { firstLevelOrOperators firstLevel { orOperators conditions { targetingType isInclude ... on CustomDataTargetingCondition { customDataIndex value valueMatchType } } } } } __typename } __typename } pageInfo { endCursor hasNextPage __typename } totalCount __typename } }",
		"variables": {
			"filter": {
				"and": [{
					"condition": {
						"field": "status",
						"operator": "IN",
						"parameters": ["ACTIVE", "DEVIATED", "USED_AS_PERSONALIZATION"]
					}
				}, {
					"condition": {
						"field": "type",
						"operator": "IN",
						"parameters": ["SERVER_SIDE", "HYBRID"]
					}
				}, {
					"condition": {
						"field": "siteCode",
						"operator": "IN",
						"parameters": ["` + siteCode + `"]
					}
				}]
			},
			"sort": [{
				"field": "id",
				"direction": "ASC"
			}]
		}
	}`
}

type ExperimentQL struct {
	ID               int             `json:"id,string"`
	Variations       []VariationQL   `json:"variations"`
	Deviations       []DeviationsQL  `json:"deviations"`
	RespoolTime      []RespoolTimeQL `json:"respoolTime"`
	TargetingSegment SegmentQL       `json:"segment"`
	types.Experiment
}

type RespoolTimeQL struct {
	VariationId string  `json:"variationId"`
	Value       float64 `json:"value"`
}

type DeviationsQL struct {
	VariationId string  `json:"variationId"`
	Value       float64 `json:"value"`
}

type ExperimentDataGraphQL struct {
	Data ExperimentGraphQL `json:"data"`
}
type ExperimentGraphQL struct {
	Experiments EdgeExperimentGraphQL `json:"experiments"`
}
type EdgeExperimentGraphQL struct {
	Edge []NodeExperimentGraphQL `json:"edges"`
}

type NodeExperimentGraphQL struct {
	Node ExperimentQL `json:"node"`
}

func (expNodeQL *NodeExperimentGraphQL) Transform() types.Experiment {
	expNode := expNodeQL.Node
	exp := expNode.Experiment
	exp.ID = expNode.ID
	//transform variations
	for _, variation := range expNode.Variations {
		exp.VariationsID = append(exp.VariationsID, variation.ID)
		variation.Variation.ID = variation.ID
		exp.Variations = append(exp.Variations, variation.Variation)
	}
	// transform segment if segment != null
	if expNode.TargetingSegment.ID != 0 {
		exp.TargetingSegmentID = expNode.TargetingSegment.ID
		expNode.TargetingSegment.Segment.ID = expNode.TargetingSegment.ID
		exp.TargetingSegment = targeting.NewSegment(&expNode.TargetingSegment.Segment)
	}
	//transform deviations
	exp.Deviations = map[string]float64{}
	for _, deviation := range expNode.Deviations {
		exp.Deviations[deviation.VariationId] = deviation.Value
	}
	//transform respoolTime
	exp.RespoolTime = map[string]float64{}
	for _, respoolTime := range expNode.RespoolTime {
		exp.RespoolTime[respoolTime.VariationId] = respoolTime.Value
	}
	return exp
}

func GetFeatureFlagsGraphQL(siteCode string, environment string) string {
	return `{
		"operationName": "getFeatureFlags",
		"query": "query getFeatureFlags($first: Int, $after: String, $filter: FilteringExpression, $sort: [SortingParameter!]) { featureFlags(first: $first, after: $after, filter: $filter, sort: $sort) { edges { node { id name site { id code isKameleoonEnabled } bypassDeviation status variations { id customJson } respoolTime { variationId value } expositionRate identificationKey featureFlagSdkLanguageType featureStatus schedules { dateStart dateEnd } segment { id name conditionsData { firstLevelOrOperators firstLevel { orOperators conditions { targetingType isInclude ... on CustomDataTargetingCondition { customDataIndex value valueMatchType } } } } } __typename } __typename } pageInfo { endCursor hasNextPage __typename } totalCount __typename } }",
		"variables": {
			"filter": {
				"and": [{
					"condition": {
						"field": "featureStatus",
						"operator": "IN",
						"parameters": ["ACTIVATED", "SCHEDULED", "DEACTIVATED"]
					}
				}, 
				{
					"condition": {
						"field": "siteCode",
						"operator": "IN",
						"parameters": ["` + siteCode + `"]
					}
				},
				{
					"condition": {
						"field": "environment.key",
						"operator": "IN",
						"parameters": ["` + environment + `"]
					}
				}]
			},
			"sort": [{
				"field": "id",
				"direction": "ASC"
			}]
		}
	}`
}

// GraphQL Helpers
type FeatureFlagQL struct {
	ID               int             `json:"id,string"`
	Variations       []VariationQL   `json:"variations"`
	RespoolTime      []RespoolTimeQL `json:"respoolTime"`
	TargetingSegment SegmentQL       `json:"segment"`
	types.FeatureFlag
}

type FeatureFlagDataGraphQL struct {
	Data FeatureFlagGraphQL `json:"data"`
}
type FeatureFlagGraphQL struct {
	FeatureFlags EdgeFeatureFlagGraphQL `json:"featureFlags"`
}
type EdgeFeatureFlagGraphQL struct {
	Edge []NodeFeatureFlagGraphQL `json:"edges"`
}

type NodeFeatureFlagGraphQL struct {
	Node FeatureFlagQL `json:"node"`
}

func (ffNodeQL *NodeFeatureFlagGraphQL) Transform() types.FeatureFlag {
	ffNode := ffNodeQL.Node
	ff := ffNodeQL.Node.FeatureFlag
	ff.ID = ffNode.ID
	// transform variations
	for _, variation := range ffNode.Variations {
		ff.VariationsID = append(ff.VariationsID, variation.ID)
		variation.Variation.ID = variation.ID
		ff.Variations = append(ff.Variations, variation.Variation)
	}
	// transform segment if segment != null
	if ffNode.TargetingSegment.ID != 0 {
		ff.TargetingSegmentID = ffNode.TargetingSegment.ID
		ffNode.TargetingSegment.Segment.ID = ffNode.TargetingSegment.ID
		ff.TargetingSegment = targeting.NewSegment(&ffNode.TargetingSegment.Segment)
	}
	// transform respoolTime
	ff.RespoolTime = map[string]float64{}
	for _, respoolTime := range ffNode.RespoolTime {
		ff.RespoolTime[respoolTime.VariationId] = respoolTime.Value
	}
	return ff
}

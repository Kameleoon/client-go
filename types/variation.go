package types

import "github.com/segmentio/encoding/json"

type Variation struct {
	ID                    int                  `json:"id"`
	SiteID                int                  `json:"siteId"`
	Name                  string               `json:"name"`
	JsCode                json.RawMessage      `json:"jsCode"`
	CssCode               json.RawMessage      `json:"cssCode"`
	IsJsCodeAfterDomReady bool                 `json:"isJsCodeAfterDomReady"`
	WidgetTemplateInput   json.RawMessage      `json:"widgetTemplateInput"`
	RedirectionStrings    string               `json:"redirectionStrings"`
	Redirection           VariationRedirection `json:"redirection"`
	ExperimentID          int                  `json:"experimentId"`
	CustomJson            json.RawMessage      `json:"customJson"`
}

type VariationRedirectionType string

const (
	VariationRedirectionGlobal    VariationRedirectionType = "GLOBAL_REDIRECTION"
	VariationRedirectionParameter VariationRedirectionType = "PARAMETER_REDIRECTION"
)

type VariationRedirection struct {
	Type                   VariationRedirectionType `json:"type"`
	Url                    string                   `json:"url"`
	Parameters             string                   `json:"parameters"`
	IncludeQueryParameters bool                     `json:"includeQueryParameters"`
}

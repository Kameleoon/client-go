package types

import "github.com/segmentio/encoding/json"

type Variation struct {
	ID         int             `json:"id,string"`
	CustomJson json.RawMessage `json:"customJson"`
}

package types

import (
	"strings"

	"github.com/Kameleoon/client-go/v2/utils"
)

type PageView struct {
	URL       string
	Title     string
	Referrers []int
}

func (v PageView) QueryEncode() string {
	var b strings.Builder
	b.WriteString("eventType=page&")
	b.WriteString(EncodeURIComponent("href", v.URL))
	b.WriteString("&title=")
	b.WriteString(v.Title)
	if len(v.Referrers) > 0 {
		b.WriteString("&referrersIndices=[")
		b.WriteString(utils.ArrayToString(v.Referrers, ","))
		b.WriteByte(']')
	}
	b.WriteString("&nonce=")
	b.WriteString(GetNonce())

	return b.String()
}

func (v PageView) DataType() DataType {
	return DataTypePageView
}

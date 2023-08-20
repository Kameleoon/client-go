package types

import (
	"fmt"
	"strings"

	"github.com/Kameleoon/client-go/v2/network"
)

const pageViewEventType = "page"

type PageView struct {
	URL       string
	Title     string
	Referrers []int
}

func (v PageView) QueryEncode() string {
	qb := network.NewQueryBuilder()
	qb.Append(network.QPEventType, pageViewEventType)
	qb.Append(network.QPHref, v.URL)
	qb.Append(network.QPTitle, v.Title)
	if len(v.Referrers) > 0 {
		qb.Append(network.QPReferrersIndices, v.encodeReferrers())
	}
	qb.Append(network.QPNonce, network.GetNonce())
	return qb.String()
}

func (v PageView) encodeReferrers() string {
	return strings.ReplaceAll(fmt.Sprint(v.Referrers), " ", ",")
}

func (v PageView) DataType() DataType {
	return DataTypePageView
}

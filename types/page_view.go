package types

import (
	"fmt"
	"strings"

	"github.com/Kameleoon/client-go/v3/utils"
)

const pageViewEventType = "page"

type PageView struct {
	duplicationSafeSendableBase
	url       string
	title     string
	referrers []int
}

func NewPageView(url string, referrers ...int) *PageView {
	return NewPageViewWithTitle(url, "", referrers...)
}
func NewPageViewWithTitle(url string, title string, referrers ...int) *PageView {
	pv := &PageView{
		url:       url,
		title:     title,
		referrers: referrers,
	}
	pv.initSendale()
	return pv
}

func (pv *PageView) dataRestriction() {}

func (pv *PageView) URL() string {
	return pv.url
}

func (pv *PageView) Title() string {
	return pv.title
}

func (pv *PageView) Referrers() []int {
	return pv.referrers
}

func (pv *PageView) QueryEncode() string {
	if len(pv.url) == 0 {
		return ""
	}
	nonce := pv.Nonce()
	if len(nonce) == 0 {
		return ""
	}
	qb := utils.NewQueryBuilder()
	qb.Append(utils.QPEventType, pageViewEventType)
	qb.Append(utils.QPHref, pv.url)
	qb.Append(utils.QPTitle, pv.title)
	if len(pv.referrers) > 0 {
		qb.Append(utils.QPReferrersIndices, pv.encodeReferrers())
	}
	qb.Append(utils.QPNonce, nonce)
	return qb.String()
}
func (pv *PageView) encodeReferrers() string {
	return strings.ReplaceAll(fmt.Sprint(pv.referrers), " ", ",")
}

func (pv *PageView) DataType() DataType {
	return DataTypePageView
}

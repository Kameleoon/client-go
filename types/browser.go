package types

import (
	"fmt"

	"github.com/Kameleoon/client-go/v2/network"
	"github.com/Kameleoon/client-go/v2/utils"
)

type BrowserType int

const (
	BrowserTypeChrome  BrowserType = 0
	BrowserTypeIE      BrowserType = 1
	BrowserTypeFirefox BrowserType = 2
	BrowserTypeSafari  BrowserType = 3
	BrowserTypeOpera   BrowserType = 4
	BrowserTypeOther   BrowserType = 5
)

const browserEventType = "staticData"

type Browser struct {
	Type    BrowserType
	Version float32
}

func (b Browser) QueryEncode() string {
	qb := network.NewQueryBuilder()
	qb.Append(network.QPEventType, browserEventType)
	qb.Append(network.QPBrowserIndex, utils.WritePositiveInt(int(b.Type)))
	if b.Version != 0 {
		qb.Append(network.QPBrowserVersion, fmt.Sprintf("%f", b.Version))
	}
	qb.Append(network.QPNonce, network.GetNonce())
	return qb.String()
}

func (b Browser) DataType() DataType {
	return DataTypeBrowser
}

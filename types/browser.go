package types

import (
	"fmt"
	"strings"

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

type Browser struct {
	Type    BrowserType
	Version float32
}

func (b Browser) QueryEncode() string {
	var sb strings.Builder
	sb.WriteString("eventType=staticData&browserIndex=")
	sb.WriteString(utils.WritePositiveInt(int(b.Type)))
	if b.Version != 0 {
		sb.WriteString("&browserVersion=")
		sb.WriteString(fmt.Sprintf("%g", b.Version))
	}
	sb.WriteString("&nonce=")
	sb.WriteString(GetNonce())
	return sb.String()
}

func (b Browser) DataType() DataType {
	return DataTypeBrowser
}

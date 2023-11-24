package cookie

import (
	"strings"
	"time"

	"github.com/Kameleoon/client-go/v3/utils"
	"github.com/valyala/fasthttp"
)

const (
	// _js_ - support to 22.08.2024, then remove
	cookieKeyJs       = "_js_"
	visitorCodeCookie = "kameleoonVisitorCode"
	cookieTTL         = 380 * 24 * time.Hour
)

func getVisitorCodeFromResponseCookie(response *fasthttp.Response) string {
	const token = visitorCodeCookie + "="
	ckBin := response.Header.PeekCookie(visitorCodeCookie)
	if ckBin == nil {
		return ""
	}
	ck := string(ckBin)
	start := strings.Index(ck, token)
	if start == -1 {
		return ""
	}
	start += len(token)
	end := start
	for (end < len(ck)) && (ck[end] != ';') {
		end++
	}
	return ck[start:end]
}

type CookieManager interface {
	IsConsentRequired() bool
	SetConsentRequired(value bool)

	GetOrAdd(request *fasthttp.Request, response *fasthttp.Response, defaultVisitorCode ...string) (string, error)

	Update(visitorCode string, legalConsent bool, response *fasthttp.Response)
}

type CookieManagerImpl struct {
	consentRequired bool
	topLevelDomain  string
}

func NewCookieManagerImpl(topLevelDomain string) *CookieManagerImpl {
	return &CookieManagerImpl{
		topLevelDomain: topLevelDomain,
	}
}

func (cm *CookieManagerImpl) IsConsentRequired() bool {
	return cm.consentRequired
}
func (cm *CookieManagerImpl) SetConsentRequired(value bool) {
	cm.consentRequired = value
}

func (cm *CookieManagerImpl) GetOrAdd(request *fasthttp.Request, response *fasthttp.Response,
	defaultVisitorCode ...string) (string, error) {
	var vc string

	if vc = getVisitorCodeFromResponseCookie(response); len(vc) > 0 {
		return vc, nil
	}

	if binaryVC := request.Header.Cookie(visitorCodeCookie); binaryVC != nil {
		vc = string(binaryVC)
		vc = strings.Replace(vc, cookieKeyJs, "", 1)
	} else {
		if len(defaultVisitorCode) > 0 {
			vc = defaultVisitorCode[0]
		} else {
			vc = utils.GenerateVisitorCode()
			if !cm.consentRequired {
				cm.add(vc, response)
			}
			return vc, nil
		}
	}

	if err := utils.ValidateVisitorCode(vc); err != nil {
		return "", err
	}
	cm.add(vc, response)
	return vc, nil
}

func (cm *CookieManagerImpl) add(visitorCode string, response *fasthttp.Response) {
	ck := fasthttp.AcquireCookie()
	defer fasthttp.ReleaseCookie(ck)
	ck.SetKey(visitorCodeCookie)
	ck.SetValue(visitorCode)
	ck.SetExpire(time.Now().Add(cookieTTL))
	ck.SetHTTPOnly(false)
	ck.SetPath("/")
	ck.SetDomain(cm.topLevelDomain)
	response.Header.SetCookie(ck)
}

func (cm *CookieManagerImpl) remove(response *fasthttp.Response) {
	if cm.consentRequired {
		response.Header.DelCookie(visitorCodeCookie)
	}
}

func (cm *CookieManagerImpl) Update(visitorCode string, consent bool, response *fasthttp.Response) {
	if consent {
		cm.add(visitorCode, response)
	} else {
		cm.remove(response)
	}
}

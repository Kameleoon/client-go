package cookie

import (
	"github.com/Kameleoon/client-go/v3/logging"
	"strings"
	"time"

	"github.com/Kameleoon/client-go/v3/managers/data"
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
	logging.Debug("CALL: getVisitorCodeFromResponseCookie(response)")
	const token = visitorCodeCookie + "="
	ckBin := response.Header.PeekCookie(visitorCodeCookie)
	var visitorCode string
	if ckBin == nil {
		visitorCode = ""
	} else {
		ck := string(ckBin)
		start := strings.Index(ck, token)
		if start == -1 {
			visitorCode = ""
		} else {
			start += len(token)
			end := start
			for (end < len(ck)) && (ck[end] != ';') {
				end++
			}
			visitorCode = ck[start:end]
		}
	}
	logging.Debug("RETURN: getVisitorCodeFromResponseCookie(response) -> (visitorCode: %s)",
		visitorCode)
	return visitorCode
}

type CookieManager interface {
	GetOrAdd(request *fasthttp.Request, response *fasthttp.Response, defaultVisitorCode ...string) (string, error)

	Update(visitorCode string, legalConsent bool, response *fasthttp.Response)
}

type CookieManagerImpl struct {
	dataManager    data.DataManager
	topLevelDomain string
}

func NewCookieManagerImpl(dataManager data.DataManager, topLevelDomain string) *CookieManagerImpl {
	cookieManagerImpl := &CookieManagerImpl{
		dataManager:    dataManager,
		topLevelDomain: topLevelDomain,
	}
	logging.Debug("CALL/RETURN: NewCookieManagerImpl(topLevelDomain: %s) -> (cookieManagerImpl)")
	return cookieManagerImpl
}

func (cm *CookieManagerImpl) GetOrAdd(request *fasthttp.Request, response *fasthttp.Response,
	defaultVisitorCode ...string) (string, error) {
	logging.Debug("CALL: CookieManagerImpl.GetOrAdd(request, response, defaultVisitorCode: %s)",
		defaultVisitorCode)
	var vc string

	if vc = getVisitorCodeFromResponseCookie(response); len(vc) > 0 {
		logging.Debug("Read visitor code %s from response %s", vc, response)
		logging.Debug("RETURN: CookieManagerImpl.GetOrAdd(request, response, defaultVisitorCode: %s) -> "+
			"(visitorCode: %s, error: <nil>)", defaultVisitorCode, vc)
		return vc, nil
	}

	if binaryVC := request.Header.Cookie(visitorCodeCookie); binaryVC != nil {
		vc = string(binaryVC)
		vc = strings.Replace(vc, cookieKeyJs, "", 1)
		logging.Debug("Read visitor code %s from request %s", vc, request)
	} else {
		if len(defaultVisitorCode) > 0 {
			vc = defaultVisitorCode[0]
			logging.Debug("Used default visitor code %s", vc)
		} else {
			vc = utils.GenerateVisitorCode()
			logging.Debug("Generated new visitor code %s", vc)
			if !cm.dataManager.IsConsentRequired() {
				cm.add(vc, response)
			}
			logging.Debug("RETURN: CookieManagerImpl.GetOrAdd(request, response, defaultVisitorCode: %s) -> "+
				"(visitorCode: %s, error: <nil>)", defaultVisitorCode, vc)
			return vc, nil
		}
	}

	err := utils.ValidateVisitorCode(vc)
	if err != nil {
		vc = ""
	} else {
		cm.add(vc, response)
	}
	logging.Debug("RETURN: CookieManagerImpl.GetOrAdd(request, response, defaultVisitorCode: %s) -> "+
		"(visitorCode: %s, error: %s)", defaultVisitorCode, vc, err)
	return vc, err
}

func (cm *CookieManagerImpl) add(visitorCode string, response *fasthttp.Response) {
	logging.Debug("CALL: CookieManagerImpl.add(visitorCode: %s, response)", visitorCode)
	ck := fasthttp.AcquireCookie()
	defer fasthttp.ReleaseCookie(ck)
	ck.SetKey(visitorCodeCookie)
	ck.SetValue(visitorCode)
	ck.SetExpire(time.Now().Add(cookieTTL))
	ck.SetHTTPOnly(false)
	ck.SetPath("/")
	ck.SetDomain(cm.topLevelDomain)
	response.Header.SetCookie(ck)
	logging.Debug("For %s was added cookies: %s", visitorCode, ck)
	logging.Debug("RETURN: CookieManagerImpl.add(visitorCode: %s, response)", visitorCode)
}

func (cm *CookieManagerImpl) remove(response *fasthttp.Response) {
	logging.Debug("CALL: CookieManagerImpl.remove(response)")
	if cm.dataManager.IsConsentRequired() {
		response.Header.DelCookie(visitorCodeCookie)
	}
	logging.Debug("RETURN: CookieManagerImpl.remove(response)")
}

func (cm *CookieManagerImpl) Update(visitorCode string, consent bool, response *fasthttp.Response) {
	logging.Debug("CALL: CookieManagerImpl.Update(visitorCode: %s, consent: %s, response)", visitorCode, consent)
	if consent {
		cm.add(visitorCode, response)
	} else {
		cm.remove(response)
	}
	logging.Debug("RETURN: CookieManagerImpl.Update(visitorCode: %s, consent: %s, response)",
		visitorCode, consent)
}

package cookie

import (
	"encoding/json"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/Kameleoon/client-go/v3/errs"
	"github.com/Kameleoon/client-go/v3/logging"
	"github.com/Kameleoon/client-go/v3/storage"
	"github.com/Kameleoon/client-go/v3/types"

	"github.com/Kameleoon/client-go/v3/managers/data"
	"github.com/Kameleoon/client-go/v3/utils"
	"github.com/valyala/fasthttp"
)

const (
	visitorCodeCookie      = "kameleoonVisitorCode"
	simulationFFDataCookie = "kameleoonSimulationFFData"
	cookieTTL              = 380 * 24 * time.Hour
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
	visitorManager storage.VisitorManager
	topLevelDomain string
}

func NewCookieManagerImpl(
	dataManager data.DataManager, visitorManager storage.VisitorManager, topLevelDomain string,
) *CookieManagerImpl {
	cookieManagerImpl := &CookieManagerImpl{
		dataManager:    dataManager,
		visitorManager: visitorManager,
		topLevelDomain: topLevelDomain,
	}
	logging.Debug(
		"CALL/RETURN: NewCookieManagerImpl(dataManager, visitorManager, topLevelDomain: %s) -> (cookieManagerImpl)",
		topLevelDomain,
	)
	return cookieManagerImpl
}

func (cm *CookieManagerImpl) Update(visitorCode string, consent bool, response *fasthttp.Response) {
	logging.Debug("CALL: CookieManagerImpl.Update(visitorCode: %s, consent: %s, response)", visitorCode, consent)
	if consent {
		cm.add(visitorCode, response)
	}
	logging.Debug("RETURN: CookieManagerImpl.Update(visitorCode: %s, consent: %s, response)", visitorCode, consent)
}

func (cm *CookieManagerImpl) GetOrAdd(request *fasthttp.Request, response *fasthttp.Response,
	defaultVisitorCode ...string) (string, error) {
	logging.Debug("CALL: CookieManagerImpl.GetOrAdd(request, response, defaultVisitorCode: %s)",
		defaultVisitorCode)
	vc, err := cm.getOrAddVisitorCode(request, response, defaultVisitorCode...)
	if err == nil {
		cm.processSimulatedVariations(request, vc)
	}
	logging.Debug("RETURN: CookieManagerImpl.GetOrAdd(request, response, defaultVisitorCode: %s) -> "+
		"(visitorCode: %s, error: %s)", defaultVisitorCode, vc, err)
	return vc, err
}

func (cm *CookieManagerImpl) getOrAddVisitorCode(
	request *fasthttp.Request, response *fasthttp.Response, defaultVisitorCode ...string,
) (string, error) {
	var vc string

	if vc = getVisitorCodeFromResponseCookie(response); len(vc) > 0 {
		logging.Debug("Read visitor code %s from response %s", vc, response)
		return vc, nil
	}

	if binaryVC := request.Header.Cookie(visitorCodeCookie); binaryVC != nil {
		vc = string(binaryVC)
		logging.Debug("Read visitor code %s from request %s", vc, request)
	} else {
		if len(defaultVisitorCode) > 0 {
			vc = defaultVisitorCode[0]
			logging.Debug("Used default visitor code %s", vc)
		} else {
			vc = utils.GenerateVisitorCode()
			logging.Debug("Generated new visitor code %s", vc)
			if !cm.dataManager.IsVisitorCodeManaged() {
				cm.add(vc, response)
			}
			return vc, nil
		}
	}

	err := utils.ValidateVisitorCode(vc)
	if err != nil {
		vc = ""
	} else if !cm.dataManager.IsVisitorCodeManaged() {
		cm.add(vc, response)
	}
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

func (cm *CookieManagerImpl) processSimulatedVariations(request *fasthttp.Request, visitorCode string) {
	svms, err := readSimulatedVariationsJson(request)
	if err == nil {
		var svs []*types.ForcedFeatureVariation
		svs, err = cm.parseSimulatedVariations(svms)
		if err == nil {
			visitor := cm.visitorManager.GetOrCreateVisitor(visitorCode)
			visitor.UpdateSimulatedVariations(svs)
			return
		}
	}
	logging.Error("Failed to process simulated variations cookie: %s", err)
}

func readSimulatedVariationsJson(request *fasthttp.Request) (svms map[string]simulatedVariationModel, err error) {
	if binarySV := request.Header.Cookie(simulationFFDataCookie); binarySV != nil {
		var unescapedSV string
		if unescapedSV, err = url.QueryUnescape(string(binarySV)); err != nil {
			return
		}
		err = json.Unmarshal([]byte(unescapedSV), &svms)
	}
	return
}

func (cm *CookieManagerImpl) parseSimulatedVariations(
	svms map[string]simulatedVariationModel,
) ([]*types.ForcedFeatureVariation, error) {
	dataFile := cm.dataManager.DataFile()
	svs := make([]*types.ForcedFeatureVariation, 0, len(svms))
	for ffKey, svm := range svms {
		sv, err := simulatedVariationFromDataFile(dataFile, ffKey, svm)
		if err != nil {
			return nil, err
		}
		svs = append(svs, sv)
	}
	return svs, nil
}

func simulatedVariationFromDataFile(
	dataFile types.DataFile, ffKey string, svm simulatedVariationModel,
) (*types.ForcedFeatureVariation, error) {
	ff, ffExists := dataFile.GetFeatureFlags()[ffKey]
	if !ffExists {
		return nil, errs.NewFeatureNotFound(ffKey)
	}
	if svm.ExperimentId == nil {
		return nil, errors.New("simulation cookie misses required field 'expId'")
	}
	experimentId := *svm.ExperimentId
	if experimentId == 0 {
		return types.NewForcedFeatureVariation(ffKey, nil, nil, true), nil
	}
	if svm.VariationId == nil {
		return nil, errors.New("simulation cookie misses field 'varId' required since the 'expId' is not zero")
	}
	variationId := *svm.VariationId
	for _, rule := range ff.GetRules() {
		if rule.GetRuleBase().ExperimentId != experimentId {
			continue
		}
		varsByExp := rule.GetRuleBase().VariationsByExposition
		for i := range varsByExp {
			varId := varsByExp[i].VariationID
			if (varId != nil) && (*varId == variationId) {
				return types.NewForcedFeatureVariation(ffKey, rule, &varsByExp[i], true), nil
			}
		}
		return nil, errs.NewFeatureVariationNotFoundWithVariationId(rule.GetRuleBase().Id, variationId)
	}
	return nil, errs.NewFeatureExperimentNotFound(experimentId)
}

type simulatedVariationModel struct {
	ExperimentId *int `json:"expId"`
	VariationId  *int `json:"varId"`
}

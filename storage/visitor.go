package storage

import (
	"sync"
	"time"

	"github.com/Kameleoon/client-go/v3/logging"
	"github.com/Kameleoon/client-go/v3/types"
)

type Visitor interface {
	LastActivityTime() time.Time
	UpdateLastActivityTime()

	LegalConsent() bool
	SetLegalConsent(consent bool)

	MappingIdentifier() *string
	SetMappingIdentifier(value *string)

	EnumerateSendableData(f func(types.Sendable) bool)
	CountSendableData() int

	UserAgent() string
	Device() *types.Device
	Browser() *types.Browser
	Cookie() *types.Cookie
	OperatingSystem() *types.OperatingSystem
	Geolocation() *types.Geolocation
	KcsHeat() *types.KcsHeat
	VisitorVisits() *types.VisitorVisits
	CustomData() DataMapStorage[int, *types.CustomData]
	PageViewVisits() DataMapStorage[string, types.PageViewVisit]
	Conversions() DataCollectionStorage[*types.Conversion]
	Variations() DataMapStorage[int, *types.AssignedVariation]

	AddData(logger logging.Logger, data ...types.Data)
	AddBaseData(logger logging.Logger, overwrite bool, data ...types.BaseData)
	AssignVariation(variation *types.AssignedVariation)
}

type VisitorImpl struct {
	mx                sync.RWMutex
	lastActivityTime  time.Time
	mappingIdentifier *string
	userAgent         string
	legalConsent      bool
	device            *types.Device
	browser           *types.Browser
	cookie            *types.Cookie
	operatingSystem   *types.OperatingSystem
	geolocation       *types.Geolocation
	kcsHeat           *types.KcsHeat
	visitorVisits     *types.VisitorVisits
	customDataMap     map[int]*types.CustomData
	pageViewVisits    map[string]types.PageViewVisit
	conversions       []*types.Conversion
	variations        map[int]*types.AssignedVariation
}

func NewVisitorImpl() *VisitorImpl {
	v := new(VisitorImpl)
	v.UpdateLastActivityTime()
	return v
}

func (v *VisitorImpl) LastActivityTime() time.Time {
	v.mx.RLock()
	defer v.mx.RUnlock()
	return v.lastActivityTime
}
func (v *VisitorImpl) UpdateLastActivityTime() {
	v.mx.Lock()
	defer v.mx.Unlock()
	v.lastActivityTime = time.Now()
}

func (v *VisitorImpl) LegalConsent() bool {
	return v.legalConsent
}
func (v *VisitorImpl) SetLegalConsent(consent bool) {
	v.legalConsent = consent
}

func (v *VisitorImpl) EnumerateSendableData(f func(types.Sendable) bool) {
	if (v.device != nil) && !f(v.device) {
		return
	}
	if (v.browser != nil) && !f(v.browser) {
		return
	}
	if (v.operatingSystem != nil) && !f(v.operatingSystem) {
		return
	}
	if (v.geolocation != nil) && !f(v.geolocation) {
		return
	}
	v.mx.RLock()
	defer v.mx.RUnlock()
	if !enumerateMap[int, *types.CustomData](v.customDataMap,
		func(cd *types.CustomData) bool { return f(cd) }) {
		return
	}
	if !enumerateMap[string, types.PageViewVisit](v.pageViewVisits,
		func(pvv types.PageViewVisit) bool { return f(pvv.PageView) }) {
		return
	}
	if !enumerateSlice[*types.Conversion](v.conversions,
		func(c *types.Conversion) bool { return f(c) }) {
		return
	}
	if !enumerateMap[int, *types.AssignedVariation](v.variations,
		func(v *types.AssignedVariation) bool { return f(v) }) {
		return
	}
}
func (v *VisitorImpl) CountSendableData() int {
	count := 0
	if v.device != nil {
		count++
	}
	if v.browser != nil {
		count++
	}
	if v.operatingSystem != nil {
		count++
	}
	if v.geolocation != nil {
		count++
	}
	v.mx.RLock()
	defer v.mx.RUnlock()
	count += len(v.customDataMap)
	count += len(v.pageViewVisits)
	count += len(v.conversions)
	count += len(v.variations)
	return count
}

func (v *VisitorImpl) MappingIdentifier() *string {
	return v.mappingIdentifier
}
func (v *VisitorImpl) SetMappingIdentifier(value *string) {
	v.mappingIdentifier = value
}

func (v *VisitorImpl) UserAgent() string {
	v.mx.RLock()
	defer v.mx.RUnlock()
	return v.userAgent
}

func (v *VisitorImpl) Device() *types.Device {
	return v.device
}

func (v *VisitorImpl) Browser() *types.Browser {
	return v.browser
}

func (v *VisitorImpl) Cookie() *types.Cookie {
	return v.cookie
}

func (v *VisitorImpl) OperatingSystem() *types.OperatingSystem {
	return v.operatingSystem
}

func (v *VisitorImpl) Geolocation() *types.Geolocation {
	return v.geolocation
}

func (v *VisitorImpl) KcsHeat() *types.KcsHeat {
	return v.kcsHeat
}

func (v *VisitorImpl) VisitorVisits() *types.VisitorVisits {
	return v.visitorVisits
}

func (v *VisitorImpl) CustomData() DataMapStorage[int, *types.CustomData] {
	return NewDataMapStorageImpl[int, *types.CustomData](&v.mx, v.customDataMap)
}

func (v *VisitorImpl) PageViewVisits() DataMapStorage[string, types.PageViewVisit] {
	return NewDataMapStorageImpl[string, types.PageViewVisit](&v.mx, v.pageViewVisits)
}

func (v *VisitorImpl) Conversions() DataCollectionStorage[*types.Conversion] {
	v.mx.RLock()
	defer v.mx.RUnlock()
	return NewDataCollectionStorageImpl[*types.Conversion](&v.mx, v.conversions)
}

func (v *VisitorImpl) Variations() DataMapStorage[int, *types.AssignedVariation] {
	return NewDataMapStorageImpl[int, *types.AssignedVariation](&v.mx, v.variations)
}

func (v *VisitorImpl) AddData(logger logging.Logger, data ...types.Data) {
	v.mx.Lock()
	defer v.mx.Unlock()
	for _, d := range data {
		v.addData(logger, true, d)
	}
}
func (v *VisitorImpl) AddBaseData(logger logging.Logger, overwrite bool, data ...types.BaseData) {
	v.mx.Lock()
	defer v.mx.Unlock()
	for _, d := range data {
		v.addData(logger, overwrite, d)
	}
}
func (v *VisitorImpl) addData(logger logging.Logger, overwrite bool, data types.BaseData) {
	dataType := data.DataType()
	switch dataType {
	case types.DataTypeUserAgent:
		v.addUserAgent(data)
	case types.DataTypeDevice:
		v.addDevice(data, overwrite)
	case types.DataTypeBrowser:
		v.addBrowser(data, overwrite)
	case types.DataTypeCookie:
		v.addCookie(data)
	case types.DataTypeOperatingSystem:
		v.addOperatingSystem(data, overwrite)
	case types.DataTypeGeolocation:
		v.addGeolocation(data, overwrite)
	case types.DataTypeKcsHeat:
		v.addKcsHeat(data)
	case types.DataTypeVisitorVisits:
		v.addVisitorVisits(data)
	case types.DataTypeCustom:
		v.addCustomData(data, overwrite)
	case types.DataTypePageView:
		v.addPageView(data)
	case types.DataTypePageViewVisit:
		v.addPageViewVisit(data)
	case types.DataTypeConversion:
		v.addConversion(data)
	case types.DataTypeAssignedVariation:
		v.addVariation(data, overwrite)
	default:
		if logger != nil {
			logger.Printf("Data has unsupported type '%s'", dataType)
		}
	}
}
func (v *VisitorImpl) addUserAgent(data types.BaseData) {
	if ua, ok := data.(types.UserAgent); ok {
		v.userAgent = ua.Value()
	}
}
func (v *VisitorImpl) addDevice(data types.BaseData, overwrite bool) {
	if d, ok := data.(*types.Device); ok && (overwrite || (v.device == nil)) {
		v.device = d
	}
}
func (v *VisitorImpl) addBrowser(data types.BaseData, overwrite bool) {
	if b, ok := data.(*types.Browser); ok && (overwrite || (v.browser == nil)) {
		v.browser = b
	}
}
func (v *VisitorImpl) addCookie(data types.BaseData) {
	if c, ok := data.(*types.Cookie); ok {
		v.cookie = c
	}
}
func (v *VisitorImpl) addOperatingSystem(data types.BaseData, overwrite bool) {
	if os, ok := data.(*types.OperatingSystem); ok && (overwrite || (v.operatingSystem == nil)) {
		v.operatingSystem = os
	}
}
func (v *VisitorImpl) addGeolocation(data types.BaseData, overwrite bool) {
	if g, ok := data.(*types.Geolocation); ok && (overwrite || (v.geolocation == nil)) {
		v.geolocation = g
	}
}
func (v *VisitorImpl) addKcsHeat(data types.BaseData) {
	if kh, ok := data.(*types.KcsHeat); ok {
		v.kcsHeat = kh
	}
}
func (v *VisitorImpl) addVisitorVisits(data types.BaseData) {
	if vv, ok := data.(*types.VisitorVisits); ok {
		v.visitorVisits = vv
	}
}
func (v *VisitorImpl) addCustomData(data types.BaseData, overwrite bool) {
	if cd, ok := data.(*types.CustomData); ok {
		if overwrite || (v.customDataMap[cd.ID()] == nil) {
			if v.customDataMap == nil {
				v.customDataMap = make(map[int]*types.CustomData, 1)
			}
			v.customDataMap[cd.ID()] = cd
		}
	}
}
func (v *VisitorImpl) addPageView(data types.BaseData) {
	if pv, ok := data.(*types.PageView); ok && (len(pv.URL()) > 0) {
		if v.pageViewVisits == nil {
			v.pageViewVisits = make(map[string]types.PageViewVisit, 1)
		}
		if pvv, contains := v.pageViewVisits[pv.URL()]; contains {
			v.pageViewVisits[pv.URL()] = pvv.Overwrite(pv)
		} else {
			v.pageViewVisits[pv.URL()] = types.PageViewVisit{PageView: pv, Count: 1}
		}
	}
}
func (v *VisitorImpl) addPageViewVisit(data types.BaseData) {
	if pvv, ok := data.(types.PageViewVisit); ok && (pvv.PageView != nil) && (pvv.PageView.URL() != "") {
		if v.pageViewVisits == nil {
			v.pageViewVisits = make(map[string]types.PageViewVisit, 1)
		}
		url := pvv.PageView.URL()
		if former, contains := v.pageViewVisits[url]; contains {
			pvv = former.Merge(pvv)
		}
		v.pageViewVisits[url] = pvv
	}
}
func (v *VisitorImpl) addConversion(data types.BaseData) {
	if c, ok := data.(*types.Conversion); ok {
		if v.conversions == nil {
			v.conversions = make([]*types.Conversion, 0, 1)
		}
		v.conversions = append(v.conversions, c)
	}
}
func (v *VisitorImpl) addVariation(data types.BaseData, overwrite bool) {
	if av, ok := data.(*types.AssignedVariation); ok {
		v.assignVariation(av, overwrite)
	}
}

func (v *VisitorImpl) AssignVariation(variation *types.AssignedVariation) {
	v.mx.Lock()
	defer v.mx.Unlock()
	v.assignVariation(variation, true)
}
func (v *VisitorImpl) assignVariation(variation *types.AssignedVariation, overwrite bool) {
	if overwrite || (v.variations[variation.ExperimentId()] == nil) {
		if v.variations == nil {
			v.variations = make(map[int]*types.AssignedVariation, 1)
		}
		v.variations[variation.ExperimentId()] = variation
	}
}

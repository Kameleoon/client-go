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

	IsUniqueIdentifier() bool
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
	CustomData() DataMapStorage[int, types.ICustomData]
	PageViewVisits() DataMapStorage[string, types.PageViewVisit]
	Conversions() DataCollectionStorage[*types.Conversion]
	Variations() DataMapStorage[int, *types.AssignedVariation]

	AddData(data ...types.Data)
	AddBaseData(overwrite bool, data ...types.BaseData)
	AssignVariation(variation *types.AssignedVariation)

	Clone() Visitor
}

type VisitorImpl struct {
	data               *visitorData
	isUniqueIdentifier bool
}

func NewVisitorImpl() *VisitorImpl {
	v := &VisitorImpl{data: new(visitorData)}
	v.UpdateLastActivityTime()
	return v
}

func cloneVisitorImpl(src *VisitorImpl) *VisitorImpl {
	v := &VisitorImpl{
		data:               src.data,
		isUniqueIdentifier: src.isUniqueIdentifier,
	}
	v.UpdateLastActivityTime()
	return v
}

func (v VisitorImpl) String() string {
	return "VisitorImpl{}"
}

func (v *VisitorImpl) LastActivityTime() time.Time {
	v.data.mx.RLock()
	defer v.data.mx.RUnlock()
	return v.data.lastActivityTime
}
func (v *VisitorImpl) UpdateLastActivityTime() {
	v.data.mx.Lock()
	defer v.data.mx.Unlock()
	v.data.lastActivityTime = time.Now()
}

func (v *VisitorImpl) IsUniqueIdentifier() bool {
	return v.isUniqueIdentifier
}

func (v *VisitorImpl) LegalConsent() bool {
	return v.data.legalConsent
}
func (v *VisitorImpl) SetLegalConsent(consent bool) {
	v.data.legalConsent = consent
}

func (v *VisitorImpl) MappingIdentifier() *string {
	return v.data.mappingIdentifier
}
func (v *VisitorImpl) SetMappingIdentifier(value *string) {
	if v.data.mappingIdentifier == nil {
		v.data.mx.Lock()
		defer v.data.mx.Unlock()
		if v.data.mappingIdentifier == nil {
			v.data.mappingIdentifier = value
		}
	}
}

func (v *VisitorImpl) EnumerateSendableData(f func(types.Sendable) bool) {
	v.data.enumerateSendableData(f)
}
func (v *VisitorImpl) CountSendableData() int {
	return v.data.countSendableData()
}

func (v *VisitorImpl) UserAgent() string {
	v.data.mx.RLock()
	defer v.data.mx.RUnlock()
	return v.data.userAgent
}

func (v *VisitorImpl) Device() *types.Device {
	logging.Debug("CALL/RETURN: VisitorImpl.Device() -> (device: %s)", v.data.device)
	return v.data.device
}

func (v *VisitorImpl) Browser() *types.Browser {
	logging.Debug("CALL/RETURN: VisitorImpl.Browser() -> (browser: %s)", v.data.browser)
	return v.data.browser
}

func (v *VisitorImpl) Cookie() *types.Cookie {
	logging.Debug("CALL/RETURN: VisitorImpl.Cookie() -> (cookie: %s)", v.data.cookie)
	return v.data.cookie
}

func (v *VisitorImpl) OperatingSystem() *types.OperatingSystem {
	logging.Debug("CALL/RETURN: VisitorImpl.OperatingSystem() -> (operatingSystem: %s)", v.data.operatingSystem)
	return v.data.operatingSystem
}

func (v *VisitorImpl) Geolocation() *types.Geolocation {
	logging.Debug("CALL/RETURN: VisitorImpl.Geolocation() -> (geolocation: %s)", v.data.geolocation)
	return v.data.geolocation
}

func (v *VisitorImpl) KcsHeat() *types.KcsHeat {
	logging.Debug("CALL/RETURN: VisitorImpl.KcsHeat() -> (kcsHeat: %s)", v.data.kcsHeat)
	return v.data.kcsHeat
}

func (v *VisitorImpl) VisitorVisits() *types.VisitorVisits {
	logging.Debug("CALL/RETURN: VisitorImpl.VisitorVisits() -> (visitorVisits: %s)", v.data.visitorVisits)
	return v.data.visitorVisits
}

func (v *VisitorImpl) CustomData() DataMapStorage[int, types.ICustomData] {
	logging.Debug("CALL: VisitorImpl.CustomData()")
	storage := NewDataMapStorageImpl(&v.data.mx, &v.data.customDataMap)
	logging.Debug("RETURN: VisitorImpl.CustomData() -> (customData: %s)", storage)
	return storage
}

func (v *VisitorImpl) PageViewVisits() DataMapStorage[string, types.PageViewVisit] {
	logging.Debug("CALL: VisitorImpl.PageViewVisits()")
	storage := NewDataMapStorageImpl(&v.data.mx, &v.data.pageViewVisits)
	logging.Debug("RETURN: VisitorImpl.PageViewVisits() -> (pageViewVisits: %s)", storage)
	return storage
}

func (v *VisitorImpl) Conversions() DataCollectionStorage[*types.Conversion] {
	logging.Debug("CALL: VisitorImpl.Conversions()")
	storage := NewDataCollectionStorageImpl(&v.data.mx, &v.data.conversions)
	logging.Debug("RETURN: VisitorImpl.Conversions() -> (conversions: %s)", storage)
	return storage
}

func (v *VisitorImpl) Variations() DataMapStorage[int, *types.AssignedVariation] {
	logging.Debug("CALL: VisitorImpl.Variations()")
	storage := NewDataMapStorageImpl(&v.data.mx, &v.data.variations)
	logging.Debug("RETURN: VisitorImpl.Variations() -> (variations: %s)", storage)
	return storage
}

func (v *VisitorImpl) AddData(data ...types.Data) {
	logging.Debug("CALL: VisitorImpl.AddData(data: %s)", data)
	v.data.mx.Lock()
	defer v.data.mx.Unlock()
	for _, d := range data {
		v.addData(true, d)
	}
	logging.Debug("RETURN: VisitorImpl.AddData(data: %s)", data)
}
func (v *VisitorImpl) AddBaseData(overwrite bool, data ...types.BaseData) {
	logging.Debug("CALL: VisitorImpl.AddBaseData(overwrite: %s, data: %s)", overwrite, data)
	v.data.mx.Lock()
	defer v.data.mx.Unlock()
	for _, d := range data {
		v.addData(overwrite, d)
	}
	logging.Debug("RETURN: VisitorImpl.AddBaseData(overwrite: %s, data: %s)", overwrite, data)
}
func (v *VisitorImpl) addData(overwrite bool, data types.BaseData) {
	logging.Debug("CALL: VisitorImpl.AddData(overwrite: %s, data: %s)", overwrite, data)
	dataType := data.DataType()
	switch dataType {
	case types.DataTypeUserAgent:
		v.data.addUserAgent(data)
	case types.DataTypeDevice:
		v.data.addDevice(data, overwrite)
	case types.DataTypeBrowser:
		v.data.addBrowser(data, overwrite)
	case types.DataTypeCookie:
		v.data.addCookie(data)
	case types.DataTypeOperatingSystem:
		v.data.addOperatingSystem(data, overwrite)
	case types.DataTypeGeolocation:
		v.data.addGeolocation(data, overwrite)
	case types.DataTypeKcsHeat:
		v.data.addKcsHeat(data)
	case types.DataTypeVisitorVisits:
		v.data.addVisitorVisits(data)
	case types.DataTypeCustom:
		v.data.addCustomData(data, overwrite)
	case types.DataTypePageView:
		v.data.addPageView(data)
	case types.DataTypePageViewVisit:
		v.data.addPageViewVisit(data)
	case types.DataTypeConversion:
		v.data.addConversion(data)
	case types.DataTypeAssignedVariation:
		v.data.addVariation(data, overwrite)
	case types.DataTypeUniqueIdentifier:
		v.setUniqueIdentifier(data)
	default:
		logging.Warning("Data has unsupported type %s", dataType)
	}
	logging.Debug("RETURN: VisitorImpl.AddData(overwrite: %s, data: %s)", overwrite, data)
}

func (v *VisitorImpl) AssignVariation(variation *types.AssignedVariation) {
	v.data.mx.Lock()
	defer v.data.mx.Unlock()
	v.data.assignVariation(variation, true)
}

func (v *VisitorImpl) setUniqueIdentifier(data types.BaseData) {
	if ui, ok := data.(*types.UniqueIdentifier); ok {
		v.isUniqueIdentifier = ui.Value()
	}
}

func (v *VisitorImpl) Clone() Visitor {
	return cloneVisitorImpl(v)
}

type visitorData struct {
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
	customDataMap     map[int]types.ICustomData
	pageViewVisits    map[string]types.PageViewVisit
	conversions       []*types.Conversion
	variations        map[int]*types.AssignedVariation
}

func (vd *visitorData) enumerateSendableData(f func(types.Sendable) bool) {
	if (vd.device != nil) && !f(vd.device) {
		return
	}
	if (vd.browser != nil) && !f(vd.browser) {
		return
	}
	if (vd.operatingSystem != nil) && !f(vd.operatingSystem) {
		return
	}
	if (vd.geolocation != nil) && !f(vd.geolocation) {
		return
	}
	vd.mx.RLock()
	defer vd.mx.RUnlock()
	if !enumerateMap[int, types.ICustomData](vd.customDataMap,
		func(cd types.ICustomData) bool { return f(cd) }) {
		return
	}
	if !enumerateMap[string, types.PageViewVisit](vd.pageViewVisits,
		func(pvv types.PageViewVisit) bool { return f(pvv.PageView) }) {
		return
	}
	if !enumerateSlice[*types.Conversion](vd.conversions,
		func(c *types.Conversion) bool { return f(c) }) {
		return
	}
	if !enumerateMap[int, *types.AssignedVariation](vd.variations,
		func(av *types.AssignedVariation) bool { return f(av) }) {
		return
	}
}
func (vd *visitorData) countSendableData() int {
	count := 0
	if vd.device != nil {
		count++
	}
	if vd.browser != nil {
		count++
	}
	if vd.operatingSystem != nil {
		count++
	}
	if vd.geolocation != nil {
		count++
	}
	vd.mx.RLock()
	defer vd.mx.RUnlock()
	count += len(vd.customDataMap)
	count += len(vd.pageViewVisits)
	count += len(vd.conversions)
	count += len(vd.variations)
	return count
}

func (vd *visitorData) addUserAgent(data types.BaseData) {
	if ua, ok := data.(types.UserAgent); ok {
		vd.userAgent = ua.Value()
	}
}
func (vd *visitorData) addDevice(data types.BaseData, overwrite bool) {
	if d, ok := data.(*types.Device); ok && (overwrite || (vd.device == nil)) {
		vd.device = d
	}
}
func (vd *visitorData) addBrowser(data types.BaseData, overwrite bool) {
	if b, ok := data.(*types.Browser); ok && (overwrite || (vd.browser == nil)) {
		vd.browser = b
	}
}
func (vd *visitorData) addCookie(data types.BaseData) {
	if c, ok := data.(*types.Cookie); ok {
		vd.cookie = c
	}
}
func (vd *visitorData) addOperatingSystem(data types.BaseData, overwrite bool) {
	if os, ok := data.(*types.OperatingSystem); ok && (overwrite || (vd.operatingSystem == nil)) {
		vd.operatingSystem = os
	}
}
func (vd *visitorData) addGeolocation(data types.BaseData, overwrite bool) {
	if g, ok := data.(*types.Geolocation); ok && (overwrite || (vd.geolocation == nil)) {
		vd.geolocation = g
	}
}
func (vd *visitorData) addKcsHeat(data types.BaseData) {
	if kh, ok := data.(*types.KcsHeat); ok {
		vd.kcsHeat = kh
	}
}
func (vd *visitorData) addVisitorVisits(data types.BaseData) {
	if vv, ok := data.(*types.VisitorVisits); ok {
		vd.visitorVisits = vv
	}
}
func (vd *visitorData) addCustomData(data types.BaseData, overwrite bool) {
	if cd, ok := data.(types.ICustomData); ok {
		if overwrite || (vd.customDataMap[cd.ID()] == nil) {
			if vd.customDataMap == nil {
				vd.customDataMap = make(map[int]types.ICustomData, 1)
			}
			vd.customDataMap[cd.ID()] = cd
		}
	}
}
func (vd *visitorData) addPageView(data types.BaseData) {
	if pv, ok := data.(*types.PageView); ok && (len(pv.URL()) > 0) {
		if vd.pageViewVisits == nil {
			vd.pageViewVisits = make(map[string]types.PageViewVisit, 1)
		}
		if pvv, contains := vd.pageViewVisits[pv.URL()]; contains {
			vd.pageViewVisits[pv.URL()] = pvv.Overwrite(pv)
		} else {
			vd.pageViewVisits[pv.URL()] = types.NewPageViewVisit(pv, 1)
		}
	}
}
func (vd *visitorData) addPageViewVisit(data types.BaseData) {
	if pvv, ok := data.(types.PageViewVisit); ok && (pvv.PageView != nil) && (pvv.PageView.URL() != "") {
		if vd.pageViewVisits == nil {
			vd.pageViewVisits = make(map[string]types.PageViewVisit, 1)
		}
		url := pvv.PageView.URL()
		if former, contains := vd.pageViewVisits[url]; contains {
			pvv = former.Merge(pvv)
		}
		vd.pageViewVisits[url] = pvv
	}
}
func (vd *visitorData) addConversion(data types.BaseData) {
	if c, ok := data.(*types.Conversion); ok {
		if vd.conversions == nil {
			vd.conversions = make([]*types.Conversion, 0, 1)
		}
		vd.conversions = append(vd.conversions, c)
	}
}
func (vd *visitorData) addVariation(data types.BaseData, overwrite bool) {
	if av, ok := data.(*types.AssignedVariation); ok {
		vd.assignVariation(av, overwrite)
	}
}

func (vd *visitorData) assignVariation(variation *types.AssignedVariation, overwrite bool) {
	if overwrite || (vd.variations[variation.ExperimentId()] == nil) {
		if vd.variations == nil {
			vd.variations = make(map[int]*types.AssignedVariation, 1)
		}
		vd.variations[variation.ExperimentId()] = variation
	}
}

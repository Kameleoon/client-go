package storage

import (
	"sync"
	"time"

	"github.com/Kameleoon/client-go/v3/logging"
	"github.com/Kameleoon/client-go/v3/types"
)

type Visitor interface {
	TimeStarted() time.Time
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
	CBScores() *types.CBScores
	VisitorVisits() *types.VisitorVisits
	CustomData() DataMapStorage[int, types.ICustomData]
	PageViewVisits() DataMapStorage[string, types.PageViewVisit]
	Conversions() DataCollectionStorage[*types.Conversion]
	Variations() DataMapStorage[int, *types.AssignedVariation]
	Personalizations() DataMapStorage[int, *types.Personalization]
	TargetedSegments() DataMapStorage[int, *types.TargetedSegment]

	AddData(data ...types.Data)
	AddBaseData(overwrite bool, data ...types.BaseData)
	AssignVariation(variation *types.AssignedVariation)

	GetForcedFeatureVariation(featureKey string) *types.ForcedFeatureVariation
	GetForcedExperimentVariation(experimentId int) *types.ForcedExperimentVariation
	ResetForcedVariation(experimentId int)
	UpdateSimulatedVariations(variations []*types.ForcedFeatureVariation)

	Clone() Visitor
}

type VisitorImpl struct {
	data               *visitorData
	isUniqueIdentifier bool
}

func NewVisitorImpl() *VisitorImpl {
	v := &VisitorImpl{data: newVisitorData()}
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

func (v *VisitorImpl) TimeStarted() time.Time {
	return v.data.timeStarted
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
	d := v.data.device
	logging.Debug("CALL/RETURN: VisitorImpl.Device() -> (device: %s)", d)
	return d
}

func (v *VisitorImpl) Browser() *types.Browser {
	b := v.data.browser
	logging.Debug("CALL/RETURN: VisitorImpl.Browser() -> (browser: %s)", b)
	return b
}

func (v *VisitorImpl) Cookie() *types.Cookie {
	c := v.data.cookie
	logging.Debug("CALL/RETURN: VisitorImpl.Cookie() -> (cookie: %s)", c)
	return c
}

func (v *VisitorImpl) OperatingSystem() *types.OperatingSystem {
	os := v.data.operatingSystem
	logging.Debug("CALL/RETURN: VisitorImpl.OperatingSystem() -> (operatingSystem: %s)", os)
	return os
}

func (v *VisitorImpl) Geolocation() *types.Geolocation {
	g := v.data.geolocation
	logging.Debug("CALL/RETURN: VisitorImpl.Geolocation() -> (geolocation: %s)", g)
	return g
}

func (v *VisitorImpl) KcsHeat() *types.KcsHeat {
	kcs := v.data.kcsHeat
	logging.Debug("CALL/RETURN: VisitorImpl.KcsHeat() -> (kcsHeat: %s)", kcs)
	return kcs
}

func (v *VisitorImpl) CBScores() *types.CBScores {
	cbs := v.data.cbscores
	logging.Debug("CALL/RETURN: VisitorImpl.CBScores() -> (cbs: %s)", cbs)
	return cbs
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

func (v *VisitorImpl) Personalizations() DataMapStorage[int, *types.Personalization] {
	logging.Debug("CALL: VisitorImpl.Personalizations()")
	storage := NewDataMapStorageImpl(&v.data.mx, &v.data.personalizations)
	logging.Debug("RETURN: VisitorImpl.Personalizations() -> (personalizations: %s)", storage)
	return storage
}

func (v *VisitorImpl) TargetedSegments() DataMapStorage[int, *types.TargetedSegment] {
	logging.Debug("CALL: VisitorImpl.TargetedSegments()")
	storage := NewDataMapStorageImpl(&v.data.mx, &v.data.targetedSegments)
	logging.Debug("RETURN: VisitorImpl.TargetedSegments() -> (targetedSegments: %s)", storage)
	return storage
}

func (v *VisitorImpl) GetForcedFeatureVariation(featureKey string) *types.ForcedFeatureVariation {
	logging.Debug("CALL: VisitorImpl.GetForcedFeatureVariation(featureKey: %s)", featureKey)
	var variation *types.ForcedFeatureVariation
	if v.data.simulatedVariations != nil {
		v.data.mx.RLock()
		variation = v.data.simulatedVariations[featureKey]
		v.data.mx.RUnlock()
	}
	logging.Debug(
		"RETURN: VisitorImpl.GetForcedFeatureVariation(featureKey: %s) -> (variation: %s)",
		featureKey, variation,
	)
	return variation
}
func (v *VisitorImpl) GetForcedExperimentVariation(experimentId int) *types.ForcedExperimentVariation {
	logging.Debug("CALL: VisitorImpl.GetForcedExperimentVariation(experimentId: %s)", experimentId)
	var variation *types.ForcedExperimentVariation
	if v.data.forcedVariations != nil {
		v.data.mx.RLock()
		variation = v.data.forcedVariations[experimentId]
		v.data.mx.RUnlock()
	}
	logging.Debug(
		"RETURN: VisitorImpl.GetForcedExperimentVariation(experimentId: %s) -> (variation: %s)",
		experimentId, variation,
	)
	return variation
}
func (v *VisitorImpl) ResetForcedVariation(experimentId int) {
	logging.Debug("CALL: VisitorImpl.ResetForcedVariation(experimentId: %s)", experimentId)
	if v.data.forcedVariations != nil {
		v.data.mx.Lock()
		delete(v.data.forcedVariations, experimentId)
		v.data.mx.Unlock()
	}
	logging.Debug("RETURN: VisitorImpl.ResetForcedVariation(experimentId: %s)", experimentId)
}
func (v *VisitorImpl) UpdateSimulatedVariations(variations []*types.ForcedFeatureVariation) {
	if (len(v.data.simulatedVariations) == 0) && (len(variations) == 0) {
		return
	}
	logging.Debug("CALL: VisitorImpl.UpdateSimulatedVariations(variations: %s)", variations)
	newSimulatedVariations := make(map[string]*types.ForcedFeatureVariation)
	for _, variation := range variations {
		newSimulatedVariations[variation.FeatureKey()] = variation
	}
	v.data.mx.Lock()
	v.data.simulatedVariations = newSimulatedVariations
	v.data.mx.Unlock()
	logging.Debug("RETURN: VisitorImpl.UpdateSimulatedVariations(variations: %s)", variations)
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
	case types.DataTypeCBScores:
		v.data.addCBScores(data, overwrite)
	case types.DataTypeVisitorVisits:
		v.data.addVisitorVisits(data, overwrite)
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
	case types.DataTypePersonalization:
		v.data.addPersonalization(data, overwrite)
	case types.DataTypeTargetedSegment:
		v.data.addTargetedSegment(data)
	case types.DataTypeForcedFeatureVariation:
		v.data.addForcedFeatureVariation(data)
	case types.DataTypeForcedExperimentVariation:
		v.data.addForcedExperimentVariation(data)
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
	mx                  sync.RWMutex
	timeStarted         time.Time
	lastActivityTime    time.Time
	mappingIdentifier   *string
	userAgent           string
	legalConsent        bool
	device              *types.Device
	browser             *types.Browser
	cookie              *types.Cookie
	operatingSystem     *types.OperatingSystem
	geolocation         *types.Geolocation
	kcsHeat             *types.KcsHeat
	cbscores            *types.CBScores
	visitorVisits       *types.VisitorVisits
	customDataMap       map[int]types.ICustomData
	pageViewVisits      map[string]types.PageViewVisit
	conversions         []*types.Conversion
	variations          map[int]*types.AssignedVariation
	personalizations    map[int]*types.Personalization
	targetedSegments    map[int]*types.TargetedSegment
	forcedVariations    map[int]*types.ForcedExperimentVariation
	simulatedVariations map[string]*types.ForcedFeatureVariation
}

func newVisitorData() *visitorData {
	return &visitorData{
		timeStarted: time.Now(),
	}
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
	if (vd.visitorVisits != nil) && !f(vd.visitorVisits) {
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
	if !enumerateMap[int, *types.AssignedVariation](vd.variations,
		func(av *types.AssignedVariation) bool { return f(av) }) {
		return
	}
	if !enumerateMap[int, *types.TargetedSegment](vd.targetedSegments,
		func(ts *types.TargetedSegment) bool { return f(ts) }) {
		return
	}
	if !enumerateSlice[*types.Conversion](vd.conversions,
		func(c *types.Conversion) bool { return f(c) }) {
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
	if vd.visitorVisits != nil {
		count++
	}
	vd.mx.RLock()
	defer vd.mx.RUnlock()
	count += len(vd.customDataMap)
	count += len(vd.pageViewVisits)
	count += len(vd.variations)
	count += len(vd.targetedSegments)
	count += len(vd.conversions)
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
func (vd *visitorData) addCBScores(data types.BaseData, overwrite bool) {
	if cbs, ok := data.(*types.CBScores); ok && (overwrite || (vd.cbscores == nil)) {
		vd.cbscores = cbs
	}
}
func (vd *visitorData) addVisitorVisits(data types.BaseData, overwrite bool) {
	if vv, ok := data.(*types.VisitorVisits); ok && (overwrite || (vd.visitorVisits == nil)) {
		vd.visitorVisits = vv.Localize(vd.timeStarted.UnixMilli())
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
func (vd *visitorData) addPersonalization(data types.BaseData, overwrite bool) {
	if p, ok := data.(*types.Personalization); ok {
		if overwrite || (vd.personalizations[p.Id()] == nil) {
			if vd.personalizations == nil {
				vd.personalizations = make(map[int]*types.Personalization, 1)
			}
			vd.personalizations[p.Id()] = p
		}
	}
}
func (vd *visitorData) addTargetedSegment(data types.BaseData) {
	if ts, ok := data.(*types.TargetedSegment); ok {
		if vd.targetedSegments == nil {
			vd.targetedSegments = make(map[int]*types.TargetedSegment, 1)
		}
		vd.targetedSegments[ts.Id()] = ts
	}
}
func (vd *visitorData) addForcedFeatureVariation(data types.BaseData) {
	if ffv, ok := data.(*types.ForcedFeatureVariation); ok {
		if vd.simulatedVariations == nil {
			vd.simulatedVariations = make(map[string]*types.ForcedFeatureVariation, 1)
		}
		vd.simulatedVariations[ffv.FeatureKey()] = ffv
	}
}
func (vd *visitorData) addForcedExperimentVariation(data types.BaseData) {
	if fev, ok := data.(*types.ForcedExperimentVariation); ok {
		if vd.forcedVariations == nil {
			vd.forcedVariations = make(map[int]*types.ForcedExperimentVariation, 1)
		}
		vd.forcedVariations[fev.Rule().GetRuleBase().ExperimentId] = fev
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

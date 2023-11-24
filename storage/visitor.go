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

	EnumerateSendableData(f func(types.Sendable) bool)
	CountSendableData() int

	UserAgent() string
	Device() *types.Device
	Browser() *types.Browser
	CustomData() DataMapStorage[int, *types.CustomData]
	PageViewVisits() DataMapStorage[string, types.PageViewVisit]
	Conversions() DataCollectionStorage[*types.Conversion]
	Variations() DataMapStorage[int, *types.AssignedVariation]

	AddData(logger logging.Logger, data ...types.Data)
	AssignVariation(variation *types.AssignedVariation)
}

type VisitorImpl struct {
	mx               sync.RWMutex
	lastActivityTime time.Time
	userAgent        string
	legalConsent     bool
	device           *types.Device
	browser          *types.Browser
	customDataMap    map[int]*types.CustomData
	pageViewVisits   map[string]types.PageViewVisit
	conversions      []*types.Conversion
	variations       map[int]*types.AssignedVariation
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
	v.mx.RLock()
	defer v.mx.RUnlock()
	count += len(v.customDataMap)
	count += len(v.pageViewVisits)
	count += len(v.conversions)
	count += len(v.variations)
	return count
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
		dataType := d.DataType()
		switch dataType {
		case types.DataTypeUserAgent:
			v.addUserAgent(d)
		case types.DataTypeDevice:
			v.addDevice(d)
		case types.DataTypeBrowser:
			v.addBrowser(d)
		case types.DataTypeCustom:
			v.addCustomData(d)
		case types.DataTypePageView:
			v.addPageView(d)
		case types.DataTypeConversion:
			v.addConversion(d)
		default:
			if logger != nil {
				logger.Printf("Data has unsupported type '%s'", dataType)
			}
		}
	}
}
func (v *VisitorImpl) addUserAgent(data types.Data) {
	if ua, ok := data.(types.UserAgent); ok {
		v.userAgent = ua.Value()
	}
}
func (v *VisitorImpl) addDevice(data types.Data) {
	if d, ok := data.(*types.Device); ok {
		v.device = d
	}
}
func (v *VisitorImpl) addBrowser(data types.Data) {
	if b, ok := data.(*types.Browser); ok {
		v.browser = b
	}
}
func (v *VisitorImpl) addCustomData(data types.Data) {
	if cd, ok := data.(*types.CustomData); ok {
		if v.customDataMap == nil {
			v.customDataMap = make(map[int]*types.CustomData, 1)
		}
		v.customDataMap[cd.ID()] = cd
	}
}
func (v *VisitorImpl) addPageView(data types.Data) {
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
func (v *VisitorImpl) addConversion(data types.Data) {
	if c, ok := data.(*types.Conversion); ok {
		if v.conversions == nil {
			v.conversions = make([]*types.Conversion, 0, 1)
		}
		v.conversions = append(v.conversions, c)
	}
}

func (v *VisitorImpl) AssignVariation(variation *types.AssignedVariation) {
	v.mx.Lock()
	defer v.mx.Unlock()
	if v.variations == nil {
		v.variations = make(map[int]*types.AssignedVariation, 1)
	}
	v.variations[variation.ExperimentId()] = variation
}

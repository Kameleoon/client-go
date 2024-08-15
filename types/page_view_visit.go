package types

import "fmt"

type PageViewVisit struct {
	PageView      *PageView
	Count         int
	LastTimestamp int64 // in milliseconds (server returns in ms as well)
}

func (pvv PageViewVisit) String() string {
	return fmt.Sprintf(
		"PageViewVisit{PageView:%v,Count:%v,LastTimestamp:%v}",
		pvv.PageView,
		pvv.Count,
		pvv.LastTimestamp,
	)
}

func (pvv PageViewVisit) Overwrite(newPageView *PageView) PageViewVisit {
	pvv.PageView = newPageView
	pvv.Count++
	return pvv
}

func (pvv PageViewVisit) Merge(other PageViewVisit) PageViewVisit {
	pvv.Count += other.Count
	if other.LastTimestamp > pvv.LastTimestamp {
		pvv.LastTimestamp = other.LastTimestamp
	}
	return pvv
}

func (pvv PageViewVisit) DataType() DataType {
	return DataTypePageViewVisit
}

package types

import (
	"fmt"
	"time"
)

type PageViewVisit struct {
	PageView      *PageView
	Count         int
	LastTimestamp int64 // in milliseconds (server returns in ms as well)
}

func NewPageViewVisit(pageView *PageView, count int, lastTimestamp ...int64) PageViewVisit {
	var lastTS int64
	if len(lastTimestamp) > 0 {
		lastTS = lastTimestamp[0]
	} else {
		lastTS = time.Now().UnixMilli()
	}
	pvv := PageViewVisit{
		PageView:      pageView,
		Count:         count,
		LastTimestamp: lastTS,
	}
	return pvv
}

func (pvv PageViewVisit) Overwrite(newPageView *PageView) PageViewVisit {
	return NewPageViewVisit(newPageView, pvv.Count+1)
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

func (pvv PageViewVisit) String() string {
	return fmt.Sprintf(
		"PageViewVisit{PageView:%v,Count:%v,LastTimestamp:%v}",
		pvv.PageView,
		pvv.Count,
		pvv.LastTimestamp,
	)
}

package types

type PageViewVisit struct {
	PageView *PageView
	Count    int
}

func (pvv PageViewVisit) Overwrite(newPageView *PageView) PageViewVisit {
	return PageViewVisit{
		PageView: newPageView,
		Count:    pvv.Count + 1,
	}
}

package kameleoon

import (
	"strings"

	"github.com/Kameleoon/client-go/types"
	"github.com/Kameleoon/client-go/utils"
)

const (
	TrackingRequestData       = "dataTracking"
	TrackingRequestExperiment = "experimentTracking"
)

type trackingRequest struct {
	Type          string
	VisitorCode   string
	VariationID   string
	ExperimentID  int
	NoneVariation bool
}

const defaultPostMaxRetries = 10

func (c *Client) postTrackingAsync(r trackingRequest) {
	req := request{
		URL:          c.buildTrackingPath(c.Cfg.TrackingURL, r),
		Method:       MethodPost,
		ContentType:  HeaderContentTypeText,
		ClientHeader: c.Cfg.TrackingVersion,
	}
	c.m.Lock()
	req.AuthToken = c.token
	c.m.Unlock()

	data := c.selectSendData(r.VisitorCode)
	c.log("Start post to tracking: %s", data)
	var sb strings.Builder
	var err error
	for _, dataCell := range data {
		for i := 0; i < len(dataCell.Data); i++ {
			if _, exist := dataCell.Index[i]; exist {
				continue
			}
			sb.WriteString(dataCell.Data[i].QueryEncode())
			sb.WriteByte('\n')
		}
		if sb.Len() == 0 {
			continue
		}
		req.BodyString = sb.String()
		sb.Reset()
		for i := defaultPostMaxRetries; i > 0; i-- {
			err = c.rest.Do(req, nil)
			if err == nil {
				break
			}
			c.log("Trials amount left: %d, error: %v", i, err)
		}
		if err != nil {
			c.log("Failed to post tracking data, error: %v", err)
			err = nil
			continue
		}
		for i := 0; i < len(dataCell.Data); i++ {
			if _, exist := dataCell.Index[i]; exist {
				continue
			}
			dataCell.Index[i] = struct{}{}
		}
	}

	c.log("Post to tracking done")
}

func (c *Client) selectSendData(visitorCode ...string) []*types.DataCell {
	var data []*types.DataCell
	if len(visitorCode) > 0 && len(visitorCode[0]) > 0 {
		if dc := c.getDataCell(visitorCode[0]); dc != nil && len(dc.Data) != len(dc.Index) {
			data = append(data, dc)
		}
		return data
	}
	for kv := range c.Data.Iter() {
		if dc, ok := kv.Value.(*types.DataCell); ok {
			if len(dc.Data) == len(dc.Index) {
				continue
			}
			data = append(data, dc)
		}
	}
	return data
}

func (c *Client) buildTrackingPath(base string, r trackingRequest) string {
	var b strings.Builder
	switch r.Type {
	case TrackingRequestData:
		b.WriteString(base)
		b.WriteString("/dataTracking?siteCode=")
		b.WriteString(c.Cfg.SiteCode)
		b.WriteString("&visitorCode=")
		b.WriteString(r.VisitorCode)
		b.WriteString("&nonce=")
		b.WriteString(types.GetNonce())
		b.WriteString("&experimentID=")
		b.WriteString(utils.WriteUint(r.ExperimentID))
		return b.String()
	case TrackingRequestExperiment:
		b.WriteString(API_SSX_URL)
		b.WriteString("/experimentTracking?siteCode=")
		b.WriteString(c.Cfg.SiteCode)
		b.WriteString("&visitorCode=")
		b.WriteString(r.VisitorCode)
		b.WriteString("&nonce=")
		b.WriteString(types.GetNonce())
		b.WriteString("&experimentID=")
		b.WriteString(utils.WriteUint(r.ExperimentID))
		if len(r.VariationID) == 0 {
			return b.String()
		}
		b.WriteString("&variationId=")
		b.WriteString(r.VariationID)
		if r.NoneVariation {
			b.WriteString("&noneVariation=true")
		}
		return b.String()
	}
	return ""
}

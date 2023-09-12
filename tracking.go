package kameleoon

import (
	"github.com/Kameleoon/client-go/v2/network"
	"github.com/Kameleoon/client-go/v2/types"
)

func (c *Client) sendTrackingRequest(visitorCode string, experimentId *int, variationId *int) {
	cells := c.selectSendData(visitorCode)
	var cell *types.DataCell
	if len(cells) > 0 {
		cell = cells[0]
	}
	unsent, lim := c.selectUnsentData(cell)
	go func() {
		sent := c.makeTrackingRequest(visitorCode, unsent, experimentId, variationId)
		if sent {
			c.markDataAsSent(cell, lim)
		}
	}()
}

func (c *Client) makeTrackingRequest(visitorCode string, data []network.QueryEncodable,
	experimentId *int, variationId *int) (sent bool) {
	ua := c.getUserAgent(visitorCode)
	token := c.token
	if (experimentId != nil) && (variationId != nil) {
		data = append(data, network.ExperimentEvent{ExperimentId: *experimentId, VariationId: *variationId})
	} else if len(data) == 0 {
		data = append(data, network.ActivityEvent{})
	}
	c.log("Start post to tracking")
	out, err := c.networkManager.SendTrackingData(visitorCode, data, ua, token, -1)
	if err != nil {
		c.log("Failed to post tracking data, error: %v", err)
		return false
	}
	c.log("Post to tracking done")
	return out
}

func (c *Client) selectUnsentData(cell *types.DataCell) ([]network.QueryEncodable, int) {
	var unsent []network.QueryEncodable
	if cell == nil {
		return unsent, -1
	}
	c.m.Lock()
	defer c.m.Unlock()
	for i := 0; i < len(cell.Data); i++ {
		if _, sent := cell.Index[i]; !sent {
			unsent = append(unsent, cell.Data[i])
		}
	}
	return unsent, len(cell.Data)
}

func (c *Client) markDataAsSent(cell *types.DataCell, lim int) {
	if (cell == nil) || (lim == -1) {
		return
	}
	c.m.Lock()
	defer c.m.Unlock()
	for i := 0; i < lim; i++ {
		if _, sent := cell.Index[i]; !sent {
			cell.Index[i] = struct{}{}
		}
	}
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

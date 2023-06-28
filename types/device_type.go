package types

import (
	"strings"
)

type DeviceType string

const (
	DeviceTypeDesktop DeviceType = "DESKTOP"
	DeviceTypePhone   DeviceType = "PHONE"
	DeviceTypeTablet  DeviceType = "TABLET"
)

type Device struct {
	Type DeviceType
}

func (device Device) QueryEncode() string {
	var b strings.Builder
	b.WriteString("eventType=staticData&deviceType=")
	b.WriteString(string(device.Type))
	b.WriteString("&nonce=")
	b.WriteString(GetNonce())
	return b.String()
}

func (device Device) DataType() DataType {
	return DataTypeDevice
}

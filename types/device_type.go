package types

import "github.com/Kameleoon/client-go/v2/network"

type DeviceType string

const (
	DeviceTypeDesktop DeviceType = "DESKTOP"
	DeviceTypePhone   DeviceType = "PHONE"
	DeviceTypeTablet  DeviceType = "TABLET"
)

const deviceEventType = "staticData"

type Device struct {
	Type DeviceType
}

func (device Device) QueryEncode() string {
	qb := network.NewQueryBuilder()
	qb.Append(network.QPEventType, deviceEventType)
	qb.Append(network.QPDeviceType, string(device.Type))
	qb.Append(network.QPNonce, network.GetNonce())
	return qb.String()
}

func (device Device) DataType() DataType {
	return DataTypeDevice
}

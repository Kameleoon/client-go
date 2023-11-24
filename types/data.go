package types

type BaseData interface {
	DataType() DataType
}

type Data interface {
	BaseData

	dataRestriction()
}

type DataType string

const (
	DataTypeAssignedVariation DataType = "EXPERIMENT"
	DataTypeCustom            DataType = "CUSTOM"
	DataTypeBrowser           DataType = "BROWSER"
	DataTypeConversion        DataType = "CONVERSION"
	DataTypeDevice            DataType = "DEVICE"
	DataTypePageView          DataType = "PAGE_VIEW"
	DataTypeUserAgent         DataType = "USER_AGENT"
)

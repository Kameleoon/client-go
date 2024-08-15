package types

type UniqueIdentifier struct {
	value bool
}

func NewUniqueIdentifier(value bool) *UniqueIdentifier {
	return &UniqueIdentifier{
		value: value,
	}
}

func (ui *UniqueIdentifier) dataRestriction() {
	// This method is required to separate external type `Data` from `BaseData` types
}

func (ui *UniqueIdentifier) Value() bool {
	return ui.value
}

func (ui *UniqueIdentifier) DataType() DataType {
	return DataTypeUniqueIdentifier
}

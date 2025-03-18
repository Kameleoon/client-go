package types

import "fmt"

type Personalization struct {
	id          int
	variationId int
}

func NewPersonalization(id int, variationId int) *Personalization {
	return &Personalization{id: id, variationId: variationId}
}

func (p *Personalization) Id() int {
	return p.id
}

func (p *Personalization) VariationId() int {
	return p.variationId
}

func (*Personalization) DataType() DataType {
	return DataTypePersonalization
}

func (p Personalization) String() string {
	return fmt.Sprintf("Personalization{id:%d,variationId:%d}", p.id, p.variationId)
}

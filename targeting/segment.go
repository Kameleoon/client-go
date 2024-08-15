package targeting

import (
	"strings"

	"github.com/Kameleoon/client-go/v3/types"
	"github.com/Kameleoon/client-go/v3/utils"
)

type Segment struct {
	ID   int
	Tree *Tree
	base *types.SegmentBase
}

func NewSegment(s *types.SegmentBase) *Segment {
	return &Segment{
		ID:   s.ID,
		Tree: NewTree(s.ConditionsData),
		base: s,
	}
}

func (s Segment) String() string {
	var b strings.Builder
	b.WriteString("\nSegment id: ")
	b.WriteString(utils.WritePositiveInt(s.ID))
	b.WriteByte('\n')
	tree := s.Tree.String()
	b.WriteString(tree)
	return b.String()
}

func (s *Segment) Data() *types.SegmentBase {
	return s.base
}

func (s *Segment) CheckTargeting(data types.TargetingDataGetter) bool {
	if s == nil || s.Tree == nil {
		return true
	}
	return s.Tree.CheckTargeting(data)
}

func (s *Segment) GetSegmentBase() *types.SegmentBase {
	return s.base
}

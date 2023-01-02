package targeting

import (
	"strings"

	"github.com/Kameleoon/client-go/v2/types"
	"github.com/Kameleoon/client-go/v2/utils"
)

type Segment struct {
	ID   int
	Tree *Tree
	s    *types.Segment
}

func NewSegment(s *types.Segment) *Segment {
	return &Segment{
		ID:   s.ID,
		Tree: NewTree(s.ConditionsData),
		s:    s,
	}
}

func (s Segment) String() string {
	var b strings.Builder
	b.WriteString("\nSegment id: ")
	b.WriteString(utils.WriteUint(s.ID))
	b.WriteByte('\n')
	tree := s.Tree.String()
	b.WriteString(tree)
	return b.String()
}

func (s Segment) Data() *types.Segment {
	return s.s
}

func (s *Segment) CheckTargeting(data GetTargetingData) bool {
	if s == nil || s.Tree == nil {
		return true
	}
	return s.Tree.CheckTargeting(data)
}

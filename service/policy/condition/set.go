package condition

import (
	"github.com/rmrobinson/jvs/service/policy"
)

type Operator int
const (
	Unknown Operator = iota
	Or
	And
)

type set struct {
	conditions []policy.Condition
	operator Operator
}

func newSet(conditions []policy.Condition, operator Operator) *set {
	return &set{
		conditions: conditions,
		operator: operator,
	}
}

func (s *set) Name() string {
	name := "("

	switch s.operator {
	case And:
		name += "AND"
	case Or:
		name += "OR"
	default:
		name += "TRUE"
	}

	for _, condition := range s.conditions {
		name += " " + condition.Name()
	}

	name += ")"
	return name
}

func (s *set) Active() bool {
	switch s.operator {
	case And:
		for _, condition := range s.conditions {
			if !condition.Active() {
				return false
			}
		}
		return true
	case Or:
		for _, condition := range s.conditions {
			if condition.Active() {
				return true
			}
		}
		return false
	default:
		return true
	}
}
func (s *set) Refresh(state policy.State) {
	for _, condition := range s.conditions {
		condition.Refresh(state)
	}
}

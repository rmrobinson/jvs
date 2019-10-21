package condition

import (
	"errors"

	"github.com/rmrobinson/jvs/service/policy"
	"github.com/rmrobinson/jvs/service/policy/pb"
)

var (
	ErrInvalidConditionState = errors.New("invalid condition state")
)

type Condition struct {
	updates chan bool
}

func New(pbc *pb.Condition) (policy.Condition, error) {
	var c policy.Condition
	var err error
	if pbc.Cron != nil {
		c, err = newCron(pbc.Cron)
		if err != nil {
			return nil, err
		}
	}
	if pbc.Set != nil {
		var ncs []policy.Condition
		for _, pse := range pbc.Set.Conditions {
			nc, err := New(pse)
			if err != nil {
				return nil, err
			}
			ncs = append(ncs, nc)
		}
		if pbc.Set.Operator == pb.Condition_Set_And {
			c = newSet(ncs, And)
		} else if pbc.Set.Operator == pb.Condition_Set_Or {
			c = newSet(ncs, Or)
		} else {
			return nil, ErrInvalidConditionState
		}
	}

	return c, nil
}
package policy

import (
"errors"

	building "github.com/rmrobinson/jvs/service/building/pb"
	"github.com/rmrobinson/jvs/service/policy/pb"
)

var (
	ErrNotImplemented = errors.New("not implemented")
)

var ExamplePolicySet = &pb.PolicySet{
	Policies: []*pb.Policy{
		{
			Name: "Test policy 1",
			Weight: 1,
			Condition: &pb.Condition{
				Name: "Multi condition",
				Set: &pb.Condition_Set{
					Operator: pb.Condition_Set_And,
					Conditions: []*pb.Condition{
						{
							Name: "Christmas 2018",
							Cron: &pb.Condition_Cron{
								Entry: "* * * * * *",
							},
						},
					},
				},
			},
			Actions: []*pb.Action{
				{
					Name: "Log action",
					Type: pb.Action_Log,
				},
			},
		},
	},
}

type Condition interface {
	Name() string
	Refresh(State)

	Active() bool
}

type Action interface {
	Execute() error
}

type Engine struct {
	policies map[string]*Policy

	timers map[string]string

	currState State
}

func NewEngine() *Engine {
	return &Engine{
		policies: map[string]*Policy{},
		timers: map[string]string{},
	}
}

type State struct {
	deviceState map[string]*building.Device

}

func (e *Engine) AddPolicy(policy *Policy) {
	e.policies[policy.name] = policy
}

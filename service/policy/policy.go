package policy

type Policy struct {
	name string
	active bool
	lastErr error

	condition Condition
	actions []Action
}

func NewPolicy(name string, condition Condition, actions []Action) *Policy {
	return &Policy{
		name: name,
		condition: condition,
		actions: actions,
	}
}

func (p *Policy) Active() bool {
	return p.active
}

func (p *Policy) Evaluate(state State) {
	p.condition.Refresh(state)

	active := p.condition.Active()
	if p.condition.Active() != p.active {
		p.active = active

		// TODO: log this
		if active {
			for _, action := range p.actions {
				if err := action.Execute(); err != nil {
					p.lastErr = err
					return
				}
			}
		}
	}
}
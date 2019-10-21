package condition

import (
	"time"

	"github.com/rmrobinson/jvs/service/policy"
	"github.com/rmrobinson/jvs/service/policy/pb"
	crontab "github.com/robfig/cron"
)

type cron struct {
	c *crontab.Cron

	active bool
}

func newCron(condition *pb.Condition_Cron) (*cron, error) {
	var loc *time.Location
	var err error
	if len(condition.Tz) > 0 {
		loc, err = time.LoadLocation(condition.Tz)
		if err != nil {
			return nil, err
		}
	} else {
		loc = time.Local
	}
	ret := &cron{
		c: crontab.NewWithLocation(loc),
	}
	err = ret.c.AddFunc(condition.Entry, ret.execute)
	if err != nil {
		return nil, err
	}

	ret.c.Start()

	return ret, nil
}

func (c *cron) Name() string {
	return "cron"
}

func (c *cron) Active() bool {
	return c.active
}

func (c *cron) Refresh(state policy.State) {
}

func (c *cron) execute() {
	c.active = true
}

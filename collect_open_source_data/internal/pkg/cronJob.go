package pkg

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
)

const (
	specSecond = "*/%d * * * * ?"
	specMinute = "0 */%d * * * ?"
	specHour   = "0 0 */%d * * ?"
	specDay    = "0 0 0 */%d * ?"
)

type CronJob struct {
	Cron    *cron.Cron
	Content context.Context
}

func NewCronJob() *CronJob {
	return &CronJob{
		Cron:    cron.New(cron.WithSeconds()),
		Content: context.TODO(),
	}
}

func (c *CronJob) AddJob(spec string, cmd cron.FuncJob) (int, error) {
	cmd()
	entryID, err := c.Cron.AddFunc(spec, cmd)
	return int(entryID), err
}

func (c *CronJob) AddJobWithSecond(t int, cmd cron.FuncJob) (int, error) {
	if t < 1 || t > 59 {
		t = 1
	}
	return c.AddJob(fmt.Sprintf(specSecond, t), cmd)
}

func (c *CronJob) AddJobWithMinute(t int, cmd cron.FuncJob) (int, error) {
	if t < 1 || t > 59 {
		t = 1
	}
	return c.AddJob(fmt.Sprintf(specMinute, t), cmd)
}

func (c *CronJob) AddJobWithHour(t int, cmd cron.FuncJob) (int, error) {
	if t < 1 || t > 23 {
		t = 1
	}
	return c.AddJob(fmt.Sprintf(specHour, t), cmd)
}

func (c *CronJob) AddJobWithDay(t int, cmd cron.FuncJob) (int, error) {
	if t < 1 || t > 31 {
		t = 1
	}
	return c.AddJob(fmt.Sprintf(specDay, t), cmd)
}

func (c *CronJob) RemoveJobWithSpec(spec string) {
}

func (c *CronJob) RemoveJob(entryID int) {
	c.Cron.Remove(cron.EntryID(entryID))
}

func (c *CronJob) Start() {
	c.Cron.Start()
}

func (c *CronJob) Stop() {
	c.Cron.Stop()
}

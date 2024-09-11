package server

import (
	"collect_open_source_data/internal/biz"
	"collect_open_source_data/internal/pkg"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type TaskServer struct {
	openS    *biz.OpenSourceInfo
	log      *log.Helper
	timerJob *pkg.CronJob
}

func NewTaskServer(info *biz.OpenSourceInfo, logger log.Logger) *TaskServer {
	return &TaskServer{
		openS:    info,
		log:      log.NewHelper(logger),
		timerJob: pkg.NewCronJob(),
	}
}

func (t *TaskServer) Start(ctx context.Context) error {
	if _, err := t.timerJob.AddJobWithHour(1, t.openS.Collect); err != nil {
		t.log.Errorf("add job failed: %v", err)
		return err
	}
	t.timerJob.Start()
	return nil
}
func (t *TaskServer) Stop(ctx context.Context) error {
	t.timerJob.Stop()
	return nil
}

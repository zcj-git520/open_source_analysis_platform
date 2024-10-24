package server

import (
	"collect_open_source_data/internal/biz"
	"collect_open_source_data/internal/conf"
	"collect_open_source_data/internal/pkg"
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

type TaskServer struct {
	openS    *biz.OpenSourceInfo
	ccf      *conf.Collect
	log      *log.Helper
	timerJob *pkg.CronJob
	cancel   context.CancelFunc
}

func NewTaskServer(info *biz.OpenSourceInfo, collect *conf.Collect, logger log.Logger) *TaskServer {
	return &TaskServer{
		openS:    info,
		ccf:      collect,
		log:      log.NewHelper(logger),
		timerJob: pkg.NewCronJob(),
	}
}

func (t *TaskServer) Start(ctx context.Context) error {
	if !t.ccf.Enable {
		return nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	t.openS.Ctx = ctx
	t.cancel = cancel
	if _, err := t.timerJob.AddJobWithDay(1, t.openS.Collect); err != nil {
		t.log.Errorf("add job failed: %v", err)
		return err
	}
	t.timerJob.Start()
	return nil
}
func (t *TaskServer) Stop(ctx context.Context) error {
	if !t.ccf.Enable {
		return nil
	}
	t.timerJob.Stop()
	t.cancel()
	return nil
}

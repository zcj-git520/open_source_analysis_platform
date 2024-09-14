package data

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"user/internal/biz"
	"user/internal/conf"
)

type casbinRuleRepo struct {
	data           *Data
	log            *log.Helper
	syncedEnforcer *casbin.SyncedEnforcer
	ModelPath      string
}

// 用户权限
func NewCasbinRuleRepo(data *Data, db *gorm.DB, conf *conf.Casbin, logger log.Logger) biz.CasbinRuleRepo {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		panic(err)
	}
	enforcer, err := casbin.NewSyncedEnforcer(conf.Model, adapter)
	if err != nil {
		panic(err)
	}
	if err = enforcer.LoadPolicy(); err != nil {
		panic(err)
	}
	return &casbinRuleRepo{
		data:           data,
		log:            log.NewHelper(logger),
		syncedEnforcer: enforcer,
		ModelPath:      conf.Model,
	}

}

func (c *casbinRuleRepo) GetModel() model.Model {
	return c.syncedEnforcer.GetModel()
}

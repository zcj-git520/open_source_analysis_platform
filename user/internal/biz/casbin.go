package biz

import (
	"github.com/casbin/casbin/v2/model"
	"github.com/go-kratos/kratos/v2/log"
)

type CasbinRuleRepo interface {
	GetModel() model.Model
	//GetAdapter() persist.Adapter
	//
	//UpdateCasbin(ctx context.Context, roleKey string, p [][]string) error
	//ClearCasbin(v int, p ...string) error
	//UpdateCasbinApi(ctx context.Context, oldPath string, newPath string, oldMethod string, newMethod string) error
	//GetPolicyPathByRoleId(roleKey string) [][]string
}

type CasbinRule struct {
	repo CasbinRuleRepo
	log  *log.Helper
}

func NewCasbinRuleUseCase(repo CasbinRuleRepo, logger log.Logger) *CasbinRule {
	return &CasbinRule{repo: repo, log: log.NewHelper(logger)}
}

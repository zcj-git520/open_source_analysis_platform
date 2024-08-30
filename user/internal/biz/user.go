package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	pb "user/api/user/v1"
)

type UserRepo interface {
	FindUserByEmail(ctx context.Context, email string) (*UserInfo, error)
	InsertUser(ctx context.Context, user *UserInfo) (int64, error)
}

type UserInfo struct {
	ID        int64  `json:"id"`
	Uid       int64  `json:"uid"`
	Status    int    `json:"status"`
	Nickname  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
	Gender    int    `json:"gender"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

type User struct {
	repo UserRepo
	log  *log.Helper
}

func NewUser(repo UserRepo, logger log.Logger) *User {
	return &User{repo: repo, log: log.NewHelper(logger)}
}

func (u *User) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	// 判断用户邮箱是否存在
	user, err := u.repo.FindUserByEmail(ctx, req.Email)
	// 用户已存在
	if err == nil && user != nil && user.Uid > 0 {
		return nil, fmt.Errorf("user already exists")
	}
	uid, err := u.repo.InsertUser(ctx, &UserInfo{
		Nickname: req.Nickname,
		Password: req.Password,
		Email:    req.Email,
		Avatar:   req.Avatar,
		Gender:   1,
	})
	return &pb.RegisterReply{
		Uid: uid,
	}, err

}

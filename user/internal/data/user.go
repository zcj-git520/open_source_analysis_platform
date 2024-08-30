package data

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"time"
	"user/internal/biz"
	"user/internal/pkg"
)

type UserInfo struct {
	ID        int64  `gorm:"primarykey;type:int" json:"id"`
	Uid       int64  `gorm:"type:int" json:"uid"`
	Status    int    `gorm:"type:TINYINT" json:"status"`
	Nickname  string `gorm:"type:varchar(50)" json:"username"`
	Password  string `gorm:"type:varchar(50)" json:"password"`
	Email     string `gorm:"type:varchar(50)" json:"email"`
	Avatar    string `gorm:"type:varchar(255)" json:"avatar"`
	Gender    int    `gorm:"type:TINYINT" json:"gender"`
	CreatedAt int64  `gorm:"type:int" json:"created_at"`
	UpdatedAt int64  `gorm:"type:int" json:"updated_at"`
}

func (UserInfo) TableName() string {
	return "user_info"
}

type userRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (u *userRepo) FindUserByEmail(ctx context.Context, email string) (*biz.UserInfo, error) {
	user := &UserInfo{}
	if err := u.data.db.Where("email = ? and status = 0", email).Find(user).Error; err != nil {
		return nil, err
	}
	return &biz.UserInfo{
		ID:        user.ID,
		Uid:       user.Uid,
		Status:    user.Status,
		Nickname:  user.Nickname,
		Password:  user.Password,
		Email:     user.Email,
		Avatar:    user.Avatar,
		Gender:    user.Gender,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

func (u *userRepo) InsertUser(ctx context.Context, user *biz.UserInfo) (int64, error) {
	sn := pkg.NewSnowflake(1, 1)
	uid := sn.NextId()
	if uid < 0 {
		return uid, fmt.Errorf("生成uid失败")
	}
	return uid, u.data.db.Create(&UserInfo{
		Uid:       uid,
		Status:    0,
		Nickname:  user.Nickname,
		Password:  user.Password,
		Email:     user.Email,
		Avatar:    user.Avatar,
		Gender:    user.Gender,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}).Error
}

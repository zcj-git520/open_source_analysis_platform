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
	ID        int64     `gorm:"primarykey;type:int" json:"id"`
	Uid       int64     `gorm:"type:int" json:"uid"`
	Phone     string    `gorm:"type:varchar(20)" json:"phone"`
	Nickname  string    `gorm:"type:varchar(50)" json:"nickname"`
	Username  string    `gorm:"type:varchar(50)" json:"username"`
	Password  string    `gorm:"type:varchar(255)" json:"password"`
	Email     string    `gorm:"type:varchar(50)" json:"email"`
	Avatar    string    `gorm:"type:varchar(255)" json:"avatar"`
	Role      []string  `gorm:"type:varchar(255)" json:"role"`
	Gender    int       `gorm:"type:TINYINT" json:"gender"`
	Status    int       `gorm:"type:TINYINT" json:"status"`
	CreatedAt time.Time `gorm:"column:created_at;type:datetime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;type:datetime" json:"updated_at"`
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

func (u *userRepo) changeType(ctx context.Context, user *UserInfo) *biz.UserInfo {
	return &biz.UserInfo{
		ID:        user.ID,
		Uid:       user.Uid,
		Status:    user.Status,
		Nickname:  user.Nickname,
		Username:  user.Username,
		Password:  user.Password,
		Email:     user.Email,
		Avatar:    user.Avatar,
		Gender:    user.Gender,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (u *userRepo) FindUserByEmail(ctx context.Context, email string) (*biz.UserInfo, error) {
	user := &UserInfo{}
	if err := u.data.db.Where("email = ? and status = 0", email).Find(user).Error; err != nil {
		return nil, err
	}
	return u.changeType(ctx, user), nil
}

func (u *userRepo) FindUserByName(ctx context.Context, name string) (*biz.UserInfo, error) {
	user := &UserInfo{}
	if err := u.data.db.Where("username = ? and status = 0", name).Find(user).Error; err != nil {
		return nil, err
	}
	return u.changeType(ctx, user), nil
}

func (u *userRepo) FindUserByPhone(ctx context.Context, phone string) (*biz.UserInfo, error) {
	user := &UserInfo{}
	if err := u.data.db.Where("phone = ? and status = 0", phone).Find(user).Error; err != nil {
		return nil, err
	}
	return u.changeType(ctx, user), nil
}

func (u *userRepo) FindUserByUid(ctx context.Context, uid int64) (*biz.UserInfo, error) {
	user := &UserInfo{}
	if err := u.data.db.Where("uid = ? and status = 0", uid).Find(user).Error; err != nil {
		return nil, err
	}
	return u.changeType(ctx, user), nil
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
		Username:  user.Username,
		Password:  user.Password,
		Email:     user.Email,
		Avatar:    user.Avatar,
		Gender:    user.Gender,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}).Error
}

func (u *userRepo) UpdateUser(ctx context.Context, user *biz.UserInfo) error {
	userInfo := &UserInfo{UpdatedAt: time.Now()}
	if user.Nickname != "" {
		userInfo.Nickname = user.Nickname
	}
	if user.Password != "" {
		userInfo.Password = user.Password
	}
	if user.Email != "" {
		userInfo.Email = user.Email
	}
	if user.Avatar != "" {
		userInfo.Avatar = user.Avatar
	}
	if user.Gender != 0 {
		userInfo.Gender = user.Gender
	}
	return u.data.db.Model(&UserInfo{}).Where("uid = ?", user.Uid).Updates(userInfo).Error
}

func (u *userRepo) SetVerifyCodeCache(ctx context.Context, email, code string, expTime time.Duration) error {
	return u.data.rdb.Set(ctx, fmt.Sprintf("user_verify_code_%s", email), code, expTime).Err()
}

func (u *userRepo) GetVerifyCodeCache(ctx context.Context, email string) (string, error) {
	return u.data.rdb.Get(ctx, fmt.Sprintf("user_verify_code_%s", email)).Result()
}

func (u *userRepo) DeleteVerifyCodeCache(ctx context.Context, email string) error {
	return u.data.rdb.Del(ctx, fmt.Sprintf("user_verify_code_%s", email)).Err()
}

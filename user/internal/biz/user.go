package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"math/rand"
	"time"
	pb "user/api/user/v1"
	"user/internal/conf"
	"user/internal/pkg"
	"user/internal/pkg/auth"
)

const (
	VerificationCodeSubject = "OSAP 注册验证码"
	VerificationCodeMsg     = `<td style="font-size:14px;color:#333;padding:24px 40px 0 40px">
                尊敬的用户您好！
                <br>
                <br>
                您的注册验证码是：<b>%d</b>，请在<b>5分钟内</b>进行验证, 过期将失效!
                <br> 
                如果该验证码不为您本人申请，请无视。
            </td>`
	Issuer   = "osap"
	Audience = "osap"
	jwtKey   = "osap-jwt"
)

type UserRepo interface {
	FindUserByEmail(ctx context.Context, email string) (*UserInfo, error)
	FindUserByPhone(ctx context.Context, phone string) (*UserInfo, error)
	InsertUser(ctx context.Context, user *UserInfo) (int64, error)
	SetVerifyCodeCache(ctx context.Context, email, code string, expTime time.Duration) error
	GetVerifyCodeCache(ctx context.Context, email string) (string, error)
	DeleteVerifyCodeCache(ctx context.Context, email string) error
	FindUserByName(ctx context.Context, name string) (*UserInfo, error)
	FindUserByUid(ctx context.Context, uid int64) (*UserInfo, error)
	UpdateUser(ctx context.Context, user *UserInfo) error
}

type UserInfo struct {
	ID        int64     `json:"id"`
	Uid       int64     `json:"uid"`
	Status    int       `json:"status"`
	Desc      string    `json:"desc"`
	Username  string    `json:"username"`
	Nickname  string    `json:"nickname"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	Avatar    string    `json:"avatar"`
	Gender    int       `json:"gender"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	repo UserRepo
	log  *log.Helper
	ac   *conf.Auth
	ec   *conf.Email
}

func NewUser(repo UserRepo, logger log.Logger, ac *conf.Auth, ec *conf.Email) *User {
	return &User{repo: repo, ac: ac, ec: ec, log: log.NewHelper(logger)}
}

func (u *User) createToken(ctx context.Context, uid int64) (string, string, error) {
	user, err := u.repo.FindUserByUid(ctx, uid)
	// 用户已存在
	if err != nil || user == nil || user.Uid == 0 {
		return "", "", fmt.Errorf("user not found")
	}
	expr := time.Now().Add(24 * time.Hour)
	token, err := auth.CreateToken(auth.Claims{
		Uid:      user.Uid,
		Nickname: user.Nickname,
		Email:    user.Email,
		Avatar:   user.Avatar,
		Gender:   user.Gender,
		Phone:    user.Phone,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Issuer,
			Subject:   Audience,
			ExpiresAt: jwt.NewNumericDate(expr),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}, u.ac.JwtKey)
	if err != nil {
		return "", "", fmt.Errorf("token create error")
	}
	return token, expr.Format("2006/01/02 15:04:05"), nil
}

func (u *User) sendVerificationCodeByEmail(ctx context.Context, to string) error {
	email := pkg.NewEmailSMTP(pkg.WithSmtpHost(u.ec.SmtpHost), pkg.WithSmtpPort(int(u.ec.SmtpPort)),
		pkg.WithSmtpUsername(u.ec.SmtpUsername), pkg.WithSmtpPassword(u.ec.SmtpPassword),
		pkg.WithFrom(u.ec.From),
		pkg.WithTo([]string{to}))
	// 生成验证码
	rand.Seed(time.Now().UnixNano())
	verificationCode := rand.Intn(999999) + 100000
	// 缓存在redis 中
	if err := u.repo.SetVerifyCodeCache(ctx, to, pkg.ToString(verificationCode), 5*time.Minute); err != nil {
		return err
	}
	// 邮件内容
	return email.SendEmailSMTP(VerificationCodeSubject, fmt.Sprintf(VerificationCodeMsg, verificationCode))
}

func (u *User) VerifyCode(ctx context.Context, req *pb.VerifyRequest) (*pb.VerifyReply, error) {
	// 检查邮箱是否合法
	if !pkg.VerifyEmailFormat(req.Email) {
		return nil, fmt.Errorf("email format error")
	}
	// 发送验证码
	result := &pb.VerifyReply{
		Success: false,
	}
	if err := u.sendVerificationCodeByEmail(ctx, req.Email); err != nil {
		return result, err
	}
	result.Success = true
	return result, nil
}

func (u *User) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	// 判断用户邮箱是否存在
	user, err := u.repo.FindUserByEmail(ctx, req.Email)
	// 用户已存在
	if err == nil && user != nil && user.Uid > 0 {
		return nil, fmt.Errorf("user already exists")
	}
	// 判断用户收集号是否存在
	user, err = u.repo.FindUserByPhone(ctx, req.Phone)
	if err == nil && user != nil && user.Uid > 0 {
		return nil, fmt.Errorf("user already exists")
	}
	code, err := u.repo.GetVerifyCodeCache(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("verification code empty")
	}
	if code != req.VerificationCode {
		return nil, fmt.Errorf("verification code error")
	}
	_ = u.repo.DeleteVerifyCodeCache(ctx, req.Email)

	pwd, err := pkg.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("password hash error")
	}
	_, err = u.repo.InsertUser(ctx, &UserInfo{
		Username: req.Username,
		Password: pwd,
		Email:    req.Email,
	})
	return &pb.RegisterReply{
		Success: true,
	}, err

}

func (u *User) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	// 判断用户名是否存在
	user, err := u.repo.FindUserByName(ctx, req.Username)
	if err != nil || user == nil || user.Uid == 0 {
		return nil, fmt.Errorf("user not exists")
	}
	if !pkg.CheckPassword(user.Password, req.Password) {
		return nil, fmt.Errorf("password error")
	}
	token, expStr, err := u.createToken(ctx, user.Uid)
	if err != nil {
		return nil, fmt.Errorf("create token error")
	}
	return &pb.LoginReply{
		Data: &pb.LoginReply_Data{
			Uid:          user.Uid,
			AccessToken:  token,
			Avatar:       user.Avatar,
			Nickname:     user.Nickname,
			Username:     user.Username,
			RefreshToken: token,
			Expires:      expStr,
			Phone:        user.Phone,
			Email:        user.Email,
			Gender:       int32(user.Gender),
			Desc:         user.Desc,
		},

		Success: true,
	}, nil
}

func (u *User) Update(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserReply, error) {
	// 根据用户uid，获取用户信息
	uid := auth.GetUid(ctx)
	if uid == 0 {
		u.log.Errorf("user not login: uid:%v", uid)
		return nil, fmt.Errorf("user not login")
	}
	// 查询uid是否存在
	user, err := u.repo.FindUserByUid(ctx, uid)
	if err != nil || user == nil || user.Uid == 0 {
		return nil, fmt.Errorf("user not exists")
	}
	userInfo := &UserInfo{
		Uid:      user.Uid,
		Desc:     req.Desc,
		Status:   int(req.Status),
		Nickname: req.Nickname,
		Password: req.Password,
		Avatar:   req.Avatar,
		Gender:   int(req.Gender),
		Phone:    req.Phone,
	}
	// 更新用户信息
	if req.Password != "" {
		if !pkg.CheckPassword(user.Password, req.Password) {
			pwd, _ := pkg.HashPassword(req.Password)
			userInfo.Password = pwd
		}
	}
	if err = u.repo.UpdateUser(ctx, userInfo); err != nil {
		return &pb.UpdateUserReply{}, fmt.Errorf("update user error:%v", err)
	}
	token, expStr, err := u.createToken(ctx, user.Uid)
	if err != nil {
		return nil, fmt.Errorf("create token error")
	}
	return &pb.UpdateUserReply{
		Data: &pb.UpdateUserReply_Data{
			Uid:          user.Uid,
			AccessToken:  token,
			Avatar:       user.Avatar,
			Nickname:     user.Nickname,
			Username:     user.Username,
			RefreshToken: token,
			Expires:      expStr,
			Phone:        user.Phone,
			Email:        user.Email,
			Gender:       int32(user.Gender),
			Desc:         user.Desc,
		},
		Success: true,
	}, nil
}

func (u *User) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	uid := auth.GetUid(ctx)
	if uid == 0 {
		u.log.Errorf("user not login: uid:%v", uid)
		return nil, fmt.Errorf("user not login")
	}
	user, err := u.repo.FindUserByUid(ctx, uid)
	if err != nil || user == nil || user.Uid == 0 {
		return nil, fmt.Errorf("user not exists")
	}
	return &pb.GetUserReply{
		Data: &pb.GetUserReply_Data{
			Uid:      user.Uid,
			Username: user.Username,
			Email:    user.Email,
			Phone:    user.Phone,
			Avatar:   user.Avatar,
			Gender:   int32(user.Gender),
			Nickname: user.Nickname,
			Desc:     user.Desc,
		},
		Success: true,
	}, nil
}

func (u *User) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserReply, error) {
	// 判断验证码是否正确
	code, err := u.repo.GetVerifyCodeCache(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("verification code empty")
	}
	if code != req.VerificationCode {
		return nil, fmt.Errorf("verification code error")
	}
	_ = u.repo.DeleteVerifyCodeCache(ctx, req.Email)
	// 获取用户是否存在
	uid := auth.GetUid(ctx)
	if uid == 0 {
		u.log.Errorf("user not login: uid:%v", uid)
		return nil, fmt.Errorf("user not login")
	}
	user, err := u.repo.FindUserByUid(ctx, uid)
	if err != nil || user == nil || user.Uid == 0 {
		return nil, fmt.Errorf("user not exists")
	}
	if err = u.repo.UpdateUser(ctx, &UserInfo{
		Status: -1,
	}); err != nil {
		return nil, fmt.Errorf("delete user error:%v", err)
	}
	return &pb.DeleteUserReply{Success: true}, nil
}

func (u *User) RefreshToken(ctx context.Context) (*pb.RefreshTokenReply, error) {
	// 根据用户uid，获取用户信息
	uid := auth.GetUid(ctx)
	if uid == 0 {
		u.log.Errorf("user not login: uid:%v", uid)
		return nil, fmt.Errorf("user not login")
	}
	// 查询uid是否存在
	user, err := u.repo.FindUserByUid(ctx, uid)
	if err != nil || user == nil || user.Uid == 0 {
		return nil, fmt.Errorf("user not exists")
	}
	token, expStr, err := u.createToken(ctx, user.Uid)
	return &pb.RefreshTokenReply{
		Data: &pb.RefreshTokenReply_Data{
			AccessToken:  token,
			RefreshToken: token,
			Expires:      expStr,
		},
		Success: true,
	}, nil
}

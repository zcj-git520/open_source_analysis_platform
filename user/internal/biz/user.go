package biz

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"math/rand"
	"time"
	pb "user/api/user/v1"
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
}

type UserInfo struct {
	ID        int64     `json:"id"`
	Uid       int64     `json:"uid"`
	Status    int       `json:"status"`
	Nickname  string    `json:"username"`
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
}

func NewUser(repo UserRepo, logger log.Logger) *User {
	return &User{repo: repo, log: log.NewHelper(logger)}
}

func (u *User) sendVerificationCodeByEmail(ctx context.Context, to string) error {
	email := pkg.NewEmailSMTP(pkg.WithSmtpHost("smtp.qq.com"), pkg.WithSmtpPort(587),
		pkg.WithSmtpUsername("osap.work@qq.com"), pkg.WithSmtpPassword(""),
		pkg.WithFrom("osap.work@qq.com"),
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
	uid, err := u.repo.InsertUser(ctx, &UserInfo{
		Nickname: req.Nickname,
		Password: pwd,
		Email:    req.Email,
	})
	return &pb.RegisterReply{
		Uid: uid,
	}, err

}

func (u *User) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	// 判断用户名是否存在
	user, err := u.repo.FindUserByName(ctx, req.Nickname)
	if err != nil || user == nil || user.Uid == 0 {
		return nil, fmt.Errorf("user not exists")
	}
	if !pkg.CheckPassword(user.Password, req.Password) {
		return nil, fmt.Errorf("password error")
	}
	token, err := auth.CreateToken(&auth.Claims{
		Uid:      user.Uid,
		Nickname: user.Nickname,
		Email:    user.Email,
		Avatar:   user.Avatar,
		Gender:   user.Gender,
		Phone:    user.Phone,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Issuer,
			Subject:   Audience,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}, jwtKey)
	if err != nil {
		return nil, fmt.Errorf("token create error")
	}
	return &pb.LoginReply{
		Token: token,
		Uid:   user.Uid,
	}, nil
}

package service

import (
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
	pb "user/api/user/v1"
	"user/internal/biz"
)

type UserService struct {
	pb.UnimplementedUserServer
	uc *biz.User
}

func NewUserService(uc *biz.User) *UserService {
	return &UserService{uc: uc}
}

func (s *UserService) Verify(ctx context.Context, req *pb.VerifyRequest) (*pb.VerifyReply, error) {
	return s.uc.VerifyCode(ctx, req)
}

func (s *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	return s.uc.Register(ctx, req)
}
func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	return s.uc.Login(ctx, req)
}
func (s *UserService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserReply, error) {
	return s.uc.Update(ctx, req)
}

func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	return s.uc.GetUser(ctx, req)
}
func (s *UserService) ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserReply, error) {
	return &pb.ListUserReply{}, nil
}

func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserReply, error) {
	return s.uc.DeleteUser(ctx, req)
}

func (s UserService) RefreshToken(ctx context.Context, req *emptypb.Empty) (*pb.RefreshTokenReply, error) {
	return s.uc.RefreshToken(ctx)
}

package service

import (
	"context"
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
	return &pb.UpdateUserReply{}, nil
}
func (s *UserService) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserReply, error) {
	return &pb.DeleteUserReply{}, nil
}
func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserReply, error) {
	return &pb.GetUserReply{}, nil
}
func (s *UserService) ListUser(ctx context.Context, req *pb.ListUserRequest) (*pb.ListUserReply, error) {
	return &pb.ListUserReply{}, nil
}

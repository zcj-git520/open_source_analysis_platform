package service

import (
	pb "collect_open_source_data/api/open_source/v1"
	"collect_open_source_data/internal/biz"
	"context"
)

type OpenSourceService struct {
	pb.UnimplementedOpenSourceServer
	uc *biz.OpenSourceInfo
}

func NewOpenSourceService(uc *biz.OpenSourceInfo) *OpenSourceService {
	return &OpenSourceService{
		uc: uc,
	}
}

func (s *OpenSourceService) GetLanguage(ctx context.Context, req *pb.LanguageRequest) (*pb.LanguageReply, error) {
	return s.uc.GetLanguage(ctx, req)
}

func (s *OpenSourceService) GetOwner(ctx context.Context, req *pb.OwnerRequest) (*pb.OwnerReply, error) {
	return s.uc.GetOwner(ctx, req)
}

func (s *OpenSourceService) GetRepo(ctx context.Context, req *pb.RepoRequest) (*pb.RepoReply, error) {
	return s.uc.GetRepo(ctx, req)
}

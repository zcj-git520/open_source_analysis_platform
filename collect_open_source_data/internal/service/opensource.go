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

func (s *OpenSourceService) GetRepoCategory(ctx context.Context, req *pb.RepoCategoryRequest) (*pb.RepoCategoryReply, error) {
	return s.uc.GetRepoCategory(ctx, req)
}

func (s *OpenSourceService) GetRepoByCategory(ctx context.Context, req *pb.RepoByCategoryRequest) (*pb.RepoByCategoryReply, error) {
	return s.uc.GetRepoByCategory(ctx, req)
}

func (s *OpenSourceService) GetRepoMeasure(ctx context.Context, req *pb.RepoMeasureRequest) (*pb.RepoMeasureReply, error) {
	return s.uc.GetRepoMeasure(ctx, req)
}

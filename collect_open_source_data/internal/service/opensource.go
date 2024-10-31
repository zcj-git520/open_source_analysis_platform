package service

import (
	pb "collect_open_source_data/api/open_source/v1"
	"collect_open_source_data/internal/biz"
	"context"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (s *OpenSourceService) RepoFav(ctx context.Context, req *pb.RepoFavRequest) (*pb.RepoFavReply, error) {
	return s.uc.RepoFav(ctx, req)
}

func (s *OpenSourceService) GetRepoFav(ctx context.Context, req *pb.RepoFavListRequest) (*pb.RepoReply, error) {
	return s.uc.GetRepoFav(ctx, req)
}

func (s *OpenSourceService) GetScreenLanguageCount(ctx context.Context, req *emptypb.Empty) (*pb.ScreenLanguageCountReply, error) {
	return s.uc.GetScreenLanguageCount(ctx, req)
}

func (s *OpenSourceService) GetScreenCategoryCount(ctx context.Context, req *emptypb.Empty) (*pb.ScreenCategoryCountReply, error) {
	return s.uc.GetScreenCategoryCount(ctx, req)
}

func (s *OpenSourceService) GetMessage(ctx context.Context, req *emptypb.Empty) (*pb.MessageReply, error) {
	return s.uc.GetMessage(ctx, req)
}

func (s *OpenSourceService) UpdateMessage(ctx context.Context, req *pb.UpdateMessageRequest) (*emptypb.Empty, error) {
	return s.uc.UpdateMessage(ctx, req)
}

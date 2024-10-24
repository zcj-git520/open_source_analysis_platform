package biz

import (
	pb "collect_open_source_data/api/open_source/v1"
	"collect_open_source_data/internal/conf"
	"collect_open_source_data/internal/domain"
	"collect_open_source_data/internal/pkg/auth"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

const (
	BaseGithubURL = "https://api.github.com/search/repositories"
)

const (
	DateTypeDay = iota
	DateTypeWeek
	DateTypeMonth
)

type OpenSourceRepo interface {
	InsertOwner(ctx context.Context, owner *domain.Owner) (int64, error)
	InsertRepo(ctx context.Context, repo *domain.RepoInfo) error
	FindOwnerByHtmlUrl(ctx context.Context, htmlUrl string) (*domain.Owner, error)
	FindRepoByName(ctx context.Context, name string) (*domain.RepoInfo, error)
	FindLanguage(ctx context.Context, name string, id int64, page *domain.Page) ([]*domain.Language, error)
	InsertLanguage(ctx context.Context) ([]*domain.Language, error)
	FindOwner(ctx context.Context, name, ownerType, email string, Id int64, page *domain.Page) ([]*domain.Owner, error)
	FindRepo(ctx context.Context, req *pb.RepoRequest, page *domain.Page) ([]*domain.RepoInfo, error)
	UpdateRepo(ctx context.Context, repo *domain.RepoInfo) error
	FindLanguageByCache(ctx context.Context, languageId int64) (*domain.Language, error)
	FindOwnerByCache(ctx context.Context, Id int64) (*domain.Owner, error)
	UpdateOwner(ctx context.Context, owner *domain.Owner) error
	FindRepoCategory(ctx context.Context, name string, id int64, page *domain.Page) ([]*domain.RepoCategory, error)
	FindRepoCategoryId(ctx context.Context, repoId, categoryId int64) bool
	AddRepoCategoryId(ctx context.Context, repoId, categoryId int64) error
	FindRepoCategoryIdByRepoId(repoId int64) bool
	FindRepoCategoryByCatId(ctx context.Context, id int64, page *domain.Page) ([]*domain.RepoCategoryId, error)
	FindRepoById(ctx context.Context, id int64) (*domain.RepoInfo, error)
	AddRepoMetrics(ctx context.Context, metrics []*domain.RepoMetrics) error
	FindRepoMetrics(ctx context.Context, metricType int, data string, page *domain.Page) ([]*domain.RepoMetricsResult, error)
	FindRepoFavor(ctx context.Context, id, uid, repoId int64, page *domain.Page) ([]*domain.RepoFav, error)
	AddRepoFavor(ctx context.Context, favorInfo *domain.RepoFav) error
	UpdateRepoFavor(ctx context.Context, favorId int64, isFavor int32) error
	UpdateRepoFaveCache(ctx context.Context, uid int64) ([]*domain.RepoFav, error)
	FindRepoFaveByCache(ctx context.Context, uid int64) ([]*domain.RepoFav, error)
}

type OpenSourceInfo struct {
	repo OpenSourceRepo
	log  *log.Helper
	ec   *conf.Email
	ccf  *conf.Collect
	Ctx  context.Context
	Page int
}

func NewOpenSourceInfo(repo OpenSourceRepo, email *conf.Email, ccf *conf.Collect, logger log.Logger) *OpenSourceInfo {
	return &OpenSourceInfo{
		repo: repo,
		ec:   email,
		ccf:  ccf,
		log:  log.NewHelper(logger),
	}
}

func (r *OpenSourceInfo) GetLanguage(ctx context.Context, req *pb.LanguageRequest) (*pb.LanguageReply, error) {
	page := &domain.Page{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
	}
	info, err := r.repo.FindLanguage(ctx, req.Name, req.ID, page)
	if err != nil {
		return nil, err
	}
	var data []*pb.LanguageInfo
	for _, item := range info {
		data = append(data, &pb.LanguageInfo{
			Id:          item.ID,
			Name:        item.Name,
			Description: item.Desc,
			ImageUrl:    item.Image,
			RepoUrl:     item.RepoRUL,
			Bio:         item.Bio,
		})
	}
	return &pb.LanguageReply{
		PageNum:   req.PageNum,
		PageSize:  req.PageSize,
		Total:     page.Total,
		Languages: data,
	}, nil

}

func (r *OpenSourceInfo) GetOwner(ctx context.Context, req *pb.OwnerRequest) (*pb.OwnerReply, error) {
	page := &domain.Page{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
	}
	info, err := r.repo.FindOwner(ctx, req.Name, req.Type, req.Email, req.ID, page)
	if err != nil {
		return nil, err
	}
	var data []*pb.OwnerInfo
	for _, item := range info {
		data = append(data, &pb.OwnerInfo{
			Id:          item.ID,
			AvatarUrl:   item.AvatarURL,
			Type:        item.Type,
			Login:       item.Login,
			HtmlUrl:     item.HtmlURL,
			Name:        item.Name,
			Email:       item.Email,
			Bio:         item.Bio,
			PublicRepos: item.PublicRepos,
			PublicGists: item.PublicGists,
			Followers:   item.Following,
			Following:   item.Followers,
			CreatedAt:   item.CreatedAt.Format(time.DateTime),
			UpdatedAt:   item.UpdatedAt.Format(time.DateTime),
		})
	}
	return &pb.OwnerReply{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
		Total:    page.Total,
		Owners:   data,
	}, nil
}

func (r *OpenSourceInfo) GetRepo(ctx context.Context, req *pb.RepoRequest) (*pb.RepoReply, error) {
	page := &domain.Page{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
	}
	info, err := r.repo.FindRepo(ctx, req, page)
	if err != nil {
		return nil, err
	}

	var data []*pb.RepoInfo
	for _, item := range info {
		data = append(data, r.repoData(ctx, item))
	}
	return &pb.RepoReply{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
		Total:    page.Total,
		Repos:    data,
	}, nil

}

func (r *OpenSourceInfo) GetRepoCategory(ctx context.Context, req *pb.RepoCategoryRequest) (*pb.RepoCategoryReply, error) {
	page := &domain.Page{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
	}

	info, err := r.repo.FindRepoCategory(ctx, req.Name, req.ID, page)
	if err != nil {
		return nil, err
	}
	var data []*pb.RepoCategoryInfo
	for _, item := range info {
		data = append(data, &pb.RepoCategoryInfo{
			Id:          item.ID,
			Name:        item.Name,
			Description: item.Desc,
			ImageUrl:    item.ImageURL})
	}
	return &pb.RepoCategoryReply{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
		Total:    page.Total,
		Category: data,
	}, nil
}

func (r *OpenSourceInfo) repoData(ctx context.Context, repoInfo *domain.RepoInfo) *pb.RepoInfo {
	ownerName := ""
	language := ""
	if owner, _ := r.repo.FindOwnerByCache(ctx, repoInfo.OwnerID); owner != nil {
		ownerName = owner.Name
	}
	if langeInfo, _ := r.repo.FindLanguageByCache(ctx, repoInfo.LanguageId); langeInfo != nil {
		language = langeInfo.Name
	}
	var topic []string
	if repoInfo.Topics != "" {
		_ = json.Unmarshal([]byte(repoInfo.Topics), &topic)
	}
	return &pb.RepoInfo{
		Id:              repoInfo.ID,
		Name:            repoInfo.Name,
		FullName:        repoInfo.FullName,
		Image:           repoInfo.Image,
		OwnerId:         repoInfo.OwnerID,
		OwnerName:       ownerName,
		Private:         repoInfo.Private,
		Desc:            repoInfo.DescZh,
		DescEn:          repoInfo.Desc,
		HtmlUrl:         repoInfo.HtmlURL,
		Homepage:        repoInfo.Homepage,
		CloneUrl:        repoInfo.CloneURL,
		StargazersCount: repoInfo.StargazersCount,
		WatchersCount:   repoInfo.WatchersCount,
		Language:        language,
		LanguageId:      repoInfo.LanguageId,
		ForksCount:      repoInfo.ForksCount,
		OpenIssuesCount: repoInfo.OpenIssuesCount,
		Topics:          topic,
		OpenIssues:      repoInfo.OpenIssues,
		Watchers:        repoInfo.Watchers,
		DefaultBranch:   repoInfo.DefaultBranch,
		Score:           repoInfo.Score,
		Size:            repoInfo.Size,
		Forks:           repoInfo.Forks,
		IsFav:           r.IsRepoFavorite(ctx, auth.GetUid(ctx), repoInfo.ID),
		CreatedAt:       repoInfo.CreatedAt.Format(time.DateTime),
		UpdatedAt:       repoInfo.UpdatedAt.Format(time.DateTime),
	}
}

func (r *OpenSourceInfo) IsRepoFavorite(ctx context.Context, uid, repoId int64) bool {
	if uid < 1 {
		return false
	}
	// 得到收藏的信息
	info, err := r.repo.FindRepoFaveByCache(ctx, uid)
	if err != nil {
		return false
	}
	for _, item := range info {
		if item.RepoID == repoId {
			return true
		}
	}
	return false
}

func (r *OpenSourceInfo) GetRepoByCategory(ctx context.Context, req *pb.RepoByCategoryRequest) (*pb.RepoByCategoryReply, error) {
	if req.Id < 1 {
		return nil, fmt.Errorf("category id must be greater than 0")
	}
	page := &domain.Page{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
	}
	info, err := r.repo.FindRepoCategoryByCatId(ctx, req.Id, page)
	if err != nil {
		return nil, err
	}

	var data []*pb.RepoInfo
	for _, item := range info {
		repoInfo, err := r.repo.FindRepoById(ctx, item.RepoID)
		if err != nil {
			continue
		}
		data = append(data, r.repoData(ctx, repoInfo))
	}
	return &pb.RepoByCategoryReply{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
		Total:    page.Total,
		Repos:    data,
	}, nil

}

func (r *OpenSourceInfo) timeConverse(ctx context.Context, dateType, num int) string {
	now := time.Now()
	switch dateType {
	case DateTypeDay:
		return now.AddDate(0, 0, -num).Format(time.DateTime)
	case DateTypeWeek:
		return now.AddDate(0, 0, -num*7).Format(time.DateTime)
	case DateTypeMonth:
		return now.AddDate(0, -num, 0).Format(time.DateTime)
	default:
		return now.Format(time.DateTime)
	}
}

func (r *OpenSourceInfo) GetRepoMeasure(ctx context.Context, req *pb.RepoMeasureRequest) (*pb.RepoMeasureReply, error) {
	page := &domain.Page{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
	}
	date := r.timeConverse(ctx, int(req.DateType), int(req.Num))
	info, err := r.repo.FindRepoMetrics(ctx, int(req.Type), date, page)
	if err != nil {
		return nil, err
	}

	var data []*pb.RepoInfo
	for _, item := range info {
		repoInfo, err := r.repo.FindRepoById(ctx, item.RepoID)
		if err != nil {
			continue
		}
		data = append(data, r.repoData(ctx, repoInfo))
	}
	return &pb.RepoMeasureReply{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
		Total:    page.Total,
		Repos:    data,
	}, nil
}

func (r *OpenSourceInfo) RepoFav(ctx context.Context, req *pb.RepoFavRequest) (*pb.RepoFavReply, error) {
	// 获取uid
	uid := auth.GetUid(ctx)
	if uid == 0 {
		return nil, fmt.Errorf("uid is empty")
	}
	if len(req.RepoIds) == 0 {
		return nil, fmt.Errorf("repo ids is empty")
	}
	// 收藏的仓库不能为空
	for _, repoId := range req.RepoIds {
		// 判断是否已经收藏
		info, _ := r.repo.FindRepoFavor(ctx, 0, uid, repoId, &domain.Page{})
		if len(info) > 0 {
			// 判断状态是否一致
			if info[0].Status != int(req.IsFav) {
				if err := r.repo.UpdateRepoFavor(ctx, info[0].ID, req.IsFav); err != nil {
					r.log.Error(ctx, "update repo favor error", err)
					continue
				}
			}
		} else {
			if err := r.repo.AddRepoFavor(ctx, &domain.RepoFav{
				UID:       uid,
				RepoID:    repoId,
				Status:    0,
				Email:     auth.GetEmail(ctx),
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}); err != nil {
				r.log.Error(ctx, "add repo favor error", err)
				continue
			}
		}
	}
	// 更新缓存
	_, _ = r.repo.UpdateRepoFaveCache(ctx, uid)
	return &pb.RepoFavReply{}, nil
}

func (r *OpenSourceInfo) GetRepoFav(ctx context.Context, req *pb.RepoFavListRequest) (*pb.RepoReply, error) {
	uid := auth.GetUid(ctx)
	if uid == 0 {
		return nil, fmt.Errorf("uid is empty")
	}
	page := &domain.Page{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
	}
	info, err := r.repo.FindRepoFavor(ctx, 0, uid, 0, page)
	if err != nil {
		return nil, err
	}
	var data []*pb.RepoInfo
	for _, item := range info {
		repoInfo, err := r.repo.FindRepoById(ctx, item.RepoID)
		if err != nil {
			continue
		}
		data = append(data, r.repoData(ctx, repoInfo))
	}
	return &pb.RepoReply{
		PageNum:  req.PageNum,
		PageSize: req.PageSize,
		Total:    page.Total,
		Repos:    data,
	}, nil

}

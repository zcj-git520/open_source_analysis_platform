package data

import (
	pb "collect_open_source_data/api/open_source/v1"
	"collect_open_source_data/internal/biz"
	"collect_open_source_data/internal/domain"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
	"time"
)

type openSourceInfoRepo struct {
	data *Data
	log  *log.Helper
}

func NewOpenSourceRepo(data *Data, logger log.Logger) biz.OpenSourceRepo {
	return &openSourceInfoRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (o *openSourceInfoRepo) InsertOwner(ctx context.Context, owner *domain.Owner) (int64, error) {
	info := o.data.db.Create(owner)
	return owner.ID, info.Error
}

func (o *openSourceInfoRepo) FindOwnerByHtmlUrl(ctx context.Context, htmlUrl string) (*domain.Owner, error) {
	ownerInfo := &domain.Owner{}
	if err := o.data.db.Where("html_url = ?", htmlUrl).First(ownerInfo).Error; err != nil {
		return nil, err
	}
	if ownerInfo.ID == 0 { // 没有找到
		return nil, nil
	}
	return ownerInfo, nil
}

func (o *openSourceInfoRepo) FindRepoByName(ctx context.Context, name string) (*domain.RepoInfo, error) {
	repoInfo := &domain.RepoInfo{}
	if err := o.data.db.Where("name = ?", name).First(repoInfo).Error; err != nil {
		return nil, err
	}
	return repoInfo, nil
}

func (o *openSourceInfoRepo) FindRepoById(ctx context.Context, id int64) (*domain.RepoInfo, error) {
	repoInfo := &domain.RepoInfo{}
	if err := o.data.db.Where("id = ?", id).First(repoInfo).Error; err != nil {
		return nil, err
	}
	return repoInfo, nil
}

func (o *openSourceInfoRepo) FindLanguage(ctx context.Context, name string, id int64, page *domain.Page) ([]*domain.Language, error) {
	var language []*domain.Language
	tx := o.data.db
	if id > 0 {
		tx = tx.Where("id = ?", id)
	}
	if name != "" {
		tx = tx.Where("`name` =  ? ", name)
	}
	tx.Find(&language).Count(&page.Total)
	err := tx.Limit(page.Limit()).Offset(page.Offset()).Find(&language).Error

	return language, err
}

func (o *openSourceInfoRepo) InsertRepo(ctx context.Context, repo *domain.RepoInfo) error {
	return o.data.db.Create(repo).Error
}

func (o *openSourceInfoRepo) InsertLanguage(ctx context.Context) ([]*domain.Language, error) {
	var info []*domain.Language
	err := o.data.db.Find(&info).Error
	return info, err
}

func (o *openSourceInfoRepo) FindOwner(ctx context.Context, name, ownerType, email string, Id int64, page *domain.Page) ([]*domain.Owner, error) {
	var ownerInfo []*domain.Owner
	tx := o.data.db
	if Id > 0 {
		tx = tx.Where("id = ?", Id)
	}
	if name != "" {
		tx = tx.Where("`name` LIKE ? ", "%"+name+"%")
	}
	if ownerType != "" {
		tx = tx.Where("`type` = ?", ownerType)
	}
	if email != "" {
		tx = tx.Where("email = ?", email)
	}
	tx.Find(&ownerInfo).Count(&page.Total)
	err := tx.Limit(page.Limit()).Offset(page.Offset()).Find(&ownerInfo).Error

	return ownerInfo, err
}

func (o *openSourceInfoRepo) filter(ctx context.Context, tx *gorm.DB, field string, filter *pb.QueryFilter) *gorm.DB {
	queryFilter := fmt.Sprintf("%s", field)
	switch filter.Op {
	case pb.QueryFilter_GT:
		tx = tx.Where(queryFilter+" > ?", filter.TargetValue)
	case pb.QueryFilter_LT:
		tx = tx.Where(queryFilter+" < ?", filter.TargetValue)
	case pb.QueryFilter_GTE:
		tx = tx.Where(queryFilter+" >= ?", filter.TargetValue)
	case pb.QueryFilter_LTE:
		tx = tx.Where(queryFilter+" <= ?", filter.TargetValue)
	}
	return tx
}

func (o *openSourceInfoRepo) FindRepo(ctx context.Context, req *pb.RepoRequest, page *domain.Page) ([]*domain.RepoInfo, error) {
	var repoInfo []*domain.RepoInfo
	tx := o.data.db
	if req.ID > 0 {
		tx = tx.Where("id = ?", req.ID)
	}
	if req.Name != "" {
		tx = tx.Where("`full_name` LIKE ? ", "%"+req.Name+"%")
	}
	if req.Desc != "" {
		tx = tx.Where("`desc` LIKE ?", "%"+req.Desc+"%")
	}
	if req.LanguageId > 0 {
		tx = tx.Where("language_id = ?", req.LanguageId)
	}
	if req.OwnerId > 0 {
		tx = tx.Where("owner_id = ?", req.OwnerId)
	}
	if len(req.Filters) > 0 {
		for _, filter := range req.Filters {
			switch filter.Field {
			case "stargazers_count":
				tx = o.filter(ctx, tx, "stargazers_count", filter)
			case "watchers_count":
				tx = o.filter(ctx, tx, "watchers_count", filter)
			case "forks_count":
				tx = o.filter(ctx, tx, "forks_count", filter)
			case "open_issues_count":
				tx = o.filter(ctx, tx, "open_issues_count", filter)
			case "open_issues":
				tx = o.filter(ctx, tx, "open_issues", filter)
			case "watchers":
				tx = o.filter(ctx, tx, "watchers", filter)
			}
		}
	}
	tx.Find(&repoInfo).Count(&page.Total)
	// 排序
	if req.Sort != nil {
		tx = tx.Order(fmt.Sprintf("%s %s", req.Sort.Field, req.Sort.Order))
	} else {
		tx = tx.Order("stargazers_count desc")
	}
	err := tx.Limit(page.Limit()).Offset(page.Offset()).Find(&repoInfo).Error
	return repoInfo, err
}

func (o *openSourceInfoRepo) UpdateRepo(ctx context.Context, repo *domain.RepoInfo) error {
	if repo == nil {
		return fmt.Errorf("repo is nil")
	}
	//repo.UpdatedAt = time.Now()
	return o.data.db.Model(&domain.RepoInfo{}).Where("id = ?", repo.ID).Updates(repo).Error
}

func (o *openSourceInfoRepo) UpdateOwner(ctx context.Context, owner *domain.Owner) error {
	if owner == nil {
		return fmt.Errorf("owner is nil")
	}
	//owner.UpdatedAt = time.Now()
	return o.data.db.Model(&domain.Owner{}).Where("id = ?", owner.ID).Updates(owner).Error
}

func (o *openSourceInfoRepo) FindLanguageByCache(ctx context.Context, languageId int64) (*domain.Language, error) {
	language := &domain.Language{}
	data, err := o.data.rdb.Get(ctx, fmt.Sprintf("opensource_languageId_%d", languageId)).Result()
	// 数据不存在，就查询数据库，并缓存到redis
	if err != nil {
		if err = o.data.db.Where("id = ?", languageId).First(language).Error; err != nil {
			return nil, err
		}
		cacheData, _ := json.Marshal(language)
		o.data.rdb.Set(ctx, fmt.Sprintf("opensource_languageId_%d", languageId), string(cacheData), 10*time.Minute)
		return language, nil
	}
	if err = json.Unmarshal([]byte(data), language); err != nil {
		return nil, err
	}
	return language, nil
}

func (o *openSourceInfoRepo) FindOwnerByCache(ctx context.Context, Id int64) (*domain.Owner, error) {
	var ownerInfo *domain.Owner
	data, err := o.data.rdb.Get(ctx, fmt.Sprintf("opensource_ownerid_%d", Id)).Result()
	if err != nil {
		if err = o.data.db.Where("id = ?", Id).First(&ownerInfo).Error; err != nil {
			return nil, err
		}
		cacheData, _ := json.Marshal(ownerInfo)
		o.data.rdb.Set(ctx, fmt.Sprintf("opensource_ownerid_%d", Id), string(cacheData), 10*time.Minute)
		return ownerInfo, nil
	}
	if err = json.Unmarshal([]byte(data), &ownerInfo); err != nil {
		return nil, err
	}
	return ownerInfo, err
}

func (o *openSourceInfoRepo) FindRepoFaveByCache(ctx context.Context, uid int64) ([]*domain.RepoFav, error) {
	var info []*domain.RepoFav
	data, err := o.data.rdb.Get(ctx, fmt.Sprintf("opensource_repoFav_%d", uid)).Result()
	if err != nil {
		return o.UpdateRepoFaveCache(ctx, uid)
	}
	if err = json.Unmarshal([]byte(data), &info); err != nil {
		return nil, err
	}
	return info, nil
}

func (o *openSourceInfoRepo) UpdateRepoFaveCache(ctx context.Context, uid int64) ([]*domain.RepoFav, error) {
	var info []*domain.RepoFav
	if err := o.data.db.Where("uid = ? and status = 0", uid).Find(&info).Error; err != nil {
		return nil, err
	}
	cacheData, _ := json.Marshal(info)
	o.data.rdb.Set(ctx, fmt.Sprintf("opensource_repoFav_%d", uid), string(cacheData), 1*time.Hour)
	return info, nil
}

func (o *openSourceInfoRepo) FindRepoCategory(ctx context.Context, name string, id int64, page *domain.Page) ([]*domain.RepoCategory, error) {
	var repoCategory []*domain.RepoCategory
	tx := o.data.db
	if id > 0 {
		tx = tx.Where("id = ?", id)
	}
	if name != "" {
		tx = tx.Where("`name` =  ? ", name)
	}
	tx.Find(&repoCategory).Count(&page.Total)
	err := tx.Limit(page.Limit()).Offset(page.Offset()).Find(&repoCategory).Error

	return repoCategory, err
}

func (o *openSourceInfoRepo) FindRepoCategoryId(ctx context.Context, repoId, categoryId int64) bool {
	var repoCategoryId *domain.RepoCategoryId
	err := o.data.db.Where("repo_id = ? and cat_id = ?", repoId, categoryId).First(&repoCategoryId).Error
	if err == nil && repoCategoryId != nil {
		return true
	}
	return false
}

func (o *openSourceInfoRepo) AddRepoCategoryId(ctx context.Context, repoId, categoryId int64) error {
	repoCategoryId := &domain.RepoCategoryId{
		RepoID: repoId,
		CatID:  categoryId,
	}
	return o.data.db.Create(repoCategoryId).Error
}

func (o *openSourceInfoRepo) FindRepoCategoryIdByRepoId(repoId int64) bool {
	var repoCategoryId *domain.RepoCategoryId
	err := o.data.db.Where("repo_id = ?", repoId).First(&repoCategoryId).Error
	if err == nil && repoCategoryId != nil {
		return true
	}
	return false
}

func (o *openSourceInfoRepo) FindRepoCategoryByCatId(ctx context.Context, id int64, page *domain.Page) ([]*domain.RepoCategoryId, error) {
	var repoCategoryId []*domain.RepoCategoryId
	tx := o.data.db
	if id > 0 {
		tx = tx.Where("cat_id = ?", id)
	}
	tx.Find(&repoCategoryId).Count(&page.Total)
	err := tx.Limit(page.Limit()).Offset(page.Offset()).Find(&repoCategoryId).Error

	return repoCategoryId, err
}

func (o *openSourceInfoRepo) AddRepoMetrics(ctx context.Context, metrics []*domain.RepoMetrics) error {
	return o.data.db.CreateInBatches(metrics, len(metrics)).Error
}

func (o *openSourceInfoRepo) FindRepoMetrics(ctx context.Context, metricType int, data string, page *domain.Page) ([]*domain.RepoMetricsResult, error) {
	var repoMetricsResult []*domain.RepoMetricsResult
	o.data.db.Table("repo_metrics").Select("repo_id, SUM(value) as total_value").
		Where(" `type` = ?", metricType).
		Where(fmt.Sprintf("date >= '%s'", data)).
		Group("repo_id").Scan(&[]*domain.RepoMetricsResult{}).Count(&page.Total)
	err := o.data.db.Table("repo_metrics").Select("repo_id, SUM(value) as total_value").
		Where(" `type` = ?", metricType).
		Where(fmt.Sprintf("date >= '%s'", data)).
		Group("repo_id").Limit(page.Limit()).
		Offset(page.Offset()).Scan(&repoMetricsResult).Error
	return repoMetricsResult, err
}

func (o *openSourceInfoRepo) FindRepoFavor(ctx context.Context, id, uid, repoId int64, page *domain.Page) ([]*domain.RepoFav, error) {
	var favorList []*domain.RepoFav
	tx := o.data.db.Where("status = 0")
	if id > 0 {
		tx = tx.Where("id = ?", id)
	}
	if uid > 0 {
		tx = tx.Where("uid = ?", uid)
	}
	if repoId > 0 {
		tx = tx.Where("repo_id = ?", repoId)
	}
	tx.Find(&favorList).Count(&page.Total)
	// 排序
	err := tx.Limit(page.Limit()).Offset(page.Offset()).Find(&favorList).Error
	return favorList, err
}

func (o *openSourceInfoRepo) AddRepoFavor(ctx context.Context, favorInfo *domain.RepoFav) error {
	return o.data.db.Create(favorInfo).Error
}

func (o *openSourceInfoRepo) UpdateRepoFavor(ctx context.Context, favorId int64, isFavor int32) error {
	return o.data.db.Model(&domain.RepoFav{}).Where("id = ?", favorId).Updates(map[string]any{"status": isFavor,
		"updated_at": time.Now()}).Error
}

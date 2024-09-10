package data

import (
	"collect_open_source_data/internal/biz"
	"collect_open_source_data/internal/domain"
	"context"
	"github.com/go-kratos/kratos/v2/log"
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

func (o *openSourceInfoRepo) FindRepo(ctx context.Context, name, desc string, languageId, ownerId int64, page *domain.Page) ([]*domain.RepoInfo, error) {
	var repoInfo []*domain.RepoInfo
	tx := o.data.db
	if name != "" {
		tx = tx.Where("`full_name` LIKE ? ", "%"+name+"%")
	}
	if desc != "" {
		tx = tx.Where("`desc` LIKE ?", "%"+desc+"%")
	}
	if languageId > 0 {
		tx = tx.Where("language_id = ?", languageId)
	}
	if ownerId > 0 {
		tx = tx.Where("owner_id = ?", ownerId)
	}
	tx.Find(&repoInfo).Count(&page.Total)
	err := tx.Limit(page.Limit()).Offset(page.Offset()).Find(&repoInfo).Error

	return repoInfo, err
}

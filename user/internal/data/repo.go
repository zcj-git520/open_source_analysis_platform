package data

import (
	"context"
	"user/internal/domain"
)

//type repoInfoRepo struct {
//	data *Data
//	log  *log.Helper
//}
//
//func NewRepoInfoRepo(data *Data, logger log.Logger) biz.RepoInfoRepo {
//	return &repoInfoRepo{
//		data: data,
//		log:  log.NewHelper(logger),
//	}
//}

func (u *userRepo) InsertOwner(ctx context.Context, owner *domain.Owner) (int64, error) {
	info := u.data.db.Create(owner)
	return owner.ID, info.Error
}

func (u *userRepo) FindOwnerByName(ctx context.Context, name string) (*domain.Owner, error) {
	ownerInfo := &domain.Owner{}
	if err := u.data.db.Where("name = ?", name).First(ownerInfo).Error; err != nil {
		return nil, err
	}
	if ownerInfo.ID == 0 { // 没有找到
		return nil, nil
	}
	return ownerInfo, nil
}

func (u *userRepo) FindRepoByName(ctx context.Context, name string) (*domain.RepoInfo, error) {
	repoInfo := &domain.RepoInfo{}
	if err := u.data.db.Where("name = ?", name).First(repoInfo).Error; err != nil {
		return nil, err
	}
	return repoInfo, nil
}

func (u *userRepo) FindLanguage(ctx context.Context, name string) (*domain.Language, error) {
	language := &domain.Language{}
	if err := u.data.db.Where("name = ?", name).First(language).Error; err != nil {
		return nil, err
	}
	return language, nil
}

func (u *userRepo) InsertRepo(ctx context.Context, repo *domain.RepoInfo) error {
	return u.data.db.Create(repo).Error
}

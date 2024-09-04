package biz

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"user/internal/domain"
	"user/internal/pkg"
)

const (
	BaseGithubURL = "https://api.github.com/search/repositories"
)

type RepoInfoRepo interface {
	InsertOwner(ctx context.Context, owner *domain.Owner) (int64, error)
	InsertRepo(ctx context.Context, repo *domain.RepoInfo) error
	FindOwnerByName(ctx context.Context, name string) (*domain.Owner, error)
	FindRepoByName(ctx context.Context, name string) (*domain.RepoInfo, error)
	FindLanguage(ctx context.Context, name string) (*domain.Language, error)
}

type RepoInfo struct {
	repo RepoInfoRepo
	log  *log.Helper
}

func NewRepoInfo(repo RepoInfoRepo, logger log.Logger) *RepoInfo {
	return &RepoInfo{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func (r *User) request(method string, url string, headers http.Header, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header = headers
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status code: %d", resp.StatusCode)
	}
	return ioutil.ReadAll(resp.Body)
}

func (r *User) getResults(search string, headers http.Header, page int, stars int) ([]*Repo, error) {
	queryParams := url.Values{}
	queryParams.Add("q", fmt.Sprintf("%s stars:<=%d", search, stars))
	queryParams.Add("page", fmt.Sprintf("%d", page))
	queryParams.Add("per_page", "50")
	queryParams.Add("sort", "stars")
	queryParams.Add("order", "desc")
	body, err := r.request("GET", fmt.Sprintf("%s?%s", BaseGithubURL, queryParams.Encode()), headers, nil)
	result := &RepoResponse{}
	err = json.Unmarshal(body, &result)
	if err != nil || len(result.Items) == 0 {
		return nil, err
	}
	return result.Items, nil
}

func (r *User) getOwnerInfo(ctx context.Context, url string, headers http.Header) (*domain.Owner, error) {
	fmt.Printf("get owner info: %s\n", url)
	body, err := r.request("GET", url, headers, nil)
	if err != nil {
		return nil, err
	}
	result := &RepoUser{}
	if err = json.Unmarshal(body, &result); err != nil || result == nil {
		return nil, err
	}
	return &domain.Owner{
		AvatarURL:   result.AvatarURL,
		Type:        result.Type,
		Login:       result.Login,
		HtmlURL:     result.HtmlURL,
		Name:        result.Name,
		Email:       pkg.ToString(result.Email),
		Bio:         result.Bio,
		PublicRepos: int64(result.PublicRepos),
		PublicGists: int64(result.PublicGists),
		Followers:   int64(result.Followers),
		Following:   int64(result.Following),
		CreatedAt:   result.CreatedAt,
		UpdatedAt:   result.UpdatedAt,
	}, nil
}

func (r *User) ParseResult(ctx context.Context, search string, headers http.Header, page int, stars int) error {
	result, err := r.getResults(search, headers, page, stars)
	if err != nil || len(result) == 0 {
		return err
	}
	for _, item := range result {
		var ownerId int64
		if item.Owner == nil || item.Owner.URL == "" {
			continue
		}
		owner, err := r.getOwnerInfo(ctx, item.Owner.URL, headers)
		if err != nil {
			r.log.Errorf("getOwnerInfo error: %v", err)
			continue
		} else {
			// 查询owner是否存在
			if info, err := r.repo.FindOwnerByName(ctx, owner.Name); err != nil || info == nil {
				// 不存在则创建
				if ownerId, err = r.repo.InsertOwner(ctx, owner); err != nil {
					r.log.Errorf("create owner error: %v", err)
					continue
				}
			} else {
				ownerId = info.ID
			}
		}

		// 查询repo是否存在
		if info, err := r.repo.FindRepoByName(ctx, item.Name); err != nil || info == nil {
			// 不存在则创建
			// 查找语言id
			var langId int64
			if lang, err := r.repo.FindLanguage(ctx, item.Language); err == nil && lang != nil {
				langId = lang.ID
			}
			topicData, _ := json.Marshal(item.Topics)
			if err = r.repo.InsertRepo(ctx, &domain.RepoInfo{
				Name:     item.Name,
				FullName: item.FullName,
				//Image:           owner.Image,
				OwnerID:         ownerId,
				Private:         item.Private,
				Desc:            item.Description,
				HtmlURL:         item.HtmlURL,
				Homepage:        item.Homepage,
				CloneURL:        item.CloneURL,
				Size:            int64(item.Size),
				StargazersCount: int64(item.StargazersCount),
				WatchersCount:   int64(item.WatchersCount),
				LanguageId:      langId,
				ForksCount:      int64(item.ForksCount),
				OpenIssuesCount: int64(item.OpenIssuesCount),
				Topics:          string(topicData),
				Forks:           int64(item.Forks),
				OpenIssues:      int64(item.OpenIssues),
				Watchers:        int64(item.Watchers),
				DefaultBranch:   item.DefaultBranch,
				Score:           int64(item.Score),
				CreatedAt:       item.CreatedAt,
				UpdatedAt:       item.UpdatedAt,
			}); err != nil {
				r.log.Errorf("create repo error: %v", err)
				return err
			}
		}

	}
	return nil
}

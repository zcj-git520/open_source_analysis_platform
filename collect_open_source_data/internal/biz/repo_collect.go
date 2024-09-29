package biz

import (
	"collect_open_source_data/internal/domain"
	"collect_open_source_data/internal/pkg"
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func (r *OpenSourceInfo) request(method string, url string, headers http.Header, body io.Reader) ([]byte, error) {
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

func (r *OpenSourceInfo) getResults(search string, headers http.Header, page int, stars int) ([]*Repo, error) {
	queryParams := url.Values{}
	queryParams.Add("q", fmt.Sprintf("%s stars:>=%d", search, stars))
	queryParams.Add("page", fmt.Sprintf("%d", page))
	queryParams.Add("per_page", "100")
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

func (r *OpenSourceInfo) getOwnerInfo(ctx context.Context, url string, headers http.Header) (*domain.Owner, error) {
	body, err := r.request("GET", url, headers, nil)
	if err != nil {
		return nil, err
	}
	result := &RepoUser{}
	if err = json.Unmarshal(body, &result); err != nil || result == nil {
		return nil, err
	}
	if result.Name == "" {
		result.Name = strings.ReplaceAll(result.HtmlURL, "https://github.com/", "")
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

func (r *OpenSourceInfo) addOwnerInfo(ctx context.Context, owner *domain.Owner) (int64, error) {
	ownerId, err := r.repo.InsertOwner(ctx, owner)
	if err != nil {
		r.log.Errorf("create owner error: %v", err)
		return ownerId, err
	}
	return ownerId, nil
}

func (r *OpenSourceInfo) updateOwnerInfo(ctx context.Context, ownerInfo, owner *domain.Owner) (int64, error) {
	update := &domain.Owner{}
	if ownerInfo.AvatarURL != owner.AvatarURL {
		update.AvatarURL = owner.AvatarURL
	}
	if ownerInfo.Bio != owner.Bio {
		update.Bio = owner.Bio
	}
	if ownerInfo.Email != owner.Email {
		update.Email = owner.Email
	}
	if ownerInfo.Followers != owner.Followers {
		update.Followers = owner.Followers
	}
	if ownerInfo.Following != owner.Following {
		update.Following = owner.Following
	}
	if ownerInfo.PublicGists != owner.PublicGists {
		update.PublicGists = owner.PublicGists
	}
	if ownerInfo.PublicRepos != owner.PublicRepos {
		update.PublicRepos = owner.PublicRepos
	}
	if update != nil {
		update.ID = ownerInfo.ID
		if err := r.repo.UpdateOwner(ctx, update); err != nil {
			r.log.Errorf("update owner error: %v", err)
			return 0, err
		}
	}
	return ownerInfo.ID, nil
}

func (r *OpenSourceInfo) OwnerInfoChange(ctx context.Context, owner *domain.Owner) (int64, error) {
	// 查询owner是否存在
	ownerInfo, err := r.repo.FindOwnerByHtmlUrl(ctx, owner.HtmlURL)
	if err != nil || ownerInfo == nil {
		// 不存在则创建
		return r.addOwnerInfo(ctx, owner)
	}
	// 更新owner信息
	return r.updateOwnerInfo(ctx, ownerInfo, owner)
}

func (r *OpenSourceInfo) addRepoInfo(ctx context.Context, item *Repo, ownerId int64, avatarURL string) error {
	var langId int64
	if langInfo, err := r.repo.FindLanguage(ctx, item.Language, 0, &domain.Page{
		PageNum:  1,
		PageSize: 1,
	}); err == nil && len(langInfo) > 0 {
		langId = langInfo[0].ID
	}
	repoImg, _ := r.getRepoImage(ctx, item.HtmlURL)
	if repoImg == "" {
		repoImg = avatarURL
	}
	tran := pkg.NewTranslateModeler(pkg.AppID, pkg.SecretKey, "en", "zh")

	topicData, _ := json.Marshal(item.Topics)
	info := &domain.RepoInfo{
		Name:            item.Name,
		FullName:        item.FullName,
		Image:           repoImg,
		OwnerID:         ownerId,
		Private:         item.Private,
		Desc:            item.Description,
		DescZh:          tran.Translate(item.Description),
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
	}
	if err := r.repo.InsertRepo(ctx, info); err != nil {
		r.log.Errorf("create repo error: %v", err)
		return err
	}
	// 更新repo分类
	go r.RepoCategoryIdChange(ctx, item.Topics, info.ID)
	return nil
}

func (r *OpenSourceInfo) updateRepoInfo(ctx context.Context, info *domain.RepoInfo, item *Repo) error {
	updateRepo := &domain.RepoInfo{}
	if int64(item.StargazersCount) != info.StargazersCount {
		updateRepo.StargazersCount = int64(item.StargazersCount)
	}
	if int64(item.WatchersCount) != info.WatchersCount {
		updateRepo.WatchersCount = int64(item.WatchersCount)
	}
	if int64(item.ForksCount) != info.ForksCount {
		updateRepo.ForksCount = int64(item.ForksCount)
	}
	if int64(item.OpenIssuesCount) != info.OpenIssuesCount {
		updateRepo.OpenIssuesCount = int64(item.OpenIssuesCount)
	}
	if int64(item.Forks) != info.Forks {
		updateRepo.Forks = int64(item.Forks)
	}
	if int64(item.OpenIssues) != info.OpenIssues {
		updateRepo.OpenIssues = int64(item.OpenIssues)
	}
	if int64(item.Watchers) != info.Watchers {
		updateRepo.Watchers = int64(item.Watchers)
	}
	if item.DefaultBranch != info.DefaultBranch {
		updateRepo.DefaultBranch = item.DefaultBranch
	}
	if int64(item.Score) != info.Score {
		updateRepo.Score = int64(item.Score)
	}
	if item.UpdatedAt != info.UpdatedAt {
		updateRepo.UpdatedAt = item.UpdatedAt
	}
	if updateRepo != nil {
		updateRepo.ID = info.ID
		updateRepo.UpdatedAt = item.UpdatedAt
		if err := r.repo.UpdateRepo(ctx, updateRepo); err != nil {
			r.log.Errorf("update repo error: %v", err)
			return err
		}
	}
	return nil
}

func (r *OpenSourceInfo) RepoInfoChange(ctx context.Context, item *Repo, ownerId int64, avatarURL string) error {
	// 查询repo是否存在
	info, err := r.repo.FindRepoByName(ctx, item.Name)
	if err != nil || info == nil {
		// 不存在则创建
		return r.addRepoInfo(ctx, item, ownerId, avatarURL)
	}
	// 存在就更新
	return r.updateRepoInfo(ctx, info, item)
}

func (r *OpenSourceInfo) ParseResult(ctx context.Context, search string, headers http.Header, page int, stars int) error {
	result, err := r.getResults(search, headers, page, stars)
	if err != nil || len(result) == 0 {
		return err
	}
	for _, item := range result {
		if item.Owner == nil || item.Owner.URL == "" {
			continue
		}
		owner, err := r.getOwnerInfo(ctx, item.Owner.URL, headers)
		if err != nil {
			r.log.Errorf("getOwnerInfo error: %v", err)
			continue
		}
		ownerId, err := r.OwnerInfoChange(ctx, owner)
		if err != nil {
			r.log.Errorf("change OwnerInfo error: %v", err)
			continue
		}
		if err = r.RepoInfoChange(ctx, item, ownerId, owner.AvatarURL); err != nil {
			r.log.Errorf("change RepoInfo error: %v", err)
			continue
		}
	}
	return nil
}

func (r *OpenSourceInfo) extractImageURLs(htmlContent string) []string {
	imageUrls := make([]string, 0)
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return imageUrls
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "img" {
			for _, attr := range n.Attr {
				if attr.Key == "src" {
					imageUrls = append(imageUrls, attr.Val)
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	return imageUrls
}

func (r *OpenSourceInfo) isImageURL(url string) bool {
	return strings.HasSuffix(url, ".png") || strings.HasSuffix(url, ".jpg") || strings.HasSuffix(url, ".jpeg")
}

func (r *OpenSourceInfo) getRepoImage(ctx context.Context, repoName string) (string, error) {
	htmlContent, err := r.request("GET", repoName, nil, nil)
	if err != nil {
		return "", err
	}

	imageUrls := r.extractImageURLs(string(htmlContent))
	for _, itemUrl := range imageUrls {
		if r.isImageURL(itemUrl) {
			if !strings.HasPrefix(itemUrl, "http") {
				itemUrl = "https://github.com" + itemUrl
			}
			return itemUrl, nil
		}
	}
	return "", nil
}

func (r *OpenSourceInfo) Collect() {
	r.Page++
	if r.Page > 10 {
		r.Page = 1
	}
	fmt.Println("========================================:  ", r.Page)
	language := []string{"Python", "JavaScript", "Java", "C", "C++", "C#", "PHP", "Ruby", "Go", "Rust", "TypeScript"}
	for _, item := range language {
		r.log.Infof("language: %v", item)
		if err := r.ParseResult(context.TODO(), "language:"+item, http.Header{
			"Authorization": []string{"token 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"},
			"Accept":        []string{"application/json"},
		}, r.Page, 1000); err != nil {
			r.log.Errorf("parse result error: %v", err)
			continue
		}
	}

}

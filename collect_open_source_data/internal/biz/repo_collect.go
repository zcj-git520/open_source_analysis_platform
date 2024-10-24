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
	"time"
)

const (
	RepoMetricStar = iota
	RepoMetricFork
	RepoMetricWatch
	RepoMetricIssue
)

const (
	RepoChangeSubject = "osap 开源仓库变更"
	RepoChangeContent = `<td style="font-size:14px;color:#333;padding:24px 40px 0 40px">
                尊敬的用户 您好！
                <br>
                <br>
                您关注的开源仓库：<b>%s</b>，发生变更!
	    		<br>
				osap地址：<b>http://192.168.40.25:5173/#/project-trends/index</b>
				<br>
				仓库地址：<b>%s</b>
				<br>
				请查看!
                <br> 
                如果您已查看，请无视。
            </td>`
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
	updateState := false
	if ownerInfo.AvatarURL != owner.AvatarURL {
		update.AvatarURL = owner.AvatarURL
		updateState = true
	}
	if ownerInfo.Bio != owner.Bio {
		update.Bio = owner.Bio
		updateState = true
	}
	if ownerInfo.Email != owner.Email {
		update.Email = owner.Email
		updateState = true
	}
	if ownerInfo.Followers != owner.Followers {
		update.Followers = owner.Followers
		updateState = true
	}
	if ownerInfo.Following != owner.Following {
		update.Following = owner.Following
		updateState = true
	}
	if ownerInfo.PublicGists != owner.PublicGists {
		update.PublicGists = owner.PublicGists
		updateState = true
	}
	if ownerInfo.PublicRepos != owner.PublicRepos {
		update.PublicRepos = owner.PublicRepos
		updateState = true
	}
	if update != nil && updateState {
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
	updateRepoState := false
	var repoMetricList []*domain.RepoMetrics
	if int64(item.StargazersCount) != info.StargazersCount {
		repoMetricList = append(repoMetricList, &domain.RepoMetrics{
			RepoID:      info.ID,
			Type:        RepoMetricStar,
			Value:       int64(item.StargazersCount) - info.StargazersCount,
			OriginValue: info.StargazersCount,
			NowValue:    int64(item.StargazersCount),
			Date:        time.Now(),
		})
		updateRepo.StargazersCount = int64(item.StargazersCount)
		updateRepoState = true
	}
	if int64(item.WatchersCount) != info.WatchersCount {
		repoMetricList = append(repoMetricList, &domain.RepoMetrics{
			RepoID:      info.ID,
			Type:        RepoMetricWatch,
			Value:       int64(item.WatchersCount) - info.WatchersCount,
			OriginValue: info.WatchersCount,
			NowValue:    int64(item.WatchersCount),
			Date:        time.Now(),
		})
		updateRepo.WatchersCount = int64(item.WatchersCount)
		updateRepoState = true
	}
	if int64(item.ForksCount) != info.ForksCount {
		repoMetricList = append(repoMetricList, &domain.RepoMetrics{
			RepoID:      info.ID,
			Type:        RepoMetricFork,
			Value:       int64(item.ForksCount) - info.ForksCount,
			OriginValue: info.ForksCount,
			NowValue:    int64(item.ForksCount),
			Date:        time.Now(),
		})
		updateRepo.ForksCount = int64(item.ForksCount)
		updateRepoState = true
	}
	if int64(item.OpenIssuesCount) != info.OpenIssuesCount {
		repoMetricList = append(repoMetricList, &domain.RepoMetrics{
			RepoID:      info.ID,
			Type:        RepoMetricIssue,
			Value:       int64(item.OpenIssuesCount) - info.OpenIssuesCount,
			OriginValue: info.OpenIssuesCount,
			NowValue:    int64(item.OpenIssuesCount),
			Date:        time.Now(),
		})
		updateRepo.OpenIssuesCount = int64(item.OpenIssuesCount)
		updateRepoState = true
	}
	if int64(item.Forks) != info.Forks {
		repoMetricList = append(repoMetricList, &domain.RepoMetrics{
			RepoID:      info.ID,
			Type:        RepoMetricFork,
			Value:       int64(item.Forks) - info.Forks,
			OriginValue: info.Forks,
			NowValue:    int64(item.Forks),
			Date:        time.Now(),
		})
		updateRepo.Forks = int64(item.Forks)
		updateRepoState = true
	}
	if int64(item.OpenIssues) != info.OpenIssues {
		repoMetricList = append(repoMetricList, &domain.RepoMetrics{
			RepoID:      info.ID,
			Type:        RepoMetricIssue,
			Value:       int64(item.OpenIssues) - info.OpenIssues,
			OriginValue: info.OpenIssues,
			NowValue:    int64(item.OpenIssues),
			Date:        time.Now(),
		})
		updateRepo.OpenIssues = int64(item.OpenIssues)
		updateRepoState = true
	}
	if int64(item.Watchers) != info.Watchers {
		updateRepo.Watchers = int64(item.Watchers)
		updateRepoState = true
	}
	if item.DefaultBranch != info.DefaultBranch {
		updateRepo.DefaultBranch = item.DefaultBranch
		updateRepoState = true
	}
	if int64(item.Score) != info.Score {
		updateRepo.Score = int64(item.Score)
		updateRepoState = true
	}
	if item.UpdatedAt != info.UpdatedAt {
		updateRepo.UpdatedAt = item.UpdatedAt
		updateRepoState = true
	}
	if int64(item.Size) != info.Size {
		// 仓库变更通知
		go r.Notify(ctx, info)
	}
	// 更新仓库指标信息
	if len(repoMetricList) > 0 {
		_ = r.repo.AddRepoMetrics(ctx, repoMetricList)
	}
	if updateRepo != nil && updateRepoState {
		// 更新仓库信息
		updateRepo.ID = info.ID
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
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
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
		return
	}
	fmt.Println("========================================:  ", r.Page)
	language := []string{"Python", "JavaScript", "Java", "C", "C++", "C#", "PHP", "Ruby", "Go", "Rust", "TypeScript"}
	for _, item := range language {
		select {
		case <-r.Ctx.Done():
			return
		default:
			r.log.Infof("language: %v", item)
			if err := r.ParseResult(r.Ctx, "language:"+item, http.Header{
				"Authorization": []string{fmt.Sprintf("token %s", r.ccf.Token)},
				"Accept":        []string{"application/json"},
			}, r.Page, 1000); err != nil {
				r.log.Errorf("parse result error: %v", err)
				continue
			}
		}
	}
}

// 仓库发生变更通知收藏该仓库的用户
func (r *OpenSourceInfo) Notify(ctx context.Context, repo *domain.RepoInfo) {
	if !r.ec.Enable {
		return
	}
	// 查询所有收藏该仓库的用户
	info, err := r.repo.FindRepoFavor(ctx, 0, 0, repo.ID, &domain.Page{
		PageNum:  1,
		PageSize: 10000000,
	})
	if err != nil {
		r.log.Errorf("change RepoInfo error: %v", err)
		return
	}
	for _, item := range info {
		// 通过uid 查询用户信息
		email := pkg.NewEmailSMTP(pkg.WithSmtpHost(r.ec.SmtpHost), pkg.WithSmtpPort(int(r.ec.SmtpPort)),
			pkg.WithSmtpUsername(r.ec.SmtpUsername), pkg.WithSmtpPassword(r.ec.SmtpPassword),
			pkg.WithFrom(r.ec.From),
			pkg.WithTo([]string{item.Email}))
		// 邮件内容
		if err = email.SendEmailSMTP(RepoChangeSubject, fmt.Sprintf(RepoChangeContent, repo.Name, repo.HtmlURL)); err != nil {
			r.log.Errorf("email: %s: send email error: %v", item.Email, err)
			continue
		}
	}
}

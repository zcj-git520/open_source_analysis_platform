package biz

import (
	pb "collect_open_source_data/api/open_source/v1"
	"collect_open_source_data/internal/domain"
	"collect_open_source_data/internal/pkg"
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	BaseGithubURL = "https://api.github.com/search/repositories"
)

var frontendRegex *regexp.Regexp
var backendRegex *regexp.Regexp
var mobileDevelopmentRegex *regexp.Regexp
var dataScienceRegex *regexp.Regexp
var gameDevelopmentRegex *regexp.Regexp
var systemRegex *regexp.Regexp
var libraryRegex *regexp.Regexp
var aiRegex *regexp.Regexp
var blockchainRegex *regexp.Regexp
var iotRegex *regexp.Regexp
var educationRegex *regexp.Regexp
var securityRegex *regexp.Regexp
var creativeRegex *regexp.Regexp

func init() {
	frontendRegex = regexp.MustCompile(strings.Join(generateRegexPattern([]string{"Frontend", "front end", "web", "webpage", "interface", "UI", "UX", "HTML", "CSS", "JavaScript", "responsive design", "React", "Vue.js", "Angular", "TypeScript", "Sass/Less", "Webpack", "Babel", "Bootstrap", "Material-UI", "Ant Design", "Foundation", "jQuery", "D3.js"}), "|"))
	backendRegex = regexp.MustCompile(strings.Join(generateRegexPattern([]string{"Backend", "back end", "server", "API", "Node.js", "Express.js", "Django", "Flask", "Ruby on Rails", "Java", "Spring", "Go", "Python", "PHP", "C#", "ASP.NET", "Rust", "Elixir", "Scala", "Kotlin", "R", "Perl", "Lua", "Groovy", "gin", "Backend", "back end", "Server", "server", "API", "application interface", "REST", "RESTful API", "GraphQL", "microservice", "Microservices", "database", "Database", "SQL", "NoSQL", "MySQL", "PostgreSQL", "MongoDB", "Redis", "ORM", "object-relational mapping", "Node.js", "Express", "Koa", "Flask", "Django", "Spring Boot", "Laravel", "Ruby on Rails", "Go", "Golang", "Java", "C#", ".NET", "PHP", "Docker", "Kubernetes", "CI/CD", "continuous integration/continuous deployment", "Git", "version control", "GitHub", "GitLab", "Bitbucket", "security", "Security", "HTTPS", "OAuth", "JWT", "encryption", "Encryption", "caching", "Caching", "Redis", "Memcached", "message queue", "Message Queue", "RabbitMQ", "Kafka", "logging", "Logging", "ELK Stack", "log analysis", "monitoring", "Monitoring", "Prometheus", "Grafana", "cloud service", "Cloud Services", "AWS", "Azure", "GCP", "cloud function", "Cloud Functions", "containerization", "Containerization", "Docker", "Kubernetes", "microservice architecture", "Microservice Architecture", "service discovery", "Service Discovery", "load balancing", "Load Balancing", "Nginx", "HAProxy", "reverse proxy", "Reverse Proxy", "API gateway", "API Gateway", "database migration", "Database Migration", "data backup", "Data Backup", "unit testing", "Unit Testing", "integration testing", "Integration Testing", "performance testing", "Performance Testing", "stress testing", "Stress Testing"}), "|"))
	mobileDevelopmentRegex = regexp.MustCompile(strings.Join(generateRegexPattern([]string{"App", "application", "mobile application", "iOS", "Android", "Flutter", "React Native", "Kotlin", "Swift", "Xamarin", "Cordova"}), "|"))
	dataScienceRegex = regexp.MustCompile(strings.Join(generateRegexPattern([]string{"Data Science", "data science", "machine learning", "deep learning", "data analysis", "data mining", "data visualization", "TensorFlow", "PyTorch", "Keras", "Scikit-learn", "Pandas", "NumPy", "Matplotlib", "Seaborn", "Jupyter Notebook", "R", "Spark", "Hadoop"}), "|"))
	gameDevelopmentRegex = regexp.MustCompile(strings.Join(generateRegexPattern([]string{"Game Development", "game development", "game engine", "Unity", "Unreal Engine", "Godot", "Cocos2d-x", "game design", "game programming", "game art", "game audio", "game physics", "game AI", "game network", "game testing", "game release"}), "|"))
	systemRegex = regexp.MustCompile(strings.Join(generateRegexPattern([]string{"System", "operating system", "command line tool", "Shell", "script", "scripting language", "Shell script", "Bash", "Zsh", "Fish", "PowerShell", "batch script", "Cygwin", "MinGW", "Cygwin", "MinGW", "Windows Subsystem for Linux", "WSL", "Linux"}), "|"))
	libraryRegex = regexp.MustCompile(strings.Join(generateRegexPattern([]string{"Library", "library", "framework", "tool", "utility program", "practical tool", "practical library", "practical framework", "practical tool library", "practical tool framework"}), "|"))
	aiRegex = regexp.MustCompile(strings.Join(generateRegexPattern([]string{"openai", "AI", "artificial intelligence", "machine learning", "deep learning", "natural language processing", "computer vision", "speech recognition", "reinforcement learning", "neural network", "deep neural network", "tensorflow", "pytorch", "ai", "model"}), "|"))
	blockchainRegex = regexp.MustCompile(strings.Join(generateRegexPattern([]string{"Blockchain", "blockchain", "cryptocurrency", "bitcoin", "ethereum", "smart contract", "consensus algorithm", "distributed ledger"}), "|"))
	iotRegex = regexp.MustCompile(strings.Join(generateRegexPattern([]string{"IoT", "Internet of Things", "device", "sensor", "network", "protocol", "communication", "edge computing", "cloud computing", "big data"}), "|"))
	educationRegex = regexp.MustCompile(strings.Join(generateRegexPattern([]string{"Education", "education", "learning", "tutorial", "course", "programming", "development", "project", "practice", "exercise", "exercise project"}), "|"))
	securityRegex = regexp.MustCompile(strings.Join(generateRegexPattern([]string{"Security", "security", "encryption", "cryptography", "authentication", "authorization", "access control", "security protocol", "security framework", "security tool"}), "|"))
	creativeRegex = regexp.MustCompile(strings.Join(generateRegexPattern([]string{"Creative", "creativity", "design", "graphic", "image", "video", "audio", "animation", "special effect", "special effect tool", "creative tool"}), "|"))
}

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
}

type OpenSourceInfo struct {
	repo OpenSourceRepo
	log  *log.Helper
	Page int
}

func NewOpenSourceInfo(repo OpenSourceRepo, logger log.Logger) *OpenSourceInfo {
	return &OpenSourceInfo{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

func generateRegexPattern(strings []string) []string {
	var patterns []string
	for _, s := range strings {
		patterns = append(patterns, regexp.QuoteMeta(s))
	}
	return patterns
}

func (r *OpenSourceInfo) updateRepoCategoryId(ctx context.Context, repoId, categoryId int64) {
	if !r.repo.FindRepoCategoryId(ctx, repoId, categoryId) {
		_ = r.repo.AddRepoCategoryId(ctx, repoId, categoryId)
	}
}

func (r *OpenSourceInfo) repoCategoryIdChange(ctx context.Context, desc string, repoId int64) {
	if frontendRegex.MatchString(desc) {
		// 前端相关逻辑
		r.updateRepoCategoryId(ctx, repoId, 1)
	}
	if backendRegex.MatchString(desc) {
		// 后端相关逻辑
		r.updateRepoCategoryId(ctx, repoId, 2)
	}
	if mobileDevelopmentRegex.MatchString(desc) {
		// 移动开发相关逻辑
		r.updateRepoCategoryId(ctx, repoId, 3)
	}
	if dataScienceRegex.MatchString(desc) {
		// 数据科学相关逻辑
		r.updateRepoCategoryId(ctx, repoId, 4)
	}
	if gameDevelopmentRegex.MatchString(desc) {
		// 游戏开发相关逻辑
		r.updateRepoCategoryId(ctx, repoId, 5)
	}
	if systemRegex.MatchString(desc) {
		// 系统相关逻辑
		r.updateRepoCategoryId(ctx, repoId, 6)
	}
	if libraryRegex.MatchString(desc) {
		// 库相关逻辑
		r.updateRepoCategoryId(ctx, repoId, 7)
	}
	if aiRegex.MatchString(desc) {
		// 人工智能相关逻辑
		r.updateRepoCategoryId(ctx, repoId, 8)
	}
	if blockchainRegex.MatchString(desc) {
		// 区块链相关逻辑
		r.updateRepoCategoryId(ctx, repoId, 9)
	}
	if iotRegex.MatchString(desc) {
		// 物联网相关逻辑
		r.updateRepoCategoryId(ctx, repoId, 10)
	}
	if educationRegex.MatchString(desc) {
		// 教育相关逻辑
		r.updateRepoCategoryId(ctx, repoId, 11)
	}
	if securityRegex.MatchString(desc) {
		// 安全相关逻辑
		r.updateRepoCategoryId(ctx, repoId, 12)
	}
	if creativeRegex.MatchString(desc) {
		// 创意相关逻辑
		r.updateRepoCategoryId(ctx, repoId, 13)
	}
}

func (r *OpenSourceInfo) RepoCategoryIdChange(ctx context.Context, topics []string, repoId int64) {
	for _, topic := range topics {
		if topic == "" {
			continue
		}
		r.repoCategoryIdChange(ctx, topic, repoId)
	}
	if !r.repo.FindRepoCategoryIdByRepoId(repoId) {
		r.updateRepoCategoryId(ctx, repoId, 14)
	}
}

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
	topicData, _ := json.Marshal(item.Topics)
	info := &domain.RepoInfo{
		Name:            item.Name,
		FullName:        item.FullName,
		Image:           repoImg,
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
		// 查找语言id
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
		Desc:            repoInfo.Desc,
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
		CreatedAt:       repoInfo.CreatedAt.Format(time.DateTime),
		UpdatedAt:       repoInfo.UpdatedAt.Format(time.DateTime),
	}
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

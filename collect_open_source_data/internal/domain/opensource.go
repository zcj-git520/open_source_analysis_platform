package domain

import (
	"time"
)

type RepoInfo struct {
	ID              int64     `gorm:"primarykey;type:int" json:"id"`
	Name            string    `gorm:"type:varchar(255)" json:"name"`           // 仓库名称
	FullName        string    `gorm:"type:varchar(255)" json:"full_name"`      // 仓库全称
	Image           string    `gorm:"type:varchar(255)" json:"image"`          // 仓库图片
	OwnerID         int64     `gorm:"type:int" json:"owner_id"`                // 仓库所有者ID
	Private         bool      `gorm:"type:tinyint" json:"private"`             // 是否私有
	Desc            string    `gorm:"type:MEDIUMTEXT" json:"desc"`             // 仓库描述
	DescZh          string    `gorm:"type:MEDIUMTEXT" json:"desc_zh"`          // 仓库中文描述
	HtmlURL         string    `gorm:"type:varchar(255)" json:"html_url"`       // 库主页 URL
	Homepage        string    `gorm:"type:varchar(255)" json:"homepage"`       // 仓库主页 URL
	CloneURL        string    `gorm:"type:varchar(255)" json:"clone_url"`      // 克隆 URL
	Size            int64     `gorm:"type:int" json:"size"`                    // 仓库大小
	StargazersCount int64     `gorm:"type:int" json:"stargazers_count"`        // 仓库星标数量
	WatchersCount   int64     `gorm:"type:int" json:"watchers_count"`          // 仓库关注者数量
	LanguageId      int64     `gorm:"type:int" json:"language_id"`             // 仓库语言ID
	ForksCount      int64     `gorm:"type:int" json:"forks_count"`             // 仓库分支数量
	OpenIssuesCount int64     `gorm:"type:int" json:"open_issues_count"`       // 仓库问题数量
	Topics          string    `gorm:"type:varchar(500)" json:"topics"`         // 仓库主题
	Forks           int64     `gorm:"type:int" json:"forks"`                   // 仓库分支数量
	OpenIssues      int64     `gorm:"type:int" json:"open_issues"`             // 仓库问题数量
	Watchers        int64     `gorm:"type:int" json:"watchers"`                // 仓库关注者数量
	DefaultBranch   string    `gorm:"type:varchar(255)" json:"default_branch"` // 默认分支
	Score           int64     `gorm:"type:int" json:"score"`                   // 仓库评分
	CreatedAt       time.Time `gorm:"type:datetime" json:"created_at"`         // 创建时间
	UpdatedAt       time.Time `gorm:"type:datetime" json:"updated_at"`         // 更新时间
}

func (RepoInfo) TableName() string {
	return "repo_info"
}

type Language struct {
	ID      int64  `gorm:"primarykey;type:int" json:"id"`
	Name    string `gorm:"type:varchar(255)" json:"name"`     // 语言名称
	Image   string `gorm:"type:varchar(255)" json:"image"`    // 语言图片
	Score   int64  `gorm:"type:int" json:"score"`             // 语言评分
	Desc    string `gorm:"type:varchar(255)" json:"desc"`     // 语言描述
	RepoRUL string `gorm:"type:varchar(255)" json:"repo_url"` // 仓库地址
	Bio     string `gorm:"type:varchar(255)" json:"bio"`      // 语言简介
}

func (Language) TableName() string {
	return "language"
}

type RepoTopic struct {
	ID     int64 `gorm:"primarykey;type:int" json:"id"`
	RepoID int64 `gorm:"type:int" json:"repo_id"` // 仓库ID
}

type Owner struct {
	ID          int64     `gorm:"primarykey;type:int" json:"id"`
	AvatarURL   string    `gorm:"type:varchar(255)" json:"avatar_url"` // 头像URL
	Type        string    `gorm:"type:varchar(255)" json:"type"`       // 类型
	Login       string    `gorm:"type:varchar(255)" json:"login"`      // 用户名
	HtmlURL     string    `gorm:"type:varchar(255)" json:"html_url"`   // 用户主页URL
	Name        string    `gorm:"type:varchar(255)" json:"name"`       // 用户名
	Email       string    `gorm:"type:varchar(255)" json:"email"`      // 用户邮箱
	Bio         string    `gorm:"type:varchar(500)" json:"bio"`        // 用户简介
	PublicRepos int64     `gorm:"type:int" json:"public_repos"`        // 用户公开仓库数量
	PublicGists int64     `gorm:"type:int" json:"public_gists"`        // 用户公开代码片段数量
	Followers   int64     `gorm:"type:int" json:"followers"`           // 用户粉丝数量
	Following   int64     `gorm:"type:int" json:"following"`           // 用户关注的人数量
	CreatedAt   time.Time `gorm:"type:datetime" json:"created_at"`     // 创建时间
	UpdatedAt   time.Time `gorm:"type:datetime" json:"updated_at"`     // 更新时间
}

func (Owner) TableName() string {
	return "owner"
}

type Page struct {
	PageNum  int32 `json:"pageNum"`  // 当前页，默认为1
	PageSize int32 `json:"pageSize"` // 分页条目数据，默认10
	Total    int64 `json:"total"`    // 查询总数量
}

func (p *Page) Offset() int {
	if p.PageSize < 1 {
		p.PageSize = 1
	}
	return int((p.PageNum - 1) * p.PageSize)
}

func (p *Page) Limit() int {
	if p.PageSize < 1 {
		p.PageSize = 10 // 当不传时，返回9999条数据
	}
	return int(p.PageSize)
}

type RepoCategory struct {
	ID       int64  `gorm:"primarykey;type:int" json:"id"`
	Name     string `gorm:"type:varchar(255)" json:"name"`      // 仓库分类名称
	ImageURL string `gorm:"type:varchar(255)" json:"image_url"` // 仓库分类图片
	Desc     string `gorm:"type:text" json:"desc"`              // 仓库分类描述
}

func (RepoCategory) TableName() string {
	return "repo_category"
}

type RepoCategoryId struct {
	ID     int64 `gorm:"primarykey;type:int" json:"id"`
	RepoID int64 `gorm:"type:int" json:"repo_id"` // 仓库ID
	CatID  int64 `gorm:"type:int" json:"cat_id"`  // 分类ID
}

func (RepoCategoryId) TableName() string {
	return "repo_category_id"
}

// 仓库指标值的增量
type RepoMetrics struct {
	ID          int64     `gorm:"primarykey;type:int" json:"id"`
	RepoID      int64     `gorm:"type:int" json:"repo_id"`      // 仓库ID
	Type        int       `gorm:"type:int" json:"type"`         // 指标类型 0: star 1: fork 2: watch 3: issue
	Value       int64     `gorm:"type:int" json:"value"`        // 增量值
	OriginValue int64     `gorm:"type:int" json:"origin_value"` // 原始值
	NowValue    int64     `gorm:"type:int" json:"now_value"`    // 当前值
	Date        time.Time `gorm:"type:datetime" json:"date"`    // 时间
}

type RepoMetricsResult struct {
	RepoID     int64
	TotalValue int64
}

func (RepoMetrics) TableName() string {
	return "repo_metrics"
}

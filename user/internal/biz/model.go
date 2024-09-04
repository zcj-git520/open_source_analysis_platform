package biz

import "time"

// Owner 表示仓库所有者信息
type Owner struct {
	Login             string `json:"login"` // 登录名
	ID                int    `json:"id"`    // 用户 ID
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`          // 用户头像 URL
	GravatarID        string `json:"gravatar_id"`         // Gravatar ID
	URL               string `json:"url"`                 // 用户 API URL
	HTMLURL           string `json:"html_url"`            // 用户主页 URL
	FollowersURL      string `json:"followers_url"`       // 用户粉丝列表 URL
	FollowingURL      string `json:"following_url"`       // 用户关注列表 URL
	GistsURL          string `json:"gists_url"`           // 用户 Gist 列表 URL
	StarredURL        string `json:"starred_url"`         // 用户收藏列表 URL
	SubscriptionsURL  string `json:"subscriptions_url"`   // 用户订阅列表 URL
	OrganizationsURL  string `json:"organizations_url"`   // 用户所属组织列表 URL
	ReposURL          string `json:"repos_url"`           // 用户仓库列表 URL
	EventsURL         string `json:"events_url"`          // 用户事件列表 URL
	ReceivedEventsURL string `json:"received_events_url"` // 用户接收的事件列表 URL
	Type              string `json:"type"`                // 用户类型（User 或 Organization）
	SiteAdmin         bool   `json:"site_admin"`          // 是否是站点管理员
}

type RepoUser struct {
	Login             string      `json:"login"`               // 用户登录名
	ID                int         `json:"id"`                  // 用户 ID
	NodeID            string      `json:"node_id"`             // 用户节点 ID
	AvatarURL         string      `json:"avatar_url"`          // 用户头像 URL
	GravatarID        string      `json:"gravatar_id"`         // Gravatar ID
	URL               string      `json:"url"`                 // 用户 API URL
	HtmlURL           string      `json:"html_url"`            // 用户主页 URL
	FollowersURL      string      `json:"followers_url"`       // 用户粉丝列表 URL
	FollowingURL      string      `json:"following_url"`       // 用户关注列表 URL
	GistsURL          string      `json:"gists_url"`           // 用户 Gist 列表 URL
	StarredURL        string      `json:"starred_url"`         // 用户收藏列表 URL
	SubscriptionsURL  string      `json:"subscriptions_url"`   // 用户订阅列表 URL
	OrganizationsURL  string      `json:"organizations_url"`   // 用户所属组织列表 URL
	ReposURL          string      `json:"repos_url"`           // 用户仓库列表 URL
	EventsURL         string      `json:"events_url"`          // 用户事件列表 URL
	ReceivedEventsURL string      `json:"received_events_url"` // 用户接收的事件列表 URL
	Type              string      `json:"type"`                // 用户类型（一般为 User）
	SiteAdmin         bool        `json:"site_admin"`          // 是否是站点管理员
	Name              string      `json:"name"`                // 用户名称
	Company           string      `json:"company"`             // 用户所在公司
	Blog              string      `json:"blog"`                // 用户博客地址
	Location          string      `json:"location"`            // 用户所在位置
	Email             interface{} `json:"email"`               // 用户邮箱（可能为 null）
	Hireable          bool        `json:"hireable"`            // 是否可雇佣
	Bio               string      `json:"bio"`                 // 用户简介
	TwitterUsername   string      `json:"twitter_username"`    // 用户推特用户名
	PublicRepos       int         `json:"public_repos"`        // 用户公开仓库数量
	PublicGists       int         `json:"public_gists"`        // 用户公开 Gist 数量
	Followers         int         `json:"followers"`           // 用户粉丝数量
	Following         int         `json:"following"`           // 用户关注数量
	CreatedAt         time.Time   `json:"created_at"`          // 用户创建时间
	UpdatedAt         time.Time   `json:"updated_at"`          // 用户更新时间
}

// License 表示许可证信息
type License struct {
	Key    string `json:"key"`     // 许可证关键字
	Name   string `json:"name"`    // 许可证名称
	SpdxID string `json:"spdx_id"` // SPDX ID
	URL    string `json:"url"`     // 许可证 URL
	NodeID string `json:"node_id"`
}

// Repo 表示单个仓库的信息
type Repo struct {
	ID                       int       `json:"id"`                          // 仓库 ID
	NodeID                   string    `json:"node_id"`                     // 仓库节点 ID
	Name                     string    `json:"name"`                        // 仓库名称
	FullName                 string    `json:"full_name"`                   // 仓库全名（所有者/仓库名）
	Private                  bool      `json:"private"`                     // 是否为私有仓库
	Owner                    *Owner    `json:"owner"`                       // 仓库所有者信息
	HtmlURL                  string    `json:"html_url"`                    // 仓库主页 URL
	Description              string    `json:"description"`                 // 仓库描述
	Fork                     bool      `json:"fork"`                        // 是否是 fork 的仓库
	URL                      string    `json:"url"`                         // 仓库 API URL
	ForksURL                 string    `json:"forks_url"`                   // 仓库 fork 列表 URL
	KeysURL                  string    `json:"keys_url"`                    // 仓库密钥列表 URL
	CollaboratorsURL         string    `json:"collaborators_url"`           // 仓库协作者列表 URL
	TeamsURL                 string    `json:"teams_url"`                   // 仓库团队列表 URL
	HooksURL                 string    `json:"hooks_url"`                   // 仓库钩子列表 URL
	IssueEventsURL           string    `json:"issue_events_url"`            // 仓库问题事件列表 URL
	EventsURL                string    `json:"events_url"`                  // 仓库事件列表 URL
	AssigneesURL             string    `json:"assignees_url"`               // 仓库指派人列表 URL
	BranchesURL              string    `json:"branches_url"`                // 仓库分支列表 URL
	TagsURL                  string    `json:"tags_url"`                    // 仓库标签列表 URL
	BlobsURL                 string    `json:"blobs_url"`                   // 仓库 Blob 对象列表 URL
	GitTagsURL               string    `json:"git_tags_url"`                // 仓库 Git 标签列表 URL
	GitRefsURL               string    `json:"git_refs_url"`                // 仓库 Git 引用列表 URL
	TreesURL                 string    `json:"trees_url"`                   // 仓库树对象列表 URL
	StatusesURL              string    `json:"statuses_url"`                // 仓库状态列表 URL
	LanguagesURL             string    `json:"languages_url"`               // 仓库语言列表 URL
	StargazersURL            string    `json:"stargazers_url"`              // 仓库星标者列表 URL
	ContributorsURL          string    `json:"contributors_url"`            // 仓库贡献者列表 URL
	SubscribersURL           string    `json:"subscribers_url"`             // 仓库订阅者列表 URL
	SubscriptionURL          string    `json:"subscription_url"`            // 仓库订阅 URL
	CommitsURL               string    `json:"commits_url"`                 // 仓库提交列表 URL
	GitCommitsURL            string    `json:"git_commits_url"`             // 仓库 Git 提交列表 URL
	CommentsURL              string    `json:"comments_url"`                // 仓库评论列表 URL
	IssueCommentURL          string    `json:"issue_comment_url"`           // 仓库问题评论列表 URL
	ContentsURL              string    `json:"contents_url"`                // 仓库内容列表 URL
	CompareURL               string    `json:"compare_url"`                 // 仓库比较 URL
	MergesURL                string    `json:"merges_url"`                  // 仓库合并请求列表 URL
	ArchiveURL               string    `json:"archive_url"`                 // 仓库归档 URL
	DownloadsURL             string    `json:"downloads_url"`               // 仓库下载列表 URL
	IssuesURL                string    `json:"issues_url"`                  // 仓库问题列表 URL
	PullsURL                 string    `json:"pulls_url"`                   // 仓库拉取请求列表 URL
	MilestonesURL            string    `json:"milestones_url"`              // 仓库里程碑列表 URL
	NotificationsURL         string    `json:"notifications_url"`           // 仓库通知列表 URL
	LabelsURL                string    `json:"labels_url"`                  // 仓库标签列表 URL
	ReleasesURL              string    `json:"releases_url"`                // 仓库发布列表 URL
	DeploymentsURL           string    `json:"deployments_url"`             // 仓库部署列表 URL
	CreatedAt                time.Time `json:"created_at"`                  // 仓库创建时间
	UpdatedAt                time.Time `json:"updated_at"`                  // 仓库更新时间
	PushedAt                 time.Time `json:"pushed_at"`                   // 仓库最后推送时间
	GitURL                   string    `json:"git_url"`                     // 仓库 Git URL
	SshURL                   string    `json:"ssh_url"`                     // 仓库 SSH URL
	CloneURL                 string    `json:"clone_url"`                   // 仓库克隆 URL
	SvnURL                   string    `json:"svn_url"`                     // 仓库 SVN URL
	Homepage                 string    `json:"homepage"`                    // 仓库主页
	Size                     int       `json:"size"`                        // 仓库大小
	StargazersCount          int       `json:"stargazers_count"`            // 仓库星标数量
	WatchersCount            int       `json:"watchers_count"`              // 仓库观察者数量
	Language                 string    `json:"language"`                    // 仓库主要语言
	HasIssues                bool      `json:"has_issues"`                  // 是否有问题
	HasProjects              bool      `json:"has_projects"`                // 是否有项目
	HasDownloads             bool      `json:"has_downloads"`               // 是否有下载
	HasWiki                  bool      `json:"has_wiki"`                    // 是否有维基
	HasPages                 bool      `json:"has_pages"`                   // 是否有页面
	HasDiscussions           bool      `json:"has_discussions"`             // 是否有讨论
	ForksCount               int       `json:"forks_count"`                 // 仓库 fork 数量
	MirrorURL                string    `json:"mirror_url"`                  // 仓库镜像 URL
	Archived                 bool      `json:"archived"`                    // 是否已归档
	Disabled                 bool      `json:"disabled"`                    // 是否已禁用
	OpenIssuesCount          int       `json:"open_issues_count"`           // 开放问题数量
	License                  License   `json:"license"`                     // 仓库许可证信息
	AllowForking             bool      `json:"allow_forking"`               // 是否允许 fork
	IsTemplate               bool      `json:"is_template"`                 // 是否是模板
	WebCommitSignoffRequired bool      `json:"web_commit_signoff_required"` // 是否需要网络提交签名
	Topics                   []string  `json:"topics"`                      // 仓库主题
	Visibility               string    `json:"visibility"`                  // 仓库可见性
	Forks                    int       `json:"forks"`                       // fork 数量
	OpenIssues               int       `json:"open_issues"`                 // 开放问题数量
	Watchers                 int       `json:"watchers"`                    // 观察者数量
	DefaultBranch            string    `json:"default_branch"`              // 默认分支
	Score                    float64   `json:"score"`                       // 得分
}

// RepoResponse 表示整个响应的结构体，包含仓库总数、是否有不完整结果以及仓库列表
type RepoResponse struct {
	TotalCount        int     `json:"total_count"`        // 仓库总数
	IncompleteResults bool    `json:"incomplete_results"` // 是否有不完整结果
	Items             []*Repo `json:"items"`              // 仓库列表
}

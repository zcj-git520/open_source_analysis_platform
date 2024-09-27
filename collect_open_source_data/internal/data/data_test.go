package data

import (
	"collect_open_source_data/internal/domain"
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"regexp"
	"strings"
	"testing"
)

//var frintend = []string{
//	"Frontend", "前端", "web", "网页", "界面", "UI", "UX", "HTML", "CSS", "JavaScript", "响应式设计", "React", "Vue.js",
//	"Angular", "TypeScript", "Sass/Less", "Webpack", "Babel", "Bootstrap", "Material-UI", "Ant Design", "Foundation",
//	"jQuery", "D3.js",
//}
//var backend = []string{"Backend", "后端", "服务器", "API", "Node.js", "Express.js", "Django", "Flask", "Ruby on Rails",
//	"Java", "Spring", "Go", "Python", "PHP", "C#", "ASP.NET", "Rust", "Elixir", "Scala", "Kotlin", "R", "Perl", "Lua", "Groovy",
//	"gin", "Backend", "后端", "Server", "服务器", "API", "应用程序接口", "REST", "RESTful API",
//	"GraphQL", "微服务", "Microservices", "数据库", "Database", "SQL", "NoSQL",
//	"MySQL", "PostgreSQL", "MongoDB", "Redis", "ORM", "对象关系映射",
//	"Node.js", "Express", "Koa", "Flask", "Django", "Spring Boot", "Laravel",
//	"Ruby on Rails", "Go", "Golang", "Java", "C#", ".NET", "PHP",
//	"Docker", "Kubernetes", "CI/CD", "持续集成/持续部署",
//	"Git", "版本控制", "GitHub", "GitLab", "Bitbucket",
//	"安全性", "Security", "HTTPS", "OAuth", "JWT", "加密", "Encryption",
//	"缓存", "Caching", "Redis", "Memcached",
//	"消息队列", "Message Queue", "RabbitMQ", "Kafka",
//	"日志", "Logging", "ELK Stack", "日志分析",
//	"监控", "Monitoring", "Prometheus", "Grafana",
//	"云服务", "Cloud Services", "AWS", "Azure", "GCP", "云函数", "Cloud Functions",
//	"容器化", "Containerization", "Docker", "Kubernetes",
//	"微服务架构", "Microservice Architecture", "服务发现", "Service Discovery",
//	"负载均衡", "Load Balancing", "Nginx", "HAProxy",
//	"反向代理", "Reverse Proxy", "API网关", "API Gateway",
//	"数据库迁移", "Database Migration", "数据备份", "Data Backup",
//	"单元测试", "Unit Testing", "集成测试", "Integration Testing",
//	"性能测试", "Performance Testing", "压力测试", "Stress Testing",
//}
//
//var mobileDevelopment = []string{"App", "应用程序", "移动应用", "iOS", "Android", "Flutter", "React Native", "Kotlin", "Swift", "Xamarin", "Cordova"}
//
//// 包含数据分析工具、机器学习算法实现等
//var dataScience = []string{"Data Science", "数据科学", "机器学习", "深度学习", "数据分析", "数据挖掘", "数据可视化", "TensorFlow",
//	"PyTorch", "Keras", "Scikit-learn", "Pandas", "NumPy", "Matplotlib", "Seaborn", "Jupyter Notebook", "R", "Spark",
//	"Hadoop"}
//
//// 涵盖各种游戏引擎和游戏项目
//var gameDevelopment = []string{"Game Development", "游戏开发", "游戏引擎", "Unity", "Unreal Engine", "Godot", "Cocos2d-x",
//	"游戏设计", "游戏编程", "游戏美术", "游戏音频", "游戏物理", "游戏AI", "游戏网络", "游戏测试", "游戏发布"}
//
//// 如操作系统相关项目、命令行工具等
//var system = []string{"System", "操作系统", "命令行工具", "Shell", "脚本", "脚本语言", "Shell脚本", "Bash", "Zsh", "Fish",
//	"PowerShell", "批处理脚本", "Cygwin", "MinGW", "Cygwin", "MinGW", "Windows Subsystem for Linux", "WSL", "Linux"}
//
//// 开源库和框架 为其他开发者提供功能模块的项目，如数据库连接库、图形处理库等
//var library = []string{"Library", "库", "框架", "工具", "实用程序", "实用工具", "实用库", "实用框架", "实用工具库", "实用工具框架"}
//
//// 包括深度学习模型、自然语言处理等项目
//var ai = []string{"AI", "人工智能", "机器学习", "深度学习", "自然语言处理", "计算机视觉", "语音识别", "强化学习", "神经网络",
//	"深度神经网络", "tensorflow", "pytorch", "ai", "发模型"}
//
//// 涉及区块链技术的实现和应用
//var blockchain = []string{"Blockchain", "区块链", "加密货币", "比特币", "以太坊", "智能合约", "共识算法", "分布式账本"}
//
//// 与物联网设备和系统相关的项目
//var iot = []string{"IoT", "物联网", "设备", "传感器", "网络", "协议", "通信", "边缘计算", "云计算", "大数据"}
//
//// 如编程教程、学习项目等
//var education = []string{"Education", "教育", "学习", "教程", "课程", "编程", "开发", "项目", "实践", "练习", "练习项目"}
//
//// 包括安全工具、加密算法实现等
//var security = []string{"Security", "安全", "加密", "密码学", "认证", "授权", "访问控制", "安全协议", "安全框架", "安全工具"}
//
//// 如图形设计工具、创意代码项目等
//var creative = []string{"Creative", "创意", "设计", "图形", "图像", "视频", "音频", "动画", "特效", "特效工具", "创意工具"}

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

var frontend []int64
var backend []int64
var mobileDevelopment []int64
var dataScience []int64
var gameDevelopment []int64
var system []int64
var library []int64
var ai []int64
var blockchain []int64
var iot []int64
var education []int64
var security []int64
var creative []int64

// 判断id是否在类别中
func isInCategory(id int64, category []int64) bool {
	for _, v := range category {
		if v == id {
			return true
		}
	}
	return false
}

func generateRegexPattern(strings []string) []string {
	var patterns []string
	for _, s := range strings {
		patterns = append(patterns, regexp.QuoteMeta(s))
	}
	return patterns
}

var db *gorm.DB

func conMysql() {
	Db, err := gorm.Open(mysql.Open("root:123456@tcp(127.0.0.1:3306)/osap_open_source?parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		os.Exit(1)
	}
	db = Db
}

func readRepoData() {
	size := 50
	num := 1
	for {
		var repo []*domain.RepoInfo
		db.Limit(size).Offset((num - 1) * size).Find(&repo)
		if len(repo) == 0 {
			println()
			return
		}
		for _, v := range repo {
			var topics []string
			err := json.Unmarshal([]byte(v.Topics), &topics)
			if err != nil {
				fmt.Println(v.ID, "json unmarshal error: ", err)
				continue
			}
			for _, topic := range topics {
				if topic == "" {
					continue
				}
				RepoCategoryIdChange(topic, v.ID)
			}
			if !FindRepoRepoId(v.ID) {
				updateRepoCategoryId(v.ID, 14)
			}
		}

		num++
	}

}

func FindRepoCategoryId(repoId, categoryId int64) bool {
	var repoCategoryId *domain.RepoCategoryId
	err := db.Where("repo_id = ? and cat_id = ?", repoId, categoryId).First(&repoCategoryId).Error
	if err == nil && repoCategoryId != nil {
		return true
	}
	return false
}

func FindRepoRepoId(repoId int64) bool {
	var repoCategoryId *domain.RepoCategoryId
	err := db.Where("repo_id = ?", repoId).First(&repoCategoryId).Error
	if err == nil && repoCategoryId != nil {
		return true
	}
	return false
}

func AddRepoCategoryId(repoId, categoryId int64) error {
	repoCategoryId := &domain.RepoCategoryId{
		RepoID: repoId,
		CatID:  categoryId,
	}
	return db.Create(repoCategoryId).Error
}

func updateRepoCategoryId(repoId, categoryId int64) {
	if !FindRepoCategoryId(repoId, categoryId) {
		_ = AddRepoCategoryId(repoId, categoryId)
	}
}

func RepoCategoryIdChange(desc string, repoId int64) {
	if frontendRegex.MatchString(desc) {
		// 前端相关逻辑
		updateRepoCategoryId(repoId, 1)
	}
	if backendRegex.MatchString(desc) {
		// 后端相关逻辑
		updateRepoCategoryId(repoId, 2)
	}
	if mobileDevelopmentRegex.MatchString(desc) {
		// 移动开发相关逻辑
		updateRepoCategoryId(repoId, 3)
	}
	if dataScienceRegex.MatchString(desc) {
		// 数据科学相关逻辑
		updateRepoCategoryId(repoId, 4)
	}
	if gameDevelopmentRegex.MatchString(desc) {
		// 游戏开发相关逻辑
		updateRepoCategoryId(repoId, 5)
	}
	if systemRegex.MatchString(desc) {
		// 系统相关逻辑
		updateRepoCategoryId(repoId, 6)
	}
	if libraryRegex.MatchString(desc) {
		// 库相关逻辑
		updateRepoCategoryId(repoId, 7)
	}
	if aiRegex.MatchString(desc) {
		// 人工智能相关逻辑
		updateRepoCategoryId(repoId, 8)
	}
	if blockchainRegex.MatchString(desc) {
		// 区块链相关逻辑
		updateRepoCategoryId(repoId, 9)
	}
	if iotRegex.MatchString(desc) {
		// 物联网相关逻辑
		updateRepoCategoryId(repoId, 10)
	}
	if educationRegex.MatchString(desc) {
		// 教育相关逻辑
		updateRepoCategoryId(repoId, 11)
	}
	if securityRegex.MatchString(desc) {
		// 安全相关逻辑
		updateRepoCategoryId(repoId, 12)
	}
	if creativeRegex.MatchString(desc) {
		// 创意相关逻辑
		updateRepoCategoryId(repoId, 13)
	}
}

func TestNewData(t *testing.T) {
	conMysql()
	readRepoData()

}

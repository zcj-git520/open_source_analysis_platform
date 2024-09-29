package biz

import (
	"context"
	"regexp"
	"strings"
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

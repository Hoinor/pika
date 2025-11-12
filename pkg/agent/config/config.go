package config

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"time"

	"gopkg.in/yaml.v3"
)

// Config Agent 配置
type Config struct {
	// 配置文件路径
	Path string `yaml:"-"`

	// 服务器配置
	Server ServerConfig `yaml:"server"`

	// Agent 配置
	Agent AgentConfig `yaml:"agent"`

	// 采集配置
	Collector CollectorConfig `yaml:"collector"`

	// 自动更新配置
	AutoUpdate AutoUpdateConfig `yaml:"auto_update"`

	// 监控配置
	Monitor MonitorConfig `yaml:"monitor"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	// 服务器地址（如：http://localhost:18888 或 https://your-server.com）
	Endpoint string `yaml:"endpoint"`

	// API Key
	APIKey string `yaml:"api_key"`

	// 是否跳过 TLS 证书验证（仅用于测试环境，生产环境不建议开启）
	InsecureSkipVerify bool `yaml:"insecure_skip_verify"`
}

// AgentConfig Agent 配置
type AgentConfig struct {
	// Agent 名称（默认使用主机名）
	Name string `yaml:"name"`
}

// CollectorConfig 采集器配置
type CollectorConfig struct {
	// 数据采集间隔（秒）
	Interval int `yaml:"interval"`

	// 心跳间隔（秒）
	HeartbeatInterval int `yaml:"heartbeat_interval"`

	// 网络采集排除的网卡列表（支持正则表达式）
	// 如果为空，默认只排除回环地址（lo, lo0）
	NetworkExclude []string `yaml:"network_exclude"`
}

// AutoUpdateConfig 自动更新配置
type AutoUpdateConfig struct {
	// 是否启用自动更新
	Enabled bool `yaml:"enabled"`

	// 检查更新间隔
	CheckInterval string `yaml:"check_interval"`
}

// MonitorConfig 监控配置
type MonitorConfig struct {
	// 是否启用监控功能
	Enabled bool `yaml:"enabled"`

	// 检测间隔（秒），默认 60 秒
	Interval int `yaml:"interval"`

	// 监控项列表
	Items []MonitorItem `yaml:"items"`
}

// MonitorItem 监控项
type MonitorItem struct {
	// 监控项名称
	Name string `yaml:"name"`

	// 监控类型: http, tcp
	Type string `yaml:"type"`

	// 目标地址
	// HTTP: 完整的 URL (如: https://example.com/api/health)
	// TCP: 地址:端口 (如: example.com:3306)
	Target string `yaml:"target"`

	// HTTP 特定配置
	HTTPConfig *HTTPMonitorConfig `yaml:"http_config,omitempty"`

	// TCP 特定配置
	TCPConfig *TCPMonitorConfig `yaml:"tcp_config,omitempty"`
}

// HTTPMonitorConfig HTTP 监控配置
type HTTPMonitorConfig struct {
	// HTTP 方法: GET, POST, PUT, DELETE 等，默认 GET
	Method string `yaml:"method"`

	// 期望的状态码，默认 200
	ExpectedStatusCode int `yaml:"expected_status_code"`

	// 期望的响应内容（关键字）
	ExpectedContent string `yaml:"expected_content,omitempty"`

	// 请求超时（秒），默认 60
	Timeout int `yaml:"timeout"`

	// 请求头
	Headers map[string]string `yaml:"headers,omitempty"`

	// 请求体
	Body string `yaml:"body,omitempty"`
}

// TCPMonitorConfig TCP 监控配置
type TCPMonitorConfig struct {
	// 连接超时（秒），默认 5
	Timeout int `yaml:"timeout"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Endpoint:           "http://localhost:8080",
			APIKey:             "",
			InsecureSkipVerify: false,
		},
		Agent: AgentConfig{
			Name: "",
		},
		Collector: CollectorConfig{
			Interval:          5,
			HeartbeatInterval: 30,
		},
		AutoUpdate: AutoUpdateConfig{
			Enabled:       true,
			CheckInterval: "10m",
		},
		Monitor: MonitorConfig{
			Enabled:  false,
			Interval: 60,
			Items:    []MonitorItem{},
		},
	}
}

// GetDefaultConfigPath 获取默认配置文件路径
func GetDefaultConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "."
	}
	return filepath.Join(homeDir, ".pika", "agent.yaml")
}

// Load 加载配置文件
func Load(path string) (*Config, error) {
	// 如果路径为空，使用默认路径
	if path == "" {
		path = GetDefaultConfigPath()
	}

	// 读取配置文件
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// 配置文件不存在，创建默认配置
			cfg := DefaultConfig()
			if err := cfg.Save(path); err != nil {
				return nil, fmt.Errorf("创建默认配置文件失败: %w", err)
			}
			return cfg, nil
		}
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析配置
	cfg := DefaultConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("配置验证失败: %w", err)
	}

	cfg.Path = path
	return cfg, nil
}

// Save 保存配置到文件
func (c *Config) Save(path string) error {
	// 如果路径为空，使用默认路径
	if path == "" {
		path = GetDefaultConfigPath()
	}

	// 确保目录存在
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}

	// 序列化配置
	data, err := yaml.Marshal(c)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	// 写入文件
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}

	return nil
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.Server.Endpoint == "" {
		return fmt.Errorf("服务器地址不能为空")
	}

	if c.Server.APIKey == "" {
		return fmt.Errorf("API Key 不能为空")
	}

	if c.Collector.Interval <= 0 {
		return fmt.Errorf("采集间隔必须大于 0")
	}

	if c.Collector.HeartbeatInterval <= 0 {
		return fmt.Errorf("心跳间隔必须大于 0")
	}

	if c.AutoUpdate.Enabled {
		if _, err := time.ParseDuration(c.AutoUpdate.CheckInterval); err != nil {
			return fmt.Errorf("更新检查间隔格式错误: %w", err)
		}
	}

	return nil
}

// GetCollectorInterval 获取采集间隔时长
func (c *Config) GetCollectorInterval() time.Duration {
	return time.Duration(c.Collector.Interval) * time.Second
}

// GetHeartbeatInterval 获取心跳间隔时长
func (c *Config) GetHeartbeatInterval() time.Duration {
	return time.Duration(c.Collector.HeartbeatInterval) * time.Second
}

// GetUpdateCheckInterval 获取更新检查间隔时长
func (c *Config) GetUpdateCheckInterval() time.Duration {
	duration, _ := time.ParseDuration(c.AutoUpdate.CheckInterval)
	return duration
}

// GetWebSocketURL 获取 WebSocket 连接地址
func (c *Config) GetWebSocketURL() string {
	u, err := url.Parse(c.Server.Endpoint)
	if err != nil {
		// 解析失败时，使用默认的 ws:// 协议
		return "ws://" + c.Server.Endpoint + "/ws/agent"
	}

	// 根据 HTTP 协议转换为对应的 WebSocket 协议
	scheme := "ws"
	if u.Scheme == "https" {
		scheme = "wss"
	}

	return fmt.Sprintf("%s://%s/ws/agent", scheme, u.Host)
}

// GetUpdateURL 获取更新检查地址
func (c *Config) GetUpdateURL() string {
	return c.Endpoint() + "/api/agent/version"
}

func (c *Config) GetDownloadURL() string {
	var filename = fmt.Sprintf("agent-%s-%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		filename += ".exe"
	}
	return c.Endpoint() + "/api/agent/downloads/" + filename
}

func (c *Config) Endpoint() string {
	u, err := url.Parse(c.Server.Endpoint)
	if err != nil {
		return c.Server.Endpoint
	}
	var endpoint = fmt.Sprintf("%s://%s", u.Scheme, u.Host)
	return endpoint
}

// GetNetworkExcludePatterns 获取网络排除的正则表达式列表
// 如果配置为空，返回默认排除规则（回环地址和常见虚拟网卡）
func (c *Config) GetNetworkExcludePatterns() ([]*regexp.Regexp, error) {
	patterns := c.Collector.NetworkExclude

	// 如果没有配置，使用默认排除规则
	if len(patterns) == 0 {
		patterns = []string{
			"^lo$",      // 回环地址
			"^lo0$",     // 回环地址
			"^docker.*", // Docker 网卡
			"^veth.*",   // 虚拟以太网接口
			"^br-.*",    // 网桥接口
		}
	}

	var regexps []*regexp.Regexp
	for _, pattern := range patterns {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return nil, fmt.Errorf("编译网络排除规则 '%s' 失败: %w", pattern, err)
		}
		regexps = append(regexps, re)
	}

	return regexps, nil
}

// ShouldExcludeNetworkInterface 检查网卡是否应该被排除
func (c *Config) ShouldExcludeNetworkInterface(interfaceName string) bool {
	patterns, err := c.GetNetworkExcludePatterns()
	if err != nil {
		// 如果正则编译失败，使用默认规则
		return interfaceName == "lo" || interfaceName == "lo0"
	}

	for _, pattern := range patterns {
		if pattern.MatchString(interfaceName) {
			return true
		}
	}

	return false
}

// GetMonitorInterval 获取监控检测间隔时长
func (c *Config) GetMonitorInterval() time.Duration {
	if c.Monitor.Interval <= 0 {
		return 60 * time.Second // 默认 60 秒
	}
	return time.Duration(c.Monitor.Interval) * time.Second
}

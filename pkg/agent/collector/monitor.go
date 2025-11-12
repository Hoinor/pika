package collector

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/dushixiang/pika/internal/protocol"
	"github.com/dushixiang/pika/pkg/agent/config"
)

// MonitorCollector 监控采集器
type MonitorCollector struct {
	cfg        *config.Config
	httpClient *http.Client
}

// NewMonitorCollector 创建监控采集器
func NewMonitorCollector(cfg *config.Config) *MonitorCollector {
	// 创建自定义的 HTTP 客户端，支持跳过 TLS 验证
	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true, // 允许自签名证书
			},
			DisableKeepAlives: true,
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// 限制重定向次数为 10
			if len(via) >= 10 {
				return fmt.Errorf("stopped after 10 redirects")
			}
			return nil
		},
	}

	return &MonitorCollector{
		cfg:        cfg,
		httpClient: httpClient,
	}
}

// Collect 采集所有监控项数据
func (c *MonitorCollector) Collect() ([]protocol.MonitorData, error) {
	if !c.cfg.Monitor.Enabled {
		return nil, nil
	}

	if len(c.cfg.Monitor.Items) == 0 {
		return nil, nil
	}

	results := make([]protocol.MonitorData, 0, len(c.cfg.Monitor.Items))

	for _, item := range c.cfg.Monitor.Items {
		var result protocol.MonitorData

		switch strings.ToLower(item.Type) {
		case "http", "https":
			result = c.checkHTTP(item)
		case "tcp":
			result = c.checkTCP(item)
		default:
			result = protocol.MonitorData{
				Name:      item.Name,
				Type:      item.Type,
				Target:    item.Target,
				Status:    "down",
				Error:     fmt.Sprintf("unsupported monitor type: %s", item.Type),
				CheckedAt: time.Now().UnixMilli(),
			}
		}

		results = append(results, result)
	}

	return results, nil
}

// checkHTTP 检查 HTTP/HTTPS 服务
func (c *MonitorCollector) checkHTTP(item config.MonitorItem) protocol.MonitorData {
	result := protocol.MonitorData{
		Name:      item.Name,
		Type:      item.Type,
		Target:    item.Target,
		CheckedAt: time.Now().UnixMilli(),
	}

	// 获取配置，使用默认值
	httpCfg := item.HTTPConfig
	if httpCfg == nil {
		httpCfg = &config.HTTPMonitorConfig{
			Method:             "GET",
			ExpectedStatusCode: 200,
			Timeout:            60,
		}
	}

	// 设置默认值
	method := httpCfg.Method
	if method == "" {
		method = "GET"
	}

	timeout := httpCfg.Timeout
	if timeout <= 0 {
		timeout = 60
	}

	expectedStatus := httpCfg.ExpectedStatusCode
	if expectedStatus == 0 {
		expectedStatus = 200
	}

	// 创建请求
	var bodyReader io.Reader
	if httpCfg.Body != "" {
		bodyReader = strings.NewReader(httpCfg.Body)
	}

	// 为请求创建带超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	// 为请求添加上下文
	req, err := http.NewRequestWithContext(ctx, method, item.Target, bodyReader)
	if err != nil {
		result.Status = "down"
		result.Error = fmt.Sprintf("create request failed: %v", err)
		return result
	}

	// 设置请求头
	if httpCfg.Headers != nil {
		for key, value := range httpCfg.Headers {
			req.Header.Set(key, value)
		}
	}

	// 发送请求并计时
	startTime := time.Now()
	resp, err := c.httpClient.Do(req)
	responseTime := time.Since(startTime).Milliseconds()
	result.ResponseTime = responseTime

	if err != nil {
		result.Status = "down"
		result.Error = fmt.Sprintf("request failed: %v", err)
		return result
	}
	defer resp.Body.Close()

	result.StatusCode = resp.StatusCode

	// 检查状态码
	if resp.StatusCode != expectedStatus {
		result.Status = "down"
		result.Error = fmt.Sprintf("status code mismatch: expected %d, got %d", expectedStatus, resp.StatusCode)
		result.Message = fmt.Sprintf("HTTP %d", resp.StatusCode)
		return result
	}

	// 检查响应内容（如果有配置）
	if httpCfg.ExpectedContent != "" {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			result.Status = "down"
			result.Error = fmt.Sprintf("read response body failed: %v", err)
			return result
		}

		bodyStr := string(body)
		if !strings.Contains(bodyStr, httpCfg.ExpectedContent) {
			result.Status = "down"
			result.Error = fmt.Sprintf("content does not contain expected string: %s", httpCfg.ExpectedContent)
			result.ContentMatch = false
			return result
		}
		result.ContentMatch = true
	}

	// 检查成功
	result.Status = "up"
	result.Message = fmt.Sprintf("HTTP %d - %dms", resp.StatusCode, responseTime)
	return result
}

// checkTCP 检查 TCP 端口
func (c *MonitorCollector) checkTCP(item config.MonitorItem) protocol.MonitorData {
	result := protocol.MonitorData{
		Name:      item.Name,
		Type:      item.Type,
		Target:    item.Target,
		CheckedAt: time.Now().UnixMilli(),
	}

	// 获取配置，使用默认值
	tcpCfg := item.TCPConfig
	timeout := 10 // 默认 10 秒
	if tcpCfg != nil && tcpCfg.Timeout > 0 {
		timeout = tcpCfg.Timeout
	}

	// 连接并计时
	startTime := time.Now()
	conn, err := net.DialTimeout("tcp", item.Target, time.Duration(timeout)*time.Second)
	responseTime := time.Since(startTime).Milliseconds()
	result.ResponseTime = responseTime

	if err != nil {
		result.Status = "down"
		result.Error = fmt.Sprintf("connection failed: %v", err)
		return result
	}
	defer conn.Close()

	// 连接成功
	result.Status = "up"
	result.Message = fmt.Sprintf("TCP connected - %dms", responseTime)
	return result
}

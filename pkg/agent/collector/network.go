package collector

import (
	"time"

	"github.com/dushixiang/pika/internal/protocol"
	"github.com/dushixiang/pika/pkg/agent/config"
	"github.com/shirou/gopsutil/v4/net"
)

// NetworkCollector 网络监控采集器
type NetworkCollector struct {
	config        *config.Config                // 配置信息
	lastStats     map[string]net.IOCountersStat // 上次采集的统计数据
	lastCollectAt time.Time                     // 上次采集时间
}

// safeDelta 计算网络计数器的增量,当出现重置或回绕时返回当前值避免溢出
func safeDelta(current, previous uint64) uint64 {
	if current >= previous {
		return current - previous
	}
	// 计数器被重置(例如接口重启)或发生回绕,只能依赖当前值
	return current
}

// calcRate 根据增量和采样间隔计算每秒速率
func calcRate(delta uint64, intervalSeconds float64) uint64 {
	if intervalSeconds <= 0 || delta == 0 {
		return 0
	}
	return uint64(float64(delta) / intervalSeconds)
}

// NewNetworkCollector 创建网络采集器
func NewNetworkCollector(cfg *config.Config) *NetworkCollector {
	return &NetworkCollector{
		config:    cfg,
		lastStats: make(map[string]net.IOCountersStat),
	}
}

// Collect 采集网络数据(计算自上次采集以来的增量)
func (n *NetworkCollector) Collect() ([]protocol.NetworkData, error) {
	now := time.Now()

	// 计算距离上次采集的时间间隔(秒)
	var intervalSeconds float64
	if !n.lastCollectAt.IsZero() {
		intervalSeconds = now.Sub(n.lastCollectAt).Seconds()
	}

	// 获取网络接口信息
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	// 创建接口信息映射
	interfaceMap := make(map[string]*protocol.NetworkData)
	for _, iface := range interfaces {
		// 使用配置中的排除规则过滤网卡
		if n.config.ShouldExcludeNetworkInterface(iface.Name) {
			continue
		}

		// 获取 IP 地址列表
		var addrs []string
		for _, addr := range iface.Addrs {
			addrs = append(addrs, addr.Addr)
		}

		interfaceMap[iface.Name] = &protocol.NetworkData{
			Interface:  iface.Name,
			MacAddress: iface.HardwareAddr,
			Addrs:      addrs,
		}
	}

	// 获取网络 IO 统计
	ioCounters, err := net.IOCounters(true)
	if err != nil {
		return nil, err
	}

	var networkDataList []protocol.NetworkData
	for _, counter := range ioCounters {
		// 使用配置中的排除规则过滤网卡
		if n.config.ShouldExcludeNetworkInterface(counter.Name) {
			continue
		}

		// 如果已有接口信息,则更新;否则创建新的
		netData := interfaceMap[counter.Name]
		if netData == nil {
			netData = &protocol.NetworkData{
				Interface: counter.Name,
			}
		}

		// 计算增量(如果是第一次采集,则使用当前值)
		lastStat, exists := n.lastStats[counter.Name]
		if exists {
			bytesSentDelta := safeDelta(counter.BytesSent, lastStat.BytesSent)
			bytesRecvDelta := safeDelta(counter.BytesRecv, lastStat.BytesRecv)
			netData.BytesSentRate = calcRate(bytesSentDelta, intervalSeconds)
			netData.BytesRecvRate = calcRate(bytesRecvDelta, intervalSeconds)
		} else {
			// 第一次采集无增量,速率保持为0
			netData.BytesSentRate = 0
			netData.BytesRecvRate = 0
		}
		netData.BytesSentTotal = counter.BytesSent
		netData.BytesRecvTotal = counter.BytesRecv

		// 保存当前统计数据用于下次计算增量
		n.lastStats[counter.Name] = counter

		networkDataList = append(networkDataList, *netData)
	}

	// 更新采集时间
	n.lastCollectAt = now

	return networkDataList, nil
}

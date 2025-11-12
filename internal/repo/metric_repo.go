package repo

import (
	"context"

	"github.com/dushixiang/pika/internal/models"
	"gorm.io/gorm"
)

type MetricRepo struct {
	db *gorm.DB
}

func NewMetricRepo(db *gorm.DB) *MetricRepo {
	return &MetricRepo{
		db: db,
	}
}

// SaveCPUMetric 保存CPU指标
func (r *MetricRepo) SaveCPUMetric(ctx context.Context, metric *models.CPUMetric) error {
	return r.db.WithContext(ctx).Create(metric).Error
}

// SaveMemoryMetric 保存内存指标
func (r *MetricRepo) SaveMemoryMetric(ctx context.Context, metric *models.MemoryMetric) error {
	return r.db.WithContext(ctx).Create(metric).Error
}

// SaveDiskMetric 保存磁盘指标
func (r *MetricRepo) SaveDiskMetric(ctx context.Context, metric *models.DiskMetric) error {
	return r.db.WithContext(ctx).Create(metric).Error
}

// SaveNetworkMetric 保存网络指标
func (r *MetricRepo) SaveNetworkMetric(ctx context.Context, metric *models.NetworkMetric) error {
	return r.db.WithContext(ctx).Create(metric).Error
}

// SaveLoadMetric 保存负载指标
func (r *MetricRepo) SaveLoadMetric(ctx context.Context, metric *models.LoadMetric) error {
	return r.db.WithContext(ctx).Create(metric).Error
}

// SaveDiskIOMetric 保存磁盘IO指标
func (r *MetricRepo) SaveDiskIOMetric(ctx context.Context, metric *models.DiskIOMetric) error {
	return r.db.WithContext(ctx).Create(metric).Error
}

// SaveGPUMetric 保存GPU指标
func (r *MetricRepo) SaveGPUMetric(ctx context.Context, metric *models.GPUMetric) error {
	return r.db.WithContext(ctx).Create(metric).Error
}

// SaveTemperatureMetric 保存温度指标
func (r *MetricRepo) SaveTemperatureMetric(ctx context.Context, metric *models.TemperatureMetric) error {
	return r.db.WithContext(ctx).Create(metric).Error
}

// SaveDockerMetric 保存Docker容器指标
func (r *MetricRepo) SaveDockerMetric(ctx context.Context, metric *models.DockerMetric) error {
	return r.db.WithContext(ctx).Create(metric).Error
}

// GetLatestDockerMetrics 获取最新的Docker容器指标列表
func (r *MetricRepo) GetLatestDockerMetrics(ctx context.Context, agentID string) ([]models.DockerMetric, error) {
	var metrics []models.DockerMetric
	// 获取每个容器的最新记录
	err := r.db.WithContext(ctx).
		Raw(`
			SELECT d1.* FROM docker_metrics d1
			INNER JOIN (
				SELECT container_id, MAX(timestamp) as max_timestamp
				FROM docker_metrics
				WHERE agent_id = ?
				GROUP BY container_id
			) d2 ON d1.container_id = d2.container_id AND d1.timestamp = d2.max_timestamp
			WHERE d1.agent_id = ?
			ORDER BY d1.name
		`, agentID, agentID).
		Scan(&metrics).Error
	return metrics, err
}

// GetLatestGPUMetrics 获取最新的GPU指标列表
func (r *MetricRepo) GetLatestGPUMetrics(ctx context.Context, agentID string) ([]models.GPUMetric, error) {
	var metrics []models.GPUMetric
	// 获取每个GPU的最新记录
	err := r.db.WithContext(ctx).
		Raw(`
			SELECT g1.* FROM gpu_metrics g1
			INNER JOIN (
				SELECT index, MAX(timestamp) as max_timestamp
				FROM gpu_metrics
				WHERE agent_id = ?
				GROUP BY index
			) g2 ON g1.index = g2.index AND g1.timestamp = g2.max_timestamp
			WHERE g1.agent_id = ?
			ORDER BY g1.index
		`, agentID, agentID).
		Scan(&metrics).Error
	return metrics, err
}

// GetLatestTemperatureMetrics 获取最新的温度指标列表
func (r *MetricRepo) GetLatestTemperatureMetrics(ctx context.Context, agentID string) ([]models.TemperatureMetric, error) {
	var metrics []models.TemperatureMetric
	// 获取每个传感器的最新记录
	err := r.db.WithContext(ctx).
		Raw(`
			SELECT t1.* FROM temperature_metrics t1
			INNER JOIN (
				SELECT sensor_key, MAX(timestamp) as max_timestamp
				FROM temperature_metrics
				WHERE agent_id = ?
				GROUP BY sensor_key
			) t2 ON t1.sensor_key = t2.sensor_key AND t1.timestamp = t2.max_timestamp
			WHERE t1.agent_id = ?
			ORDER BY t1.sensor_key
		`, agentID, agentID).
		Scan(&metrics).Error
	return metrics, err
}

// SaveHostMetric 保存主机信息指标（只保留最新的一条记录）
func (r *MetricRepo) SaveHostMetric(ctx context.Context, metric *models.HostMetric) error {
	// 使用事务确保原子性
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先删除该 agent 的所有旧记录
		if err := tx.Where("agent_id = ?", metric.AgentID).Delete(&models.HostMetric{}).Error; err != nil {
			return err
		}
		// 插入新记录
		return tx.Create(metric).Error
	})
}

// GetLatestHostMetric 获取最新的主机信息
func (r *MetricRepo) GetLatestHostMetric(ctx context.Context, agentID string) (*models.HostMetric, error) {
	var metric models.HostMetric
	err := r.db.WithContext(ctx).
		Where("agent_id = ?", agentID).
		Order("timestamp DESC").
		First(&metric).Error
	if err != nil {
		return nil, err
	}
	return &metric, nil
}

// DeleteOldMetrics 删除指定时间之前的所有指标数据
func (r *MetricRepo) DeleteOldMetrics(ctx context.Context, beforeTimestamp int64) error {
	// 批量大小
	batchSize := 1000

	// 定义要清理的表（Host 信息只保留最新的，不需要清理）
	tables := []interface{}{
		&models.CPUMetric{},
		&models.MemoryMetric{},
		&models.DiskMetric{},
		&models.NetworkMetric{},
		&models.LoadMetric{},
		&models.DiskIOMetric{},
		&models.GPUMetric{},
		&models.TemperatureMetric{},
		&models.DockerMetric{},
		&models.MonitorMetric{},
	}

	// 对每个表进行分批删除
	for _, table := range tables {
		for {
			// 分批删除，避免长事务
			result := r.db.WithContext(ctx).
				Where("timestamp < ?", beforeTimestamp).
				Limit(batchSize).
				Delete(table)

			if result.Error != nil {
				return result.Error
			}

			// 如果删除的行数少于批量大小，说明已经删除完毕
			if result.RowsAffected < int64(batchSize) {
				break
			}
		}
	}

	return nil
}

// AggregatedCPUMetric CPU聚合指标
type AggregatedCPUMetric struct {
	Timestamp    int64   `json:"timestamp"`
	AvgUsage     float64 `json:"avgUsage"`
	LogicalCores int     `json:"logicalCores"`
}

// GetCPUMetrics 获取聚合后的CPU指标（始终返回聚合数据）
// interval: 聚合间隔，单位秒（如：60表示1分钟，3600表示1小时）
func (r *MetricRepo) GetCPUMetrics(ctx context.Context, agentID string, start, end int64, interval int) ([]AggregatedCPUMetric, error) {
	var metrics []AggregatedCPUMetric

	query := `
		SELECT
			CAST(FLOOR(timestamp / ?) * ? AS BIGINT) as timestamp,
			AVG(usage_percent) as avg_usage,
			MAX(logical_cores) as logical_cores
		FROM cpu_metrics
		WHERE agent_id = ? AND timestamp >= ? AND timestamp <= ?
		GROUP BY 1
		ORDER BY timestamp ASC
	`

	intervalMs := int64(interval * 1000)
	err := r.db.WithContext(ctx).
		Raw(query, intervalMs, intervalMs, agentID, start, end).
		Scan(&metrics).Error

	return metrics, err
}

// AggregatedMemoryMetric 内存聚合指标
type AggregatedMemoryMetric struct {
	Timestamp int64   `json:"timestamp"`
	AvgUsage  float64 `json:"avgUsage"`
	Total     uint64  `json:"total"`
}

// GetMemoryMetrics 获取聚合后的内存指标（始终返回聚合数据）
func (r *MetricRepo) GetMemoryMetrics(ctx context.Context, agentID string, start, end int64, interval int) ([]AggregatedMemoryMetric, error) {
	var metrics []AggregatedMemoryMetric

	query := `
		SELECT
			CAST(FLOOR(timestamp / ?) * ? AS BIGINT) as timestamp,
			AVG(usage_percent) as avg_usage,
			MAX(total) as total
		FROM memory_metrics
		WHERE agent_id = ? AND timestamp >= ? AND timestamp <= ?
		GROUP BY 1
		ORDER BY timestamp ASC
	`

	intervalMs := int64(interval * 1000)
	err := r.db.WithContext(ctx).
		Raw(query, intervalMs, intervalMs, agentID, start, end).
		Scan(&metrics).Error

	return metrics, err
}

// AggregatedDiskMetric 磁盘聚合指标
type AggregatedDiskMetric struct {
	Timestamp  int64   `json:"timestamp"`
	MountPoint string  `json:"mountPoint"`
	AvgUsage   float64 `json:"avgUsage"`
	Total      uint64  `json:"total"`
}

// GetDiskMetrics 获取聚合后的磁盘指标（始终返回聚合数据）
func (r *MetricRepo) GetDiskMetrics(ctx context.Context, agentID string, start, end int64, interval int) ([]AggregatedDiskMetric, error) {
	var metrics []AggregatedDiskMetric

	query := `
		SELECT
			CAST(FLOOR(timestamp / ?) * ? AS BIGINT) as timestamp,
			mount_point,
			AVG(usage_percent) as avg_usage,
			MAX(total) as total
		FROM disk_metrics
		WHERE agent_id = ? AND timestamp >= ? AND timestamp <= ?
		GROUP BY 1, mount_point
		ORDER BY timestamp ASC, mount_point
	`

	intervalMs := int64(interval * 1000)
	err := r.db.WithContext(ctx).
		Raw(query, intervalMs, intervalMs, agentID, start, end).
		Scan(&metrics).Error

	return metrics, err
}

// AggregatedNetworkMetric 网络聚合指标
type AggregatedNetworkMetric struct {
	Timestamp   int64   `json:"timestamp"`
	AvgSentRate float64 `json:"avgSentRate"`
	AvgRecvRate float64 `json:"avgRecvRate"`
}

// GetNetworkMetrics 获取聚合后的网络指标（合并所有网卡接口）
func (r *MetricRepo) GetNetworkMetrics(ctx context.Context, agentID string, start, end int64, interval int) ([]AggregatedNetworkMetric, error) {
	var metrics []AggregatedNetworkMetric

	query := `
		SELECT
			CAST(FLOOR(timestamp / ?) * ? AS BIGINT) as timestamp,
			AVG(bytes_sent_rate) as avg_sent_rate,
			AVG(bytes_recv_rate) as avg_recv_rate
		FROM network_metrics
		WHERE agent_id = ? AND timestamp >= ? AND timestamp <= ?
		GROUP BY 1
		ORDER BY timestamp ASC
	`

	intervalMs := int64(interval * 1000)
	err := r.db.WithContext(ctx).
		Raw(query, intervalMs, intervalMs, agentID, start, end).
		Scan(&metrics).Error

	return metrics, err
}

// AggregatedLoadMetric 负载聚合指标
type AggregatedLoadMetric struct {
	Timestamp int64   `json:"timestamp"`
	AvgLoad1  float64 `json:"avgLoad1"`
	AvgLoad5  float64 `json:"avgLoad5"`
	AvgLoad15 float64 `json:"avgLoad15"`
}

// GetLoadMetrics 获取聚合后的负载指标（始终返回聚合数据）
func (r *MetricRepo) GetLoadMetrics(ctx context.Context, agentID string, start, end int64, interval int) ([]AggregatedLoadMetric, error) {
	var metrics []AggregatedLoadMetric

	query := `
		SELECT
			CAST(FLOOR(timestamp / ?) * ? AS BIGINT) as timestamp,
			AVG(load1) as avg_load1,
			AVG(load5) as avg_load5,
			AVG(load15) as avg_load15
		FROM load_metrics
		WHERE agent_id = ? AND timestamp >= ? AND timestamp <= ?
		GROUP BY 1
		ORDER BY timestamp ASC
	`

	intervalMs := int64(interval * 1000)
	err := r.db.WithContext(ctx).
		Raw(query, intervalMs, intervalMs, agentID, start, end).
		Scan(&metrics).Error

	return metrics, err
}

// GetLatestCPUMetric 获取最新的CPU指标
func (r *MetricRepo) GetLatestCPUMetric(ctx context.Context, agentID string) (*models.CPUMetric, error) {
	var metric models.CPUMetric
	err := r.db.WithContext(ctx).
		Where("agent_id = ?", agentID).
		Order("timestamp DESC").
		First(&metric).Error
	if err != nil {
		return nil, err
	}
	return &metric, nil
}

// GetLatestMemoryMetric 获取最新的内存指标
func (r *MetricRepo) GetLatestMemoryMetric(ctx context.Context, agentID string) (*models.MemoryMetric, error) {
	var metric models.MemoryMetric
	err := r.db.WithContext(ctx).
		Where("agent_id = ?", agentID).
		Order("timestamp DESC").
		First(&metric).Error
	if err != nil {
		return nil, err
	}
	return &metric, nil
}

// GetLatestDiskMetrics 获取最新的磁盘指标（所有挂载点）
func (r *MetricRepo) GetLatestDiskMetrics(ctx context.Context, agentID string) ([]models.DiskMetric, error) {
	// 先获取最新时间戳
	var latestTimestamp int64
	err := r.db.WithContext(ctx).
		Model(&models.DiskMetric{}).
		Where("agent_id = ?", agentID).
		Select("MAX(timestamp)").
		Scan(&latestTimestamp).Error

	if err != nil {
		return nil, err
	}

	// 获取该时间戳的所有磁盘数据
	var metrics []models.DiskMetric
	err = r.db.WithContext(ctx).
		Where("agent_id = ? AND timestamp = ?", agentID, latestTimestamp).
		Find(&metrics).Error

	return metrics, err
}

// GetLatestNetworkMetrics 获取最新的网络指标（所有网卡）
func (r *MetricRepo) GetLatestNetworkMetrics(ctx context.Context, agentID string) ([]models.NetworkMetric, error) {
	// 先获取最新时间戳
	var latestTimestamp int64
	err := r.db.WithContext(ctx).
		Model(&models.NetworkMetric{}).
		Where("agent_id = ?", agentID).
		Select("MAX(timestamp)").
		Scan(&latestTimestamp).Error

	if err != nil {
		return nil, err
	}

	// 获取该时间戳的所有网络数据
	var metrics []models.NetworkMetric
	err = r.db.WithContext(ctx).
		Where("agent_id = ? AND timestamp = ?", agentID, latestTimestamp).
		Find(&metrics).Error

	return metrics, err
}

// GetLatestLoadMetric 获取最新的负载指标
func (r *MetricRepo) GetLatestLoadMetric(ctx context.Context, agentID string) (*models.LoadMetric, error) {
	var metric models.LoadMetric
	err := r.db.WithContext(ctx).
		Where("agent_id = ?", agentID).
		Order("timestamp DESC").
		First(&metric).Error
	if err != nil {
		return nil, err
	}
	return &metric, nil
}

// SaveMonitorMetric 保存监控指标
func (r *MetricRepo) SaveMonitorMetric(ctx context.Context, metric *models.MonitorMetric) error {
	return r.db.WithContext(ctx).Create(metric).Error
}

// GetMonitorMetrics 获取监控指标列表
func (r *MetricRepo) GetMonitorMetrics(ctx context.Context, agentID, monitorName string, start, end int64) ([]models.MonitorMetric, error) {
	var metrics []models.MonitorMetric
	query := r.db.WithContext(ctx).
		Where("agent_id = ? AND timestamp >= ? AND timestamp <= ?", agentID, start, end)

	// 如果指定了监控项名称，则只查询该监控项
	if monitorName != "" {
		query = query.Where("name = ?", monitorName)
	}

	err := query.Order("timestamp ASC").Find(&metrics).Error
	return metrics, err
}

// GetLatestMonitorMetrics 获取最新的监控指标（每个监控项的最新一条）
func (r *MetricRepo) GetLatestMonitorMetrics(ctx context.Context, agentID string) ([]models.MonitorMetric, error) {
	// 先获取该 agent 下所有监控项名称
	var names []string
	err := r.db.WithContext(ctx).
		Model(&models.MonitorMetric{}).
		Where("agent_id = ?", agentID).
		Distinct("name").
		Pluck("name", &names).Error
	if err != nil {
		return nil, err
	}

	// 对每个监控项获取最新的一条记录
	var metrics []models.MonitorMetric
	for _, name := range names {
		var metric models.MonitorMetric
		err := r.db.WithContext(ctx).
			Where("agent_id = ? AND name = ?", agentID, name).
			Order("timestamp DESC").
			First(&metric).Error
		if err == nil {
			metrics = append(metrics, metric)
		}
	}

	return metrics, nil
}

// GetMonitorMetricsByName 获取指定监控项的历史数据
func (r *MetricRepo) GetMonitorMetricsByName(ctx context.Context, agentID, monitorName string, start, end int64, limit int) ([]models.MonitorMetric, error) {
	var metrics []models.MonitorMetric
	query := r.db.WithContext(ctx).
		Where("agent_id = ? AND name = ? AND timestamp >= ? AND timestamp <= ?", agentID, monitorName, start, end).
		Order("timestamp DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	err := query.Find(&metrics).Error
	return metrics, err
}

// AggregatedDiskIOMetric 磁盘IO聚合指标
type AggregatedDiskIOMetric struct {
	Timestamp       int64   `json:"timestamp"`
	Device          string  `json:"device"`
	AvgReadRate     float64 `json:"avgReadRate"`     // 平均读取速率(字节/秒)
	AvgWriteRate    float64 `json:"avgWriteRate"`    // 平均写入速率(字节/秒)
	TotalReadBytes  uint64  `json:"totalReadBytes"`  // 总读取字节数
	TotalWriteBytes uint64  `json:"totalWriteBytes"` // 总写入字节数
}

// GetDiskIOMetrics 获取聚合后的磁盘IO指标
func (r *MetricRepo) GetDiskIOMetrics(ctx context.Context, agentID string, start, end int64, interval int) ([]AggregatedDiskIOMetric, error) {
	var metrics []AggregatedDiskIOMetric

	// 计算速率需要根据时间差来计算，这里简化处理，直接计算平均值
	query := `
		SELECT
			CAST(FLOOR(timestamp / ?) * ? AS BIGINT) as timestamp,
			device,
			AVG(read_bytes_rate) as avg_read_rate,
			AVG(write_bytes_rate) as avg_write_rate,
			CASE
				WHEN MAX(read_bytes) >= MIN(read_bytes) THEN MAX(read_bytes) - MIN(read_bytes)
				ELSE MAX(read_bytes)
			END as total_read_bytes,
			CASE
				WHEN MAX(write_bytes) >= MIN(write_bytes) THEN MAX(write_bytes) - MIN(write_bytes)
				ELSE MAX(write_bytes)
			END as total_write_bytes
		FROM disk_io_metrics
		WHERE agent_id = ? AND timestamp >= ? AND timestamp <= ?
		GROUP BY 1, device
		ORDER BY timestamp ASC, device
	`

	intervalMs := int64(interval * 1000)
	err := r.db.WithContext(ctx).
		Raw(query, intervalMs, intervalMs, agentID, start, end).
		Scan(&metrics).Error

	return metrics, err
}

// AggregatedGPUMetric GPU聚合指标
type AggregatedGPUMetric struct {
	Timestamp      int64   `json:"timestamp"`
	AvgUtilization float64 `json:"avgUtilization"`
	AvgMemoryUsed  float64 `json:"avgMemoryUsed"`
	AvgTemperature float64 `json:"avgTemperature"`
	AvgPowerDraw   float64 `json:"avgPowerDraw"`
}

// GetGPUMetrics 获取聚合后的GPU指标
func (r *MetricRepo) GetGPUMetrics(ctx context.Context, agentID string, start, end int64, interval int) ([]AggregatedGPUMetric, error) {
	var metrics []AggregatedGPUMetric

	query := `
		SELECT
			CAST(FLOOR(timestamp / ?) * ? AS BIGINT) as timestamp,
			AVG(utilization) as avg_utilization,
			AVG(memory_used) as avg_memory_used,
			AVG(temperature) as avg_temperature,
			AVG(power_draw) as avg_power_draw
		FROM gpu_metrics
		WHERE agent_id = ? AND timestamp >= ? AND timestamp <= ?
		GROUP BY 1
		ORDER BY timestamp ASC
	`

	intervalMs := int64(interval * 1000)
	err := r.db.WithContext(ctx).
		Raw(query, intervalMs, intervalMs, agentID, start, end).
		Scan(&metrics).Error

	return metrics, err
}

// AggregatedTemperatureMetric 温度聚合指标
type AggregatedTemperatureMetric struct {
	Timestamp      int64   `json:"timestamp"`
	SensorKey      string  `json:"sensorKey"`
	SensorLabel    string  `json:"sensorLabel"`
	AvgTemperature float64 `json:"avgTemperature"`
}

// GetTemperatureMetrics 获取聚合后的温度指标
func (r *MetricRepo) GetTemperatureMetrics(ctx context.Context, agentID string, start, end int64, interval int) ([]AggregatedTemperatureMetric, error) {
	var metrics []AggregatedTemperatureMetric

	query := `
		SELECT
			CAST(FLOOR(timestamp / ?) * ? AS BIGINT) as timestamp,
			sensor_key,
			sensor_label,
			AVG(temperature) as avg_temperature
		FROM temperature_metrics
		WHERE agent_id = ? AND timestamp >= ? AND timestamp <= ?
		GROUP BY 1, sensor_key, sensor_label
		ORDER BY timestamp ASC, sensor_key
	`

	intervalMs := int64(interval * 1000)
	err := r.db.WithContext(ctx).
		Raw(query, intervalMs, intervalMs, agentID, start, end).
		Scan(&metrics).Error

	return metrics, err
}

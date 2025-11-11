import { useEffect } from 'react';
import { Form, Input, Switch, InputNumber, Button, Space, App, Card, Select, Divider } from 'antd';
import { Save } from 'lucide-react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import type { AlertConfig } from '../../types';
import { getAlertConfigsByAgent, createAlertConfig, updateAlertConfig } from '../../api/alert';
import { getAgents } from '../../api/agent';
import { getNotificationChannels } from '../../api/notification-channel';
import { getErrorMessage } from '../../lib/utils';

// 获取渠道类型的中文标签
const getChannelTypeLabel = (type: string): string => {
    const labels: Record<string, string> = {
        dingtalk: '钉钉',
        wecom: '企业微信',
        feishu: '飞书',
        webhook: '自定义Webhook',
        email: '邮件',
    };
    return labels[type] || type;
};

const AlertSettings = () => {
    const [form] = Form.useForm();
    const { message: messageApi } = App.useApp();
    const queryClient = useQueryClient();

    // 获取探针列表
    const { data: agentsData } = useQuery({
        queryKey: ['agents'],
        queryFn: getAgents,
    });
    const agents = agentsData?.items || [];

    // 获取通知渠道列表
    const { data: channels = [] } = useQuery({
        queryKey: ['notificationChannels'],
        queryFn: getNotificationChannels,
    });

    // 获取全局告警配置
    const { data: configsData, isLoading: configLoading } = useQuery({
        queryKey: ['alertConfigs', 'global'],
        queryFn: () => getAlertConfigsByAgent('global'),
    });

    const configId = configsData && configsData.length > 0 ? configsData[0].id : null;

    // 设置表单默认值
    useEffect(() => {
        if (configsData && configsData.length > 0) {
            const config = configsData[0];
            form.setFieldsValue(config);
        } else if (!configLoading) {
            form.setFieldsValue({
                name: '全局告警配置',
                enabled: true,
                agentIds: [],
                notificationChannelIds: [],
                rules: {
                    cpuEnabled: true,
                    cpuThreshold: 80,
                    cpuDuration: 60,
                    memoryEnabled: true,
                    memoryThreshold: 80,
                    memoryDuration: 60,
                    diskEnabled: true,
                    diskThreshold: 85,
                    diskDuration: 60,
                    networkEnabled: false,
                    networkDuration: 60,
                },
            });
        }
    }, [configsData, configLoading, form]);

    // 创建/更新 mutation
    const createMutation = useMutation({
        mutationFn: (config: AlertConfig) => createAlertConfig(config),
        onSuccess: () => {
            messageApi.success('告警配置创建成功');
            queryClient.invalidateQueries({ queryKey: ['alertConfigs', 'global'] });
        },
        onError: (error: unknown) => {
            messageApi.error(getErrorMessage(error, '保存配置失败'));
        },
    });

    const updateMutation = useMutation({
        mutationFn: ({ id, config }: { id: string; config: AlertConfig }) => updateAlertConfig(id, config),
        onSuccess: () => {
            messageApi.success('告警配置更新成功');
            queryClient.invalidateQueries({ queryKey: ['alertConfigs', 'global'] });
        },
        onError: (error: unknown) => {
            messageApi.error(getErrorMessage(error, '保存配置失败'));
        },
    });

    const handleSubmit = async () => {
        try {
            const values = await form.validateFields();
            const alertConfig: AlertConfig = {
                ...values,
                agentId: 'global',
            };

            if (configId) {
                updateMutation.mutate({ id: configId, config: alertConfig });
            } else {
                createMutation.mutate(alertConfig);
            }
        } catch (error) {
            // 表单验证失败
        }
    };

    return (
        <div>
            <Form form={form}>
                <Space direction="vertical" className="w-full">
                    <Card title="基本信息" type="inner">
                        <Form.Item label="配置名称" name="name" rules={[{ required: true }]}>
                            <Input placeholder="例如：全局告警配置" />
                        </Form.Item>
                        <Form.Item label="启用告警" name="enabled" valuePropName="checked">
                            <Switch checkedChildren="开启" unCheckedChildren="关闭" />
                        </Form.Item>
                        <Form.Item label="监控范围" name="agentIds" tooltip="留空表示监控所有探针">
                            <Select
                                mode="multiple"
                                placeholder="留空监控所有探针"
                                allowClear
                                options={agents.map((agent) => ({
                                    label: agent.name,
                                    value: agent.id,
                                }))}
                            />
                        </Form.Item>
                        <Form.Item
                            label="通知渠道"
                            name="notificationChannelIds"
                            rules={[{ required: true, message: '请选择至少一个通知渠道' }]}
                        >
                            <Select
                                mode="multiple"
                                placeholder="选择通知渠道"
                                options={channels
                                    .filter((ch) => ch.enabled)
                                    .map((ch) => ({
                                        label: getChannelTypeLabel(ch.type),
                                        value: ch.type,
                                    }))}
                            />
                        </Form.Item>
                    </Card>

                    <Divider orientation="left">告警规则</Divider>

                    {[
                        { key: 'cpu', title: 'CPU 告警规则', thresholdLabel: 'CPU 使用率阈值 (%)' },
                        { key: 'memory', title: '内存告警规则', thresholdLabel: '内存使用率阈值 (%)' },
                        { key: 'disk', title: '磁盘告警规则', thresholdLabel: '磁盘使用率阈值 (%)' },
                    ].map((rule) => (
                        <Card key={rule.key} title={rule.title} type="inner">
                            <Form.Item noStyle shouldUpdate>
                                {({ getFieldValue }) => {
                                    const enabled = getFieldValue(['rules', `${rule.key}Enabled`]);
                                    return (
                                        <div className="flex items-center gap-8">
                                            <Form.Item
                                                label="开关"
                                                name={['rules', `${rule.key}Enabled`]}
                                                valuePropName="checked"
                                                className="mb-0"
                                            >
                                                <Switch />
                                            </Form.Item>
                                            <Form.Item
                                                label={rule.thresholdLabel}
                                                name={['rules', `${rule.key}Threshold`]}
                                                className="mb-0"
                                            >
                                                <InputNumber
                                                    min={0}
                                                    max={100}
                                                    style={{ width: '100%' }}
                                                    disabled={!enabled}
                                                />
                                            </Form.Item>
                                            <Form.Item
                                                label="持续时间（秒）"
                                                name={['rules', `${rule.key}Duration`]}
                                                className="mb-0"
                                            >
                                                <InputNumber min={1} max={3600} style={{ width: '100%' }} disabled={!enabled} />
                                            </Form.Item>
                                        </div>
                                    );
                                }}
                            </Form.Item>
                        </Card>
                    ))}

                    <div className="flex justify-end pt-4">
                        <Button
                            type="primary"
                            icon={<Save size={16} />}
                            loading={createMutation.isPending || updateMutation.isPending}
                            onClick={handleSubmit}
                        >
                            保存配置
                        </Button>
                    </div>
                </Space>
            </Form>
        </div>
    );
};

export default AlertSettings;

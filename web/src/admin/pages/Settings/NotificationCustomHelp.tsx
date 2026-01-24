import React from 'react';
import {Collapse, type CollapseProps} from "antd";

const NotificationCustomHelp = () => {

    const items: CollapseProps['items'] = [
        {
            key: 'help',
            label: '请求体格式说明',
            children: <div className={'space-y-3 text-sm'}>
                {/* JSON 格式说明 */}
                <div className={'space-y-1'}>
                    <strong>1. JSON 格式 (默认)：</strong>
                    <div className={'text-gray-600 text-xs'}>
                        发送 <code className={'bg-gray-100 px-1 rounded'}>application/json</code> 格式的数据
                    </div>
                    <pre className={'border p-2 rounded-md text-xs mt-1 bg-gray-50'}>
                                    {JSON.stringify({
                                        "msg_type": "text",
                                        "text": {"content": "告警消息内容"},
                                        "agent": {
                                            "id": "agent-id",
                                            "name": "探针名称",
                                            "hostname": "主机名",
                                            "ip": "192.168.1.1"
                                        },
                                        "alert": {
                                            "type": "cpu",
                                            "level": "warning",
                                            "status": "firing",
                                            "message": "CPU使用率过高",
                                            "threshold": 80,
                                            "actualValue": 85.5,
                                            "firedAt": 1234567890000,
                                            "resolvedAt": 0
                                        }
                                    }, null, 2)}
                                </pre>
                </div>

                {/* Form 表单格式说明 */}
                <div className={'space-y-1'}>
                    <strong>2. Form 表单格式：</strong>
                    <div className={'text-gray-600 text-xs'}>
                        发送 <code
                        className={'bg-gray-100 px-1 rounded'}>application/x-www-form-urlencoded</code> 格式的数据
                    </div>
                    <div className={'border p-2 rounded-md text-xs mt-1 bg-gray-50'}>
                        <div className={'font-semibold mb-1'}>包含以下字段：</div>
                        <div className={'grid grid-cols-2 gap-x-4 gap-y-1'}>
                            <div>• <code>message</code> - 告警消息</div>
                            <div>• <code>agent_id</code> - 探针ID</div>
                            <div>• <code>agent_name</code> - 探针名称</div>
                            <div>• <code>agent_hostname</code> - 主机名</div>
                            <div>• <code>agent_ip</code> - IP地址(合并)</div>
                            <div>• <code>agent_ipv4</code> - IPv4 地址</div>
                            <div>• <code>agent_ipv6</code> - IPv6 地址</div>
                            <div>• <code>alert_type</code> - 告警类型</div>
                            <div>• <code>alert_level</code> - 告警级别</div>
                            <div>• <code>alert_status</code> - 告警状态</div>
                            <div>• <code>alert_message</code> - 告警详情</div>
                            <div>• <code>threshold</code> - 阈值</div>
                            <div>• <code>actual_value</code> - 当前值</div>
                            <div>• <code>fired_at</code> - 触发时间(格式化)</div>
                            <div>• <code>resolved_at</code> - 恢复时间(格式化)</div>

                        </div>
                    </div>
                </div>

                {/* 自定义模板说明 */}
                <div className={'space-y-1'}>
                    <strong>3. 自定义模板：</strong>
                    <div className={'text-gray-600 text-xs'}>
                        支持变量替换，Content-Type 为 <code
                        className={'bg-gray-100 px-1 rounded'}>text/plain</code>
                    </div>
                    <div className={'border p-2 rounded-md text-xs mt-1 bg-gray-50'}>
                        <div className={'font-semibold mb-1'}>可用变量：</div>
                        <div className={'grid grid-cols-2 gap-x-4 gap-y-1'}>
                            <div>• <code>{`{{message}}`}</code> - 告警消息</div>
                            <div>• <code>{`{{agent.id}}`}</code> - 探针ID</div>
                            <div>• <code>{`{{agent.name}}`}</code> - 探针名称</div>
                            <div>• <code>{`{{agent.hostname}}`}</code> - 主机名</div>
                            <div>• <code>{`{{agent.ip}}`}</code> - IP地址(合并)</div>
                            <div>• <code>{`{{agent.ipv4}}`}</code> - IPv4 地址</div>
                            <div>• <code>{`{{agent.ipv6}}`}</code> - IPv6 地址</div>
                            <div>• <code>{`{{alert.type}}`}</code> - 告警类型</div>
                            <div>• <code>{`{{alert.level}}`}</code> - 告警级别</div>
                            <div>• <code>{`{{alert.status}}`}</code> - 告警状态</div>
                            <div>• <code>{`{{alert.message}}`}</code> - 告警消息</div>
                            <div>• <code>{`{{alert.threshold}}`}</code> - 阈值</div>
                            <div>• <code>{`{{alert.actualValue}}`}</code> - 当前值</div>
                            <div>• <code>{`{{alert.firedAt}}`}</code> - 触发时间(格式化)</div>
                            <div>• <code>{`{{alert.resolvedAt}}`}</code> - 恢复时间(格式化)</div>
                        </div>
                        <div className={'mt-2 pt-2 border-t'}>
                            <div className={'font-semibold mb-1'}>示例：</div>
                            <pre className={'text-xs'}>
                                            {`{
  "alert": "{{alert.message}}",
  "host": "{{agent.hostname}}",
  "level": "{{alert.level}}"
}`}
                                        </pre>
                        </div>
                    </div>
                </div>
            </div>,
        },
    ];

    return <Collapse
        bordered={false}
        items={items}
    />;
};

export default NotificationCustomHelp;
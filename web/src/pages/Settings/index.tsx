import {Tabs} from 'antd';
import {Bell, MessageSquare} from 'lucide-react';
import AlertSettings from './AlertSettings';
import NotificationChannels from './NotificationChannels';
import {PageHeader} from "../../components";

const Settings = () => {
    const items = [
        {
            key: 'channels',
            label: (
                <span className="flex items-center gap-2">
                    <MessageSquare size={16}/>
                    通知渠道
                </span>
            ),
            children: <NotificationChannels/>,
        },
        {
            key: 'alert',
            label: (
                <span className="flex items-center gap-2">
                    <Bell size={16}/>
                    告警规则
                </span>
            ),
            children: <AlertSettings/>,
        },
    ];

    return (
        <div className={'space-y-6'}>
            <PageHeader
                title="系统设置"
                description="CONFIGURATION"
            />
            <Tabs defaultActiveKey="channels"
                  tabPosition={'left'}
                  items={items}
            />
        </div>
    );
};

export default Settings;

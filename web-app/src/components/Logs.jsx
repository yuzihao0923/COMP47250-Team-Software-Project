import React from 'react';

import { List, Typography } from 'antd';

const { Text } = Typography;

export default function Logs(props) {

    const { logsTitle, logsBackgroundColor, logsData} = props

    return (
        <div>
            <h2 className='mb-3 text-gray-500 font-medium'>{logsTitle}</h2>
            <div className={`${logsBackgroundColor} py-5 px-5 mb-10 max-h-60 overflow-y-auto overflow-x-hidden`}>
                <List
                    itemLayout='horizontal'
                    dataSource={logsData}
                    renderItem={(item) => (
                        <List.Item>
                            {item.trim().startsWith('[ERROR]') ? <Text type="danger">{item}</Text>
                                :
                                item.trim().startsWith('[WARNING]') ? <Text type="warning">{item}</Text> : item}
                        </List.Item>
                    )}
                />
            </div>
        </div>
    )
}

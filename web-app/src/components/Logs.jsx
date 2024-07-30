import React, { useRef } from 'react';
import { Button, List, Typography } from 'antd';

const { Text } = Typography;

export default function Logs(props) {
    const { logsTitle, logsData, style } = props;
    const downloadLinkRef = useRef(null);

    const exportLogs = () => {
        const file = new Blob([logsData.join('\n')], { type: 'text/plain' });
        const url = URL.createObjectURL(file);
        downloadLinkRef.current.href = url;
        downloadLinkRef.current.download = `${logsTitle}.txt`;
        downloadLinkRef.current.click();
        URL.revokeObjectURL(url);
    };

    return (
        <div style={{style, marginBottom: '20px'}}>
            <div className='flex justify-between items-center mb-3'>
                <h2 style={{ color: 'white' }}>{logsTitle}</h2>
                <Button ghost onClick={exportLogs} style={{ color: 'white', borderColor: 'rgba(255, 255, 255, 0.3)' }}>Export Logs</Button>
            </div>
            <div className="py-3 px-4" style={{ backgroundColor: 'rgba(0, 0, 0, 0.7)', overflowY: 'auto', maxHeight: '300px', border: '1px solid rgba(255, 255, 255, 0.2)' }}>
                <List
                    itemLayout='horizontal'
                    dataSource={logsData}
                    renderItem={(item) => (
                        <List.Item>
                            <Text style={{ color: 'cyan', fontSize: '18px' }}>
                                {item.trim().startsWith('[ERROR]') ? <Text type="danger">{item}</Text>
                                    :
                                    item.trim().startsWith('[WARNING]') ? <Text type="warning">{item}</Text> : item}
                            </Text>
                        </List.Item>
                    )}
                />
            </div>
        </div>
    )
}

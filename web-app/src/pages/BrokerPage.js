// src/pages/BrokerPage.js
import React, { useState, useEffect } from 'react';
import messageService from '../services/message';

const BrokerPage = () => {
    const [logs, setLogs] = useState([]);

    useEffect(() => {
        const fetchLogs = async () => {
            try {
                const brokerLogs = await messageService.getBrokerLogs();
                setLogs(brokerLogs);
            } catch (error) {
                console.error('Fetch broker logs failed:', error);
            }
        };
        fetchLogs();
    }, []);

    return (
        <div>
            <h2>Broker Console</h2>
            <ul>
                {logs.map((log, index) => (
                    <li key={index}>{log}</li>
                ))}
            </ul>
        </div>
    );
};

export default BrokerPage;

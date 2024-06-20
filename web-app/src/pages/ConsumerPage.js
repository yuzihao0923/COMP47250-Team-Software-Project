// src/pages/ConsumerPage.js
import React, { useState, useEffect } from 'react';
import messageService from '../services/message';

const ConsumerPage = () => {
    const [messages, setMessages] = useState([]);

    useEffect(() => {
        const fetchMessages = async () => {
            try {
                const msgs = await messageService.consumeMessages();
                setMessages(msgs);
            } catch (error) {
                console.error('Fetch messages failed:', error);
            }
        };
        fetchMessages();
    }, []);

    return (
        <div>
            <h2>Consumer</h2>
            <ul>
                {messages.map((msg, index) => (
                    <li key={index}>{msg}</li>
                ))}
            </ul>
        </div>
    );
};

export default ConsumerPage;

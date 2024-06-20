// src/pages/ProducerPage.js
import React, { useState } from 'react';
import messageService from '../services/message';

const ProducerPage = () => {
    const [message, setMessage] = useState('');

    const handleSendMessage = async () => {
        try {
            await messageService.sendMessage(message);
            alert('Message sent!');
            setMessage('');
        } catch (error) {
            console.error('Send message failed:', error);
            alert('Send message failed');
        }
    };

    return (
        <div>
            <h2>Producer</h2>
            <textarea
                value={message}
                onChange={(e) => setMessage(e.target.value)}
                placeholder="Type your message here"
            />
            <button onClick={handleSendMessage}>Send Message</button>
        </div>
    );
};

export default ProducerPage;

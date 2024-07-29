import React from 'react';
import { Button } from 'antd';
import { LogoutOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { useDispatch, useSelector } from 'react-redux';
import { logout as reduxLogout } from '../store/userSlice';

export default function Head() {
    const { username } = useSelector(state => state.user);
    const navigate = useNavigate();
    const dispatch = useDispatch();

    function logout() {
        dispatch(reduxLogout());
        navigate('/');
    }

    return (
        <div className='bg-gray-900 border-b border-gray-700 shadow-md sticky top-0 z-50'>
            <header className='flex justify-between items-center max-w-6xl mx-auto h-20 max-h-20'>
                <div className='flex-shrink-0'>
                    <h1 className='hollow-fluorescent-edge text-2xl font-bold'>Hi, {username}</h1>
                </div>
                <div className='h-full flex items-center space-x-6 flex-nowrap'>
                    <Button type="default" ghost size={'middle'} icon={<LogoutOutlined />} onClick={logout} style={{ color: 'white', borderColor: 'rgba(255, 255, 255, 0.3)' }}>
                        Logout
                    </Button>
                </div>
            </header>
        </div>
    );
}

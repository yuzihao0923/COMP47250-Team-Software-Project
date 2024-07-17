import React from 'react';
import { Outlet, useNavigate } from 'react-router';
import Head from '../components/Header';
import { Menu } from 'antd';
import { SettingOutlined, BugOutlined, UserOutlined } from '@ant-design/icons';

const items = [
    {
        label: 'Components Logs',
        key: '/',
        icon: <BugOutlined />,
    },
    {
        label: 'Profile',
        key: 'profile',
        icon: <UserOutlined />
    },
    {
        label: 'Settings',
        key: 'settings',
        icon: <SettingOutlined />,
    },
];

export default function Home() {

    const navigate = useNavigate()

    function onMenuClick(e) {
        // console.log(location);
        // console.log('click ', e);
        if (e.key === '/') {
            navigate('/')
        } else {
            navigate(`/${e.key}`)
        }
    }
    return (
        <div>
            <Head />
            <div className='max-w-6xl mx-auto mt-6'>
                
                <Menu 
                    mode="horizontal" 
                    items={items} 
                    theme='dark' 
                    style={{
                        height: 80,
                        display: 'flex',
                        alignItems: 'center',
                        backgroundColor: '#272b2c', 
                        color: '#fff'
                    }}
                    defaultSelectedKeys={['/']}
                    onClick={onMenuClick}
                />
               
            </div>
            <Outlet />
        </div>
    );
}

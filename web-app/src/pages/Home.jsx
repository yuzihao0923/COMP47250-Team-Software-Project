import React from 'react';
import { Outlet } from 'react-router';
import Head from '../components/Header';
// import { Menu } from 'antd';
// import { AppstoreOutlined, BugOutlined } from '@ant-design/icons';

// const items = [
//     {
//         label: 'Components Logs',
//         key: 'logs',
//         icon: <BugOutlined />,
//     },
//     {
//         label: 'Jasmine',
//         key: 'app',
//         icon: <AppstoreOutlined />,
//     },
// ];

export default function Home() {
    // function onMenuClick(e) {
    //     console.log('click ', e);
    // }

    return (
        <div>
            <Head />
            <div className='max-w-6xl mx-auto mt-6'>
                {/* 
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
                    defaultSelectedKeys={['logs']}
                    onClick={onMenuClick}
                />
                */}
            </div>
            <Outlet />
        </div>
    );
}

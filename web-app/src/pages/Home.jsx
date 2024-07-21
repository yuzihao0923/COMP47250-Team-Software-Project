import React, { useEffect, useState } from 'react';
import { Outlet, useLocation, useNavigate } from 'react-router';
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
    const location = useLocation()
    const [menuKey, setMenuKey] = useState(null)

    function onMenuClick(e) {
        // console.log(location);
        // console.log('click ', e);
        if (e.key === '/') {
            navigate('/')
        } else {
            navigate(`/${e.key}`)
        }
    }

    useEffect(()=>{
        const currentPath = location.pathname
        if(currentPath === '/'){
            setMenuKey(currentPath)
        }else{
            setMenuKey(currentPath.replace('/',''))
        }
    },[location.pathname])

    if(menuKey === null){
        return null;
    }

    return (
        <div>
            <Head />
            <div className='max-w-6xl mx-auto mt-6'>
                
                <Menu 
                    mode="horizontal" 
                    items={items} 
                    style={{
                        height: 80,
                        paddingLeft: '1.25rem',
                        paddingRight: '1.25rem',
                        display: 'flex',
                        alignItems: 'center',
                        backgroundColor: 'rgb(233, 228, 240)'
                    }}
                    selectedKeys={[path]}
                    onClick={onMenuClick}
                />
               
            </div>
            <Outlet />
        </div>
    );
}

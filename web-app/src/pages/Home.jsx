import React, { useState, useEffect } from 'react';
import { Outlet, useLocation, useNavigate } from 'react-router';
import Head from '../components/Header';
import { Menu } from 'antd';
import { SettingOutlined, BugOutlined, UserOutlined } from '@ant-design/icons';

const iconStyle = {
    color: '#0ff',
    textShadow: '0 0 8px #0ff',
    fontSize: '24px',
};

const items = [
    {
        label: 'Components Logs',
        key: '/',
        icon: <BugOutlined style={iconStyle} />,
        className: 'hollow-fluorescent-edge',
    },
    {
        label: 'Profile',
        key: 'profile',
        icon: <UserOutlined style={iconStyle} />,
        className: 'hollow-fluorescent-edge',
    },
    {
        label: 'Settings',
        key: 'settings',
        icon: <SettingOutlined style={iconStyle} />,
        className: 'hollow-fluorescent-edge',
    },
];

export default function Home() {
    const navigate = useNavigate();
    const location = useLocation();
    const [menuKey, setMenuKey] = useState(null);
    const [isNavOpen, setIsNavOpen] = useState(false);

    useEffect(() => {
        const currentPath = location.pathname;
        if (currentPath === '/') {
            setMenuKey(currentPath);
        } else {
            setMenuKey(currentPath.replace('/', ''));
        }
    }, [location.pathname]);

    const onMenuClick = (e) => {
        if (e.key === '/') {
            navigate('/');
        } else {
            navigate(`/${e.key}`);
        }
        setIsNavOpen(false);
    };

    const toggleNav = () => {
        setIsNavOpen(!isNavOpen);
    };

    const closeNav = (e) => {
        if (isNavOpen && e.target === e.currentTarget) {
            setIsNavOpen(false);
        }
    };

    if (menuKey === null) {
        return null;
    }

    return (
        <div onClick={closeNav} style={{ backgroundColor: '#121212', minHeight: '100vh', width: '100%' }}>
            <Head />
            <div style={{
                position: 'absolute',
                top: '20px',
                left: isNavOpen ? '256px' : '10px',
                cursor: 'pointer',
                zIndex: '1000'
            }} onClick={toggleNav}>
                <div style={{
                    width: '0',
                    height: '0',
                    borderTop: '15px solid transparent',
                    borderBottom: '15px solid transparent',
                    borderLeft: `15px solid ${isNavOpen ? '#121212' : '#fff'}`,
                }}></div>
            </div>
            <div style={{
                width: isNavOpen ? '250px' : '10px',
                height: '20vh',
                position: 'absolute',
                left: isNavOpen ? '0px' : '-230px',
                backgroundColor: '#000',
                transition: 'left 0.3s ease-in-out',
                paddingTop: '60px',
                boxSizing: 'border-box'
            }}>
            <Menu 
                mode="inline"
                items={items}
                style={{
                    height: '100%',
                    color: '#fff', // This keeps text color as white
                    borderRight: 'none',
                    backgroundColor: '#000', // Background color for the entire Menu
                    fontSize: '17px',
                }}
                selectedKeys={[menuKey]}
                onClick={onMenuClick}
                itemStyle={{
                    backgroundColor: '#000' // Normal item background color
                }}
                selectedItemStyle={{
                    backgroundColor: '#808080', // Gray background for selected item
                    color: '#000' // White text color for selected item for visibility
                }}
            />
            </div>
            <Outlet />
        </div>
    );
}

import React from 'react'
import { Button } from 'antd';
import { LogoutOutlined } from '@ant-design/icons';
// import { Link, useNavigate } from 'react-router-dom'
import { useNavigate } from 'react-router-dom'
import { useDispatch, useSelector } from 'react-redux';
import { logout as reduxLogout } from '../store/userSlice'

export default function Head() {
    const { username } = useSelector(state => state.user)
    const navigate = useNavigate()
    const dispatch = useDispatch()

    function logout() {
        dispatch(reduxLogout())
        navigate('/')
    }

    return (
        <div className='bg-[#2a2f33] border-b border-slate-400 shadow-md sticky top-0 z-50'>
            <header className='flex justify-between items-center max-w-6xl mx-auto h-20 max-h-20'>
                <div className='flex-shrink-0'>
                    <h1 className='text-white font-bold text-2xl'>Hi, {username}</h1>  {/* 加粗并保持一致的字体风格 */}
                </div>
                <div className='h-full flex items-center space-x-6 flex-nowrap'>
                    {/* <Link to="/change-password" className='text-white text-lg hover:text-blue-300 font-bold'>Change password</Link> */}
                    <Button type="primary" size={'middle'} icon={<LogoutOutlined />} onClick={logout} style={{ backgroundColor: '#81d4fa', borderColor: '#57c3f3', fontWeight: 'bold' }}>
                        Logout  {/* 加粗按钮文本 */}
                    </Button>
                </div>
            </header>
        </div>
    )
}

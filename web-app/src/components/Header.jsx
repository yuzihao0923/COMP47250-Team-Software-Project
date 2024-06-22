import React from 'react'
import { Button } from 'antd';
import { LogoutOutlined } from '@ant-design/icons';
import { Link, useNavigate } from 'react-router-dom'
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
        <div className='bg-white border-b border-slate-200 shadow-md sticky top-0 z-50'>
            <header className='flex justify-between items-center max-w-6xl mx-auto h-20 max-h-20'>
                <div className='flex-shrink-0'>
                    <h1 className='text-black font-light text-2xl'>Hi, {username}</h1>
                </div>
                <div className='h-full flex items-center space-x-6 flex-nowrap '>
                    <Link className='contents lg:text-sm md:text-base sm:text-lg '>Change password</Link>
                    <Button type="primary" size={'middle'} icon={<LogoutOutlined />} onClick={logout}>
                        Logout
                    </Button>
                </div>
            </header>
        </div>

    )
}
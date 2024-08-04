import React from 'react';
import { BarChartOutlined } from '@ant-design/icons';

export default function Card(props) {
    const {
        logoBackground = 'bg-gray-800',
        logo = <BarChartOutlined style={{ color: 'white', fontSize: '16px' }} />,
        data = 0,
        dataTitle = 'Default Title',
        titleColor = 'text-white'
    } = props;

    return (
        <div className='flex h-36 hover:shadow-lg transition duration-300 ease-in-out bg-black text-white'>
            <div className={`flex items-center justify-center ${logoBackground} w-1/4 rounded-l-lg card-icon`}>
                <span style={{ fontSize: '28px' }}> {}
                    {logo}
                </span>
            </div>
            <div className='flex flex-1 flex-col justify-center items-center text-center rounded-r-lg' style={{ background: '#1a1a1a' }}>
                <span className='text-3xl font-bold mb-3'>{data}</span>
                <h3 className={`text-base font-bold ${titleColor}`}>{dataTitle}</h3>
            </div>
        </div>
    )
}

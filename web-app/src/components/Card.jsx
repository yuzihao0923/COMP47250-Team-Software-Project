import React from 'react'
import { BarChartOutlined } from '@ant-design/icons';

export default function Card(props) {

    const {
        logoBackground = 'bg-pink-200',
        logo = <BarChartOutlined />,
        data = 0,
        dataTitle = 'Default Title'
    } = props;

    return (
        <div className='flex h-36 '>
            <div className={`flex items-center justify-center ${logoBackground} w-1/4 rounded-l-lg`}>
                <span>
                    {logo}
                </span>
            </div>
            <div className='bg-fuchsia-200 flex flex-1 flex-col justify-center items-center text-center rounded-r-lg'>
                <span className='text-3xl font-bold text-neutral-500 mb-3'>{data}</span>
                <h3 className='text-base font-light text-neutral-500'>{dataTitle}</h3>
            </div>
        </div>
    )
}

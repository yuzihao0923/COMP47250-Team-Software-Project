import React from 'react'

export default function Card(props) {

    const { logoBackground, logo, data, dataTitle } = props

    return (
        <div className='flex h-56 '>
            <div className={`flex items-center justify-center ${logoBackground} w-1/4`}>
                <span>
                    {logo}
                </span>
            </div>
            <div className='bg-slate-100 flex flex-1 flex-col justify-center items-center text-center'>
                <span className='text-3xl font-bold'>{data}</span>
                <h3 className='text-base font-light'>{dataTitle}</h3>
            </div>
        </div>
    )
}

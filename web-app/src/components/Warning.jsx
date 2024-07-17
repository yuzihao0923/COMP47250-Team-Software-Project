import React from 'react'
import '../css/Warning.css'

export default function Warning(props) {

    const { warningNumber } = props

    return (

        <div className='alert'>
            {/* {alert && <span className="text-white text-2xl font-bold">Warning</span>} */}
            <span className='alert-message'>The broker has received {warningNumber} messages in the last 3 seconds</span>
        </div>
    )
}

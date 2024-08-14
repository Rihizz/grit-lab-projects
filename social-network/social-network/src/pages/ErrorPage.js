import React from 'react'
import { Link } from 'react-router-dom'
import '../styles/ErrorPage.css'

const ErrorPage = ({errorType}) => {
    console.log('ErrorPage', errorType)
    return (
        <div className='error'>
        <h1>ErrorPage</h1>
        <h2 className='errorType'>{errorType}</h2>
        <h2 className='errorMessage'>Oops! Something went wrong. Leave this URL immediately!</h2>
        {/* refresh home page */}
        <button className='errorButton'><a href='/'>Home</a></button>
        </div>

    )
}

export default ErrorPage
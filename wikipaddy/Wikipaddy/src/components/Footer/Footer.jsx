import React from 'react'
import './Footer.css'

const Footer = ({theme}) => {
  return (
    <div className={`footer ${theme}`}>
      <div className="waves">
        <div className={`wave ${theme}`} id='wave1'></div>
        <div className={`wave ${theme}`} id='wave2'></div>
        <div className={`wave ${theme}`} id='wave3'></div>
        <div className={`wave ${theme}`} id='wave4'></div>
      </div>
      <div className={`footer-content ${theme}`}>
        <p>Copyright © 2024 Wikipaddy. All rights reserved.</p>
      </div>    
    </div>
  )
}

export default Footer
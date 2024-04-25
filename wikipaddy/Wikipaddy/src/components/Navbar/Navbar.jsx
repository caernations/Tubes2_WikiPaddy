import React, { useState, useEffect } from 'react';
import './Navbar.css';
import light from '../../assets/sun.png';
import dark from '../../assets/moon.png';

const Navbar = ({ theme, setTheme }) => {
  const [isScrolled, setIsScrolled] = useState(false);

  useEffect(() => {
    const handleScroll = () => {
      const scrollPosition = window.scrollY;
      const windowHeight = window.innerHeight;
      const pageHeight = document.body.clientHeight;
      const twentyPercentHeight = pageHeight * 0.2;

      setIsScrolled(scrollPosition > twentyPercentHeight);
    };

    window.addEventListener('scroll', handleScroll);
    return () => {
      window.removeEventListener('scroll', handleScroll);
    };
  }, []);

  const toggle_mode = () => {
    theme === 'light' ? setTheme('dark') : setTheme('light');
  };

  return (
    <div className={`navbar ${theme} ${isScrolled ? 'scrolled' : ''}`}>
      <h3 className={`logo ${theme}`}>Wikipaddy</h3>

      <ul>
        <li className={`navlink ${theme}`}>Home</li>
        <li className={`navlink ${theme}`}>About Us</li>
        <li className={`navlink ${theme}`}>Information</li>
        <li className={`navlink ${theme}`}>How to Use</li>
      </ul>

      <img onClick={toggle_mode} src={theme === 'light' ? light : dark} alt="" className='toggle'/>
    </div>
  );
};

export default Navbar;
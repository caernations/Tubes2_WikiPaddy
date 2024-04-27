import React, { useState, useEffect } from 'react';
import { Link } from 'react-scroll';
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
      <li className={`navlink ${theme}`}>
          <Link to='project' spy={true} smooth={true} duration={500}>Home</Link>
        </li>
        <li className={`navlink ${theme}`}>
          <Link to='about' spy={true} smooth={true} duration={500}>About Us</Link>
        </li>
        <li className={`navlink ${theme}`}>
          <Link to='what' spy={true} smooth={true} duration={500} offset={80}>Information</Link>
        </li>
      </ul>

      <img onClick={toggle_mode} src={theme === 'light' ? light : dark} alt="" className='toggle'/>
    </div>
  );
};

export default Navbar;
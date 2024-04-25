import React, { useState, useEffect } from 'react';
import './About.css';
import photo from '../../assets/photo1.jpg';

const About = ({ theme }) => {
  const [titleFontSize, setTitleFontSize] = useState('100px');

  useEffect(() => {
    const handleScroll = () => {
      const scrollPercentage = (window.scrollY / (document.documentElement.scrollHeight - window.innerHeight)) * 200;
      const newSize = 200 - scrollPercentage * 1.6
      ;
      setTitleFontSize(`${newSize}px`);
    };

    window.addEventListener('scroll', handleScroll);

    return () => {
      window.removeEventListener('scroll', handleScroll);
    };
  }, []);

  const flipCard = (index) => {
    const cards = document.querySelectorAll('.card');
    cards.forEach((card, i) => {
      if (i === index) {
        card.classList.toggle('is-flipped');
      } else {
        card.classList.remove('is-flipped');
      }
    });
  };

  return (
    <div className={`about ${theme}`}>
      <div className={`about-title ${theme}`}>
        <h1 style={{ fontSize: titleFontSize }}>About Us</h1>
      </div>
      <div className='about-container'>
        <div className='card' onClick={() => flipCard(0)}>
          <div className={`front ${theme}`}>
            <div className='card-img'>
              <img src={photo}/>
            </div>
            <div className={`card-text ${theme}`}>
              <h3>13522001</h3>
              <p>Stima Mantap</p>
            </div>
          </div>
          <div className={`back ${theme}`}>
            <h3>Back</h3>
            <p>Back of the card content</p>
          </div>
        </div>
        <div className='card' onClick={() => flipCard(1)}>
          <div className={`front ${theme}`}>
            <div className='card-img'>
              <img src={photo}/>
            </div>
            <div className={`card-text ${theme}`}>
              <h3>13522002</h3>
              <p>Stima Excellent</p>
            </div>
          </div>
          <div className={`back ${theme}`}>
            <h3>Back</h3>
            <p>Back of the card content</p>
          </div>
        </div>
        <div className='card' onClick={() => flipCard(2)}>
          <div className={`front ${theme}`}>
            <div className='card-img'>
              <img src={photo}/>
            </div>
            <div className={`card-text ${theme}`}>
              <h3>13522003</h3>
              <p>Stima Awesome</p>
            </div>
          </div>
          <div className={`back ${theme}`}>
            <h3>Back</h3>
            <p>Back of the card content</p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default About;
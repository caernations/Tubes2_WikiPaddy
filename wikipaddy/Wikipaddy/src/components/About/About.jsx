import React, { useState, useEffect } from 'react';
import './About.css';
import photo1 from '../../assets/photo1.jpg';
import photo2 from '../../assets/photo2.jpg';
import photo3 from '../../assets/photo3.jpg';

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
              <img src={photo1}/>
            </div>
            <div className={`card-text ${theme}`}>
              <h3>13522140</h3>
              <p>Yasmin Farisah Salma</p>
            </div>
          </div>
          <div className={`back ${theme}`}>
            <p>stream ttpd !!</p>
          </div>
        </div>
        <div className='card' onClick={() => flipCard(1)}>
          <div className={`front ${theme}`}>
            <div className='card-img'>
              <img src={photo2}/>
            </div>
            <div className={`card-text ${theme}`}>
              <h3>13522142</h3>
              <p>Farhan Raditya Aji</p>
            </div>
          </div>
          <div className={`back ${theme}`}>
            <p>
              si bobi sifatnya angkuh <br></br>
              menyala abangkuhhh 
            </p>
          </div>
        </div>
        <div className='card' onClick={() => flipCard(2)}>
          <div className={`front ${theme}`}>
            <div className='card-img'>
              <img src={photo3}/>
            </div>
            <div className={`card-text ${theme}`}>
              <h3>13522156</h3>
              <p>Jason Fernando</p>
            </div>
          </div>
          <div className={`back ${theme}`}>
            <p>
              ikan hiu makan kedondong<br></br>
              stima i love u dong
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default About;
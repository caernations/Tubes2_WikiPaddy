import React, { useEffect, useRef } from 'react';
import './What.css';
import { motion, useScroll, useTransform } from "framer-motion";
import wikipedia from './../../assets/wikipedia-logo.png';

const TextWrapper = ({children}) => {
  const textRef = useRef(null);
  
  const { scrollYProgress } = useScroll({
    target: textRef,
    offset: ["start end", "end start"],
  });

  const scrollYProgressWithOffset = useTransform(scrollYProgress, value => value );

  const opacity = useTransform(scrollYProgressWithOffset, [1, 0.8, 0], [1, 0.8, 0]);
  const x = useTransform(scrollYProgressWithOffset, [1, 0.4, 0], [0, 0, 1000]);
  // const colorChange = useTransform(
  //   scrollYProgressWithOffset,
  //   [0, 0.2, 0.4, 0.6, 0.8, 1],
  //   [
  //     "hsla(180, 7%, 75%, 0.9)",
  //     "hsla(180, 7%, 75%, 0.9)",
  //     "#ffffff",
  //     "#ffffff",
  //     "#ffffff",
  //     "#ffffff",
  //   ]
  // );

  return (
    <div ref={textRef}>
      <motion.p style={{ opacity, x}}>{children}</motion.p>
    </div>
  );
};

const LogoWrapper = ({photo}) => {
  const textRef = useRef(null);
  
  const { scrollYProgress } = useScroll({
    target: textRef,
    offset: ["start end", "end start"],
  });

  const scrollYProgressWithOffset = useTransform(scrollYProgress, value => value );

  const opacity = useTransform(scrollYProgressWithOffset, [1, 0.8, 0], [1, 1, 1]);
  const x = useTransform(scrollYProgressWithOffset, [1, 0.5, 0], [0, 0, 500]);

  return (
    <div ref={textRef}>
      <motion.img src={wikipedia} style={{ opacity, x}}>{photo}</motion.img>
    </div>
  );
};

const What = ({ theme }) => {
  const textRef = useRef(null);
  const howToUseRef = useRef(null);

  useEffect(() => {
    const observer = new IntersectionObserver(
      ([entry]) => {
        if (entry.isIntersecting) {
          textRef.current.style.opacity = 0;
        } else {
          textRef.current.style.opacity = 1;
        }
      },
      { threshold: 0.2 } 
    );

    if (howToUseRef.current) {
      observer.observe(howToUseRef.current);
    }

    return () => {
      if (howToUseRef.current) {
        observer.unobserve(howToUseRef.current);
      }
    };
  }, []);

  return (
    <div className={`what ${theme}`}>
      <div ref={textRef} className={`text ${theme}`} style={{ transition: 'opacity 0.5s' }}>
        <TextWrapper><h1>What</h1></TextWrapper>
        <TextWrapper><h1>Is</h1></TextWrapper>
        <TextWrapper><h1>Wikipaddy</h1></TextWrapper>
        <TextWrapper><h1>?</h1></TextWrapper>
        <div className={`text-p ${theme}`}>
          <TextWrapper>
            <p>
              Wikipaddy adalah sebuah projek dimana pengguna memberikan dua halaman Wikipedia yang berbeda dan wikipaddy akan mencari jalan terpendek untuk sampai dari satu halaman ke halaman lainnya hanya dengan mengikuti tautan-tautan dalam artikel tersebut. Tujuan dari proyek ini adalah untuk menguji kecepatan algoritma IDS dan BFS dalam menavigasi dan mencari informasi di Wikipedia.
            </p>
          </TextWrapper>
        </div>
      </div>
      <div ref={howToUseRef} className={`how-to-use ${theme}`}>
        <div className='logo-wiki'>
          <LogoWrapper style={{zIndex: '20'}}>
          </LogoWrapper>
        </div>
        <div className={`how-to-use-text ${theme}`}>
              <p>
                1. Masukan 2 judul wikipedia yang berbeda <br></br>
                2. Pilih Algoritma antara IDS dan BFS <br></br>
                3. Mulai pencarian dengan submit
              </p>
        </div>
      </div>
    </div>
  );
};

export default What;
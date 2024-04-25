import React, {useEffect, useState} from 'react'
import Navbar from './components/Navbar/Navbar'
import Project from './components/Projects/Project'
import About from './components/About/About'
import What from './components/What/What'
import Footer from './components/Footer/Footer'
import './index.css'

const App = () => {

  const current_theme = localStorage.getItem('current_theme');
  const [theme, setTheme] = useState(current_theme ? 
    current_theme : 'light');
  

  useEffect(()=>{
    localStorage.setItem('current_theme', theme);
  })

  return (
    <div className={`container ${theme}`}>
      <Navbar theme={theme} setTheme={setTheme}/>
      <Project theme={theme} setTheme={setTheme}/>
      <About theme={theme} setTheme={setTheme}/>
      <What theme={theme} setTheme={setTheme}/>
      <Footer theme={theme} setTheme={setTheme}/>
    </div>
  )
}

export default App
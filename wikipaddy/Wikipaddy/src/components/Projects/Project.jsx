import React from 'react';
import './Project.css';

function handleButtonClick(e) {
  const buttons = document.querySelectorAll('.input-algorithm-box button');

  buttons.forEach((button) => {
    button.classList.remove('active');
  });

  e.currentTarget.classList.add('active');
}

const Project = ({ theme, awal, akhir }) => {
  return (
    <div className={`project ${theme}`}>
      <div className={`title ${theme}`}>
        <p>Wikipaddy</p>
      </div>
      <div className='project-container'>
        <form>
          <div className={`input-box ${theme}`}>
            <input className='awal' type='text' value={awal} placeholder='Input Start Text' />
          </div>
          <div className={`project-text ${theme}`}>
            <h2>To</h2>
          </div>
          <div className={`input-box ${theme}`}>
            <input className='akhir' type='text' value={akhir} placeholder='Input End Text'/>
          </div>
          <div className={`project-text ${theme}`}>
            <h2>Select Algorithm</h2>
          </div>
          <div className={`input-algorithm-box ${theme}`}>
            <button type='button' id='BFS' onClick={handleButtonClick}>BFS</button>
            <button type='button' id='IDS' onClick={handleButtonClick}>IDS</button>
          </div>
          <div className={`submit-box ${theme}`}>
            <button type='submit'>SUBMIT</button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default Project;
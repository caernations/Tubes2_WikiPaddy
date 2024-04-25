import React from 'react';
import './Project.css';

const Project = ({ theme, awal, akhir }) => {
  return (
    <div className={`project ${theme}`}>
      <div className={`title ${theme}`}>
        <p>Wikipaddy</p>
      </div>
      <div className='project-container'>
        <form>
          <div className={`input-box ${theme}`}>
            <input className='awal' type='text' value={awal} />
          </div>
          <div className={`project-text ${theme}`}>
            <h2>To</h2>
          </div>
          <div className={`input-box ${theme}`}>
            <input className='akhir' type='text' value={akhir} />
          </div>
          <div className={`project-text ${theme}`}>
            <h2>Select Algorithm</h2>
          </div>
          <div className={`input-algorithm-box ${theme}`}>
            <button>BFS</button>
            <button>IDS</button>
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
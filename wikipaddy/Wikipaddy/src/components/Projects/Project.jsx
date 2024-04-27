import React, {useState} from 'react';
import './Project.css';

const Project = ({ theme, awal, akhir, bfs, ids, bibfs }) => {
  const [suggestionsAwal, setSuggestionsAwal] = useState([]);
  const [suggestionsAkhir, setSuggestionsAkhir] = useState([]);
  const [inputAwal, setInputAwal] = useState(awal);
  const [inputAkhir, setInputAkhir] = useState(akhir);
  const [inputAwalStyle, setInputAwalStyle] = React.useState({});
  const [inputAkhirStyle, setInputAkhirStyle] = React.useState({});
  const [activeButton, setActiveButton] = useState(''); 
  const [showResult, setShowResult] = useState(false);
  const [isExiting, setIsExiting] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [results, setResult] = useState('');

  function handleButtonClick(e) {
    const buttons = document.querySelectorAll('.input-algorithm-box button');
  
    buttons.forEach((button) => {
      button.classList.remove('active');
    });
  
    e.currentTarget.classList.add('active');
    setActiveButton(e.currentTarget.id);
  }

  const handleInputChangeAwal = async (e) => {
    const query = e.target.value;
    setInputAwal(query);
  
    if (query) {
      const response = await fetch(`https://en.wikipedia.org/w/api.php?origin=*&action=opensearch&search=${query}&limit=7&namespace=0&format=json`);
      const data = await response.json();
  
      setSuggestionsAwal(data[1]);
      setInputAwalStyle({ borderRadius: '15px 15px 0 0', background: 'white'});
    } else {
      setSuggestionsAwal([]);
      setInputAwalStyle({});
    }
  };
  
  const handleInputChangeAkhir = async (e) => {
    const query = e.target.value;
    setInputAkhir(query);
  
    if (query) {
      const response = await fetch(`https://en.wikipedia.org/w/api.php?origin=*&action=opensearch&search=${query}&limit=7&namespace=0&format=json`);
      const data = await response.json();
  
      setSuggestionsAkhir(data[1]);
      setInputAkhirStyle({ borderRadius: '15px 15px 0 0', background: 'white' });
    } else {
      setSuggestionsAkhir([]);
      setInputAkhirStyle({});
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsLoading(true);
    const formattedInputAwal = inputAwal.replace(/ /g, '_'); 
    const formattedInputAkhir = inputAkhir.replace(/ /g, '_');
    const response = await fetch(`http://localhost:8080/${activeButton}?start_title=${formattedInputAwal}&end_title=${formattedInputAkhir}`);
    const data = await response.json();
    setIsLoading(false);
    setShowResult(true);
    setResult(data);
    console.log(data);
  };

  const handleExit = () => {
    setIsExiting(true);
    setTimeout(() => setShowResult(false), 500);
    window.location.reload();
  };

  return (
    <div className={`project ${theme}`}>
      <div className={`title ${theme}`}>
        <p>Wikipaddy</p>
      </div>
      <div className='project-container'>
        <form>
          <div className={`input-box ${theme}`} style={inputAwalStyle}>
            <input className='awal' type='text' value={inputAwal} placeholder='Input Start Text' onChange={handleInputChangeAwal} />
            <div className='suggestions'>
              {suggestionsAwal.map((suggestion, index) => (
                <p key={index} onClick={() => {setInputAwal(suggestion); setSuggestionsAwal([]); setInputAwalStyle({});}}>{suggestion}</p>
              ))}
            </div>
          </div>
          <div className={`project-text ${theme}`}>
            <h2>To</h2>
          </div>
          <div className={`input-box ${theme}`} style={inputAkhirStyle}>
            <input className='akhir' type='text' value={inputAkhir} placeholder='Input End Text' onChange={handleInputChangeAkhir} />
            <div className='suggestions'>
              {suggestionsAkhir.map((suggestion, index) => (
                <p key={index} onClick={() => {setInputAkhir(suggestion); setSuggestionsAkhir([]); setInputAkhirStyle({});}}>{suggestion}</p>
              ))}
            </div>
          </div>
          <div className={`project-text ${theme}`}>
            <h2>Select Algorithm</h2>
          </div>
          <div className={`input-algorithm-box ${theme}`}>
            <button type='button' id='BFS' value={bfs} onClick={handleButtonClick}>BFS</button>
            <button type='button' id='IDS' value={ids} onClick={handleButtonClick}>IDS</button>
            <button type='button' id='BI-BFS' value={bibfs} onClick={handleButtonClick}>BI-BFS</button>
          </div>
          <div className={`submit-box ${theme}`}>
          <button type='submit' onClick={handleSubmit}>
            {isLoading ? 'Loading...' : 'SUBMIT'}
          </button>
          </div>
        </form>
        {showResult && (
          <div className={`project-result ${isExiting ? 'slide-down' : ''}`}>
            <div className={`exit-button ${theme}`}>
              <button type='button' onClick={handleExit}>X</button>
            </div>
            <div className={`result-box ${theme}`}>
              <div className='result-text'>
                <div className={`result-algorithm ${theme}`}>{activeButton}</div>
                <div className={`result-start-text ${theme}`}>{inputAwal}</div>
                <div className={`result-path ${theme}`}>
                  <ul className={`result-list ${theme}`}>
                    {results.path.map((path, index) => (
                      <li key={index}><a href={path}>{path}</a></li>
                    ))}
                  </ul>
                </div>
                <div className={`result-end-text ${theme}`}>{inputAkhir}</div>
                <div className={`result-time ${theme}`}>
                  <p>TIME : {results.waktu_eksekusi} Seconds</p>
                  <p>DEPTH : {results.kedalaman}</p>
                  <p>TOTAL LINKS : {results.total}</p>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
    
  );
};

export default Project;
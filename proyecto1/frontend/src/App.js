import './App.css';
import Table from './components/Table'
import SecondTable from './components/SecondTable';
import { Progress, Segment } from 'semantic-ui-react'
import { update } from 'react-spring';
import Tree from './components/Tree';
import { useState, useEffect } from 'react';




function App() {
  let maquina = localStorage.getItem('maquinaseleccionada');
  let kill = localStorage.getItem('kill');
  let pid = localStorage.getItem('pid');
  if (pid == null) {
    pid = 0;
  }
  const [ram, setRam] = useState(0);
  const [cpu, setCpu] = useState(0);
  useEffect(() => {
    const interval = setInterval(() => {
      fetch(`http://35.245.67.156:4000/uso/${maquina}/${kill}/${pid}`)
        .then(response => response.json())
        .then(data => {
          const ramValue = parseFloat(data[data.length - 1].ram);
          const cpuValue = parseFloat(data[data.length - 1].cpu);
          setRam(ramValue);
          setCpu(cpuValue);
          localStorage.setItem('kill', false);

        })
        .catch(error => console.log(error));
    }, 5000);
    return () => clearInterval(interval);
  }, []);
  const opcionesMaquinas = [
    { nombre: 'Maquina 1', numero: '1' },
    { nombre: 'Maquina 2', numero: '2' },
    { nombre: 'Maquina 3', numero: '3' },
    { nombre: 'Maquina 4', numero: '4' },
  ];

  const [maquinaSeleccionada, setMaquinaSeleccionada] = useState(1);

  const handleSeleccionMaquina = (event) => {
    localStorage.setItem('maquinaseleccionada', event.target.value);
    setMaquinaSeleccionada(event.target.value);
  };
  
  
  const handleKillClick = (e) => {

    // Almacenar el valor de PID en el Local Storage
    localStorage.setItem('kill', true);

    let maquina = localStorage.getItem('maquinaseleccionada');
    let kill = true;
    let pid = localStorage.getItem('pid');
    if (pid == null) {
      pid = 0;
    }


    fetch(`http://35.245.67.156:4000/kill/${maquina}/${kill}/${pid}`)
        .then(response => response.json())
        .catch(error => console.log(error));
    // Limpiar el valor del input
  }

  const handleInputChange = (e) => {
    localStorage.setItem('pid', e.target.value );
  }


  return (

    
    <div className="container">
      <h1 className='titulo-h1'>Proyecto 1 Sopes 1</h1>
      
      <div className='container-graphics'>
        <Segment inverted>
        <h1>RAM</h1>
          <Progress percent={ram} inverted color='green' progress />
          <h1>CPU</h1>
          <Progress percent={cpu} inverted color='red' progress />
        </Segment>
      </div>
      <br></br>
      <br></br>
      <br></br>
      <br></br>
      <br></br>
      <br></br>
      
      <div className='container-table-1'>
        <div>
          <h1>Selecciona una Máquina</h1>
          <select value={maquinaSeleccionada} onChange={handleSeleccionMaquina}>
            {opcionesMaquinas.map((maquina) => (
              <option key={maquina.numero} value={maquina.numero}>
                {maquina.nombre}
              </option>
            ))}
          </select>
          <p>Has seleccionado: Maquina {maquinaSeleccionada}</p>
        </div>
        <br></br>
        <div>
          <label htmlFor="pidInput">PID:</label>
          <input
            type="text"
            id="pidInput"
            onChange={handleInputChange}
          />
          <button onClick={handleKillClick}>Kill</button>
        </div>
      </div>
      <div className='container-table-2'>
        <SecondTable/>
      </div>
      <br></br>
      <br></br>
      <br></br>
      <br></br>
      <br></br>
      
      
    </div>
  );
}

export default App;

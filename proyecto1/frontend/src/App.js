import './App.css';
import Table from './components/Table'
import SecondTable from './components/SecondTable';
import { Progress, Segment } from 'semantic-ui-react'
import { update } from 'react-spring';
import Tree from './components/Tree';
import { useState, useEffect } from 'react';




function App() {
  let maquina = '1';

  const [ram, setRam] = useState(0);
  const [cpu, setCpu] = useState(0);
  useEffect(() => {
    const interval = setInterval(() => {
      fetch('http://35.245.67.156:4000/uso')
        .then(response => response.json())
        .then(data => {
          const ramValue = parseFloat(data[data.length - 1].ram);
          const cpuValue = parseFloat(data[data.length - 1].cpu);
          setRam(ramValue);
          setCpu(cpuValue);
        })
        .catch(error => console.log(error));
    }, 50000);
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
    setMaquinaSeleccionada(event.target.value);

  };
  
  return (

    
    <div className="container">
      <h1 className='titulo-h1'>Proyecto 1 Sopes 1</h1>
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

      
      <div className='container-table-1'>
        <Table/>
      </div>
      <div className='container-table-2'>
        <SecondTable/>
      </div>
      <br></br>
      <br></br>
      <br></br>
      <br></br>
      <br></br>
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
      <Tree/>
    </div>
  );
}

export { maquinaSeleccionada }; // Exporta maquinaSeleccionada de forma individual
export default App;

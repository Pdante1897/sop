const express = require('express');
const mysql = require('mysql2');
const cors = require('cors'); 
const app = express();
app.use(cors()); 

const connection = mysql.createConnection({
  host: '10.150.0.2',
  user: 'root',       
  password: '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', 
  database: 'proyecto_1' 
});

connection.connect((error) => {
  if (error) {
    console.error('Error al conectar a la base de datos:', error);
  } else {
    console.log('Conexión exitosa a la base de datos MySQL');
  }
});


let maquinaGlobal;
let killGlobal;
let pidGlobal;


app.get('/proceso/:maquina/', (req, res) => {
  const maquina = req.params.maquina; 
  const maquina2 = req.params.maquina; 

  // Extrae el valor de la variable "maquina" de la URL
  connection.query("select * from proceso where maquina = ? and proces = (select Max(proces) from proceso where maquina = ?) order by id ", [maquina, maquina2], (error, results) => {
    if (error) {
      console.error('Error al realizar la consulta:', error);
      res.status(500).send('Error al realizar la consulta');
    } else {
      res.send(results);
    }
  });
});

app.get('/uso/:maquina', (req, res) => {
  const maquina = req.params.maquina; // Extrae el valor de la variable "maquina" de la URL

  connection.query("SELECT * FROM uso WHERE maquina = ? ORDER BY id desc limit  1", [maquina], (error, results) => {
    if (error) {
      console.error('Error al realizar la consulta:', error);
      res.status(500).send('Error al realizar la consulta');
    } else {
      res.send(results);
    }
  });
});

app.get('/tarea/:maquina', (req, res) => {
  const maquina = req.params.maquina; // Extrae el valor de la variable "maquina" de la URL
  connection.query("SELECT * FROM tarea WHERE maquina = ? ORDER BY id desc limit  1", [maquina], (error, results) => {
    if (error) {
      console.error('Error al realizar la consulta:', error);
      res.status(500).send('Error al realizar la consulta');
    } else {
      res.send(results);
      
    }
  });
  
});


app.get('/hijo/:maquina', (req, res) => {
  const maquina = req.params.maquina; // Extrae el valor de la variable "maquina" de la URL
  connection.query('SELECT * FROM hijo', (error, results) => {
    if (error) {
      console.error('Error al realizar la consulta:', error);
      res.status(500).send('Error al realizar la consulta');
    } else {
      res.send(results);
    }
  });
});

app.get('/kill/:maquina/:kill/:pid', (req, res) => {
  const maquina = req.params.maquina; // Extrae el valor de la variable "maquina" de la URL
  const kill = req.params.kill;
  const pid = req.params.pid;
  maquinaGlobal = maquina;
  killGlobal = kill;
  pidGlobal = pid;
  console.log(kill);
  console.log(pid);
  
  res.status(200).json({ message: `El kill guardado` });
});



// Configuración de la conexión a la base de datos

// Middleware para analizar JSON en solicitudes POST
app.use(express.json());

// Endpoint para insertar un proceso
app.post('/insertar_proceso/:maquina/:proces', (req, res) => {
  const maquina = req.params.maquina;
  const proces = req.params.proces;
  console.log(req.body);

  // Suponiendo que el JSON tiene una propiedad "Procesos" que es un arreglo de objetos
  const datosAInsertar = req.body.Procesos;

  const values = datosAInsertar.map(({ Pid, Nombre, Usuario, Estado, Ram }) =>
    [Estado, Pid, Nombre, Usuario, Ram, maquina, proces]
  );

  const sql_command = `INSERT INTO proceso(estado, pid, name, user, ram, maquina, proces) VALUES ?`;
  connection.query(sql_command, [values], (error, results) => {
    if (error) {
      console.error('Error al insertar proceso:', error);
      res.status(500).send('Error al insertar proceso');
    } else {
      const lastId = results.insertId;
      res.status(200).json({ message: `El id del último proceso ingresado es: ${lastId}` });
    }
  });
});

// Endpoint para insertar un hijo
app.post('/insertar_hijo', (req, res) => {
  const { pid_padre, pid_hijo, name } = req.body;
  const sql_command = `INSERT INTO hijo(pid_padre, pid_hijo, name) VALUES(?, ?, ?)`;
  connection.query(sql_command, [pid_padre, pid_hijo, name], (error, results) => {
    if (error) {
      console.error('Error al insertar hijo:', error);
      res.status(500).send('Error al insertar hijo');
    } else {
      const lastId = results.insertId;
      res.status(200).json({ message: `El id del último hijo ingresado es: ${lastId}` });
    }
  });
});

// Endpoint para insertar un uso
app.post('/insertar_uso/:maquina', (req, res) => {
  const { ram, cpu } = req.body;
  const maquina = req.params.maquina;
  let data = {};
  if (maquinaGlobal == maquina) {
    data = {
      pid: pidGlobal,
      kill: killGlobal
    }
  }
  

  const sql_command = `INSERT INTO uso(ram, cpu, maquina) VALUES(?, ?, ?)`;
  connection.query(sql_command, [ram, cpu, maquina], (error, results) => {
    if (error) {
      console.error('Error al insertar uso:', error);
      res.status(500).send('Error al insertar uso');
    } else {
      const lastId = results.insertId;
      pidGlobal = 0;
      killGlobal = false;
      res.status(200).json({ message: `El id del último uso ingresado es: ${lastId}`, data: data });
    }
  });
});

// Endpoint para insertar una tarea
app.post('/insertar_tarea/:maquina', (req, res) => {
  const maquina = req.params.maquina;
  const { running, sleeping, zombie, stopped, total } = req.body;
  const sql_command = `INSERT INTO tarea(running, sleeping, zombie, stopped, total, maquina) VALUES(?, ?, ?, ?, ?, ?)`;
  connection.query(sql_command, [running, sleeping, zombie, stopped, total, maquina], (error, results) => {
    if (error) {
      console.error('Error al insertar tarea:', error);
      res.status(500).send('Error al insertar tarea');
    } else {
      const lastId = results.insertId;
      res.status(200).json({ message: `El id del último Task ingresado es: ${lastId}` });
    }
  });
});



app.listen(4000, () => {
  console.log('API escuchando en el puerto 4000');
});
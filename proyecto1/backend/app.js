const express = require('express');
const mysql = require('mysql');
const cors = require('cors'); 
const app = express();
app.use(cors()); 

const connection = mysql.createConnection({
  host: 'localhost',
  user: 'root',       
  password: 'admin', 
  database: 'password' 
});

connection.connect((error) => {
  if (error) {
    console.error('Error al conectar a la base de datos:', error);
  } else {
    console.log('Conexión exitosa a la base de datos MySQL');
  }
});

app.get('/proceso', (req, res) => {
  connection.query('SELECT * FROM proceso', (error, results) => {
    if (error) {
      console.error('Error al realizar la consulta:', error);
      res.status(500).send('Error al realizar la consulta');
    } else {
      res.send(results);
    }
  });
});

app.get('/uso', (req, res) => {
  connection.query('SELECT * FROM uso', (error, results) => {
    if (error) {
      console.error('Error al realizar la consulta:', error);
      res.status(500).send('Error al realizar la consulta');
    } else {
      res.send(results);
      connection.query('TRUNCATE usos_modificados;', (error, results) => {
        if (error) {
          console.error('Error al realizar la consulta:', error);
        }
      });
    }
  });
});

app.get('/tarea', (req, res) => {
  connection.query('SELECT * FROM tarea', (error, results) => {
    if (error) {
      console.error('Error al realizar la consulta:', error);
      res.status(500).send('Error al realizar la consulta');
    } else {
      res.send(results);
    }
  });
});

app.get('/hijo', (req, res) => {
  connection.query('SELECT * FROM hijo', (error, results) => {
    if (error) {
      console.error('Error al realizar la consulta:', error);
      res.status(500).send('Error al realizar la consulta');
    } else {
      res.send(results);
    }
  });
});

app.listen(4000, () => {
  console.log('API escuchando en el puerto 4000');
});
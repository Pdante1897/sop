import React from 'react';

// Ahora puedes usar maquinaSeleccionada en tu componente.
class Row extends React.Component {
    render() {
        const { proceso } = this.props;
        return (
            <tr>
                <td><p>{proceso.running}</p></td>
                <td><p>{proceso.sleeping}</p></td>
                <td><p>{proceso.zombie}</p></td>
                <td><p>{proceso.stopped}</p></td>
                <td><p>{proceso.total}</p></td>
            </tr>
        );
    }
}

class Table extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            procesos: []
        };
        this.update = this.update.bind(this);
    }

    componentDidMount() {
        this.interval = setInterval(() => this.update(), 1000);
    }

    componentWillUnmount() {
        clearInterval(this.interval);
    }

    update() {  
        let maquina = localStorage.getItem('maquinaseleccionada');

        console.log(maquina);
        fetch(`http://35.245.67.156:4000/tarea/${maquina}`, {
            method: 'GET',
            mode: 'cors',
        })
        .then(response => response.json())
        .then(data => {
            if (Array.isArray(data)) {
                this.setState({ procesos: data });
            } else {
                console.error("Error: Data is not an array.");
            }
        })
        .catch(error => {
            console.error(error);
        });
    }

    render() {
        const { procesos } = this.state;
        return (
            <div>
                <table>
                    <thead>
                        <tr>
                            <th scope="col">Running</th>
                            <th scope="col">Sleeping</th>
                            <th scope="col">Zombie</th>
                            <th scope="col">Stopped</th>
                            <th scope="col">Total</th>
                        </tr>
                    </thead>
                    <tbody>
                        {procesos.map((proceso, index) => <Row key={index} proceso={proceso} />)}
                    </tbody>
                </table>
            </div>
        );
    }
}

export default Table;
import React from 'react';

// Ahora puedes usar maquinaSeleccionada en tu componente.
class Row extends React.Component {
    render() {
        const { proceso } = this.props;
        return (
            <tr>
                <td><p>{proceso.pid}</p></td>
                <td><p>{proceso.name}</p></td>
                <td><p>{proceso.user}</p></td>
                <td><p>{proceso.estado}</p></td>
                <td><p>{proceso.ram}</p></td>
            </tr>
        );
    }
}

class SecondTable extends React.Component {
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

        fetch('http://35.245.67.156:4000/proceso/${maquina}', {
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
                            <th scope="col">PID</th>
                            <th scope="col">Name</th>
                            <th scope="col">User</th>
                            <th scope="col">State</th>
                            <th scope="col">RAM</th>
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

export default SecondTable;
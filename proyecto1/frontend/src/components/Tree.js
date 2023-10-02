import { useState, useEffect } from 'react';
import App from '../App';

function Tree() {
  const [parentData, setParentData] = useState([]);
  const [childData, setChildData] = useState([]);
  const [activeItem, setActiveItem] = useState(null);
  let maquina = App.maquina;
  const fetchData = () => {
    fetch('http://35.245.67.156:4000/proceso/${maquina}')
      .then(response => response.json())
      .then(data => setParentData(data));

    fetch('http://35.245.67.156:4000/hijo/${maquina}')
      .then(response => response.json())
      .then(data => setChildData(data));
  };

  useEffect(() => {
    fetchData();
    const intervalId = setInterval(fetchData, 1000);

    return () => {
      clearInterval(intervalId);
    };
  }, []);

  const getChildren = parentId => {
    return childData.filter(child => child.pid_padre === parentId);
  };

  const toggleAccordionItem = itemId => {
    if (itemId === activeItem) {
      setActiveItem(null);
    } else {
      setActiveItem(itemId);
    }
  };

  return (
    <div>
      {parentData.map(parent => (
        <div className='container-tree' key={parent.pid}>
          <button onClick={() => toggleAccordionItem(parent.pid)}>
            {parent.name}
          </button>
          {activeItem === parent.pid && (
            <div>
              <p>{parent.estado} - {parent.ram}MB</p>
              <ul>
                {getChildren(parent.pid).map(child => (
                  <li key={child.idHijos}>{child.name}</li>
                ))}
              </ul>
            </div>
          )}
        </div>
      ))}
    </div>
  );
}

export default Tree;
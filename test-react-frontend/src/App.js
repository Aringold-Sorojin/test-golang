import React, {useState, useEffect} from 'react';
import './App.css';
import axios from 'axios';

function App(props) {
  const [count, setCount] = useState('0');
  const [increment, setIncrement] = useState(0)

  useEffect(() => {
    setCount(readCounter());
    console.log(count)
  }, []);

  function readCounter() {
    axios.post('http://127.0.0.1:8080/readCounter').then(res => {
      setCount(res.data.result);
    }).catch(e => {
      console.log(e);
    })
    return 0;
  }

  function incrementByOne() {
    axios.post('http://127.0.0.1:8080/incrementByOne').then(res => {
      setCount(res.data.result);
    }).catch(e => {
      console.log(e);
    })
  }

  function incrementByN() {
    axios.post('http://127.0.0.1:8080/incrementByN', {increment: increment}).then(res => {
      setCount(res.data.result);
    }).catch(e => {
      console.log(e);
    })
  }

  return (
    <div className="App">
      <div className="min-h-screen bg-gray-100 flex items-center justify-center">
        <div className="bg-white p-8 rounded-lg shadow-md">
          <h1 className="text-2xl font-bold mb-4">Massa Counter</h1>
          <p className="text-xl mb-4">Current Count: </p>
          <div className="flex mb-4">
            <h1>{count}</h1>
            <input
              type="text"
              onChange={(e) => setIncrement(parseInt(e.target.value))}
              className="border rounded px-2 py-1 mr-2"
            />
            <button
              onClick={incrementByN}
              className="bg-blue-500 text-white px-4 py-2 rounded"
            >
              Increment By N
            </button>
            <button
              onClick={incrementByOne}
              className="bg-blue-500 text-white px-4 py-2 rounded"
            >
              Increment By One
            </button>
          </div>
          {props.error && <p className="text-red-500">{props.error}</p>}
        </div>
      </div>
    </div>
  );
}

export default App;

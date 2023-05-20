import React, { useState } from 'react';
import axios from 'axios';
import './Login.css';

axios.defaults.headers.post['Access-Control-Allow-Origin'] = '*';

function Login({ onLogin }) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [errorMessage, setErrorMessage] = useState('');

  const handleSubmit = async (e) => {
    e.preventDefault();
    setErrorMessage('');

    try {
      const response = await axios.get('http://127.0.0.1:8080/login', {
        headers: {
          Authorization: `Basic ${btoa(`${username}:${password}`)}`,
          Accept: '*/*',
        },
      });
      const jwt = response.data.token;

      localStorage.setItem('jwt', jwt);
      localStorage.setItem('username', username);

      onLogin();
    } catch (error) {
      setErrorMessage('Login failed. Please check your credentials.');
    }
  };

  return (
    <div className="login-container">
      <h2>Login</h2>
      {errorMessage && <p className="error-message">{errorMessage}</p>}
      <form onSubmit={handleSubmit}>
        <div>
          <label>Username:</label>
          <input
            type="text"
            value={username}
            placeholder='username'
            onChange={(e) => setUsername(e.target.value)}
          />
        </div>
        <div>
          <label>Password:</label>
          <input
            type="password"
            value={password}
            placeholder='password'
            onChange={(e) => setPassword(e.target.value)}
          />
        </div>
        <button type="submit">Login</button>
      </form>
    </div>
  );
}

export default Login;

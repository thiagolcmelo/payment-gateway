import React, { useState } from 'react';
import axios from 'axios';
import './Dashboard.css'; // Import the CSS file for styling

const Dashboard = ({ onLogout }) => {
  const [isLoggedIn, setIsLoggedIn] = useState(true); // Replace with your login state logic
  const [shopper, setShopper] = useState('');
  const [currency, setCurrency] = useState('USD');
  const [amount, setAmount] = useState('');
  const [errorMessage, setErrorMessage] = useState('');
  const [submissions, setSubmissions] = useState([]);

  const handleFormSubmit = (event) => {
    event.preventDefault();
    // Reset error message
    setErrorMessage('');

    // Shopper validation
    if (!shopper) {
        setErrorMessage('Please select a shopper.');
        return;
    }

    // Amount validation
    if (!amount || isNaN(amount) || Number(amount) <= 0) {
        setErrorMessage('Please enter a valid positive number for the amount.');
        return;
    }

    const cur_shopper = hardCodedShoppers.find((s) => s.id === parseInt(shopper));
    const purchate_time = new Date().toISOString().slice(0, -1);
    
    const url = 'http://127.0.0.1:8080/payment';
    const data = {
        "amount": parseFloat(amount),
        "currency": currency,
        "purchate_time": purchate_time,
        "validation_method": "push",
        "card": {
            "number": cur_shopper["card"]["number"],
            "name": cur_shopper["card"]["name"],
            "expire_month": cur_shopper["card"]["expire_month"],
            "expire_year": cur_shopper["card"]["expire_year"],
            "cvv": cur_shopper["card"]["cvv"],
        },
        "metadata": "shopper 0"
    };
    const token = localStorage.getItem('jwt');
    const config = {
      headers: {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json',
        Accept: '*/*',
      },
    };

    axios.post(url, data, config).then((response) => {
        const id = response.data.id;
        const status = response.data.status;
        const bank_message = response.data.bank_message

        // Create new submission object
        const newSubmission = {
            id: id,
            shopper: cur_shopper.display,
            currency,
            amount,
            status: status,
            bank_message: bank_message,
            timestamp: purchate_time,
        };

        // Update submissions state with the new submission
        setSubmissions([newSubmission, ...submissions]);

        // Clear form inputs
        setShopper('');
        setCurrency('USD');
        setAmount('');
    }).catch((error) => {
      console.log('Error fetching submission data:', error);
      setErrorMessage('Could not request payment');
    });
  };

  const handleCheckSubmission = (submission) => {
    const token = localStorage.getItem('jwt');
    const config = {
      headers: {
        Authorization: `Bearer ${token}`,
        'Content-Type': 'application/json',
        Accept: '*/*',
      },
    };
    axios.get(`http://127.0.0.1:8080/payment/${submission.id}`, config)
      .then((response) => {
        const status = response.data.status;
        const bank_message = response.data.bank_message;
  
        // Find the index of the submission in the submissions array
        const submissionIndex = submissions.findIndex((item) => item.id === submission.id);
  
        // Update the submission with the new data
        const updatedSubmissions = [...submissions];
        updatedSubmissions[submissionIndex].status = status;
        updatedSubmissions[submissionIndex].bank_message = bank_message;
  
        // Update the submissions state with the updated array
        setSubmissions(updatedSubmissions);
      })
      .catch((error) => {
        console.log('Error fetching submission data:', error);
      });
  };

  const handleLogout = () => {
    localStorage.setItem('jwt', '');
    setIsLoggedIn(false);
    onLogout();
  };

  const originalString = localStorage.getItem('username');
  const id = originalString.slice(-1);
  const firstLetter = originalString.charAt(0).toUpperCase();
  const header = "Welcome, " + firstLetter + originalString.slice(1, -1) + " " + id + " Ltd.";

  const hardCodedShoppers = [
    {
        "id": 0,
        "name": "shopper 0",
        "card": {
            "number": "1111-2222-3333-4444",
            "name": "shopper 0",
            "expire_month": 10,
            "expire_year": 2050,
            "cvv": 123
        },
        "display": "shopper 0 - ****-****-****-4444",
    },
    {
        "id": 1,
        "name": "shopper 1",
        "card": {
            "number": "5555-6666-7777-8888",
            "name": "shopper 1",
            "expire_month": 10,
            "expire_year": 2040,
            "cvv": 456
        },
        "display": "shopper 1 - ****-****-****-8888",
    },
    {
        "id": 2,
        "name": "shopper 2",
        "card": {
            "number": "9999-1010-1111-1212",
            "name": "shopper 2",
            "expire_month": 3,
            "expire_year": 2045,
            "cvv": 789
        },
        "display": "shopper 2 - ****-****-****-1212",
    },
    {
        "id": 4,
        "name": "shopper 4",
        "card": {
            "number": "1313-1414-1515-1616",
            "name": "shopper 5",
            "expire_month": 1,
            "expire_year": 2070,
            "cvv": 987
        },
        "display": "shopper 4 - ****-****-****-1616",
    },
    {
        "id": 5,
        "name": "shopper 5",
        "card": {
            "number": "1717-1818-1919-2020",
            "name": "shopper 5",
            "expire_month": 1,
            "expire_year": 2070,
            "cvv": 987
        },
        "display": "shopper 5 - ****-****-****-2020",
    }];

  return (
    <div className="dashboard-container">
      <header className="dashboard-header">
        <h1 className="dashboard-title">{header}</h1>
        <div className="header-right">
          {isLoggedIn && (
            <button className="logout-button" onClick={handleLogout}>Logout</button>
          )}
        </div>
      </header>
      <form className="form-container" onSubmit={handleFormSubmit}>
        <label htmlFor="shopper" className="form-label">Shopper:</label>
        <select id="shopper" className="form-input" value={shopper} onChange={(e) => setShopper(e.target.value)}>
            <option value="">Select Shopper</option>
            {hardCodedShoppers.map((shopper) => (
                <option key={shopper.id} value={shopper.id}>
                {shopper.display}
                </option>
            ))}
        </select>

        <label htmlFor="currency" className="form-label">Currency:</label>
        <select id="currency" className="form-input" value={currency} onChange={(e) => setCurrency(e.target.value)}>
          <option value="USD">USD</option>
          <option value="EUR">EUR</option>
          <option value="GBP">GBP</option>
        </select>

        <label htmlFor="amount" className="form-label">Amount:</label>
        <input type="text" id="amount" className="form-input" value={amount} onChange={(e) => setAmount(e.target.value)} />

        <button type="submit" className="submit-button">Request Payment</button>
      </form>

      {errorMessage && <p style={{ color: 'red' }}>{errorMessage}</p>}

      <h2>Previous Requests</h2>

      <table className="submission-table">
        <thead>
          <tr>
            <th>Shopper</th>
            <th>Currency</th>
            <th>Amount</th>
            <th>Date</th>
            <th>Status</th>
            <th>Bank Message</th>
            <th>Action</th>
          </tr>
        </thead>
        <tbody>
          {submissions.map((submission) => (
            <tr key={submission.id}>
              <td>{submission.shopper}</td>
              <td>{submission.currency}</td>
              <td>{submission.amount}</td>
              <td>{submission.timestamp}</td>
              <td>{submission.status}</td>
              <td>{submission.bank_message}</td>
              <td>
                <button onClick={() => handleCheckSubmission(submission)}>Check</button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};

export default Dashboard;
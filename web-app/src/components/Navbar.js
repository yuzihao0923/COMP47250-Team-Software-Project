// src/components/Navbar.js
import React from 'react';
import { Link } from 'react-router-dom';

const Navbar = () => {
    return (
        <nav>
            <ul>
                <li>
                    <Link to="/producer">Producer</Link>
                </li>
                <li>
                    <Link to="/consumer">Consumer</Link>
                </li>
                <li>
                    <Link to="/broker">Broker Console</Link>
                </li>
                <li>
                    <Link to="/login">Login</Link>
                </li>
            </ul>
        </nav>
    );
};

export default Navbar;

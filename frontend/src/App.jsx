import React from 'react';
import {
    BrowserRouter as Router,
    Routes,
    Route,
} from "react-router-dom";

import Navbar from 'react-bootstrap/Navbar';
import Nav from 'react-bootstrap/Nav';
import Form from 'react-bootstrap/Form';
import FormControl from 'react-bootstrap/FormControl';
import Button from 'react-bootstrap/Button';
import './App.css';

import 'bootstrap/dist/css/bootstrap.min.css';
import CoffeeList from "./CoffeeList.jsx"
import Admin from "./Admin.jsx";

function App(){
    return (
        <div className="App">
            <Navbar bg="light" expand="lg">
                <Navbar.Brand href="#home">Coffee Shop</Navbar.Brand>
                <Navbar.Toggle aria-controls="basic-navbar-nav" />
                <Navbar.Collapse id="basic-navbar-nav">
                    <Nav className="mr-auto">
                        <Nav.Link href="#home">Home</Nav.Link>
                        <Nav.Link href="#link">Link</Nav.Link>
                    </Nav>
                    <Form inline="true">
                        <div className="search-container">
                            <FormControl type="text" placeholder="Search" className="mr-sm-2" />
                            <Button variant="outline-success">Search</Button>
                        </div>

                    </Form>
                </Navbar.Collapse>
            </Navbar>
            <Routes>
                <Route path="/admin" element={<Admin/>}/>
                <Route path="/" element={<CoffeeList/>}/>
            </Routes>

        </div>
    )
}

export default App
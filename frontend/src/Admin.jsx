import React, { useState } from 'react';
import Form from 'react-bootstrap/Form';
import Col from 'react-bootstrap/Col';
import Row from 'react-bootstrap/Row';
import Button from 'react-bootstrap/Button';
import Container from 'react-bootstrap/Container';
import Toast from './Toast.jsx';

import axios from 'axios';

const Admin = () => {
    const [validated, setValidated] = useState(false);
    const [id, setId] = useState('');
    const [file, setFile] = useState(null);
    const [buttonDisabled, setButtonDisabled] = useState(false);
    const [toastShow, setToastShow] = useState(false);
    const [toastText, setToastText] = useState('');

    const handleSubmit = (event) => {
        event.preventDefault();

        const form = event.currentTarget;
        if (!form.checkValidity()) {
            event.stopPropagation();
            setValidated(true);
            return;
        }

        setButtonDisabled(true);
        setToastShow(false);

        const data = new FormData();
        data.append('file', file);
        data.append('id', id);

        axios
            .post(`http://localhost:9091/images/${id}/${file.name}`, data, {
                headers: { 'content-type': `multipart/form-data; boundary=${data._boundary}` },
            })
            .then((res) => {
                console.log(res);
                let toastText = '';
                if (res.status === 200) {
                    toastText = 'Uploaded file';
                } else {
                    toastText = `Unable to upload file. Error: ${res.statusText}`;
                }

                setButtonDisabled(false);
                setToastShow(true);
                setToastText(toastText);
            })
            .catch((error) => {
                console.log('Error: ' + error);
                setButtonDisabled(false);
                setToastShow(true);
                setToastText(`Unable to upload file. ${error}`);
            });
    };
    console.log(file)
    const changeHandler = (event) => {
        if (event.target.name === 'file') {
            setFile(event.target.files[0]);
            setToastShow(false);
            return;
        }

        setId(event.target.value);
        setToastShow(false);
    };

    return (
        <div>
            <h1 style={{ marginBottom: '40px' }}>Admin</h1>
            <Container className="text-left">
                <Form noValidate validated={validated} onSubmit={handleSubmit}>
                    <Form.Group as={Row} controlId="productID">
                        <Form.Label column sm="2">
                            Product ID:
                        </Form.Label>
                        <Col sm="6">
                            <Form.Control
                                type="text"
                                name="id"
                                placeholder=""
                                required
                                style={{ width: '80px' }}
                                value={id}
                                onChange={changeHandler}
                            />
                            <Form.Text className="text-muted">Enter the product id to upload an image for</Form.Text>
                            <Form.Control.Feedback type="invalid">Please provide a product ID.</Form.Control.Feedback>
                        </Col>
                        <Col sm="4">
                            <Toast show={toastShow} message={toastText} />
                        </Col>
                    </Form.Group>
                    <Form.Group as={Row}>
                        <Form.Label column sm="2">
                            File:
                        </Form.Label>
                        <Col sm="10">
                            <Form.Control
                                type="file"
                                name="file"
                                placeholder=""
                                required
                                onChange={changeHandler}
                            />
                            <Form.Text className="text-muted">Image to associate with the product</Form.Text>
                            <Form.Control.Feedback type="invalid">Please select a file to upload.</Form.Control.Feedback>
                        </Col>
                    </Form.Group>
                    <Button type="submit" disabled={buttonDisabled}>
                        Submit form
                    </Button>
                </Form>
            </Container>
        </div>
    );
};

export default Admin;

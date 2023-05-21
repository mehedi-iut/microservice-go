import React, { useState, useEffect } from 'react';
import Toast from 'react-bootstrap/Toast';

const Toaster = (props) => {
    const [show, setShow] = useState(false);

    const hide = () => {
        setShow(false);
    };

    useEffect(() => {
        setShow(props.show);
    }, [props.show]);

    return (
        <Toast onClose={hide} show={show} delay={3000} autohide>
            <Toast.Header>
                <strong className="mr-auto">File Upload</strong>
            </Toast.Header>
            <Toast.Body>{props.message}</Toast.Body>
        </Toast>
    );
};

export default Toaster;

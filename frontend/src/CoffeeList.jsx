import React, {useState, useEffect} from "react";
import Table from "react-bootstrap/Table"
import axios from 'axios'

function CoffeeList(){
    const [products, setProducts] = useState([])

    const readData = () =>{
        axios
            .get("http://localhost:9090" + '/products')
            .then(function(response){
                console.log(typeof(response.data))
                console.log(response.data)
                setProducts(response.data)
            })
            .catch(function(error){
                console.log(error)
            })
    }

    useEffect(()=>{
        readData()
    }, [])

    const getProducts = () =>{
        let table = []
        for (let i=0; i<products.length; i++){
            table.push(
                <tr>
                    <td>{products[i].name}</td>
                    <td>{products[i].price}</td>
                    <td>{products[i].description}</td>
                </tr>
            )
        }
        return table
    }

    // const getProducts = () =>{
    //     return products.map((product, i) =>(
    //         <tr key={i}>
    //             <td>{product.name}</td>
    //             <td>{product.price}</td>
    //             <td>{product.sku}</td>
    //         </tr>
    //     ))
    // }

    return (
        <div>
            <h1 style={{marginBottom: '40px'}}>Menu</h1>
            <Table>
                <thead>
                    <tr>
                        <th>Name</th>
                        <th>Price</th>
                        <th>Description</th>
                    </tr>
                </thead>
                <tbody>{getProducts()}</tbody>
            </Table>
        </div>
    )
}

export default CoffeeList
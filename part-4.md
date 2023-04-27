# Micro-service in go
## Part-4

In this part, we want to add more rest api method. In the handler **ServeHTTP** method
```go
if r.Method == http.MethodPost {
    p.addProduct(rw, r)
    return
}
```
Now we will implement **addProduct** method
```go
func(p *Products) addProduct(rw http.ResponseWriter, r *http.Request){
    p.l.Println("Handle POST Products")
}
```
Now we need to manupulate the data we get from post request and convert it into our own **Product** object. So to do that we need to navigate to the data folder and modify the **products.go**

```go
func(p *Product) FromJSON(r io.Reader) error {
    e := json.NewDecoder(r)
    return e.Decode(p)
}
```
So we are adding a method **FromJSON** to the **Product** struct. **FromJSON** take **io.Reader** which is our **http.Request**. create **NewDecoder** with the **io.Reader**. then we can call **Decode** method with destination format. in our case **Product** with **p**

```go
prod := &data.Product{}
err := prod.FromJSON(r.Body)
if err != nil {
    http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
}
p.l.Printf("Prod: %#v", prod)
```

In Go programming, it is common to use a buffered reader when receiving HTTP requests. This is because Go does not read everything from the client immediately upon receiving the request. Instead, it buffers some of the data, which can be a large amount of data such as a 10GB file.

Purpose:

The purpose of using a buffered reader is to allow the progressive reading of the incoming data, rather than having to buffer all of the data into memory and allocate memory for the entire request. This helps to prevent memory overload and improve the overall performance of the application.

Conclusion:

In conclusion, using a buffered reader is a recommended approach when receiving HTTP requests in Go programming, especially when dealing with large amounts of data. This approach allows for a more efficient handling of incoming data, by progressive reading, and helps to prevent memory overload.

Now, we need to save the data that is converted by docoder to our productList.
So we create **AddProduct** function to do that.
```go
func AddProduct(p *Product){

}
```

but first we need to generate the **ID**. whether we add to **product** struct or database, we need to generate the unique **ID**
```go
func getNextID()int{
    lp := productList[len(productList)-1]
    return lp.ID+1
}
```
Now we will set the product in productList with updated ID
```go
func AddProduct(p *Product){
    p.ID = getNextID()
    productList = append(productList, p)
}
```
Now in the product handler, we will call **AddProduct** to add the data from post body to **productList**
```go
prod := &data.Product{}
err := prod.FromJSON(r.Body)
if err != nil {
    http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
}
data.AddProduct(prod)
```

### About REST API Method
POST, PUT, and PATCH are HTTP methods used to transfer data to a server.

* **POST**: is used to submit data to the server to be processed, usually resulting in the creation of a new resource or updating an existing resource. It allows for sending data in the body of the request.

* **PUT**: is used to send data to the server to update an existing resource. It requires the client to include the entire updated representation of the resource in the request body. Unlike POST, PUT is idempotent, which means multiple identical PUT requests should result in the same state of the resource.

* **PATCH**: is used to partially update an existing resource. It allows for sending only the changes to the resource, rather than sending the entire representation. PATCH is useful when you need to update a resource but don't need to resend all of its data.

In summary, use POST for creating a resource, PUT for replacing a resource, and PATCH for updating a resource.

Now will add **PUT** method. to implement **PUT** method, we need to first extract the **id** of the resource from the **URI**. but go standard library doesn't provide any method to do that, we need to extract that manually. in this case go framework like **Gin**, **Gorilla** help

```go
if r.Method == http.MethodPut {
    // expect the id in the URI
    p := r.URL.Path
}
```

so, here we extract the url path from **http.Request** and assign to variable **p** and we will use **regex** on p to extract the **ID**

```go
if r.Method == http.MethodPut {
    // expect the id in the URI
    reg := regexp.MustCompile('/([0-9])+')
    g := reg.FindAllStringSubmatch(r.URL.Path, -1)

    if len(g) != 1{
        http.Error(rw, "Invalid URI", http.StatusBadRequest)
        return
    }

    if len(g[0]) != 2 {
        http.Error(rw, "Invalid URI", http.StatusBadRequest)
        return
    }

    idString := g[0][1]
    id, err := strconv.Atoi(idString)

    if err != nil {
        http.Error(rw, "Invalid URI", http.StatusBadRequest)
        return
    }

    p.l.Println("got id", id)

    p.updateProducts(id, rw, r)
    return
}
```

Above code is only to get **ID** from the URL. so we need to use go framework to handle that, because URL contains many complex thing that we need extract, we wil be very complex if we do that manually

Now we got our **ID**, we wil use that to update the **ProductList** array. To do that we create another method for **Products**, **updateProducts**.

```go
func (p Products) updateProducts(id int, rw http.ResponseWriter, r *http.Request){
    p.l.Println("Handle PUT Product")

    prod := &data.Product{}

    err := prod.FromJSON(r.Body)
    if err != nil {
        http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
    }

    data.UpdateProduct(id, prod)
}
```

Now we need to crete a function **UpdateProduct** in the data module. To do that we need to first find the product. we will create **findProduct** function
```go
// create a structured error
var ErrProductNotFound = fmt.Errorf("Product not founc")

func findProduct(id int) (*Product, int, error){
    for i, p := range productList{
        if p.ID == id{
            return p, i, nil
        }
    }
    return nil, -1, ErrProductNotFound
}
```
Now we will implement the **UpdateProduct** function
```go
func UpdateProduct(id int, p *Product) error {
    _, pos, err := findProduct(id)
    if err != nil{
        return err
    }

    p.ID = id
    productList[pos] = p
    return nil
}
```

In the **products** handler, we need to update that
```go
func (p Products) updateProducts(id int, rw http.ResponseWriter, r *http.Request){
    p.l.Println("Handle PUT Product")

    prod := &data.Product{}

    err := prod.FromJSON(r.Body)
    if err != nil {
        http.Error(rw, "Unable to unmarshal json", http.StatusBadRequest)
    }

    err = data.UpdateProduct(id, prod)
    if err == data.ErrProductNotFound{
        http.Error(rw, "Product not found", http.StatusNotFound)
        return
    }

    if err != nil {
        http.Error(rw, "Product not Found", http.StatusInternalServerError)
        return
    }

}
```
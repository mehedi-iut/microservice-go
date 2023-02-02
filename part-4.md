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

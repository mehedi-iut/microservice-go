# Micro-service in go
## Part-8
## Swagger Client

From previous episode code, we can now add swagger client in our project. but why we need swagger client?

*the function of the Swagger client in Go is to automate the generation of client libraries for RESTful APIs based on the Swagger specification, making it easier for developers to interact with the APIs and reducing the amount of boilerplate code they need to write.*

Now to create swagger client first we need to create directory in root called **client**. 
```bash
mkdir client
cd client
```

Now we need to run swagger command to generate client
``` bash
swagger generate client -f ../microservice-go/swagger.yaml -A product-api
```

you can also use ```$(go env GOPATH)/bin/swagger``` to generate the client.



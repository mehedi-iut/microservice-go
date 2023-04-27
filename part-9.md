# Micro-service in go
## Part-9
## CORS

In this blog we will discuss about CORS (Cross Origin Resource Sharing)

From previous code block, there is a frontend code written in React. This is running in localhost:3000 and our api server running in localhost:9090. So whenever I try to access the my api server from frontend, I will get **CORS** error. because they are not in same origin, frontend is running in port **3000** and api is running in **9090**. In the cors i need to add **localhost:3000**. For more information follow [link](https://medium.com/@baphemot/understanding-cors-18ad6b478e2b) [link](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)

To add CORS in our code, In **main.go**, first we need to import gorilla cors

```go
gohandlers "github.com/gorilla/handlers"
```

then we need to define cors handler
```go
ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"http://localhost:3000"}))

s := http.Server{
    Addr: "localhost:9090",
    Handler: ch(sm),
    Errorlog: l,
    ReadTimeout: 5 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout: 120 * time.Second,
}
```

if we run the frontend code by navigating to frontend folder and run ```yarn start``` we will get an error due to nodejs upgrade and some package depricated. To solve the issue you can follow this [link](https://gankrin.org/how-to-fix-error-digital-envelope-routinesunsupported-in-node-js-or-react/)


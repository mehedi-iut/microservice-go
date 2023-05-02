# Microservice in go
## gRPC Error messages in Unary RPCs

hello welcome to another episode of building microservices with go and today we're gonna look at is gRPC error handling that's right error handling in gRPC because it's not the same as rest and I think this one's going to be confusing for some folks now let's take a quick dig into this, so error handling in gRPC, now the key thing that you need to remember is you're not using the HTTP transport in the same way that you would with a restful approach so the restful approach you're issuing HTTP requests and you're using HTTP status codes to measure the response so you're gonna have things like a 200 or you can have things like a 404 for not found where you can follow that standard approach now the thing with gRPC is the gRPC doesn't use HTTP status codes so an error in gRPC or a failed request could still return a HTTP 200 because the transport was decoupled from the actual messages itself with gRPC. what you're doing is you're encoding your errors into your protobufs so you do have the concept of error handling and you do have the concept of error codes. now error codes they don't correspond to HTTP status code, so don't think that an OK in gRPC is going to be the same as HTTP 200 it's actually code 0 but you do have the case of general errors status deadlines you know internal status error etc [link](https://grpc.io/docs/guides/error/)and we're going to take a look at how we can implement these and how we can build these and use them and it's it's pretty straightforward

Now let's look at the unary operator **getRate** below
```go
func(c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error){
    c.log.Info("Handle request for GetRate", "base", rr.GetBase(), "dest", rr.GetDestination())

    rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
    if err != nil{
        return nil, err
    }

    return &protos.RateResponse{Base: rr.Base, Destination, Rate: rate}, nil
}
```

now what we want to do here is let's say we want to do some validation so what I want to do is I want to validate that the currency in the destination are not the same as base. so we can write code like below to handle that
```go
if rr.Base == rr.Destination{
    return nil, fmt.Errorf("Base can not be the same as the destination")
}
```

now if you were doing that in **REST** what you would do is you would return an HTTP status code, you probably run something like precondition failed or something to indicate that validation wasn't working and we can do something similar in gRPC so rather than just getting this error message back what we can actually do is use gRPC status because actually under the the hood what is being returned here is an RPC status type. so let's take a look at that. so what we can do to create a more rich error message with gRPC is we can use the status package and **Errorf** so we can do something like below
```go
if rr.Base == rr.Destination{
    err := status.Errorf(
        codes.InvalidArgument,
        "Base currency %s can not be the same as the destination currency %s",
        rr.Base.String(),
        rr.Destination.String(),
    )

    return nil, err
}
```
Now, I can return **err** because **status.Errorf** if you dig into this what you actually have is sort of an underlying error message
```go
func Errorf(c codes.Code, format string, a ...interface{}) error {
    return Error(c, fmt.Sprintf(format, a...))
}
```
but all of the things in the status package in builtin **status.go** are using in the above object which is the RPC status and we'll come back to this in a little bit but for now we're just using this the status utility allows me to do things like create error messages I can set the codes and program it handles all of the serialization for me and it also means that my status object is compatible with error so I've got that now one of the nice things about the gRPc error message is that I can add metadata to it so for example I might have another error message or I might have some data that I want to add. but now if we access **err.** we will not get anything. So we shouldn't go in this route, we need to change ```status.Errorf``` which create error message directly which will compatible with gRPC error we need to use ```status.NewF``` which will return *status*. now we will access other method in *err* variable. if we look into gRPC **status.proto** file in [link](https://github.com/googleapis/googleapis/blob/master/google/rpc/status.proto) and if we look into ```message status``` we will find it has *code*, *message*, *details* with type **google.protobuf.Any**. so in the *details* we can add any protocol buffer as collection. now I can use *details* in *err* variable to get new copy of our *status*, where we can add our own metadata.
```go
err, wde := err.WithDetails(rr)
if wde != nil{
    return nil, wde
}
```
now our new *err* is ```*status.Status```, so we can't return *err*, this isn't error anymore. but if you look at the methods which are on a status object you'll see that one is called **Err()** and what that's going to do is that will return you a go error which is basically a wrap of the gRPC status object encapsulated in the error format so that it can be serialized over the wire and we can get that back.

Now, we can run our *currency* server ```go run main.go``` and test with **grpcurl**
```bash
grpcurl --plaintext -d '{"Base": "USD", "Destination": "USD"}' localhost:9092 Currency/GetRate
```
and we will get the Error that we define above

now let's take a look actually how we can handle this in the client and let's see how we can decode this error message in a client
we basically have our product-api which calls the currency api which is a gRPC API this is where we are calling **GetRate** so we're calling that unary operator in gRPC and what we're doing is we're just getting the error message and we're just passing it on

```go
resp, err := p.currency.GetRate(context.Background(), rr)
if err != nil{
    return -1, err
}
```

now if you look at what you get from the go's gRPC framework from the Code is generated from, you get back ```*protos RateResponse``` and an ```error``` just a plain go error now we know that this is actually a rich gRPC error because we're sending it from currency as we defined above. we're constructing it and were sending it so what we can actually do is we can try and convert this basic go error and deserialize it into a rich gRPC error.

so the first thing we're gonna check if it's not nil then what we can do is we can try and convert it into a gRPC error and again we use that status package
```go
resp, err := p.currency.GetRate(context.Background(), rr)
if err != nil{
    if s, ok := status.FromError(err); ok {
        return -1, fmt.Errorf("Unable to get rate from currency server")
    }
}
```

now we could put some additional information in here so let's have a look what we could say is we could say well we want be the base rate or we can get that from s and we could say well we've got this with details so if we want to get the base rate we can get that from the metadata because we know it's encoded and now so let's get let's get our metadata so I'm going to say metadata and I'm going to do ```md := s.Details()[0]``` that returns me a collection be the first item so I can get that from my first item my collection now it is not typed it's just returning interface, so we can't access our base or destination currency using ```md.``` so you do have to have an understanding of what's kind of going on with these error messages I could just run through things and look at things with with reflection or something like that but I wrote both sides of this so I can actually put a little bit more more information in there you could also package this up into a client if you were going to build a client for for the consumer rather than allowing people to have to implement it themselves but we can cast it quite simply we know it is ```*protos.RateRequest``` I'm going to show you the the art of the possible and then we can do *md* and we can get all of the various different things so I can get the base currency and destination currency
```go
resp, err := p.currency.GetRate(context.Background(), rr)
if err != nil{
    if s, ok := status.FromError(err); ok {
        md := s.Details()[0].(*protos.RateRequest)
        return -1, fmt.Errorf("Unable to get rate from currency server, base: %s, dest: %s", md.Base.String(), md.Destination.String())
    }
    return -1, err
}
```

but what we haven't done yet is we haven't looked at that error code because it's possible that a number of different types of error might be returned from the server and you might be getting sort of different handling on on each of those so in the instance that we get that Invalid argument because we're getting that you know the two currencies which are the same the base and the destination we can handle that so rather than just having this kind of generic case here we can handle that differently
```**products.go**```
```go
resp, err := p.currency.GetRate(context.Background(), rr)
if err != nil{
    if s, ok := status.FromError(err); ok {
        md := s.Details()[0].(*protos.RateRequest)
        if s.Code() == codes.InvalidArgument{
            return -1, fmt.Errorf("Unable to get rate from currency server, destination and base currencies can not be the same, base: %s, dest: %s", md.Base.String(), md.Destination.String())
        }
        return -1, fmt.Errorf("Unable to get rate from currency server, base: %s, dest: %s", md.Base.String(), md.Destination.String())
    }
    return -1, err
}
```

run our **products.go** and run ```curl "localhost:9090/products/1?currency=EUR"``` we will get the rich gRPC error



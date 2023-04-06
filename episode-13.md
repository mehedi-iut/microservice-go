# Microservice in Go
## gRPC

gRPC is an open-source remote procedure call (RPC) framework developed by Google. It allows clients to call methods on a server as if they were calling a local object, abstracting away the details of network communication.

gRPC uses the Protocol Buffers binary serialization format to efficiently encode and transmit data between clients and servers. It also supports multiple programming languages, making it easy to build distributed systems that span multiple platforms and technologies.

One of the main advantages of gRPC is its speed and efficiency. It is designed to be lightweight and performant, making it ideal for building high-performance distributed systems. Additionally, gRPC supports bi-directional streaming, which allows for more complex interactions between clients and servers.

In gRPC, the intention behind this is that you're still using standard protocols in this instance it's gonna be HTTP but rather than JSON what you're gonna use is a binary based message protocol called **protobufs**

so with **protobufs** because they're binary based there are obviously quicker and faster serialize and send it over the wire and with them also what you do is you end up defining these interfaces these proto files and because you define a proto file anybody can generate a client based off your proto file

now in a proto file you define **services** and you define **methods** for **services** and you define **messages** for **methods**

```
syntax = "proto3";
option go_package="/currency";

service Currency {
    rpc GetRate(RateRequest) returns (RateResponse);
}

message RateRequest {
    string Base = 1;
    string Destination = 2;
}

message RateResponse {
    float Rate = 1;
}
```

Now, protobufs has its own type, so that it can be language agnostic.[link](https://protobuf.dev/programming-guides/proto3/#scalar)

From the protobufs we need to generate code for specific language For example go, to do that we need to use **protoC** to generate golang code from protobufs. [protoc installation](https://grpc.io/docs/protoc-installation/)

we also need **grpc** module and plugins.
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```


To generate go code, create a **Makefile** and add the below code
```bash
.PHONY: protos

protos:
    protoc -I protos/ --go_out=. --go-grpc_out=. protos/currency.proto
```

Now run ```make protos``` in the terminal

now if we look into the generated code we can find a interface which we need to implement using struct to have gRPC server
```go
type CurrencyServer interface {
    GetRate(context.Context, *RateRequest) (*RateResponse, error)
}
```
in this interface called **CurrencyServer** which has one method called **GetRate**. This interface is essentially a contract that any implementation of the **CurrencyServer** must follow. The **GetRate** method takes in a context object and a **RateRequest** object and returns a **RateResponse** object and an **error**.

Now, in order to use this interface, you need to create an implementation of it. This implementation will define how the **GetRate** method will work. Once you have your implementation, you can use a code-generated helper function, called **RegisterCurrencyServer**, to map your implementation to the **gRPC server**. This is similar to mapping routes to a server in HTTP.

When you register your implementation with the gRPC server, it will know how to handle incoming requests for the GetRate method. When a client makes a request to the server, the server will call the appropriate method on your implementation of the CurrencyServer interface to process the request and return the response.


Now in the **main.go**
```go
func main(){
    log := hclog.Default()

    gs := grpc.NewServer()
    cs := server.NewCurrency(log)

    protos.RegisterCurrencyServer(gs, cs)
}
```

First we create gRPC server and try to register currency service using ```protos.RegisterCurrencyService```. but to register currency service in gRPC server, we need to create our currency service ```cs := server.NewCurrency(log)```

now we need to implement the currency server using **CurrencyServer** interface

```go server.go
package server
import (
    "context"
    "github.com/hashicorp/go-hclog"
    protos "currency/protos/currency"
)

type Currency struct {
    log hclog.Logger
    protos.UnimplementedCurrencyServer
}

func NewCurrency(l hclog.Logger) *Currency {
    return &Currency{log: l}
}

func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error){
    c.log.Info("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())
    return &protos.RateResponse{Rate: 0.5}, nil
}
```

Now, in the main.go, we need to create gRPC server ```gs.Serve()``` which is similar to http.ListenAndSeve(). but gs.Serve() take **net.Listen** as a parameter
```go
l, err := net.Listen("tcp", ":9092")
if err != nil {
    log.Error("Unable to listen", "error", err)
    os.Exit(1)
}
gs.Serve(l)
```

Now we can run the main.go but we can't test it using *curl*, as curl use REST with JSON.

To test the gRPC we need either write unit test or use *grpcurl*
Here we will use *grpcurl* [link](https://github.com/fullstorydev/grpcurl)

```bash
grpcurl --plaintext localhost:9092 list
```

we will get error ```Failed to list services: server does not support the reflection API```

To resolve this issue, we need to use gRPC module reflection method
```
reflection.Register(gs)
```

but reflection shouldn't be use in production env

In gRPC, reflection is used to provide an additional service called the Reflection service. The Reflection service allows clients to query a gRPC server for information about the gRPC services it exposes.

Now if we want to see method of the services,
```grpcurl --plaintext localhost:9092 list <service-name>```

we can describe method
```grpcurl --plaintext localhost:9092 describe Currency.GetRate```
here, *Currency* is service and *GetRate* is method

we can describe message
```grpcurl --plaintext localhost:9092 describe .RateRequest```

here .RateRequest is message
we will get output like below:
```
RateRequest is a message:
message RateRequest {
    string Base = 1 [json_name = "base"];
    string Destination = 2 [json_name = "destination"];
}
```
Now, we can send data in json format with the field name as json_name above.
if you can't find the ```[json_name = "base"]``` then we need to find that in generated code in **RateRequest** part and find json name. if we want to call *Currency.GetRate* method for data payloadd, we need call *Currency.GetRate*
```grpcurl --plaintext -d '{"base": "GBP", "destination": "USD"}' localhost:9092 Currency.GetRate```



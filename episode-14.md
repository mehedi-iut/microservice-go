# Microservice in Go
## gRPC client

In the previous episode, we create a dummy conversion value of 0.5. But in this episode, we want to use real data from bank. To do that, we need to add **enum** in *proto* file

#### What is enum?

In gRPC, an enum (short for enumeration) is a data type that represents a set of named values. It allows you to define a finite set of values that a field can take, which makes your code more readable and less prone to errors.

In gRPC, enums are typically used to define message fields that can take a limited set of values. For example, you might define an enum to represent the type of a message, or the status of an operation.

Here's an example of how to define an enum in a gRPC service definition file:

```
enum Status {
  OK = 0;
  CANCELLED = 1;
  UNKNOWN = 2;
}
```
In this example, we define an enum called Status that has three values: OK, CANCELLED, and UNKNOWN. Each value is assigned an integer value, starting from 0.

Once you have defined an enum, you can use it as the type for a field in a message definition:

```
message MyMessage {
  Status status = 1;
}
```
In this example, we define a message called MyMessage that has a field called status of type Status. This field can take one of the three values defined in the Status enum.

Using enums in your gRPC service can make your code more expressive and less prone to errors, since it helps ensure that the values passed between client and server are consistent and well-defined.

documentation [link](https://protobuf.dev/reference/protobuf/proto3-spec/#enum_definition)
#### define enum for currency
in the proto file we need to add 
```
enum Currencies {
    EUR=0;
}
```

So our porto file looks like this 
```
syntax = "proto3";

service Currency {
    // GetRate returns the exchange rate for the two provided currency codes
    rpc GetRate(RateRequest) returns (RateResponse);
}

// RateRequest defines the request for a GetRate call
message RateRequest {
    // Base is the base currency code for the rate
    string Base = 1;
    // Destination is the destination currency code for the rate
    string Destination = 2;
}

// RateResponse is the response from a GetRate call, it contains
// rate which is a floating point number and can be used to convert between the
// two currencies specified in the request.
message RateResponse {
    float rate = 1;
}

// Currencies is an enum which represents the allowed currencies for the API
enum Currencies {
    EUR=0;
}
```

Now we want to add the all other currency here with EUR,

```
// Currencies is an enum which represents the allowed currencies for the API
enum Currencies {
  EUR=0;
  USD=1;
  JPY=2;
  BGN=3;
  CZK=4;
  DKK=5;
  GBP=6;
  HUF=7;
  PLN=8;
  RON=9;
  SEK=10;
  CHF=11;
  ISK=12;
  NOK=13;
  HRK=14;
  RUB=15;
  TRY=16;
  AUD=17;
  BRL=18;
  CAD=19;
  CNY=20;
  HKD=21;
  IDR=22;
  ILS=23;
  INR=24;
  KRW=25;
  MXN=26;
  MYR=27;
  NZD=28;
  PHP=29;
  SGD=30;
  THB=31;
  ZAR=32;
}
```


Now we will change the **RateRequest** base and destination from *string* to *enum*

```
// RateRequest defines the request for a GetRate call
message RateRequest {
    // Base is the base currency code for the rate
    Currencies Base = 1;
    // Destination is the destination currency code for the rate
    Currencies Destination = 2;
}
```
so when we change RateRequest from string to enum, we now need to provide the value that is define in the enum list, we supply value other than enum list value it will show error.

Now we need generate golang code using protoc. we have *Makefile* to do that.
```makefile
.PHONY: protos

protos:
	protoc -I protos/ --go_out=. --go-grpc_out=. protos/currency.proto
```

so we can run 
```
make protos
```

Now, we want to call our gRPC service from our product-api. to do that, we need to use **CurrencyClient** interface that can be found in golang generated code.

First we need to import the generated code in the product-api
```go
import (
    protos "microservice-go/currency/currency"
)

```

we need to instantiate new *CurrencyClient* using ```protos.NewCurrencyClient()``` which accept **grpc.ClientConnInterface**
Documentation of gRPC client is [here](https://grpc.io/docs/languages/go/basics/#client)

we need to copy code from the above documentation
```go
conn, err := grpc.Dial("localhost:9092")
if err != nil{
    panic(err)
}
defer conn.Close()
```

now we need to pass the **conn** to the ```protos.NewCurrencyClient```

```go
cc := protos.NewCurrencyClient(conn)
```

Now we create our client. we want to use this currency conversion everytime when we call the other API

to do that we will pass the gRPC client to the handler

```go
ph := handlers.NewProducts(l, v, cc)
```

Now in the **NewProducts** function we need to add the gRPC argument

```go
import (
    protos "microservice-go/currency/currency"
)

type Products struct {
    l *log.Logger
    v *data.Validation
    cc protos.CurrencyClient
}

func NewProducts(l *log.Logger, v *data.Validation, cc protos.CurrencyClient) *Products {
    return &Products{l, v, cc}
}
```

Now in the **ListSingle** method, we need to add our conversion method. because when we query the product, we want to convert that to the desire currency.

in the *get.go* file, we need to add our gRPC currency client

in the above code, we define **cc**, which is a currency service. which has **GetRate** method. we need to call that method. the **GetRate** method accept context, which we can pass as ```context.Background()``` and **RateRequest**. First we need to construct **RateRequest** with base and destination currency from the enum we define in the protos file.
```go
rr := &protos.RateRequest{
    Base: protos.Currencies(protos.Currencies_value["EUR"]),
    Destination: protos.Currencies(protos.Currencies_value["GBP"]),
}

resp, err := p.cc.GetRate(context.Background(), rr)
if err != nil{
    p.l.Println("[Error] error getting new rate", err)
    data.ToJSON(&GenericError{Message: err.Error()}, rw)
    return
}

p.l.Printf("Resp %#v", resp)

prod.Price = prod.Price * resp.Rate
```

in the gRPC, error is different from the traditional error like 2xx, 3xx, 4xx, 5xx. that's why we just add log there

as we will get a ratio between base and destination in the **resp**. we just multiply it with the original price to get the converted price.

Now, if we run the code, we will get error
```panic: grpc: no trnasport security set ( use grpc.WithInsecure() explicitly or set credentials)```

because *grpc.Dial* will use http/2 with https, we need to specify **opts** in the *grpc.Dial* with *grpc.WithInsecure()*

```go
conn, err := grpc.Dial("localhost:9092", grpc.WithInsecure())
```

now, if we run main.go and run ```curl localhost:9090/products/1```

we will get converted value.
here we are running product-api main.go and curl the rest endpoint. but internally it will call our gRPC client and convert the rate.
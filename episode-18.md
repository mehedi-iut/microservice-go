# Microservice in go
## bi-directional gRPC


In our products API, when it calls for **getProduct**, we make a call to our currency API to get the currency. As you may remember from the previous episode, the currency API returns interest rates from the ECB and the European Central Bank. However, we make a call to that currency API with every single request for every product. This is not very efficient because the rate doesn't change very often. In fact, the ECB API that we are using doesn't change very often at all. So it's not efficient to keep calling that service when the rate isn't changing.

We want to implement a more efficient method where the currency service can notify us when there is a change in the rate. To achieve this, we can use a push-based model where the currency service sends us notifications when there is a change in the rate. We can use bi-directional streaming for this pub/sub model.



Bi-directional streaming is a special type of Remote Procedure Call (RPC) in gRPC. If you look at the getRate method here, it's known as a unary method, which has a simple request-response model. There are no open connections or anything like that. It works exactly the same way as a standard HTTP handler would work.

With gRPC and streaming, we can define a method like this. Let's use the subscribeRates method name and add the stream keyword to the method signature. The stream keyword tells the protobuffer and ultimately the gRPC generation that we are creating a streaming API.

We can stream messages from the client to the server by using the subscribeRates method with a RateRequest message. We can also stream messages back to the client using a RateResponse message. We don't have to do streaming both ways; we can choose to do one-way streaming by having a single request and a streaming response or just streaming requests and no streaming response at all.

To enable bi-directional streaming, we need to rebuild our protobufs. ```make protos```

```
service Currency {
    rpc GetRate(RateRequest) returns (RateResponse);
    rpc SubscribeRates(stream RateRequest) returns (stream RateResponse);
}
```

Now our **CurrencyServer** *interface* has new method **SubscribeRates** and we have to implement that method in our code **currency.go**

in the generated code for **CurrencyServer** interface
```go
type CurrencyServer interface {
    GetRate(context.Context, *RateRequest) (*RateResponse, error)
    SubscribeRates(Currency_SubscribeRatesServer) error
}
```

in the **currency.go** we now add the **SubscribeRates** method
```go
func (c *Currency) SubscribeRates(src protos.Currency_SubscribeRatesServer) error {
    return src.Send(&protos.RateResponse{Rate: 12.1})
}
```

in the above code *{Rate: 12.1}* is the dummy value. 
subscriberRates server has two methods that we are going to be using, the *send* and *receive* message now both of these are kind of gonna be handled in in this instance so let's deal with just the send method first

From the above code, I just kind of do this like this I'm only ever gonna get one message that's gonna be sent because in order to be able to continually read and continually write to my client I've got a block inside of this method it's not something which I just immediately return from so let's just refactor that a little bit to to take that into account

```go
// This method handles bidirectional streaming, allowing the server to send messages to the client continuously
func (c *Currency) SubscribeRates(src protos.Currency_SubscribeRatesServer) error {
  for {
    // Send a message to the client
    err := src.Send(&protos.RateResponse{Rate: 12.1})
    if err != nil {
      // If there's an error, return it
      return err
    }
    // Sleep for 5 seconds before sending the next message
    time.Sleep(5 * time.Second)
  }    
}

```

This code defines a method called SubscribeRates which handles bidirectional streaming. This means that the server can send messages to the client continuously.

The src parameter is a Currency_SubscribeRatesServer object, which has two methods that we'll be using: Send and Recv.

In this code, we're only using the Send method, which sends a RateResponse message to the client every 5 seconds.

We've added a loop to continuously send messages, and we're checking for any errors that might occur while sending the message. If an error occurs, we simply return it.

Overall, this code allows the server to send messages to the client continuously, making it a good solution for applications that require real-time updates.

Now we can run our code ```go run main.go``` and to test we will **grpcurl** using following command
```bash
grpcurl --plaintext --msg-template -d @ localhost:9092 Currency/SubscribeRates
```

we will get output every 5 seconds

So how do we handle requests from the client in our service? We can use the **Receive** method on the **CurrencySubscriberServer**, which returns a **RateRequest** and an **error**.
```go
func (c *Currency) SubscribeRates(src protos.Currency_SubscribeRatesServer) error {
  for {
    // Receive a message from the client
    request, err := src.Recv()
    if err != nil {
      return err
    }

    // Handle the received message
    // ...

    // Sleep for some time before receiving the next message
    time.Sleep(5 * time.Second)
  }    
}
```

The **Receive** method is a blocking method, which means that it will block until the client sends a message. It's similar to using a **net** package connection in the standard library. To handle all incoming messages from the client, we can use a for loop.

Inside the loop, we first receive a message using **src.Recv()**, which returns a **RateRequest** and an error. We then check if there was an error and return it if there was. Next, we can handle the received message as required. Finally, we sleep for some time before receiving the next message.

```go
func (c *Currency) SubscribeRates(src protos.Currency_SubscribeRatesServer) error {
  
  for {
    rr, err := src.Recv()
    // this technically not error, 
    // it happens when client close the connection
    if err == io.EOF {
        c.log.Info("Client has closed connection")
        break
    }

    if err != nil{
        c.log.Error("Unable to read from client", "error", err)
        break
    }
  }
  
  for {
    // Send a message to the client
    err := src.Send(&protos.RateResponse{Rate: 12.1})
    if err != nil {
      // If there's an error, return it
      return err
    }
    // Sleep for 5 seconds before sending the next message
    time.Sleep(5 * time.Second)
  }    
}
```

Now the problem is we got two **for** loops and we will block in two places, So we need to use **go func()** around the receive so that we only block in server part

```go
func (c *Currency) SubscribeRates(src protos.Currency_SubscribeRatesServer) error {
  
  go func(){
    for {
        rr, err := src.Recv()
        // this technically not error, 
        // it happens when client close the connection
        if err == io.EOF {
            c.log.Info("Client has closed connection")
            break
        }

        if err != nil{
            c.log.Error("Unable to read from client", "error", err)
            break
        }
    }
  }()
  
  
  for {
    // Send a message to the client
    err := src.Send(&protos.RateResponse{Rate: 12.1})
    if err != nil {
      // If there's an error, return it
      return err
    }
    // Sleep for 5 seconds before sending the next message
    time.Sleep(5 * time.Second)
  }    
}
```

Now we need to handle the bound message in receiver part. so we just added ```c.log.Info("Handle client request", "request", rr)```

```go
func (c *Currency) SubscribeRates(src protos.Currency_SubscribeRatesServer) error {
  
  go func(){
    for {
        rr, err := src.Recv()
        // this technically not error, 
        // it happens when client close the connection
        if err == io.EOF {
            c.log.Info("Client has closed connection")
            break
        }

        if err != nil{
            c.log.Error("Unable to read from client", "error", err)
            break
        }
        // added this inbound message
        c.log.Info("Handle client request", "rquest", rr)
    }
  }()
  
  
  for {
    // Send a message to the client
    err := src.Send(&protos.RateResponse{Rate: 12.1})
    if err != nil {
      // If there's an error, return it
      return err
    }
    // Sleep for 5 seconds before sending the next message
    time.Sleep(5 * time.Second)
  }    
}
```
Now we want to test our **SubscribeRates**. first we need to find the payload we need to provide
```bash
grpcurl --plaintext --msg-template -d @ localhost:9092 describe Currency.SubscribeRates
```
above command show us the stream RateRequest and stream RateResponse
```bash
grpcurl --plaintext --msg-template -d @ localhost:9092 describe .RateRequest
```
above command shows you 
```
Message template:
{
    "Base": "EUR",
    "Destination": "EUR"
}
```
so we need to provide payload in the above format

```bash
grpcurl --plaintext --msg-template -d @ localhost:9092 Currency/SubscribeRates
```
in the above command we added ```-d @```, so it is expecting input from user, so we can paste in terminal after running the above command
```
{
    "Base": "EUR",
    "Destination": "EUR"
}
```
we will see the result in the terminal where we run ```go run main.go```




# Microservice in go
## gRPC Error handling in gRPC Bidirectional streams

In our **currency.go** file in *currency* service, we have our **SubscribeRates** method. this is our bidirectional streaming method. Now in our **currency.proto** file we have below method
```
service Currency {
    rpc GetRate(RateRequest) returns (RateResponse);
    rpc SubscribeRates(stream RateRequest) returns (stream RateResponse);
}
```
so we have ```stream RateRequest``` and ```stream RateResponse```. now in the instance that an error occurs in either the request from the client after processing or when you're sanding the response how do you handle the the error messages and the key thing is it's not that easy because actually here you've got ```rr, err := src.Recv()``` now you could kind of do some stuff you want to do some processing maybe save some
things onto the database or something like that there is no way of sending a direct message back to the caller that sent the the message again which is actually just happening in ```p.client.Send(rr)``` where we're subscribing for updates in our products go but you can't send a response if you return an error what you're actually doing is you're terminating the connection so you're stopping both the inbound on the outbound streaming connections when you return an error from this handler that's
why we block with this code ```rr, err := src.Recv()```.we went into that in sort of a little detail on the last two episodes but we do need a way of sending a message back to the user so how do we do it so you think well or it's a stream why can't I just send an error message back from to the straight to the client you've got to think about streams in gRPC is kind of like I think of it like a pipe so it's a constant flow of water it's one way only so you've got one-way stream for messages which are coming from the client and you've got a one-way
stream for messages which are going to the client so there's no way to interrupt that stream so when we want to send an error message the only way we can do it is we've got to send it in the outbound stream so let's say I processing some stuff in here let's have a look what we got well we've got this subscribe so what we're doing here is we're subscribing for currency rates and we don't have any validation here we're literally just popping these things straight onto a queue but what if we want to do
some validation so what we're gonna do is we're going to validate the subscription doesn't already exist before we just add it on to the queue so let's just run some code through there and look see what that looks like so we're just gonna say check that subscription does not exist so because our subscription is just a flat list we can just loop over it so we can just write below code
```go
// check that subscription doesn't exists
for _, v := range rrs {
    if v.Base == rr.Base && v.Destination == rr.Destination {
        // subscription exists return errors
    }
}
```
above, ```rr.Base``` and ```rr.Destination``` is subscription request

we want to return an error message just letting them know that subscription already exists well we've already said that we can't return an error from this function because what that's going to do is that's going to to break the connection it's going to stop the client being able to send any more messages and we don't want that it's not a terminal
error it's a it's more of a warning in a sense I suppose in this case but what we can do is we can use the send channel because this is a bi-directional streaming API. So we can send an error message back but there's a small problem and the small problem is that if you look at the type for the send back it's just a response so it's just this message here
```
message RateResponse {
    Currencies Base = 1;
    Currencies Destination = 2;
    double Rate = 3;
}
```
in **currency.protos** file

and actually we want this message to be something a little bit smaller we want it to potentially be a rate request. we want to potentially be a an error message so how do we kind of deal with that so what we can do is we can add new message in **currency.proto**
```
message StreamingRateResponse {
    oneof message {
        RateResponse rate_response = 1;
        google.rpc.Status error = 2;
    }
}
```

here we create new message in proto file, it contains **oneof** which help us to return one of the message type. because we're using this **oneof** what's going to happen is that if I set rate response if there's a value set for error it'll get a param wiped out and if there's a value for error the same with rate response you can only have one of these one of the others

now if we run ```make protos```, we will get error, because, in the above code, we try to use gRPC status, but we didn't import it. so we need to import it, but we can't import like go.mod, we need to download the **Status** file from internet and put it inside folder **google/rpc**. so we need to manually add these files into a location where where they're expected. so in the **currency.proto** we need to add 
```
import "google/rpc/status.proto";
```
now we need to crete folder **google/rpc** and create **status.proto** in that path. we can find the **status.proto** in this [link](https://github.com/googleapis/googleapis/blob/master/google/rpc/status.proto). copy the file content and paste it in the **status.proto** file.

now if we run ```make protos``` we will get another error, ```google/protobuf/any.proto: File not found.```, so we need to create **any.proto** in the path **google/protobuf/**. we can find the **any.proto** in this [link](https://github.com/protobuf-net/protobuf-net/blob/main/src/protobuf-net.Reflection/google/protobuf/any.proto)

now if we run the ```make protos```, it will run successfully.

Now, we need to update our **Currency** service in **currency.proto** file to use our **StreamingRateResponse**
```
service Currency {
    rpc GetRate(RateRequest) returns (RateResponse);
    rpc SubscribeRates(stream RateRequest) returns (stream StreamingRateResponse);
}
```

we need to again build with ```make protos```

now, if we getback to our code, where we lef off in **SubscribeRates**
```go
// check that subscription doesn't exists
for _, v := range rrs {
    if v.Base == rr.Base && v.Destination == rr.Destination {
        // subscription exists return errors
        s := status.Newf(
            codes.AlreadyExists,
            "Unable to subscribe for currency as subscription already exists"
        )

        // add the original request as metadata
        s, err := s.WithDetails(rr)
        if err != nil {
            c.log.Error("Unable to add metadata to error", "error", err)
            break
        }
        src.Send(
            &protos.StreamingRateResponse{
                Message: &protos.StreamingRateResponse_Error{
                    Error: s.Proto(),
                },
            },
        )
        break
    }
}

rrs = append(rrs, rr)
c.subscriptions[src] = rrs
```
above code has one problem, so in the for loop, we are checking subscription exists or not, if it exists, it will break the loop, and it will always **append** to **rrs** but it shouldn't, as we already have subscription exists. so we need to fix that


```go
// check that subscription doesn't exists
var validationError *status.Status
for _, v := range rrs {
    if v.Base == rr.Base && v.Destination == rr.Destination {
        // subscription exists return errors
        validationError = status.Newf(
            codes.AlreadyExists,
            "Unable to subscribe for currency as subscription already exists"
        )

        // add the original request as metadata
        validationError, err = validationError.WithDetails(rr)
        if err != nil {
            c.log.Error("Unable to add metadata to error", "error", err)
            break
        }
        break
    }
    // if a validation error return error and continue
    if validationError != nil {
        src.Send(
            &protos.StreamingRateRequest{
                Message: &protos.StreamingRateResponse_Error{
                    Error: validationError.Proto(),
                },
            },
        )
        continue
    }
    // all ok
    rrs = append(rrs, rr)
    c.subscription[src] = rrs
}
```

Now, we need to handle error message in client side in product-api **products.go**

from client side we are orginating error, because we send subscription that already exists. but error message isn't coming back here though. what's going to happen is we're gonna get the error message back in the receive loop which we are dealing in the **handleUpdates()** method. so this handle update is monitoring for subscription updates from the currency server. so this is getting those error message. but now ```sub.Recv()``` is not a simple receive rate response. it's a streaming rate response and streaming rate response can be one of a number of different types.

```**products.go**```
```go
for {
    rr, err := sub.Recv()

    if grpcError := rr.GetError(); grpcError != nil {
        p.log.Error("Error subscribing for rates", "error", grpcError)
        continue
    }

    if resp := rr.GetRateResponse(); resp != nil {
        p.log.Info("Received updated rate from server", "dest", resp.GetDestination().String())

        if err != nil {
            p.log.Error("Error receiving message", "error", err)
            return
        }

        p.rates[resp.Destination.String()] = resp.Rate
    }
}
```

in the above code, we check if we get Error regarding subscribe rates, and we show error, continue the loop and start listening for the next message, then if we didn't get the subscribe rate we are sure that it will be other type error, as we have one of the two types.
so this is now you see how we can handle our inbound stream so our server to client stream, how we can handle messages which might be of a number of different types it might be a rate response a standard sort of
application message or it might be an error this is the only way we can send errors when we're running a gRPC bi-directional stream it's by sending it back down the like the client stream.

Now we need to fix the standard error response, that is sending right now to our changed response. in the **currency.go** in **handleUpdates**. and it is now ```err = k.Send(&protos.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: r})``` and changed to below

```go
err = k.Send(
    &protos.StreamingRateResponse{
        Message: &protos.StreamingRateResponse_RateResponse{
            RateResponse: &protos.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: r},
        },
    },
)
```

Now we can run the currency and product-api service using ```go run main.go``` and then run ```curl "localhost:9090/products/1?currency=GBP"``` and now if we run the curl command again with *GBP*, i.e we are again subscribing to the same currency. we will get error in *products-api* service output.


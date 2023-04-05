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

From the protobufs we need to generate code for specific language example go, to do that we need to use **protoC** to generate golang code from protobufs.

we also need **grpc** module 
# gRPC / Protobuf Go Examples


## gRPC Stream example

This is to demonstrate an asynchronouse source that sends a stream of messages to the
peer. The peer sends messages back (on its own terms) asynchronously


Compile protobuf in the base 'grpc' directory with the following command:

```bash
    protoc -Iexample/ --go_out=plugins=grpc:example example/example.proto
```


### grpc-stream-source.go


### grpc-stream-consumer.go
 
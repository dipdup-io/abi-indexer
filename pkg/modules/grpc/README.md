# gRPC

Metadata gRPC API exposes on 7778 port. Link to [proto file](/pkg/modules/grpc/proto/metadata.proto).

## Endpoints

First of all you have to handshake with the server. `HelloService` is responsible for it.

```protobuf
service HelloService {
    rpc Hello(HelloRequest) returns (HelloResponse);
}

```

* `Hello` - handshake endpoint which requests personal identity for subscriber. It has to be called first. Personal identity will be used in others requests.

```protobuf
message HelloRequest {}

message HelloResponse {
    string id = 1;
}
```


Metadata service implements following endpoints:

```protobuf
service MetadataService {
    rpc SubscribeOnMetadata(DefaultRequest) returns (stream Metadata);
    rpc UnsubscribeFromMetadata(DefaultRequest) returns (Message);

    rpc GetMetadata(GetMetadataRequest) returns (Metadata);
    rpc ListMetadata(ListMetadataRequest) returns (ListMetadataResponse);
    rpc GetMetadataByMethodSinature(GetMetadataByMethodSinatureRequest) returns (ListMetadataResponse);
    rpc GetMetadataByTopic(GetMetadataByTopicRequest) returns (ListMetadataResponse);
}
```

* `SubscribeOnMetadata` - subscribes on new metadata receiving events.

```protobuf
message DefaultRequest {
    string id = 1;
}

// stream of Metadata
message Metadata {
    string address = 1;
    bytes metadata = 2;
    bytes json_schema = 3;
}

```

* `UnsubscribeFromMetadata` - unsubscribes from metadata stream

```protobuf
message DefaultRequest {
    string id = 1;
}
message Message {
    string message = 1;
}
```

* `GetMetadata` - receives ABI by contract address.

```protobuf
message GetMetadataRequest {
    string id = 1;
    string address = 2;
}
message Metadata {
    string address = 1;
    bytes metadata = 2;
    bytes json_schema = 3;
}
```

* `ListMetadata` - receives all ABIs with pagination and sorting.

```protobuf
enum SortOrder {
    ASC = 0;
    DESC = 1;
}

message Page {
    uint64 limit = 1;
    uint64 offset = 2;
    SortOrder order = 3;
}

message ListMetadataRequest {
    string id = 1;
    Page page = 2;
}

message ListMetadataResponse {
    repeated Metadata metadata = 1;
}
```

* `GetMetadataByMethodSinature` - receives all metadata contains certain method signature with sorting and pagination.

```protobuf
message GetMetadataByMethodSinatureRequest {
    string id = 1;
    Page page = 2;
    string signature = 3;
}
```


* `GetMetadataByTopic` - receives all metadata contains certain method signature with sorting and pagination.

```protobuf
message GetMetadataByTopicRequest {
    string id = 1;
    Page page = 2;
    string topic = 3;
}
``` 

## Usage

There are server and client modules in the package.

To run server write the following code which can be found [here](/cmd/indexer/main.go):

```go
grpcModule, err := grpc.NewServer(cfg.GRPC.Server, storage)
if err != nil {
    log.Panic().Err(err).Msg("creating grpc module")
    cancel()
    return
}

metadataIndexer.Subscribe(grpcModule.Subscriber, messages.TopicMetadata)
grpcModule.Start(ctx)

// your code here

if err := grpcModule.Close(); err != nil {
    log.Panic().Err(err).Msg("closing grpc server")
}
```

To create client module write the following code:

```go
grpcClient := grpc.NewClient(cfg.GRPC.Client)                       // create module
grpcClient.Subscribe(yourModule.Subscriber, messages.TopicMetadata) // subscribe on internal events

if err := grpcClient.Connect(ctx); err != nil {                     // create connection to server
    log.Panic().Err(err).Msg("GetMetadata")
    cancel()
    return
}

grpcClient.Start(ctx)                                               // listening for server events

data, err := grpcClient.GetMetadata(ctx, "0x...")                   // receiving metadata by gRPC
if err != nil {
    log.Panic().Err(err).Msg("GetMetadata")
    cancel()
    return
}

data, err := grpcClient.ListMetadata(ctx, 10, 0, pb.SortOrder_ASC)  // receiving list metadata by gRPC
if err != nil {
    log.Panic().Err(err).Msg("ListMetadata")
    cancel()
    return
}

// your code here

if err := grpcClient.Close(); err != nil {
	log.Panic().Err(err).Msg("closing grpc client")
}
```
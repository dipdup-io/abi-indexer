# gRPC

Metadata gRPC API exposes on 7778 port. Link to [proto file](/pkg/modules/grpc/proto/metadata.proto).

## Endpoints


Metadata service implements following endpoints:

```protobuf
service MetadataService {
    rpc SubscribeOnMetadata(DefaultRequest) returns (stream SubscriptionMetadata);
    rpc UnsubscribeFromMetadata(UnsubscribeRequest) returns (UnsubscribeResponse);

    rpc GetMetadata(GetMetadataRequest) returns (Metadata);
    rpc ListMetadata(ListMetadataRequest) returns (ListMetadataResponse);
    rpc GetMetadataByMethodSinature(GetMetadataByMethodSinatureRequest) returns (ListMetadataResponse);
    rpc GetMetadataByTopic(GetMetadataByTopicRequest) returns (ListMetadataResponse);
}
```

* `SubscribeOnMetadata` - subscribes on new metadata receiving events.

```protobuf
message DefaultRequest {}

// stream of Metadata

message SubscriptionMetadata {
    SubscribeResponse subscription = 1;
    Metadata metadata = 2;
}

message SubscribeResponse {
    uint64 id = 1;
}

message Metadata {
    string address = 1;
    bytes metadata = 2;
    bytes json_schema = 3;
}

```

* `UnsubscribeFromMetadata` - unsubscribes from metadata stream

```protobuf
message UnsubscribeRequest {
    uint64 id = 1;
}

message UnsubscribeResponse {
    uint64 id = 1;
    Message response = 2;
}
message Message {
    string message = 1;
}
```

* `GetMetadata` - receives ABI by contract address.

```protobuf
message GetMetadataRequest {
    string address = 1;
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
    Page page = 1;
}

message ListMetadataResponse {
    repeated Metadata metadata = 1;
}
```

* `GetMetadataByMethodSinature` - receives all metadata contains certain method signature with sorting and pagination.

```protobuf
message GetMetadataByMethodSinatureRequest {
    Page page = 1;
    string signature = 2;
}
```


* `GetMetadataByTopic` - receives all metadata contains certain method signature with sorting and pagination.

```protobuf
message GetMetadataByTopicRequest {
    Page page = 1;
    string topic = 2;
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

if err := grpcClient.Connect(ctx); err != nil {                     // create connection to server
    log.Panic().Err(err).Msg("Connect")
    return
}

grpcClient.Start(ctx)                                                          // listening for server events
id, err := grpcClient.SubscribeOnMetadata(ctx, yourModule.Subscriber)          // subscribe on internal events. retruns subscription id which required on unsubscribe.
if err != nil {
    log.Panic().Err(err).Msg("SubscribeOnMetadata")
    return
}

data, err := grpcClient.GetMetadata(ctx, "0x...")                   // receiving metadata by gRPC
if err != nil {
    log.Panic().Err(err).Msg("GetMetadata")
    return
}

data, err := grpcClient.ListMetadata(ctx, 10, 0, pb.SortOrder_ASC)  // receiving list metadata by gRPC
if err != nil {
    log.Panic().Err(err).Msg("ListMetadata")
    return
}

// your code here

if err := grpcClient.UnsubscribeFromMetadata(ctx, yourModule.Subscriber, id); err != nil {
    log.Panic().Err(err).Msg("UnsubscribeFromMetadata")
    return
}

if err := grpcClient.Close(); err != nil {
	log.Panic().Err(err).Msg("closing grpc client")
}
```
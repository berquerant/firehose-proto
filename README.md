# firehose-proto

# Initialize

``` bash
git submodule update --init
make docker
```

# Add proto

Create a new directory and add a proto file to it.

Compile `example/example.proto` by

``` bash
make example/example.pb.go
```

`example/example.pb.go` should belong to package example.

# Test

``` bash
make dev
```

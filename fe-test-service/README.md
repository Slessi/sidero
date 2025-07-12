Simple Log Server Stub
---

This package contains the simple Golang HTTP server that is using gRPC gateway to serve mock Talos logs.
OpenAPIv2 spec is stored in `openapiv2` folder.
OpenAPIv3 spec is stored in `openapiv3` folder.

To build and run the server run:

```
make server
```

This will create binaries in the `_out` folder.

To run the server run one of the executables in the `_out` folder, depending on your OS type and Architecture.
For example for linux it should be:

```
_out/server-linux-amd64
```

The server will start listening on HTTP on port 12000.

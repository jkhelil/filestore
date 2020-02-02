# Filestore client and server
Filestore is a simple file store service (HTTP server and a command line client) that stores plain-text files

## Directory structure
[client/](./client) Command line client
[server/](./server) HTTP server store service
[docker/](./docker) Containing Dockerfiles for client and server

 ## Building The project
1. Build all (client and server)
```bash
make all
```
2. Build the server only
```bash
make linux-server
```
3. Build the client only
```bash
make linux-client
```
Other Build option are available for Mac os users.

## Start a filestore server locally
```bash
STORE=[directory path acting as volume to bind mount inside the server container and serving as the server store]
docker run -v $STORE:/store -p 9090:9090 emircs/filestore-server:latest
```

## Use the client cli
- Download a client cli release for mac os user and add it to you path
```bash
curl -L -O https://github.com/jkhelil/filestore/releases/download/v0.0.1/store-darwin-amd64
```
- Download a client cli for linux user and add it to you path
```bash
curl -L -O https://github.com/jkhelil/filestore/releases/download/v0.0.1/store-linux-amd64
```

1. Add a file to the store
```bash
store add test.txt
```
2. Remove file from the store
```bash
store rm test.txt
```
3. Update file in the store
```bash
store update test.txt
```
4. List file in the store
```bash
store list
```
# kargos-backend

An backend server that provides REST API to frontend and stores data from `kargos-agent`s.

## Features
- gRPC server for recepting data from each `kargos-agent`s.
- Kubernetes API client for retrieving data from the cluster. (using [client-go](https://github.com/kubernetes/client-go))
- DB handler for storing K8s information and data from `kargos-agent`s.
- Providing REST API for frontend.

## Docker
![enter image description here](https://img.shields.io/docker/pulls/kargos/kargos-backend)
Example execution command:
```
docker run -it kargos/kargos-backend
```
> It is highly suggested to use `-p` option to bind host ports into your environment. If you are using Docker-compose, it is suggested to use `expose`.
> 
The backend container will start listening for REST API and gRPC communications in following ports:
- 9000: for REST API (`kargos-frontend`)
- 50001: for gRPC (`kargos-agent`)

### MongoDB
Kargos-backend uses [MongoDB](https://hub.docker.com/_/mongo) in order to store and retrieve data. Therefore there **MUST** be an MongoDB instance (a containered one or just the native one) that shall be running for `kargos-backend`.

### Environment Variables:
In order for agent to work properly, each environment variables must be set in accordance to the backend server that is listening for the agent's data.
> Any environment variables not marked *(optional)* is mandatory arguments for the container to execute.

- `MONGODB_LISTEN_ADDR`: The server IP where MongoDB is listening at.
- `MONGODB_LISTEN_PORT`: The server port where MongoDB is listening at.
- `GRPC_LISTEN_PORT`: The port that will be listening for `kargos-agent`s. This defaults to `50001` if not explicitly set.

## REST API
If you are willing to use REST API for your own purpose, the documents on each endpoints will be discussed in *wiki* section.

> As of 10 Feb 2023, this is still in construction.

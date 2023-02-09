# kargos-agent
An agent that retrieves information from nodes.

## Features
- Container detection using containerd client.
- Child process extraction and metric calculation.
- Sending data to backend for data storage.

## Docker
![enter image description here](https://img.shields.io/docker/pulls/kargos/kargos-agent)
Example execution command:
```
docker run -it --env SERVER_IP="The_server_IP" --env SERVER_PORT="The_server_port" --env GRPC_DELAY="The_interval" --pid=host --volume=/run:/run --volume=/sys:/sys isukim/kargos-agent
```
In order for agent to work properly, each environment variables must be set in accordance to the backend server that is listening for the agent's data.

### Disclaimer
In order to communicate with the host to gather metrics and data, this docker container must have settings of:
- `--pid=host`: For retrieving other container's processes.
- `--volume=/run:/run`: For providing `/run/containerd/containerd.sock` to gather information about each containers. This might be different between users.
- `--volume=/sys:/sys`: For retrieving metrics, such as CPU and RAM usage, for each processes properly.

### Environment Variables:
> Any environment variables not marked *(optional)* is mandatory arguments for the container to execute.
- `SERVER_IP`: The server IP that this agent will be sending information to. 
- `SERVER_PORT`: The server port that this agent will be sending information to.
- `CONTAINERD_SOCK`: The containerd socket. This defaults to `/run/containerd/containerd.sock` if not explicitly set. *(Optional)*
- `GRPC_DELAY`: The interval of sending data to backend server. This defaults to 60 seconds if not explicitly set. *(Optional)*
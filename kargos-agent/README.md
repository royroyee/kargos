# kargos-agent
An agent for [Kargos](https://github.com/boanlab/kargos)

## Docker
Execute this using docker by following command:
```
docker run -it --env SERVER_IP="The_server_IP" --env SERVER_PORT="The_server_port" --env GRPC_DELAY="The_interval" --pid=host --volume=/run:/run --volume=/sys:/sys --volume=/dev:/dev --volume=/etc:/etc --network=host isukim/kargos-agent
```

## Kubernetes
To deploy an agent into your Kubernetes cluster, it is highly suggested to deploy kargos-agent as a DaemonSet. We have included an example of DaemonSet spec in `DaemonSet.yaml`. Please refer to the file for more information

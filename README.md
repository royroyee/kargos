# Kargos

<p align="center"><img src="https://user-images.githubusercontent.com/88774925/218298174-d067dcca-1a9b-4519-ab8e-e2f451966ce9.png" height="400" width="700"></p>

**Kubernetes Management and Monitoring Dashboard**

## Description
Kargos allows you to monitor your Kubernetes cluster's resources, performance with ease.


## Features
- Kubernetes Monitoring
	- Controllers
	- Resources
	- Metrics
	- Events
- Detailed Pod Monitoring
	- Processes
	- Resource usage per processes


## Installation
Kargos supports easy installation and deployment of the ecosystem using:
```bash
kubectl create namespace kargos
kubectl apply -f https://raw.githubusercontent.com/boanlab/kargos/main/kargos.yaml
```

If you would like more diverse preferences, you can use `Docker` to pull each components to your system and set it up for yourself. 
- [Docker Hub](https://hub.docker.com/u/kargos)
- [Kargos-agent](./kargos-agent): For probing data from each nodes.
- [Kargos-backend](./kargos-backend): For storing data and processing data from infra and providing REST API for frontend.
- [Kargos-frontend](./kargos-backend): For providing visual information to the user.




## Contributors
- [Isu Kim](https://github.com/isu-kim)
- [Younghwan Kim](https://github.com/royroyee)
- [Junha Kim](https://github.com/kim-wnsgk)

## License
MIT License

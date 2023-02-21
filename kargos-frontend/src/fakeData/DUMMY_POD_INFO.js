import { FormText } from "@themesberg/react-bootstrap";

export default {
    "podInfo": [
    {
        "index": 1,
        "type": "Image",
        "value": "ubuntu",
    },
    {
        "index": 2,
        "type": "Node",
        "value": "k8s-worker-1",
    },
    {
        "index": 3,
        "type": "IP",
        "value": "10.0.20.112",
    },
    {
        "index": 4,
        "type": "Ports",
        "value": "80/TCP, 8080/TCP",
    },
    {
        "index": 5,
        "type": "Volumes",
        "value": "ubuntu-pv",
    },
    {
        "index": 6,
        "type": "Controlled by",
        "value": "ubuntu-replicaset",
    },
    ],
    "podLog": "Kargos: 2023/02/11 08:24:55 Pod Data stored successfully\nKargos: 2023/02/11 08:24:55 received data from agent kargos-agent-8rf7j\n\Kargos: 2023/02/11 08:24:56 Pod Data stored successfully\nKargos: 2023/02/11 08:24:56 received data from agent kargos-agent-sb8f7\nKargos: 2023/02/11 08:25:00 Pod Data stored successfully\nKargos: 2023/02/11 08:25:00 received data from agent kargos-agent-8rf7j\nKargos: 2023/02/11 08:25:01 Pod Data stored successfully\nKargos: 2023/02/11 08:25:01 received data from agent kargos-agent-sb8f7",
    "containers": [
        {
            "id": "6b1e813f-08f9-4ca2-93c3-417d7c1d51a4",
            "image": "ubuntu/ubuntu:latest",
            "node": "k8s-worker-1",
            "processes": [
                {
                    "name": "/bin/bash",
                    "status": "sleeping",
                    "PID": 2467323,
                    "cpu": Math.floor(Math.random() * 1000),
                    "ram": Math.floor(Math.random() * 1000),
                },
                {
                    "name": "gcc",
                    "status": "sleeping",
                    "PID": 2467324,
                    "cpu": Math.floor(Math.random() * 1000),
                    "ram": Math.floor(Math.random() * 1000),
                },
                {
                    "name": "a.out",
                    "status": "running",
                    "PID": 2467323,
                    "cpu": Math.floor(Math.random() * 1000),
                    "ram": Math.floor(Math.random() * 1000),
                }
            ]
        },
        {
            "id": "2eccfa5f-2d82-44de-93c2-5fa60178a0eb",
            "image": "ubuntu/ubuntu:latest",
            "node": "k8s-worker-2",
            "processes": [
                {
                    "name": "/bin/bash",
                    "status": "sleeping",
                    "PID": 2467321,
                    "cpu": Math.floor(Math.random() * 1000),
                    "ram": Math.floor(Math.random() * 1000),
                },
                {
                    "name": "java",
                    "status": "sleeping",
                    "PID": 2467324,
                    "cpu": Math.floor(Math.random() * 1000),
                    "ram": Math.floor(Math.random() * 1000),
                },
                {
                    "name": "java -jar test.jar",
                    "status": "running",
                    "PID": 2467123,
                    "cpu": Math.floor(Math.random() * 1000),
                    "ram": Math.floor(Math.random() * 1000),
                }
            ]
        },
    ]
};
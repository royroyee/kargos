export default {
    "name": "deployment-1",
    "controllerInfo": [
    {
        "type": "Labels",
        "value": "app=ubuntu",
    },
    {
        "type": "Limits",
        "value": "cpu: 2, memory: 4Gi",
    },
    {
        "type": "Environment",
        "value": "<none>",
    },
    {
        "type": "Mounts",
        "value": "nas1-shared",
    },
    {
        "type": "Volumes",
        "value": "ubuntu-pv",
    },
    {
        "type": "Controlled by",
        "value": "ubuntu-replicaset",
    },
    ],
    "templateContainers": [
        {
            "name": "ubuntu-modified",
            "info": [
                {
                    "type": "Image",
                    "value": "ubuntu/ubuntu:latest"
                },
                {
                    "type": "Port",
                    "value": "22/TCP"
                },
                {
                    "type": "Host Port",
                    "value": "0/TCP"
                },
                {
                    "type": "Command",
                    "value": "\"/bin/bash\", \"-c\", \"while true; do sleep 30; done;\""
                },
            ]
        },
        {
            "name": "ubuntu-apache",
            "info": [
                {
                    "type": "Image",
                    "value": "ubuntu/ubuntu:latest"
                },
                {
                    "type": "Port",
                    "value": "22/TCP, 80/TCP"
                },
                {
                    "type": "Host Port",
                    "value": "0/TCP"
                },
                {
                    "type": "Command",
                    "value": "/bin/bash\", \"-c\", \"service apache2 restart; while true; do sleep 30; done;\""
                },
            ]
        },
    ],
    "volumes": [
        {
            "name": "main-volume",
            "info": [
                {
                    "type": "Type",
                    "value": "PersistentVolumeClaim"
                },
                {
                    "type": "ClaimName",
                    "value": "ubuntu-pvc"
                },
                {
                    "type": "ReadOnly",
                    "value": "false"
                },
            ]
        },
        {
            "name": "volume-shared",
            "info": [
                {
                    "type": "Type",
                    "value": "PersistentVolumeClaim"
                },
                {
                    "type": "ClaimName",
                    "value": "ubuntu-pvc-shared"
                },
                {
                    "type": "ReadOnly",
                    "value": "false"
                },
            ]
        },
    ],
    "conditions": [
        {
            "type": "Progressing",
            "status": "True",
            "reason": "NewReplicaSetAvailable"
        },
        {
            "type": "Available",
            "status": "True",
            "reason": "MinimumReplicasAvailable"
        }
    ],
};

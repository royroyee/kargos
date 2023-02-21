import { FormText } from "@themesberg/react-bootstrap";

/**
 * Generate fake cpu usage percentages from 0 to 100 for last 24 hours.
 * @returns 24 items that ranges from 0 to 100.
 */
function generateRandomPercent() {
    var usage = [];
    for (var i = 0; i < 24; i++) {
        usage.push(Math.floor(Math.random() * 100))
    }
    return usage
}

export default {
    "nodeInfo": [
    {
        "index": 1,
        "type": "OS",
        "value": "Ubuntu 22.04 LTS",
    },
    {
        "index": 1,
        "type": "Hostname",
        "value": "k8s-worker-1",
    },
    {
        "index": 2,
        "type": "IP",
        "value": "172.25.244.1",
    },
    {
        "index": 3,
        "type": "Kubernetes Version",
        "value": "v1.24.1",
    },
    {
        "index": 4,
        "type": "Containerd Version",
        "value": "1.6.7",
    },
    {
        "index": 5,
        "type": "Running Containers",
        "value": "20",
    },
    {
        "index": 6,
        "type": "CPU Cores",
        "value": "40 Cores",
    },
    {
        "index": 7,
        "type": "RAM Capacity",
        "value": "256 GB",
    },
    {
        "index": 8,
        "type": "Status",
        "value": "Ready",
    },
    ],
    "nodeLog": "Jan 09 03:44:07 cloud-06 kubelet[1231]: E0109 03:44:07.062950    1231 eviction_manager.go:560] \"Eviction manager: cannot evict a critical pod\" pod=\"kube-system/etcd-cloud-06\"\n\
Jan 09 03:44:07 cloud-06 kubelet[1231]: E0109 03:44:07.062973    1231 eviction_manager.go:560] \"Eviction manager: cannot evict a critical pod\" pod=\"kube-system/kube-controller-manager-cloud-06\"\n\
Jan 09 03:44:07 cloud-06 kubelet[1231]: E0109 03:44:07.062993    1231 eviction_manager.go:560] \"Eviction manager: cannot evict a critical pod\" pod=\"kube-system/kube-scheduler-cloud-06\"\n\
Jan 09 03:44:07 cloud-06 kubelet[1231]: E0109 03:44:07.063012    1231 eviction_manager.go:560] \"Eviction manager: cannot evict a critical pod\" pod=\"kube-system/kube-apiserver-cloud-06\"\n\
Jan 09 03:44:07 cloud-06 kubelet[1231]: E0109 03:44:07.063033    1231 eviction_manager.go:560] \"Eviction manager: cannot evict a critical pod\" pod=\"kube-system/coredns-64897985d-jbn98\"\n\
Jan 09 03:44:07 cloud-06 kubelet[1231]: E0109 03:44:07.063054    1231 eviction_manager.go:560] \"Eviction manager: cannot evict a critical pod\" pod=\"kube-system/calico-node-nfznn\"\n\
Jan 09 03:44:07 cloud-06 kubelet[1231]: E0109 03:44:07.063071    1231 eviction_manager.go:560] \"Eviction manager: cannot evict a critical pod\" pod=\"kube-system/kube-proxy-sv9sd\"\n\
Jan 09 03:44:07 cloud-06 kubelet[1231]: I0109 03:44:07.063095    1231 eviction_manager.go:390] \"Eviction manager: unable to evict any pods from the node\"\n\
Jan 09 03:44:17 cloud-06 kubelet[1231]: I0109 03:44:17.111819    1231 eviction_manager.go:338] \"Eviction manager: attempting to reclaim\" resourceName=\"ephemeral-storage\"",
    "cpuPercents": generateRandomPercent(),
    "ramPercents": generateRandomPercent(),
    "diskPercents": generateRandomPercent(),
    "networkPercents": generateRandomPercent(),
};
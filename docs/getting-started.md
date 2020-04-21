# Tencent Cloud Controller Manager User Guide

The CCM provides a route-based Pod network and a LoadBalancer for each Service with a LoadBalancer type.

## Prerequisites

* A K8s cluster with version 1.10, 1.12, 1.14, 1.16 on a VPC network
* Each node should have a node name same as its IP address, or CCM can't initialize them. Using `--hostname-override` flag of `kubelet` is recommended.
* A route table is required for the Pod network. We build `route-ctl` to support the route table creation. See [route-ctl](https://github.com/TencentCloud/tencentcloud-cloud-controller-manager/tree/master/route-ctl).

## Tencent Cloud CCM Installation

**WARNING**: New workloads won't be scheduled until some nodes are initialized by CCM. Workloads already scheduled won't be affected.

### Update the control plane configuration

Clear the flag `--cloud-provider` of `kube-apiserver` and `kube-controller-manager`. If the cluster never uses a in-tree cloud provider, this flag should be empty. You can found Pod manifests in `/etc/kubernetes/manifests` on each master in a regular K8s cluster.

The sample manifests are [here](https://raw.githubusercontent.com/TencentCloud/tencentcloud-cloud-controller-manager/master/docs/example-manifests/out-of-tree/kube-apiserver.yaml).

More details are in [the official documents](https://kubernetes.io/docs/tasks/administer-cluster/running-cloud-controller/#running-cloud-controller-manager).

### Update the kubelet configuration

1. Set flag `--cloud-provider` of `kubelet` to `external`. If you run kubelet through `systemd`, you can edit the unit file `/etc/systemd/system/kubelet.service`, or `/etc/systemd/system/kubelet.service.d/10-kubeadm.conf` for `kubeadm`.
2. Set flag `--node-status-update-frequency` to `30s` to increase the kubelet status report frequency. A smaller frequency may lead to update failure of node status. 

### Deploy Tencent Cloud CCM

#### 1.You need to determine the following parameters and replace all the placeholders in the manifest before deploying. 

The placeholders below are parts of a Secret. All the values should be encoded via `base64`.
You can run the following command to get the encrypted context. Notice that the `-n` is required.

```shell script
echo -n "<Plain Text>" | base64
```

| Parameter Placeholder | Description | Value |
| ---- | ---- | ---- |
| `<REGION>` | The region your CVMs assisted | All region IDs(with a prefix `ap-`) could be found in section `Region List` of the [document](https://intl.cloud.tencent.com/document/api/213/31574) |
| `<SECRET_ID> & <SECRET_KEY>` | Identity to access the Tencent Cloud API | Following the [document](https://intl.cloud.tencent.com/document/product/598/34228) |
| `<CLUSTER_NETWORK_ROUTE_TABLE_NAME>` | Route table name of the Pod network | The route table must be created via the utility `route-ctl`. See [route-ctl](https://github.com/TencentCloud/tencentcloud-cloud-controller-manager/tree/master/route-ctl) |
| `<TENCENTCLOUD_CLOUD_CONTROLLER_MANAGER_VPC_ID>` | ID of the current VPC Network | It can be found on [TencentCloud VPC Console](https://console.cloud.tencent.com/vpc/vpc) , usually has a prefix `vpc-`. |


| Flag | Description | Value |
| ---- | ---- | ---- |
| `--cloud-provider` | Identify of the current CCM. | Must be `tencentcloud` |
| `--allocate-node-cidrs` | Approve CCM to allocate Pod CIDR for each node. | Should be enabled if a Pod network is desired. |
| `--configure-cloud-routes` | Allow CCM to create routes for Pod traffic among nodes. | Should be enabled if a Pod network is desired. |
| `--cluster-cidr` | Cluster CIDR(a.k.a Pod CIDR). | Should be same as the CIDR associated with <CLUSTER_NETWORK_ROUTE_TABLE_NAME>. e.g. `192.168.0.0/20` |


#### 2. Choose a container network plugin

If `--configure-cloud-routes` of CCM is enabled, the `kubernet` plugin is recommended to handle traffic among Pods along with the VPC network.
You can follow the [document](https://kubernetes.io/docs/concepts/extend-kubernetes/compute-storage-net/network-plugins/#kubenet) to set it up.
In addition, you need to add some `iptables rules` to accept traffic forwarded between `cbr0` and `eth0`.


#### 3. Install CCM

For cluster which version belows v1.16, run the command below to install CCM.
```shell script
kubectl apply -f https://raw.githubusercontent.com/TencentCloud/tencentcloud-cloud-controller-manager/master/docs/example-manifests/out-of-tree/cloud-controller-manager.yaml
```

For cluster which version is v1.16 or above, use the following manifest instead.
```shell script
kubectl apply -f https://raw.githubusercontent.com/TencentCloud/tencentcloud-cloud-controller-manager/master/docs/example-manifests/out-of-tree/cloud-controller-manager-v1.16.yaml
```

### Verify the installation

1. Wait until all nodes ready. It may take a few minutes.
2. Deploy [the sample of Service](https://github.com/TencentCloud/tencentcloud-cloud-controller-manager/blob/master/docs/resources/service/README.md).

# Tencent Cloud Controller Manager User Guide

## Prerequisites

* A K8s cluster with version 1.10, 1.12, or 1.14 using VPC network
* Each node should have a node name same as its IP address, or CCM can't initialize them. Using `--hostname-override` flag of `kubelet` is recommended.

## Tencent Cloud CCM Installation

**WARNING**: New workloads won't be scheduled until some nodes are initialized by CCM. Workloads already scheduled won't be affected.

### Update the control plane configuration

Clear the flag `--cloud-provider` of `kube-apiserver` and `kube-controller-manager`. If the cluster never uses a in-tree cloud provider, this flag should be empty. You can found Pod manifests in `/etc/kubernetes/manifests` on each master in a regular K8s cluster.

The sample manifests are [here](https://raw.githubusercontent.com/TencentCloud/tencentcloud-cloud-controller-manager/master/docs/example-manifests/out-of-tree/kube-apiserver.yaml).

More details are in [the official documents](https://kubernetes.io/docs/tasks/administer-cluster/running-cloud-controller/#running-cloud-controller-manager)

### Update the kubelet configuration

1. Set flag `--cloud-provider` of `kubelet` to `external`. If you run kubelet through `systemd`, you can edit the unit file `/etc/systemd/system/kubelet.service`, or `/etc/systemd/system/kubelet.service.d/10-kubeadm.conf` for `kubeadm`.
2. Set flag `--node-status-update-frequency` to `30s` to increase the kubelet status report frequency. A smaller frequency may lead to update failure of node status. 

### Deploy Tencent Cloud CCM

The following parameters should be determined before deploying. These placeholders are part of a Secret. All the values should be encoded via `base64`.

| Parameter Placeholder | Description | Value |
| ---- | ---- | ---- |
| <REGION> | The region your CVMs assisted | All region IDs(with a prefix `ap-`) could be found in section `Region List` of the [document](https://intl.cloud.tencent.com/document/api/213/31574) |
| <SECRET_ID> & <SECRET_KEY> | Identity to access the Tencent Cloud API | Following the [document](https://intl.cloud.tencent.com/document/product/598/34228) |
| <CLUSTER_NETWORK_ROUTE_TABLE_NAME> | ID of the route table of Pod network | It can be found on [TencentCloud Route Console](https://console.cloud.tencent.com/vpc/route) , usually has a prefix `rtb-`. |
| <TENCENTCLOUD_CLOUD_CONTROLLER_MANAGER_VPC_ID> | ID of the current VPC Network | It can be found on [TencentCloud VPC Console](https://console.cloud.tencent.com/vpc/vpc) , usually has a prefix `vpc-`. |


| Flag | Description | Value |
| ---- | ---- | ---- |
| --cloud-provider | Identify of the current CCM. The value is fixed. | `tencentcloud` |
| --allocate-node-cidrs | Approve CCM to allocate Pod CIRD for each node. | Disabled by setting to `false` |
| --cluster-cidr | Cluster CIDR(a.k.a Pod CIDR). The subnet must be created before using. | e.g. `192.168.0.0/20` |
| --configure-cloud-routes | Allow CCM to create routes for Pods. | Disabled by setting to `false` |


To deploy CCM,
```shell script
kubectl apply -f https://raw.githubusercontent.com/TencentCloud/tencentcloud-cloud-controller-manager/master/docs/example-manifests/out-of-tree/cloud-controller-manager.yaml
```

### Verify the installation

1. Wait until all nodes ready.
2. Deploy [the sample of Service](https://github.com/TencentCloud/tencentcloud-cloud-controller-manager/blob/master/docs/resources/service/README.md).

# Kubernetes Cloud Controller Manager for Tencent Cloud

`tencentcloud-cloud-controller-manager` is the Kubernetes cloud controller manager implementation for Tencent Cloud Container Service. Read more about cloud controller managers [here](https://kubernetes.io/docs/tasks/administer-cluster/running-cloud-controller/).

**WARNING**: this project is a work in progress and may not be production ready. Currently, only Kubernetes 1.10.x is supported.

## Requirements

At the current state of Kubernetes, running cloud controller manager requires a few things. Please read through the requirements carefully as they are critical to running cloud controller manager on a Kubernetes cluster on Tencent Cloud.

### --cloud-provider=external
All `kubelet`s in your cluster **MUST** set the flag `--cloud-provider=external`. `kube-apiserver` and `kube-controller-manager` must **NOT** set the flag `--cloud-provider` which will default them to use no cloud provider natively.

**WARNING**: setting `--cloud-provider=external` will taint all nodes in a cluster with `node.cloudprovider.kubernetes.io/uninitialized`, it is the responsibility of cloud controller managers to untaint those nodes once it has finished initializing them. This means that most pods will be left unscheduable until the cloud controller manager is running.

In the future, `--cloud-provider=external` will be the default. Learn more about the future of cloud providers in Kubernetes [here](https://github.com/kubernetes/community/blob/master/contributors/design-proposals/cloud-provider/cloud-provider-refactoring.md).

### Kubernetes node names must match the instance private ipv4 ip
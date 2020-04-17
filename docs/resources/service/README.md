# Service

Tencent CCM could provide a load balancer（CLB）for each Service with `type: LoadBalancer`. 

**WARNING**:The `spec.sessionAffinity` is not supported.

## Annotations

You can also change the CLB configuration by annotating Service.
 
[Annotations](https://raw.githubusercontent.com/TencentCloud/tencentcloud-cloud-controller-manager/master/docs/resources/service/annotations.md)

## Sample

To deploy the sample,

```shell script
kubectl apply -f https://raw.githubusercontent.com/TencentCloud/tencentcloud-cloud-controller-manager/master/docs/resources/service/smaple.yaml
```

To verify the result,

```shell script
❯ kubectl --kubeconfig=tke.kubeconf get po
NAME                         READY   STATUS    RESTARTS   AGE
nginx-574b87c764-hgt6d       1/1     Running   0          80s

❯ kubectl --kubeconfig=tke.kubeconf get svc
NAME                TYPE           CLUSTER-IP     EXTERNAL-IP      PORT(S)          AGE
nginx               LoadBalancer   172.18.0.136   106.55.71.219    80:32158/TCP     103s
```

You will found that a public IP address is allocated for the Service nginx in the column `EXTERNAL-IP`.

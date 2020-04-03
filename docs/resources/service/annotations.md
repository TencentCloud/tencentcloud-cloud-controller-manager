# Annotations of Service

| Annotation | Description | Value |
| ---- | ---- | ---- |
| service.beta.kubernetes.io/tencentcloud-loadbalancer-kind | Type of the load balancer will be created for the Service   |  `classic` for L4 CLB, and  `application` for L7 CLB which natively supports HTTP(S) traffic  |
| service.beta.kubernetes.io/tencentcloud-loadbalancer-type | Scope where the load balancer could be accessed from | `public` when readable from Internet  or `private` from the same VPC |
| service.beta.kubernetes.io/tencentcloud-loadbalancer-type-internal-subnet-id | ID of the subnet in which the `private` CLB to be created | You can found it on Tencent VPC Console. It usually has a prefix `subnet-` |
| service.beta.kubernetes.io/tencentcloud-loadbalancer-name | Name of the CLB to be created |  |

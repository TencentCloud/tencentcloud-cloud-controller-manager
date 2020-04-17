# route-ctl

中文文档在[这里](https://github.com/TencentCloud/tencentcloud-cloud-controller-manager/blob/master/route-ctl/README_zhCN.md)

`route-ctl` is a utility to build a route table for the Pod network. The route table is created in the same VPC in which CVMs occupied. 
The route table is required if you would like to use Tencent Cloud Controller as a out-of-tree cloud provider in your Kubernetes cluster. 


## Build

Before getting into code, you must install the go toolchain. The [official document](https://golang.org/doc/install#install) is provided as a reference. 
Then, make sure that the environment variable `GOPATH` is set correctly in current terminal and,
```
mkdir -p "${GOPATH}/src/github.com/tencentcloud/"
git clone https://github.com/tencentcloud/tencentcloud-cloud-controller-manager.git "${GOPATH}/src/github.com/tencentcloud/tencentcloud-cloud-controller-manager"
cd ${GOPATH}/src/github.com/tencentcloud/tencentcloud-cloud-controller-manager/route-ctl
go build -v
```

## Usage

The authority identity should be provided before running any command. 
You can found `SecretID(QCloudSecretId)` and `SecretKey(QCloudSecretKey)` by following the [guide](https://intl.cloud.tencent.com/document/product/598/34228). 
The `QCloudCcsAPIRegion` specifies the region `route-ctl` will communicate with. An available region list is provided [here](https://intl.cloud.tencent.com/document/api/213/31574).

```
export QCloudSecretId=************************************
export QCloudSecretKey=********************************
export QCloudCcsAPIRegion=ap-shanghai
```

### Create route table
```
./route-ctl route-table create --route-table-cidr-block 10.10.0.0/16 --route-table-name route-table-test --vpc-id vpc-********
```

Notice that the CIDR specified by `--route-table-cidr-block` must not overlap any subnet CIDR or cluster CIDR in the same VPC. 
If some CIDR in the VPC is overlapped and `--ignore-cidr-conflict` is enabled, the routes to the overlapped subnet will also be overridden. It is **DANGEROUS**.

### List all route tables
```
./route-ctl route-table list
```

### Delete route table
```
./route-ctl route-table delete --route-table-name route-table-test
```

### Create a route
```
./route-ctl route create --destination-cidr-block 10.10.1.0/24 --route-table-name route-table-test --gateway-ip 192.168.1.4
```

### List all routes in the specific route table
```
./route-ctl route list --route-table-name route-table-test
```

### Delete a route
```
./route-ctl route delete --destination-cidr-block 10.10.1.0/24 --route-table-name route-table-test --gateway-ip 192.168.1.4
```

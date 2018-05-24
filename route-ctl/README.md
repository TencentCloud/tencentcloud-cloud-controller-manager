# route-ctl

`route-ctl` 是一个用于在腾讯云 vpc 环境创建路由表的小工具，典型的应用场景包括

* 在腾讯云从云服务器手工搭建 kubernetes 且集群网络使用路由方案时，可以使用此工具创建 kubernetes 用来建立 pod 网络的子网，使得集群中的 pod 通信能够和腾讯云 vpc 打通


## 编译

将此项目 clone 到 GOPATH 下，假设 GOPATH 为 /root/go

```
mkdir -p /root/go/src/github.com/dbdd4us/
git clone https://github.com/dbdd4us/tencentcloud-cloud-controller-manager.git /root/go/src/github.com/dbdd4us/tencentcloud-cloud-controller-manager
cd /root/go/src/github.com/dbdd4us/tencentcloud-cloud-controller-manager/route-ctl
go build -v
```

## 使用

**注意**:  routetable 的 cidr 不能与对应 vpc 的 cidr、其他集群的 cidr 以及已经存在的子网的 cidr 存在冲突

使用前需要通过环境变量设置 `QCloudSecretId`、 `QCloudSecretKey` 以及 `QCloudCcsAPIRegion`

```
export QCloudSecretId=************************************
export QCloudSecretKey=********************************
export QCloudCcsAPIRegion=ap-shanghai
```

### 创建路由表
```
./route-ctl create --route-table-cidr-block 10.10.0.0/16 --route-table-name route-table-test --vpc-id vpc-********
```

### 查看现存的路由表
```
./route-ctl list
```

### 删除指定路由表
```
./route-ctl delete --route-table-name route-table-test
```

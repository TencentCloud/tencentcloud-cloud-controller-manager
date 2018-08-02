# route-ctl

`route-ctl` 是一个用于在腾讯云 vpc 环境创建路由表的小工具，典型的应用场景包括

* 在腾讯云从云服务器手工搭建 kubernetes 且集群网络使用路由方案时，可以使用此工具创建 kubernetes 用来建立 pod 网络的子网，使得集群中的 pod 通信能够和腾讯云 vpc 打通


## 编译

将此项目 clone 到 GOPATH 下，假设 GOPATH 为 /root/go

```
mkdir -p /root/go/src/github.com/dbdd4us/
git clone https://github.com/tencentcloud/tencentcloud-cloud-controller-manager.git /root/go/src/github.com/tencentcloud/tencentcloud-cloud-controller-manager
cd /root/go/src/github.com/tencentcloud/tencentcloud-cloud-controller-manager/route-ctl
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
./route-ctl route-table create --route-table-cidr-block 10.10.0.0/16 --route-table-name route-table-test --vpc-id vpc-********
```

当通过 `route-ctl` 创建路由表时， `route-ctl` 会先检查所要创建的路由表的 `cidr` 是否和现存的网络设置冲突，具体的检查包括下面四项：

1. route table 所在 vpc 的 cidr
2. route table 所在 vpc 的子网路由的 cidr
3. route table 所在 vpc 的 ccs 集群的集群网络 cidr
4. route table 所在 vpc 的其它通过 route-ctl 创建的 route table 的 cidr

___Note:___ 当存在 `cidr` 冲突时，`route-ctl` 支持通过 `--ignore-cidr-conflict` 选项忽略冲突进行创建，需要注意的是，通过 route table 创建的具体路由条目在 vpc 内拥有最高的匹配优先级，忽略 `cidr` 冲突进行创建可能会导致现存的网络出现问题，___请谨慎使用___。

### 查看现存的路由表
```
./route-ctl route-table list
```

### 删除指定路由表
```
./route-ctl route-table delete --route-table-name route-table-test
```

### 创建路由
```
./route-ctl route create --destination-cidr-block 10.10.1.0/24 --route-table-name route-table-test --gateway-ip 192.168.1.4
```

### 查看现存的路由
```
./route-ctl route list --route-table-name route-table-test
```

### 删除指定路由
```
./route-ctl route delete --destination-cidr-block 10.10.1.0/24 --route-table-name route-table-test --gateway-ip 192.168.1.4
```

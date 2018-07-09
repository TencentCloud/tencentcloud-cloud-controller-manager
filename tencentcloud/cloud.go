package tencentcloud

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/dbdd4us/qcloudapi-sdk-go/cvm"
	"github.com/dbdd4us/qcloudapi-sdk-go/ccs"
	"github.com/dbdd4us/qcloudapi-sdk-go/clb"
	"github.com/dbdd4us/qcloudapi-sdk-go/common"

	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/cloudprovider"
	"k8s.io/kubernetes/pkg/controller"
)

const (
	providerName = "tencentcloud"
)

var (
	CloudInstanceNotFound = errors.New("tencentcloud instance not found")
)

func init() {
	cloudprovider.RegisterCloudProvider(providerName, NewCloud)
}

func NewCloud(config io.Reader) (cloudprovider.Interface, error) {
	var c Config
	if config != nil {
		cfg, err := ioutil.ReadAll(config)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(cfg, &c); err != nil {
			return nil, err
		}
	}

	if c.Region == "" {
		c.Region = os.Getenv("TENCENTCLOUD_CLOUD_CONTROLLER_MANAGER_REGION")
	}
	if c.VpcId == "" {
		c.VpcId = os.Getenv("TENCENTCLOUD_CLOUD_CONTROLLER_MANAGER_VPC_ID")
	}
	if c.SecretId == "" {
		c.SecretId = os.Getenv("TENCENTCLOUD_CLOUD_CONTROLLER_MANAGER_SECRET_ID")
	}
	if c.SecretKey == "" {
		c.SecretKey = os.Getenv("TENCENTCLOUD_CLOUD_CONTROLLER_MANAGER_SECRET_KEY")
	}

	if c.ClusterRouteTable == "" {
		c.ClusterRouteTable = os.Getenv("TENCENTCLOUD_CLOUD_CONTROLLER_MANAGER_CLUSTER_ROUTE_TABLE")
	}

	return &Cloud{config: c}, nil
}

type Cloud struct {
	config Config

	kubeClient kubernetes.Interface

	cvm   *cvm.Client
	cvmV3 *cvm.Client
	ccs   *ccs.Client
	clb   *clb.Client
}

type Config struct {
	Region string `json:"region"`
	VpcId  string `json:"vpc_id"`

	SecretId  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`

	ClusterRouteTable string `json:"cluster_route_table"`
}

// Initialize provides the cloud with a kubernetes client builder and may spawn goroutines
// to perform housekeeping activities within the cloud provider.
func (cloud *Cloud) Initialize(clientBuilder controller.ControllerClientBuilder) {
	cloud.kubeClient = clientBuilder.ClientOrDie("tencentcloud-cloud-provider")
	cvmClient, err := cvm.NewClient(
		common.Credential{SecretId: cloud.config.SecretId, SecretKey: cloud.config.SecretKey},
		common.Opts{Region: cloud.config.Region},
	)
	if err != nil {
		panic(err)
	}
	cloud.cvm = cvmClient
	cvmV3Client, err := cvm.NewClient(
		common.Credential{SecretId: cloud.config.SecretId, SecretKey: cloud.config.SecretKey},
		common.Opts{Region: cloud.config.Region, Host: cvm.CvmV3Host, Path: cvm.CvmV3Path},
	)
	if err != nil {
		panic(err)
	}
	cloud.cvmV3 = cvmV3Client
	ccsClient, err := ccs.NewClient(
		common.Credential{SecretId: cloud.config.SecretId, SecretKey: cloud.config.SecretKey},
		common.Opts{Region: cloud.config.Region},
	)
	if err != nil {
		panic(err)
	}
	cloud.ccs = ccsClient
	clbClient, err := clb.NewClient(
		common.Credential{SecretId: cloud.config.SecretId, SecretKey: cloud.config.SecretKey},
		common.Opts{Region: cloud.config.Region},
	)
	if err != nil {
		panic(err)
	}
	cloud.clb = clbClient
	return
}

// LoadBalancer returns a balancer interface. Also returns true if the interface is supported, false otherwise.
func (cloud *Cloud) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	return cloud, true
}

// Instances returns an instances interface. Also returns true if the interface is supported, false otherwise.
func (cloud *Cloud) Instances() (cloudprovider.Instances, bool) {
	return cloud, true
}

// Zones returns a zones interface. Also returns true if the interface is supported, false otherwise.
func (cloud *Cloud) Zones() (cloudprovider.Zones, bool) {
	return nil, false
}

// Clusters returns a clusters interface.  Also returns true if the interface is supported, false otherwise.
func (cloud *Cloud) Clusters() (cloudprovider.Clusters, bool) {
	return nil, false
}

// Routes returns a routes interface along with whether the interface is supported.
func (cloud *Cloud) Routes() (cloudprovider.Routes, bool) {
	return cloud, true
}

// ProviderName returns the cloud provider ID.
func (cloud *Cloud) ProviderName() string {
	return providerName
}

// HasClusterID returns true if a ClusterID is required and set
func (cloud *Cloud) HasClusterID() bool {
	return false
}

package tencentcloud

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/dbdd4us/qcloudapi-sdk-go/metadata"
	"github.com/dbdd4us/qcloudapi-sdk-go/cvm"
	"github.com/dbdd4us/qcloudapi-sdk-go/ccs"
	"github.com/dbdd4us/qcloudapi-sdk-go/common"

	"k8s.io/client-go/kubernetes"
	"k8s.io/kubernetes/pkg/cloudprovider"
	"k8s.io/kubernetes/pkg/controller"
)

const (
	providerName = "qcloud"
)

func init() {
	cloudprovider.RegisterCloudProvider(providerName, NewCloud)
}

func NewCloud(config io.Reader) (cloudprovider.Interface, error) {
	cfg, err := ioutil.ReadAll(config)
	if err != nil {
		return nil, err
	}

	var c Config
	if err := json.Unmarshal(cfg, &c); err != nil {
		return nil, err
	}

	return &Cloud{config: c}, nil
}

type Cloud struct {
	config Config

	kubeClient kubernetes.Interface

	metadata *metadata.MetaData
	cvm      *cvm.Client
	ccs      *ccs.Client
}

type Config struct {
	Region string `json:"region"`
	Zone   string `json:"zone"`

	VpcId     string `json:"vpcId"`
	SecretId  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`

	ClusterRouteTable string `json:"cluster_route_table"`
}

// Initialize provides the cloud with a kubernetes client builder and may spawn goroutines
// to perform housekeeping activities within the cloud provider.
func (cloud *Cloud) Initialize(clientBuilder controller.ControllerClientBuilder) {
	cloud.kubeClient = clientBuilder.ClientOrDie("tencentcloud-cloud-provider")
	cloud.metadata = metadata.NewMetaData(http.DefaultClient)
	cvmClient, err := cvm.NewClient(
		common.Credential{SecretId: cloud.config.SecretId, SecretKey: cloud.config.SecretKey},
		common.Opts{Region: cloud.config.Region},
	)
	if err != nil {
		panic(err)
	}
	cloud.cvm = cvmClient
	ccsClient, err := ccs.NewClient(
		common.Credential{SecretId: cloud.config.SecretId, SecretKey: cloud.config.SecretKey},
		common.Opts{Region: cloud.config.Region},
	)
	if err != nil {
		panic(err)
	}
	cloud.ccs = ccsClient
	return
}

// LoadBalancer returns a balancer interface. Also returns true if the interface is supported, false otherwise.
func (cloud *Cloud) LoadBalancer() (cloudprovider.LoadBalancer, bool) {
	return nil, false
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

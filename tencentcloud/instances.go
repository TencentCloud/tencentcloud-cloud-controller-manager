package tencentcloud

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/dbdd4us/qcloudapi-sdk-go/cvm"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/cloudprovider"
)

// NodeAddresses returns the addresses of the specified instance.
// TODO(roberthbailey): This currently is only used in such a way that it
// returns the address of the calling instance. We should do a rename to
// make this clearer.
func (cloud *Cloud) NodeAddresses(ctx context.Context, name types.NodeName) ([]v1.NodeAddress, error) {
	node, err := cloud.getInstanceByInstancePrivateIp(string(name))
	if err != nil {
		return []v1.NodeAddress{}, err
	}
	addresses := make([]v1.NodeAddress, len(node.PrivateIPAddresses)+len(node.PublicIPAddresses))
	for idx, ip := range node.PrivateIPAddresses {
		addresses[idx] = v1.NodeAddress{Type: v1.NodeInternalIP, Address: ip}
	}
	for idx, ip := range node.PublicIPAddresses {
		addresses[len(node.PrivateIPAddresses)+idx] = v1.NodeAddress{Type: v1.NodeExternalIP, Address: ip}
	}
	return addresses, nil
}

// NodeAddressesByProviderID returns the addresses of the specified instance.
// The instance is specified using the providerID of the node. The
// ProviderID is a unique identifier of the node. This will not be called
// from the node whose nodeaddresses are being queried. i.e. local metadata
// services cannot be used in this method to obtain nodeaddresses
func (cloud *Cloud) NodeAddressesByProviderID(ctx context.Context, providerID string) ([]v1.NodeAddress, error) {
	id := strings.TrimPrefix(providerID, fmt.Sprintf("%s://", providerName))
	parts := strings.Split(id, "/")
	if len(parts) == 3 {
		instance, err := cloud.getInstanceByInstanceID(parts[2])
		if err != nil {
			return []v1.NodeAddress{}, err
		}
		addresses := make([]v1.NodeAddress, len(instance.PrivateIPAddresses)+len(instance.PublicIPAddresses))
		for idx, ip := range instance.PrivateIPAddresses {
			addresses[idx] = v1.NodeAddress{Type: v1.NodeInternalIP, Address: ip}
		}
		for idx, ip := range instance.PublicIPAddresses {
			addresses[len(instance.PrivateIPAddresses)+idx] = v1.NodeAddress{Type: v1.NodeExternalIP, Address: ip}
		}
		return addresses, nil
	}
	return []v1.NodeAddress{}, errors.New(fmt.Sprintf("invalid format for providerId %s", providerID))
}

// ExternalID returns the cloud provider ID of the node with the specified NodeName.
// Note that if the instance does not exist or is no longer running, we must return ("", cloudprovider.InstanceNotFound)
func (cloud *Cloud) ExternalID(ctx context.Context, nodeName types.NodeName) (string, error) {
	node, err := cloud.getInstanceByInstancePrivateIp(string(nodeName))
	if err != nil {
		return "", err
	}

	return node.InstanceID, nil
}

// InstanceID returns the cloud provider ID of the node with the specified NodeName.
func (cloud *Cloud) InstanceID(ctx context.Context, nodeName types.NodeName) (string, error) {
	node, err := cloud.getInstanceByInstancePrivateIp(string(nodeName))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("/%s/%s", node.Placement.Zone, node.InstanceID), nil
}

// InstanceType returns the type of the specified instance.
func (cloud *Cloud) InstanceType(ctx context.Context, name types.NodeName) (string, error) {
	return providerName, nil
}

// InstanceTypeByProviderID returns the type of the specified instance.
func (cloud *Cloud) InstanceTypeByProviderID(ctx context.Context, providerID string) (string, error) {
	return providerName, nil
}

// AddSSHKeyToAllInstances adds an SSH public key as a legal identity for all instances
// expected format for the key is standard ssh-keygen format: <protocol> <blob>
func (cloud *Cloud) AddSSHKeyToAllInstances(ctx context.Context, user string, keyData []byte) error {
	return cloudprovider.NotImplemented
}

// CurrentNodeName returns the name of the node we are currently running on
// On most clouds (e.g. GCE) this is the hostname, so we provide the hostname
func (cloud *Cloud) CurrentNodeName(ctx context.Context, hostname string) (types.NodeName, error) {
	return types.NodeName(""), cloudprovider.NotImplemented
}

// InstanceExistsByProviderID returns true if the instance for the given provider id still is running.
// If false is returned with no error, the instance will be immediately deleted by the cloud controller manager.
func (cloud *Cloud) InstanceExistsByProviderID(ctx context.Context, providerID string) (bool, error) {
	return true, nil
}

func (cloud *Cloud) getInstanceByInstancePrivateIp(privateIp string) (*cvm.InstanceInfo, error) {
	instances, err := cloud.cvm.DescribeInstances(&cvm.DescribeInstancesArgs{
		Version: cvm.DefaultVersion,
		Filters: &[]cvm.Filter{cvm.NewFilter(cvm.FilterNamePrivateIpAddress, privateIp)},
	})
	if err != nil {
		return nil, err
	}
	for _, instance := range instances.InstanceSet {
		if instance.VirtualPrivateCloud.VpcID != cloud.config.VpcId {
			continue
		}
		for _, ip := range instance.PrivateIPAddresses {
			if ip == privateIp {
				return &instance, nil
			}
		}
	}
	return nil, CloudInstanceNotFound
}

func (cloud *Cloud) getInstanceByInstanceID(instanceID string) (*cvm.InstanceInfo, error) {
	instances, err := cloud.cvm.DescribeInstances(&cvm.DescribeInstancesArgs{
		Version: cvm.DefaultVersion,
		Filters: &[]cvm.Filter{cvm.NewFilter(cvm.FilterNameInstanceId, instanceID)},
	})
	if err != nil {
		return nil, err
	}
	for _, instance := range instances.InstanceSet {
		if instance.VirtualPrivateCloud.VpcID != cloud.config.VpcId {
			continue
		}
		if instance.InstanceID == instanceID {
			return &instance, nil
		}
	}
	return nil, CloudInstanceNotFound
}

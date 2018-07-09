package tencentcloud

import (
	"context"
	"errors"
	"k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/cloudprovider"

	"github.com/dbdd4us/qcloudapi-sdk-go/clb"
	"github.com/dbdd4us/qcloudapi-sdk-go/cvm"
)

const (
	// classic clb or application clb
	ServiceAnnotationLoadBalancerKind = "service.beta.kubernetes.io/tencentcloud-loadbalancer-kind"
	LoadBalancerKindClassic           = "classic"
	LoadBalancerKindApplication       = "application"

	// public network based clb or private network based clb
	ServiceAnnotationLoadBalancerType = "service.beta.kubernetes.io/tencentcloud-loadbalancer-type"
	LoadBalancerTypePublic            = "public"
	LoadBalancerTypePrivate           = "private"

	// subnet id for private network based clb
	ServiceAnnotationLoadBalancerTypeInternalSubnetId = "service.beta.kubernetes.io/tencentcloud-loadbalancer-type-internal-subnet-id"
)

var (
	ErrCloudLoadBalancerNotFound = errors.New("LoadBalancer not found")

	ClbLoadBalancerTypePublic  = 2
	ClbLoadBalancerTypePrivate = 3

	ClbLoadBalancerKindClassic     = 0
	ClbLoadBalancerKindApplication = 1

	ClbLoadBalancerListenerProtocolHTTP  = 1
	ClbLoadBalancerListenerProtocolHTTPS = 4
	ClbLoadBalancerListenerProtocolTCP   = 2
	ClbLoadBalancerListenerProtocolUDP   = 3
)

func (cloud *Cloud) GetLoadBalancer(ctx context.Context, clusterName string, service *v1.Service) (status *v1.LoadBalancerStatus, exists bool, err error) {
	loadBalancerName := cloudprovider.GetLoadBalancerName(service)

	loadBalancer, err := cloud.getLoadBalancerByName(loadBalancerName)
	if err != nil {
		if err == ErrCloudLoadBalancerNotFound {
			return nil, false, nil
		}
		return nil, false, err
	}

	ingresses := make([]v1.LoadBalancerIngress, len(loadBalancer.LoadBalancerVips))

	for i, vip := range loadBalancer.LoadBalancerVips {
		ingresses[i] = v1.LoadBalancerIngress{IP: vip}
	}

	return &v1.LoadBalancerStatus{
		Ingress: ingresses,
	}, true, nil
}

func (cloud *Cloud) EnsureLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) (*v1.LoadBalancerStatus, error) {
	if service.Spec.SessionAffinity != v1.ServiceAffinityNone {
		return nil, errors.New("SessionAffinity is not supported currently")
	}

	// TODO check if kubernetes has already do validate

	// 1. ensure loadbalancer created
	err := cloud.ensureLoadBalancerInstance(ctx, clusterName, service)
	if err != nil {
		return nil, err
	}
	// 2. ensure loadbalancer listener created
	err = cloud.ensureLoadBalancerListeners(ctx, clusterName, service)
	if err != nil {
		return nil, err
	}
	// 3. ensure right hosts is bounded to loadbalancer
	err = cloud.ensureLoadBalancerBackends(ctx, clusterName, service, nodes)
	if err != nil {
		return nil, err
	}

	loadBalancer, err := cloud.getLoadBalancerByName(cloudprovider.GetLoadBalancerName(service))
	if err != nil {
		return nil, err
	}

	ingresses := make([]v1.LoadBalancerIngress, len(loadBalancer.LoadBalancerVips))

	for i, vip := range loadBalancer.LoadBalancerVips {
		ingresses[i] = v1.LoadBalancerIngress{IP: vip}
	}

	return &v1.LoadBalancerStatus{
		Ingress: ingresses,
	}, nil
}

func (cloud *Cloud) UpdateLoadBalancer(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) error {
	return cloud.ensureLoadBalancerBackends(ctx, clusterName, service, nodes)
}

func (cloud *Cloud) EnsureLoadBalancerDeleted(ctx context.Context, clusterName string, service *v1.Service) error {
	_, err := cloud.getLoadBalancerByName(cloudprovider.GetLoadBalancerName(service))
	if err != nil {
		if err == ErrCloudLoadBalancerNotFound {
			return nil
		}
	}

	return cloud.deleteLoadBalancer(ctx, clusterName, service)
}

func (cloud *Cloud) getLoadBalancerByName(name string) (*clb.LoadBalancer, error) {
	// we don't need to check loadbalancer kind here because ensureLoadBalancerInstance will ensure the kind is right
	forward := -1
	response, err := cloud.clb.DescribeLoadBalancers(&clb.DescribeLoadBalancersArgs{
		LoadBalancerName: &name,
		Forward:          &forward,
	})
	if err != nil {
		return nil, err
	}
	if len(response.LoadBalancerSet) < 1 {
		return nil, ErrCloudLoadBalancerNotFound
	}

	return &response.LoadBalancerSet[0], nil
}

func (cloud *Cloud) ensureLoadBalancerInstance(ctx context.Context, clusterName string, service *v1.Service) error {
	loadBalancerName := cloudprovider.GetLoadBalancerName(service)

	loadBalancer, err := cloud.getLoadBalancerByName(loadBalancerName)
	if err != nil {
		if err != ErrCloudLoadBalancerNotFound {
			return err
		}
		loadBalancer, err = cloud.createLoadBalancer(ctx, clusterName, service)
		if err != nil {
			return err
		}
	}

	loadBalancerDesiredKind, ok := service.Annotations[ServiceAnnotationLoadBalancerKind]
	if !ok || (loadBalancerDesiredKind != LoadBalancerKindClassic && loadBalancerDesiredKind != LoadBalancerKindApplication) {
		loadBalancerDesiredKind = LoadBalancerKindApplication
	}
	loadBalancerDesiredType, ok := service.Annotations[ServiceAnnotationLoadBalancerType]
	if !ok || (loadBalancerDesiredType != LoadBalancerTypePrivate && loadBalancerDesiredType != LoadBalancerTypePublic) {
		loadBalancerDesiredType = LoadBalancerTypePublic
	}

	// don't check subnet id because clb could bound to instance in differnet subnet
	//var loadBalancerDesiredSubnetId string
	//if loadBalancerDesiredType == LoadBalancerTypePrivate {
	//	loadBalancerDesiredSubnetId, ok = service.Annotations[ServiceAnnotationLoadBalancerTypeInternalSubnetId]
	//	if !ok {
	//   return errors.New("Can not find subnetId")
	// }
	//}

	needRecreate := false
	switch {
	case loadBalancerDesiredType == LoadBalancerTypePublic && loadBalancerDesiredKind == LoadBalancerKindApplication:
		if !(loadBalancer.LoadBalancerType == ClbLoadBalancerTypePublic && loadBalancer.Forward == ClbLoadBalancerKindApplication && loadBalancer.UniqVpcId == cloud.config.VpcId) {
			needRecreate = true
		}
	case loadBalancerDesiredType == LoadBalancerTypePublic && loadBalancerDesiredKind == LoadBalancerKindClassic:
		if !(loadBalancer.LoadBalancerType == ClbLoadBalancerTypePublic && loadBalancer.Forward == ClbLoadBalancerKindClassic && loadBalancer.UniqVpcId == cloud.config.VpcId) {
			needRecreate = true
		}
	case loadBalancerDesiredType == LoadBalancerTypePrivate && loadBalancerDesiredKind == LoadBalancerKindApplication:
		if !(loadBalancer.LoadBalancerType == ClbLoadBalancerTypePrivate && loadBalancer.Forward == ClbLoadBalancerKindApplication && loadBalancer.UniqVpcId == cloud.config.VpcId) {
			needRecreate = true
		}
	case loadBalancerDesiredType == LoadBalancerTypePrivate && loadBalancerDesiredKind == LoadBalancerKindClassic:
		if !(loadBalancer.LoadBalancerType == ClbLoadBalancerTypePrivate && loadBalancer.Forward == ClbLoadBalancerKindClassic && loadBalancer.UniqVpcId == cloud.config.VpcId) {
			needRecreate = true
		}
	default:
		needRecreate = true
	}

	if needRecreate {
		if err := cloud.deleteLoadBalancer(ctx, clusterName, service); err != nil {
			return err
		}
		if loadBalancer, err = cloud.createLoadBalancer(ctx, clusterName, service); err != nil {
			return err
		}
	}

	return nil
}

func (cloud *Cloud) ensureLoadBalancerListeners(ctx context.Context, clusterName string, service *v1.Service) error {
	loadBalancerName := cloudprovider.GetLoadBalancerName(service)

	loadBalancer, err := cloud.getLoadBalancerByName(loadBalancerName)
	if err != nil {
		return err
	}

	switch loadBalancer.Forward {
	case ClbLoadBalancerKindClassic:
		return cloud.ensureClassicLoadBalancerListeners(ctx, clusterName, service, loadBalancer)
	case ClbLoadBalancerKindApplication:
		return cloud.ensureApplicationLoadBalancerListeners(ctx, clusterName, service, loadBalancer)
	default:
		return errors.New("Unsupported loadbalancer kind")
	}
}

func (cloud *Cloud) ensureClassicLoadBalancerListeners(ctx context.Context, clusterName string, service *v1.Service, loadBalancer *clb.LoadBalancer) error {
	response, err := cloud.clb.DescribeLoadBalancerListeners(&clb.DescribeLoadBalancerListenersArgs{
		LoadBalancerId: loadBalancer.LoadBalancerId,
	})
	if err != nil {
		return err
	}

	loadBalancerListeners := response.ListenerSet

	usedListenerIds := []string{}

	createdServicePortNames := []string{}

	findOneListenerValid := func(port v1.ServicePort) (listenerId string) {
		listenerId = ""

		for _, listener := range loadBalancerListeners {
			if listener.LoadBalancerPort == port.Port && listener.InstancePort == port.NodePort && cloud.mapClbProtoToServicePortProto(listener.Protocol) == port.Protocol {
				return listener.UnListenerId
			}
		}

		return
	}

	for _, port := range service.Spec.Ports {
		listenerId := findOneListenerValid(port)
		if listenerId != "" {
			createdServicePortNames = append(createdServicePortNames, port.Name)
			usedListenerIds = append(usedListenerIds, listenerId)
		}
	}

	listenersToCreate := []clb.CreateListenerOpts{}

	for _, port := range service.Spec.Ports {
		ensured := false

		for _, portName := range createdServicePortNames {
			if port.Name == portName {
				ensured = true
				break
			}
		}

		if !ensured {
			listenersToCreate = append(listenersToCreate, clb.CreateListenerOpts{
				LoadBalancerPort: port.Port,
				InstancePort:     port.NodePort,
				Protocol:         cloud.mapServicePortProtoClbProto(port.Protocol),
				ListenerName:     &port.Name,
			})
		}
	}

	listenersToDelete := []string{}

	for _, listener := range loadBalancerListeners {
		used := false

		for _, usedListenerId := range usedListenerIds {
			if listener.UnListenerId == usedListenerId {
				used = true
			}
		}

		if !used {
			listenersToDelete = append(listenersToDelete, listener.UnListenerId)
		}
	}
	if len(listenersToDelete) > 0 {
		result, err := clb.WaitUntilDone(
			func() (clb.AsyncTask, error) {
				return cloud.clb.DeleteLoadBalancerListeners(
					loadBalancer.LoadBalancerId,
					listenersToDelete,
				)
			},
			cloud.clb,
		)
		if err != nil {
			return err
		}
		if result != clb.TaskSuccceed {
			return errors.New("task is not succeed")
		}
	}

	if len(listenersToCreate) > 0 {
		result, err := clb.WaitUntilDone(
			func() (clb.AsyncTask, error) {
				return cloud.clb.CreateLoadBalancerListeners(&clb.CreateLoadBalancerListenersArgs{
					LoadBalancerId: loadBalancer.LoadBalancerId,
					Listeners:      listenersToCreate,
				})
			},
			cloud.clb,
		)
		if err != nil {
			return err
		}
		if result != clb.TaskSuccceed {
			return errors.New("task is not succeed")
		}

	}

	return nil
}

func (cloud *Cloud) ensureApplicationLoadBalancerListeners(ctx context.Context, clusterName string, service *v1.Service, loadBalancer *clb.LoadBalancer) error {
	response, err := cloud.clb.DescribeForwardLBListeners(&clb.DescribeForwardLBListenersArgs{
		LoadBalancerId: loadBalancer.LoadBalancerId,
	})
	if err != nil {
		return err
	}

	loadBalancerListeners := response.ListenerSet

	usedListenerIds := make([]string, 0)

	createdServicePortNames := make([]string, 0)

	findOneListenerValid := func(port v1.ServicePort) (listenerId string) {
		listenerId = ""

		for _, listener := range loadBalancerListeners {
			if listener.LoadBalancerPort == int(port.Port) && cloud.mapClbProtoToServicePortProto(listener.Protocol) == port.Protocol {
				return listener.ListenerId
			}
		}

		return
	}

	for _, port := range service.Spec.Ports {
		listenerId := findOneListenerValid(port)
		if listenerId != "" {
			// TODO check if port name is unique
			createdServicePortNames = append(createdServicePortNames, port.Name)
			usedListenerIds = append(usedListenerIds, listenerId)
		}
	}

	listenersToCreate := make([]clb.CreateFourthLayerListenerOpts, 0)

	for _, port := range service.Spec.Ports {
		ensured := false

		for _, portName := range createdServicePortNames {
			if port.Name == portName {
				ensured = true
				break
			}
		}

		if !ensured {
			listenersToCreate = append(listenersToCreate, clb.CreateFourthLayerListenerOpts{
				LoadBalancerPort: int(port.Port),
				Protocol:         cloud.mapServicePortProtoClbProto(port.Protocol),
			})
		}
	}

	listenersToDelete := make([]string, 0)

	for _, listener := range loadBalancerListeners {
		used := false

		for _, usedListenerId := range usedListenerIds {
			if listener.ListenerId == usedListenerId {
				used = true
			}
		}

		if !used {
			listenersToDelete = append(listenersToDelete, listener.ListenerId)
		}
	}

	for _, unusedListener := range listenersToDelete {
		result, err := clb.WaitUntilDone(
			func() (clb.AsyncTask, error) {
				return cloud.clb.DeleteForwardLBListener(&clb.DeleteForwardLBListenerArgs{
					LoadBalancerId: loadBalancer.LoadBalancerId,
					ListenerId:     unusedListener,
				})
			},
			cloud.clb,
		)
		if err != nil {
			return err
		}
		if result != clb.TaskSuccceed {
			return errors.New("task is not succeed")
		}
	}

	if len(listenersToCreate) > 0 {
		result, err := clb.WaitUntilDone(
			func() (clb.AsyncTask, error) {
				return cloud.clb.CreateForwardLBFourthLayerListeners(&clb.CreateForwardLBFourthLayerListenersArgs{
					LoadBalancerId: loadBalancer.LoadBalancerId,
					Listeners:      listenersToCreate,
				})
			},
			cloud.clb,
		)
		if err != nil {
			return err
		}
		if result != clb.TaskSuccceed {
			return errors.New("task is not succeed")
		}
	}

	return nil
}

func (cloud *Cloud) ensureLoadBalancerBackends(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node) error {
	loadBalancerName := cloudprovider.GetLoadBalancerName(service)

	loadBalancer, err := cloud.getLoadBalancerByName(loadBalancerName)
	if err != nil {
		return err
	}

	switch loadBalancer.Forward {
	case ClbLoadBalancerKindClassic:
		return cloud.ensureClassicLoadBalancerBackends(ctx, clusterName, service, nodes, loadBalancer)
	case ClbLoadBalancerKindApplication:
		return cloud.ensureApplicationLoadBalancerBackends(ctx, clusterName, service, nodes, loadBalancer)
	default:
		return errors.New("task is not succeed")
	}
}

func (cloud *Cloud) ensureClassicLoadBalancerBackends(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node, loadBalancer *clb.LoadBalancer) error {
	backends, err := cloud.describeLoadBalancerListenersBackends(loadBalancer.LoadBalancerId)
	if err != nil {
		return err
	}

	nodeLanIps := []string{}
	for _, node := range nodes {
		nodeLanIps = append(nodeLanIps, node.Name)
	}

	instancesInMultiVpc, err := cloud.describeInstancesByMultiLanIp(nodeLanIps)
	if err != nil {
		return err
	}

	instances := []cvm.InstanceInfo{}

	for _, instance := range instancesInMultiVpc {
		if instance.VirtualPrivateCloud.VpcID == cloud.config.VpcId {
			instances = append(instances, instance)
		}
	}

	backendsToAdd := []string{}
	backendsToDelete := []string{}

	for _, instance := range instances {
		found := false

		for _, backend := range backends {
			if backend.UnInstanceId == instance.InstanceID {
				found = true
			}
		}

		if !found {
			backendsToAdd = append(backendsToAdd, instance.InstanceID)
		}
	}

	for _, backend := range backends {
		found := false

		for _, instance := range instances {
			if instance.InstanceID == backend.UnInstanceId {
				found = true
			}
		}

		if !found {
			backendsToDelete = append(backendsToDelete, backend.UnInstanceId)
		}
	}

	backendToRegister := []clb.RegisterInstancesOpts{}
	backendToDeRegister := []string{}

	for _, backendToAdd := range backendsToAdd {
		backendToRegister = append(backendToRegister, clb.RegisterInstancesOpts{
			InstanceId: backendToAdd,
		})
	}

	for _, backendToDelete := range backendsToDelete {
		backendToDeRegister = append(backendToDeRegister, backendToDelete)
	}

	if len(backendToRegister) > 0 {
		result, err := clb.WaitUntilDone(
			func() (clb.AsyncTask, error) {
				return cloud.clb.RegisterInstancesWithLoadBalancer(&clb.RegisterInstancesWithLoadBalancerArgs{
					LoadBalancerId: loadBalancer.LoadBalancerId,
					Backends:       backendToRegister,
				})
			}, cloud.clb,
		)
		if err != nil {
			return err
		}
		if result != clb.TaskSuccceed {
			return errors.New("task is not succeed")
		}
	}

	if len(backendToDeRegister) > 0 {
		result, err := clb.WaitUntilDone(
			func() (clb.AsyncTask, error) {
				return cloud.clb.DeregisterInstancesFromLoadBalancer(
					loadBalancer.LoadBalancerId,
					backendToDeRegister,
				)
			}, cloud.clb,
		)
		if err != nil {
			return err
		}
		if result != clb.TaskSuccceed {
			return errors.New("task is not succeed")
		}
	}

	return nil
}

func (cloud *Cloud) ensureApplicationLoadBalancerBackends(ctx context.Context, clusterName string, service *v1.Service, nodes []*v1.Node, loadBalancer *clb.LoadBalancer) error {
	nodeLanIps := []string{}
	for _, node := range nodes {
		nodeLanIps = append(nodeLanIps, node.Name)
	}

	instancesInMultiVpc, err := cloud.describeInstancesByMultiLanIp(nodeLanIps)
	if err != nil {
		return err
	}

	instances := []cvm.InstanceInfo{}
	for _, instance := range instancesInMultiVpc {
		if instance.VirtualPrivateCloud.VpcID == cloud.config.VpcId {
			instances = append(instances, instance)
		}
	}

	response, err := cloud.clb.DescribeForwardLBBackends(&clb.DescribeForwardLBBackendsArgs{
		LoadBalancerId: loadBalancer.LoadBalancerId,
	})
	if err != nil {
		return err
	}

	forwardListeners := response.Data

	// remove unused listener first
	for _, port := range service.Spec.Ports {
		// find listener match this service port
		forwardListener := new(clb.ForwardLBListener)
		for _, listener := range forwardListeners {
			if listener.LoadBalancerPort == int(port.Port) && cloud.mapClbProtoToServicePortProto(listener.Protocol) == port.Protocol {
				forwardListener = &listener
				break
			}
		}

		if forwardListener == nil {
			return errors.New("Can not find loadbalancer listener for this service port")
		}

		backendsToDelete := make([]clb.ForwardLBListenerBackend, 0)

		for _, backend := range forwardListener.Backends {

			found := false

			for _, instance := range instances {
				if backend.UnInstanceId == instance.InstanceID && backend.Port == int(port.NodePort) {
					found = true
				}
			}

			if !found {
				backendsToDelete = append(backendsToDelete, backend)
			}
		}

		backendToDeRegister := make([]clb.DeregisterInstancesWithForwardLBFourthListenerBackendOpts, 0)

		for _, backendToDelete := range backendsToDelete {
			backendToDeRegister = append(backendToDeRegister, clb.DeregisterInstancesWithForwardLBFourthListenerBackendOpts{
				InstanceId: backendToDelete.UnInstanceId,
				Port:       backendToDelete.Port,
			})
		}

		if len(backendToDeRegister) > 0 {
			result, err := clb.WaitUntilDone(
				func() (clb.AsyncTask, error) {
					return cloud.clb.DeregisterInstancesFromForwardLBFourthListener(&clb.DeregisterInstancesFromForwardLBFourthListenerArgs{
						LoadBalancerId: loadBalancer.LoadBalancerId,
						ListenerId:     forwardListener.ListenerId,
						Backends:       backendToDeRegister,
					})
				}, cloud.clb,
			)
			if err != nil {
				return err
			}
			if result != clb.TaskSuccceed {
				return errors.New("task is not succeed")
			}
		}
	}

	// then add listener needed
	for _, port := range service.Spec.Ports {
		// find listener match this service port
		forwardListener := new(clb.ForwardLBListener)
		for _, listener := range forwardListeners {
			if listener.LoadBalancerPort == int(port.Port) && cloud.mapClbProtoToServicePortProto(listener.Protocol) == port.Protocol {
				forwardListener = &listener
				break
			}
		}

		if forwardListener == nil {
			return errors.New("Can not find loadbalancer listener for this service port")
		}

		backendsToAdd := make([]string, 0)

		for _, instance := range instances {
			found := false

			for _, backend := range forwardListener.Backends {
				if backend.UnInstanceId == instance.InstanceID && backend.Port == int(port.NodePort) {
					found = true
				}
			}

			if !found {
				backendsToAdd = append(backendsToAdd, instance.InstanceID)
			}
		}

		backendToRegister := make([]clb.RegisterInstancesWithForwardLBFourthListenerBackendOpts, 0)

		for _, backendToAdd := range backendsToAdd {
			backendToRegister = append(backendToRegister, clb.RegisterInstancesWithForwardLBFourthListenerBackendOpts{
				InstanceId: backendToAdd,
				Port:       int(port.NodePort),
			})
		}

		if len(backendToRegister) > 0 {
			result, err := clb.WaitUntilDone(
				func() (clb.AsyncTask, error) {
					return cloud.clb.RegisterInstancesWithForwardLBFourthListener(&clb.RegisterInstancesWithForwardLBFourthListenerArgs{
						LoadBalancerId: loadBalancer.LoadBalancerId,
						ListenerId:     forwardListener.ListenerId,
						Backends:       backendToRegister,
					})
				}, cloud.clb,
			)
			if err != nil {
				return err
			}
			if result != clb.TaskSuccceed {
				return errors.New("task is not succeed")
			}
		}
	}

	return nil
}

func (cloud *Cloud) createLoadBalancer(ctx context.Context, clusterName string, service *v1.Service) (*clb.LoadBalancer, error) {
	loadBalancerName := cloudprovider.GetLoadBalancerName(service)

	args := clb.CreateLoadBalancerArgs{
		VpcId:            &cloud.config.VpcId,
		LoadBalancerName: &loadBalancerName,
	}

	loadBalancerDesiredKind, ok := service.Annotations[ServiceAnnotationLoadBalancerKind]
	if !ok || (loadBalancerDesiredKind != LoadBalancerKindClassic && loadBalancerDesiredKind != LoadBalancerKindApplication) {
		loadBalancerDesiredKind = LoadBalancerKindApplication
	}
	loadBalancerDesiredType, ok := service.Annotations[ServiceAnnotationLoadBalancerType]
	if !ok || (loadBalancerDesiredType != LoadBalancerTypePrivate && loadBalancerDesiredType != LoadBalancerTypePublic) {
		loadBalancerDesiredType = LoadBalancerTypePublic
	}

	switch loadBalancerDesiredKind {
	case LoadBalancerKindApplication:
		args.Forward = &ClbLoadBalancerKindApplication
	case LoadBalancerKindClassic:
		args.Forward = &ClbLoadBalancerKindClassic
	default:
		args.Forward = &ClbLoadBalancerKindApplication
	}

	switch loadBalancerDesiredType {
	case LoadBalancerTypePrivate:
		args.LoadBalancerType = ClbLoadBalancerTypePrivate
	case LoadBalancerTypePublic:
		args.LoadBalancerType = ClbLoadBalancerTypePublic
	default:
		args.LoadBalancerType = ClbLoadBalancerTypePublic
	}

	if loadBalancerDesiredType == LoadBalancerTypePrivate {
		loadBalancerDesiredSubnetId, ok := service.Annotations[ServiceAnnotationLoadBalancerTypeInternalSubnetId]
		if !ok {
			return nil, errors.New("Subnet must be specified for private loadbalancer")
		}
		args.SubnetId = &loadBalancerDesiredSubnetId
	}

	result, err := clb.WaitUntilDone(
		func() (clb.AsyncTask, error) {
			return cloud.clb.CreateLoadBalancer(&args)
		},
		cloud.clb,
	)
	if err != nil {
		return nil, err
	}
	if result != clb.TaskSuccceed {
		return nil, errors.New("task is not succeed")
	}
	return cloud.getLoadBalancerByName(loadBalancerName)
}

func (cloud *Cloud) deleteLoadBalancer(ctx context.Context, clusterName string, service *v1.Service) error {
	loadBalancerName := cloudprovider.GetLoadBalancerName(service)
	loadBalancer, err := cloud.getLoadBalancerByName(loadBalancerName)
	if err != nil {
		if err == ErrCloudLoadBalancerNotFound {
			return nil
		}
		return err
	}

	result, err := clb.WaitUntilDone(
		func() (clb.AsyncTask, error) {
			return cloud.clb.DeleteLoadBalancers([]string{loadBalancer.LoadBalancerId})
		},
		cloud.clb,
	)
	if err != nil {
		return err
	}
	if result != clb.TaskSuccceed {
		return errors.New("task is not succeed")
	}

	return nil
}

func (cloud *Cloud) describeLoadBalancerListenersBackends(loadBalancerId string) ([]clb.LoadBalancerBackends, error) {
	backends := []clb.LoadBalancerBackends{}

	offset := 0
	limit := 100

	for {
		response, err := cloud.clb.DescribeLoadBalancerBackends(loadBalancerId, offset, limit)
		if err != nil {
			return []clb.LoadBalancerBackends{}, err
		}
		for _, backend := range response.BackendSet {
			backends = append(backends, backend)
		}

		if len(backends) < response.TotalCount {
			offset = len(backends)
		} else {
			break
		}
	}
	return backends, nil
}

func (cloud *Cloud) describeInstancesByMultiLanIp(ips []string) ([]cvm.InstanceInfo, error) {
	instances := []cvm.InstanceInfo{}

	offset := 0
	limit := 100

	ipsParas := make([]interface{}, len(ips))

	for idx, ip := range ips {
		ipsParas[idx] = ip
	}

	for {
		response, err := cloud.cvmV3.DescribeInstances(&cvm.DescribeInstancesArgs{
			Version: cvm.DefaultVersion,
			Filters: &[]cvm.Filter{{Name: cvm.FilterNamePrivateIpAddress, Values: ipsParas}},
			Offset:  &offset,
			Limit:   &limit,
		})
		if err != nil {
			return []cvm.InstanceInfo{}, err
		}
		for _, instance := range response.InstanceSet {
			instances = append(instances, instance)
		}

		if len(instances) < response.TotalCount {
			offset = len(instances)
		} else {
			break
		}
	}

	return instances, nil
}

func (cloud *Cloud) mapClbProtoToServicePortProto(proto int) v1.Protocol {
	switch proto {
	case ClbLoadBalancerListenerProtocolHTTP, ClbLoadBalancerListenerProtocolTCP, ClbLoadBalancerListenerProtocolHTTPS:
		return v1.ProtocolTCP
	case ClbLoadBalancerListenerProtocolUDP:
		return v1.ProtocolUDP
	default:
		return v1.ProtocolTCP
	}
}

func (cloud *Cloud) mapServicePortProtoClbProto(proto v1.Protocol) int {
	switch proto {
	case v1.ProtocolTCP:
		return ClbLoadBalancerListenerProtocolTCP
	case v1.ProtocolUDP:
		return ClbLoadBalancerListenerProtocolUDP
	default:
		return ClbLoadBalancerListenerProtocolTCP
	}
}

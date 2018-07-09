package clb

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/dbdd4us/qcloudapi-sdk-go/cvm"
)

func TestLoadBalanceBackends(t *testing.T) {
	client, err := NewClientFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	createLoadBalancerArgs := CreateLoadBalancerArgs{
		LoadBalancerType: LoadBalancerTypePrivateNetwork,
	}
	lb, err := client.CreateLoadBalancer(&createLoadBalancerArgs)
	if err != nil {
		t.Fatal(err)
	}

	dealId := lb.DealIds[0]
	lbId := lb.UnLoadBalancerIds[dealId][0]

	describeArgs := DescribeLoadBalancersArgs{
		LoadBalancerIds: &[]string{lbId},
	}

	for {
		time.Sleep(time.Second * 1)
		describeResponse, err := client.DescribeLoadBalancers(&describeArgs)
		if err != nil {
			continue
		}
		if len(describeResponse.LoadBalancerSet) > 0 {
			break
		}
	}

	cvmClient, err := cvm.NewClientFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	describeInstanceArgs := &cvm.DescribeInstancesArgs{
		Version: "2017-03-12",
	}

	instances, err := cvmClient.DescribeInstances(describeInstanceArgs)
	if err != nil {
		t.Fatal(err)
	}

	if len(instances.InstanceSet) <= 0 {
		t.Fatal("no enough instance for test")
	}

	instanceId := instances.InstanceSet[0].InstanceID

	registerArgs := RegisterInstancesWithLoadBalancerArgs{
		LoadBalancerId: lbId,
		Backends: []RegisterInstancesOpts{
			{
				InstanceId: instanceId,
			},
		},
	}

	registerResponse, err := client.RegisterInstancesWithLoadBalancer(&registerArgs)
	if err != nil {
		t.Fatal(err)
	}

	task := NewTask(registerResponse.RequestId)
	task.WaitUntilDone(context.Background(), client)

	describeResponse, err := client.DescribeLoadBalancerBackends(lbId, 0, 20)
	if err != nil {
		t.Fatal(err)
	}

	in := false

	for _, backend := range describeResponse.BackendSet {
		if backend.UnInstanceId == instanceId {
			in = true
			break
		}
	}

	if !in {
		t.Fatal(in)
	}

	modifyBackendArgs := ModifyLoadBalancerBackendsArgs{
		LoadBalancerId: lbId,
		Backends: []ModifyBackendOpts{
			{
				InstanceId: instanceId,
				Weight:     int(rand.Intn(100)),
			},
		},
	}

	_, err = WaitUntilDone(
		func() (AsyncTask, error) {
			return client.ModifyLoadBalancerBackends(&modifyBackendArgs)
		},
		client,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = WaitUntilDone(
		func() (AsyncTask, error) {
			return client.DeregisterInstancesFromLoadBalancer(lbId, []string{instanceId})
		},
		client,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = WaitUntilDone(
		func() (AsyncTask, error) {
			return client.DeleteLoadBalancers([]string{lbId})
		},
		client,
	)
	if err != nil {
		t.Fatal(err)
	}

}

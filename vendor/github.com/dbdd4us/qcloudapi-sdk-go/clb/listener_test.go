package clb

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestLoadBalancerListeners(t *testing.T) {
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

	createListenerArgs := CreateLoadBalancerListenersArgs{
		LoadBalancerId: lbId,
		Listeners: []CreateListenerOpts{
			{
				LoadBalancerPort: 9000 + rand.Int31n(1000),
				InstancePort:     9000 + rand.Int31n(1000),
				Protocol:         LoadBalanceListenerProtocolUDP,
			},
		},
	}

	createListenerResponse, err := client.CreateLoadBalancerListeners(&createListenerArgs)
	if err != nil {
		t.Fatal(err)
	}

	task := NewTask(createListenerResponse.RequestId)
	task.WaitUntilDone(context.Background(), client)

	describeLoadBalancerListenersArgs := DescribeLoadBalancerListenersArgs{
		LoadBalancerId: lbId,
	}

	lbListeners, err := client.DescribeLoadBalancerListeners(&describeLoadBalancerListenersArgs)
	if err != nil {
		t.Fatal(err)
	}

	newName := fmt.Sprintf("lb-listener-v-%d", rand.Int())

	modifyListenerArgs := ModifyLoadBalancerListenerArgs{
		LoadBalancerId: lbId,
		ListenerId:     lbListeners.ListenerSet[0].UnListenerId,
		ListenerName:   &newName,
	}

	_, err = WaitUntilDone(
		func() (AsyncTask, error) {
			return client.ModifyLoadBalancerListener(&modifyListenerArgs)
		},
		client,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = WaitUntilDone(
		func() (AsyncTask, error) {
			return client.DeleteLoadBalancerListeners(lbId, createListenerResponse.ListenerIds)
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

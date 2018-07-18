package clb

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestLoadBalancer(t *testing.T) {
	client, err := NewClientFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	loadBalancerType := LoadBalancerTypePrivateNetwork

	createArgs := CreateLoadBalancerArgs{
		LoadBalancerType: loadBalancerType,
	}

	createResponse, err := client.CreateLoadBalancer(&createArgs)
	if err != nil {
		t.Fatal(err)
	}

	dealId := createResponse.DealIds[0]
	lbId, ok := createResponse.UnLoadBalancerIds[dealId]
	if !ok {
		t.Fatal(err)
	}

	describeArgs := DescribeLoadBalancersArgs{
		LoadBalancerIds: &[]string{lbId[0]},
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

	newName := fmt.Sprintf("test-lb-v-%d", rand.Int())

	modifyArgs := ModifyLoadBalancerAttributesArgs{
		LoadBalancerId:   lbId[0],
		LoadBalancerName: &newName,
	}

	_, err = WaitUntilDone(
		func() (AsyncTask, error) {
			return client.ModifyLoadBalancerAttributes(&modifyArgs)
		},
		client,
	)
	if err != nil {
		t.Fatal(err)
	}

	_, err = WaitUntilDone(
		func() (AsyncTask, error) {
			return client.DeleteLoadBalancers(lbId)
		},
		client,
	)
	if err != nil {
		t.Fatal(err)
	}
}

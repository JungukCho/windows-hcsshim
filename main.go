package main

import (
	"encoding/json"
	"fmt"

	hcn "github.com/Microsoft/hcsshim/hcn"
)

func createEndPointReq(actionType hcn.ActionType, ipAddress string, port string) (*hcn.PolicyEndpointRequest, error) {
	in := hcn.AclPolicySetting{
		Protocols: "6",
		Action:    actionType,

		Direction:      hcn.DirectionTypeIn,
		LocalAddresses: ipAddress,
		//RemoteAddresses:
		LocalPorts: port,
		//RemotePorts:     "80",
		//RuleType: hcn.RuleTypeSwitch,
		Priority: 100,
	}

	rawJSON, err := json.Marshal(in)
	if err != nil {
		return nil, err
	}

	inPolicy := hcn.EndpointPolicy{
		Type:     hcn.ACL,
		Settings: rawJSON,
	}

	out := hcn.AclPolicySetting{
		Protocols: "6",
		Action:    actionType,
		Direction: hcn.DirectionTypeOut,
		//RuleType:  hcn.RuleTypeSwitch,
		Priority: 100,
	}

	rawJSON, err = json.Marshal(out)
	if err != nil {
		return nil, err
	}
	outPolicy := hcn.EndpointPolicy{
		Type:     hcn.ACL,
		Settings: rawJSON,
	}

	endpointRequest := hcn.PolicyEndpointRequest{
		Policies: []hcn.EndpointPolicy{inPolicy, outPolicy},
	}

	return &endpointRequest, nil

}

func printEndPointInfo(endpointID string) {
	eps, err := hcn.ListEndpointsOfNetwork(endpointID)
	if err != nil {
		fmt.Printf("cannot list %s endpoint due to err %v\n", endpointID, err)
	}

	for _, ep := range eps {
		fmt.Printf("%+v\n", ep)
	}
}
func printAllEPs() {
	eps, err := hcn.ListEndpoints()
	if err != nil {
		fmt.Printf("err %v\n", err)
	}

	for idx, ep := range eps {
		fmt.Printf("idx : %d %s %s %v\n", idx, ep.Id, ep.Name, ep.IpConfigurations)
		printEndPointInfo(ep.Id)
		for policyIdx, policy := range ep.Policies {
			jsonString, err := policy.Settings.MarshalJSON()
			if err != nil {
				fmt.Println("error ", err)
				continue
			}
			fmt.Println(policyIdx, " type :", policy.Type, " ", string(jsonString))
		}
	}
}

func checkSetPolicy() {
	supportedFeatures := hcn.GetSupportedFeatures()
	err := hcn.SetPolicySupported()
	if supportedFeatures.SetPolicy && err != nil {
		fmt.Println(err)
	}
	if !supportedFeatures.SetPolicy && err == nil {
		fmt.Println(err)
	}
}

func main() {
	checkSetPolicy()
	printAllEPs()

	// this is endpoint ID of nginx
	endpointID := "edffc876-3051-46ec-80b3-820dda20aa14"

	// there are two type of functions - What is the difference between hcsshim and hcn?
	// eps, err := hcsshim.GetHNSEndpointByID(endpointID)
	eps, err := hcn.GetEndpointByID(endpointID)
	if err != nil {
		fmt.Printf("Cannot find %s %v\n", endpointID, err)
		return
	}

	// simple test - first block traffic and then allow traffic.
	// #1. Block
	// endpointRequest, err := createEndPointReq(hcn.ActionTypeBlock, eps.IpConfigurations[0].IpAddress, "80")
	// #2. Allow
	endpointRequest, err := createEndPointReq(hcn.ActionTypeAllow, eps.IpConfigurations[0].IpAddress, "80")

	if err != nil {
		fmt.Println(err)
		return
	}

	jsonString, err := json.Marshal(endpointRequest)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("ACLS JSON:\n%s \n", jsonString)

	// #1. Block
	//err = eps.ApplyPolicy(hcn.RequestTypeAdd, *endpointRequest)

	// #2. Allow
	err = eps.ApplyPolicy(hcn.RequestTypeUpdate, *endpointRequest)

	// Delete operation does not work
	// err = eps.ApplyPolicy(hcn.RequestTypeRemove, *endpointRequest)

	if err != nil {
		fmt.Println(err)
		return
	}

}

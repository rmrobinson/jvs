package main

import (
	"fmt"

	"github.com/davecgh/go-spew/spew"
	"github.com/golang/protobuf/jsonpb"
	"github.com/rmrobinson/jvs/service/policy"
	"github.com/rmrobinson/jvs/service/policy/condition"
	"github.com/rmrobinson/jvs/service/policy/pb"
)

func parsePolicies(policySetStr string) (*policy.Engine, error) {
	var policySet pb.PolicySet
	err := jsonpb.UnmarshalString(policySetStr, &policySet)
	if err != nil {
		return nil, err
	}

	e := policy.NewEngine()

	for _, p := range policySet.Policies {
		c, err := condition.New(p.Condition)
		if err != nil  {
			return nil, err
		}

		np := policy.NewPolicy(p.Name, c, nil)
		e.AddPolicy(np)
	}

	return e, nil
}

func main() {
	m := jsonpb.Marshaler{}
	examplePolicy, err := m.MarshalToString(policy.ExamplePolicySet)
	if err != nil {
		fmt.Printf("Error marshaling example policy set: %s\n", err.Error())
		return
	}

	p, err := parsePolicies(examplePolicy)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
	} else {
		spew.Dump(p)
	}
}

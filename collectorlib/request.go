package collectorlib

import (
	"encoding/json"
	"io"

	"github.com/hashicorp/consul/consul/structs"
)

type Request struct {
	Root          structs.CheckServiceNodes
	TargetCheckID string
}

func ParseRequest(input io.Reader) (Request, error) {
	decoder := json.NewDecoder(input)

	var root structs.CheckServiceNodes
	if err := decoder.Decode(&root); err != nil {
		return Request{}, err
	}

	var defaultTargetCheckID string
	if len(root) > 0 {
		defaultTargetCheckID = "service:" + root[0].Service.ID
	}
	return Request{
		Root:          root,
		TargetCheckID: defaultTargetCheckID,
	}, nil
}

func (req Request) IPsByTag(tag string) []string {
	var ips []string
	for _, check := range req.Root {
		hasTag := false
		if tag == "*" {
			hasTag = true
		} else {
			for _, t := range check.Service.Tags {
				if t == tag {
					hasTag = true
				}
			}
		}
		if !hasTag {
			continue
		}

		for _, c := range check.Checks {
			if c.CheckID == req.TargetCheckID && c.Status == "passing" {
				ip := FindIPFromOutput(c.Output)
				if ip != "" {
					ips = append(ips, ip)
				}
			}
		}
	}

	return ips
}

package collectorlib

import (
	"fmt"
	"strings"
)

type Domain struct {
	FQDN string
	Tag  string
}

func NewDomain(option string) (*Domain, error) {
	pair := strings.Split(option, ":")
	switch len(pair) {
	case 1:
		return &Domain{FQDN: pair[0], Tag: "*"}, nil
	case 2:
		return &Domain{FQDN: pair[0], Tag: pair[1]}, nil
	default:
		return nil, fmt.Errorf("Invalid domain format: %s", option)
	}
}

func NewDomains(options []string) ([]*Domain, error) {
	var ret []*Domain
	for _, o := range options {
		d, err := NewDomain(o)
		if err != nil {
			return nil, err
		}
		ret = append(ret, d)
	}

	return ret, nil
}

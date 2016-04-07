package collectorlib

import (
	"strings"
	"testing"
)

var testJson = `[
  {
    "Node": {
      "Node": "nginx001",
      "Address": "10.1.10.101"
    },
    "Service": {
      "ID": "nginx",
      "Service": "nginx",
      "Tags": ["nginx", "lb-a"],
      "Port": 80
    },
    "Checks": [
      {
	"Node": "nginx001",
	"CheckID": "service:nginx",
	"Name": "Service 'nginx' check",
	"Status": "passing",
	"Notes": "",
	"Output": "ipaddr\t192.0.2.101",
	"ServiceID": "redis",
	"ServiceName": "redis"
      },
      {
	"Node": "foobar",
	"CheckID": "serfHealth",
	"Name": "Serf Health Status",
	"Status": "passing",
	"Notes": "",
	"Output": "",
	"ServiceID": "",
	"ServiceName": ""
      }
    ]
  },
  {
    "Node": {
      "Node": "nginx002",
      "Address": "10.1.10.102"
    },
    "Service": {
      "ID": "nginx",
      "Service": "nginx",
      "Tags": ["nginx", "lb-a"],
      "Port": 80
    },
    "Checks": [
      {
	"Node": "nginx002",
	"CheckID": "service:nginx",
	"Name": "Service 'nginx' check",
	"Status": "passing",
	"Notes": "",
	"Output": "ipaddr\t192.0.2.102",
	"ServiceID": "redis",
	"ServiceName": "redis"
      },
      {
	"Node": "foobar",
	"CheckID": "serfHealth",
	"Name": "Serf Health Status",
	"Status": "passing",
	"Notes": "",
	"Output": "",
	"ServiceID": "",
	"ServiceName": ""
      }
    ]
  },
  {
    "Node": {
      "Node": "nginx003",
      "Address": "10.1.10.103"
    },
    "Service": {
      "ID": "nginx",
      "Service": "nginx",
      "Tags": ["nginx", "lb-b"],
      "Port": 80
    },
    "Checks": [
      {
	"Node": "nginx003",
	"CheckID": "service:nginx",
	"Name": "Service 'nginx' check",
	"Status": "passing",
	"Notes": "",
	"Output": "ipaddr\t192.0.2.103",
	"ServiceID": "redis",
	"ServiceName": "redis"
      },
      {
	"Node": "foobar",
	"CheckID": "serfHealth",
	"Name": "Serf Health Status",
	"Status": "passing",
	"Notes": "",
	"Output": "",
	"ServiceID": "",
	"ServiceName": ""
      }
    ]
  }
]
`

func TestParseRequest(t *testing.T) {
	input := strings.NewReader(testJson)

	req, err := ParseRequest(input)
	if err != nil {
		t.Fatalf("Error parsing: %s", err.Error())
	}

	if size := len(req.Root); size != 3 {
		t.Errorf("req.Root size expected 3, got %d", size)
	}

	tags := req.Root[0].Service.Tags
	if tags[0] != "nginx" || tags[1] != "lb-a" {
		t.Errorf("req.Root[0].Service.Tags expected [\"nginx\", \"lb-a\"], got %v", tags)
	}

	tags = req.Root[2].Service.Tags
	if tags[0] != "nginx" || tags[1] != "lb-b" {
		t.Errorf("req.Root[2].Service.Tags expected [\"nginx\", \"lb-b\"], got %v", tags)
	}
}

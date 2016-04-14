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
	"Output": "ipaddr:192.0.2.101\n",
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
	"Output": "ipaddr:192.0.2.102\n",
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
	"Output": "ipaddr:192.0.2.103\n",
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
      "Node": "nginx004",
      "Address": "10.1.10.104"
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
	"Status": "failing",
	"Notes": "",
	"Output": "ipaddr:192.0.2.104\n",
	"ServiceID": "redis",
	"ServiceName": "redis"
      },
      {
	"Node": "foobar",
	"CheckID": "serfHealth",
	"Name": "Serf Health Status",
	"Status": "failing",
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

	if size := len(req.Root); size != 4 {
		t.Errorf("req.Root size expected 4, got %d", size)
	}

	tags := req.Root[0].Service.Tags
	if tags[0] != "nginx" || tags[1] != "lb-a" {
		t.Errorf("req.Root[0].Service.Tags expected [\"nginx\", \"lb-a\"], got %v", tags)
	}

	tags = req.Root[2].Service.Tags
	if tags[0] != "nginx" || tags[1] != "lb-b" {
		t.Errorf("req.Root[2].Service.Tags expected [\"nginx\", \"lb-b\"], got %v", tags)
	}

	ips := req.IPsByTag("lb-a")
	if len(ips) != 2 || ips[0] != "192.0.2.101" || ips[1] != "192.0.2.102" {
		t.Errorf("req.IPsByTag(\"lb-a\") expected [\"192.0.2.101\", \"192.0.2.102\"], got %v", ips)
	}

	ips = req.IPsByTag("lb-b")
	if len(ips) != 1 || ips[0] != "192.0.2.103" {
		t.Errorf("req.IPsByTag(\"lb-b\") expected [\"192.0.2.103\"], got %v", ips)
	}

	ips = req.IPsByTag("*")
	if len(ips) != 3 {
		t.Errorf("req.IPsByTag(\"*\") expected lb-a + lb-b, got %v", ips)
	}
}

func TestRequestCheckID(t *testing.T) {
	testJson2 := `[
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
	"Output": "ipaddr:192.0.2.101",
	"ServiceID": "redis",
	"ServiceName": "redis"
      },
      {
	"Node": "nginx001",
	"CheckID": "service:nginx_another_check",
	"Name": "Service 'nginx' check",
	"Status": "passing",
	"Notes": "",
	"Output": "ipaddr:192.0.2.111",
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
	input := strings.NewReader(testJson2)
	req, err := ParseRequest(input)
	if err != nil {
		t.Fatalf("Error parsing: %s", err.Error())
	}

	ips := req.IPsByTag("*")
	if len(ips) != 1 || ips[0] != "192.0.2.101" {
		t.Errorf("Check service:nginx expected to have IP [\"192.0.2.101\"], got %v", ips)
	}

	req.TargetCheckID = "service:nginx_another_check"
	ips = req.IPsByTag("*")
	if len(ips) != 1 || ips[0] != "192.0.2.111" {
		t.Errorf("Check service:nginx_another_check expected to have IP [\"192.0.2.111\"], got %v", ips)
	}
}

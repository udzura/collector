package collectorlib

import (
	"testing"
)

func TestFindIPFromOutput(t *testing.T) {
	str := "ipaddr:192.0.2.100"
	expected := "192.0.2.100"
	found := FindIPFromOutput(str)
	if found != expected {
		t.Errorf("Expect %s found, got %s", expected, found)
	}

	str = "ipaddr:192.0.2.101\n"
	expected = "192.0.2.101"
	found = FindIPFromOutput(str)
	if found != expected {
		t.Errorf("Expect %s found, got %s", expected, found)
	}

	str = "foo:bar buz\txxx:10.0.0.10\tipaddr:192.0.2.102\n"
	expected = "192.0.2.102"
	found = FindIPFromOutput(str)
	if found != expected {
		t.Errorf("Expect %s found, got %s", expected, found)
	}

	str = "foo:bar buz\txxx:10.0.0.10\n"
	expected = ""
	found = FindIPFromOutput(str)
	if found != expected {
		t.Errorf("Expect %s found, got %s", expected, found)
	}
}

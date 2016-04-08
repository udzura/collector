package collectorlib

import "testing"

func TestNewDomain(t *testing.T) {
	d1, _ := NewDomain("udzura.example.com")
	if d1.FQDN != "udzura.example.com" || d1.Tag != "*" {
		t.Errorf("Domain{udzura.example.com, *} expected, got %v", d1)
	}

	d2, _ := NewDomain("udzura2.example.com:lb-001")
	if d2.FQDN != "udzura2.example.com" || d2.Tag != "lb-001" {
		t.Errorf("Domain{udzura2.example.com, lb-001} expected, got %v", d2)
	}

	d3, err := NewDomain("udzura2.example.com:lb-001:invalid")
	if d3 != nil || err == nil {
		t.Errorf("Expected error, got %v, %v", d3, err)
	}
}

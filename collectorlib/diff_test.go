package collectorlib

import "testing"

func TestDiffChanged(t *testing.T) {
	d1 := NewDiff(
		[]string{"192.168.0.1", "192.168.0.2"},
		[]string{"192.168.0.1"},
	)
	if !d1.IsChanged() {
		t.Errorf("d1.IsChanged() should be true, but not: %v", d1)
	}

	d2 := NewDiff(
		[]string{"192.168.0.1", "192.168.0.2"},
		[]string{"192.168.0.1", "192.168.0.3"},
	)
	if !d2.IsChanged() {
		t.Errorf("d1.IsChanged() should be true, but not: %v", d2)
	}

	d3 := NewDiff(
		[]string{"192.168.0.1", "192.168.0.2"},
		[]string{"192.168.0.1", "192.168.0.2"},
	)
	if d3.IsChanged() {
		t.Errorf("d1.IsChanged() should be false, but not: %v", d3)
	}
}

func TestDiffToString(t *testing.T) {
	var expected string

	d1 := NewDiff(
		[]string{"192.168.0.1", "192.168.0.2"},
		[]string{"192.168.0.1"},
	)
	expected = ` [192.168.0.1]
-[192.168.0.2]`
	if d1.ToString() != expected {
		t.Errorf("d1.ToString() expected:\n%s\n\ngot:\n%s", expected, d1.ToString())
	}

	d2 := NewDiff(
		[]string{"192.168.0.1", "192.168.0.2"},
		[]string{"192.168.0.1", "192.168.0.3"},
	)
	expected = `+[192.168.0.3]
 [192.168.0.1]
-[192.168.0.2]`
	if d2.ToString() != expected {
		t.Errorf("d2.ToString() expected:\n%s\n\ngot:\n%s", expected, d2.ToString())
	}

	d3 := NewDiff(
		[]string{"192.168.0.1", "192.168.0.2"},
		[]string{"192.168.0.1", "192.168.0.2"},
	)
	expected = ` [192.168.0.1 192.168.0.2]`
	if d3.ToString() != expected {
		t.Errorf("d3.ToString() expected:\n%s\n\ngot:\n%s", expected, d3.ToString())
	}

	d4 := NewDiff(
		[]string{},
		[]string{"192.168.0.1", "192.168.0.2"},
	)
	expected = `+[192.168.0.1 192.168.0.2]`
	if d4.ToString() != expected {
		t.Errorf("d4.ToString() expected:\n%s\n\ngot:\n%s", expected, d4.ToString())
	}
}

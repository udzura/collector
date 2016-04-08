package collectorlib

import (
	"fmt"
	"sort"
	"strings"
)

type Diff struct {
	IPsBefore []string
	IPsAfter  []string

	added, deleted, existsBoth []string
}

func NewDiff(before, after []string) *Diff {
	d := &Diff{
		IPsBefore: before,
		IPsAfter:  after,
	}
	d.initialize()
	return d
}

func (d *Diff) initialize() {
	sort.Strings(d.IPsBefore)
	sort.Strings(d.IPsAfter)
	d.detectChange()
}

func (d *Diff) detectChange() {
	d.deleted = append([]string{}, d.IPsBefore...)
	d.added = append([]string{}, d.IPsAfter...)
	d.existsBoth = []string{}

	for _, i1 := range d.IPsAfter {
		for idx, i2 := range d.deleted {
			if i1 == i2 {
				d.deleted = append(d.deleted[:idx], d.deleted[idx+1:]...)
				d.existsBoth = append(d.existsBoth, i1)
			}
		}
	}

	for _, i1 := range d.existsBoth {
		for idx, i2 := range d.added {
			if i1 == i2 {
				d.added = append(d.added[:idx], d.added[idx+1:]...)
			}
		}
	}
}

func (d *Diff) IsChanged() bool {
	return len(d.added) != 0 || len(d.deleted) != 0
}

func (d *Diff) ToString() string {
	s := ""
	if len(d.added) > 0 {
		s += fmt.Sprintf("+%v\n", d.added)
	}
	if len(d.existsBoth) > 0 {
		s += fmt.Sprintf(" %v\n", d.existsBoth)
	}
	if len(d.deleted) > 0 {
		s += fmt.Sprintf("-%v\n", d.deleted)
	}
	s = strings.TrimSuffix(s, "\n")

	return s
}

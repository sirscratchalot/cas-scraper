package main

import (
	"testing"
)

func TestRegexp(t *testing.T) {
	casNumbers := []string{"64-17-5", "7440-22-4"}
	for _, cas := range casNumbers {
		if !casFormat.MatchString(cas) {
			t.Errorf("Cas regexp did not match: %s", cas)
		}
	}
}

package main

import "testing"

func TestDangerDir(t *testing.T) {
	if !isDangerDir("/") {
		t.Fail()
	}
	if !isDangerDir("/bin") {
		t.Fail()
	}
	if !isDangerDir("/bin/") {
		t.Fail()
	}
}

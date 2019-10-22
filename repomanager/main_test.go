package main

import "testing"

func Test_clone(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
		{
			name: "fdsfd",
			args: args{"https://github.com/Kuangcp/GoBase"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			println(tt.name, t)
		})
	}
}
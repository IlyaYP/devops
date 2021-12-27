package main

import (
	"github.com/IlyaYP/devops/storage/inmemory"
	"testing"
)

func Test_testStore(t *testing.T) {
	type args struct {
		st *inmemory.Storage
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testStore(tt.args.st)
		})
	}
}

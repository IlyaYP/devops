package internal

import "testing"

func TestSend(t *testing.T) {
	type args struct {
		endpoint string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Send(tt.args.endpoint)
		})
	}
}

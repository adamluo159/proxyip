package proxyip

import "testing"

func Test_startProxyIpServer(t *testing.T) {
	type args struct {
		addr string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			startProxyIpServer(tt.args.addr)
		})
	}
}

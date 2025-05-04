package netutil

import (
	"testing"
)

func TestIsMyActiveHostIp(t *testing.T) {
	type args struct {
		ip string
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "valid IP",
			args: args{
				ip: "172.30.1.96",
			},
			want: true,
		},
		{
			name: "localhost IP",
			args: args{
				ip: "127.0.0.1",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IsMyActiveHostIp(tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsMyActiveHostIp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("IsMyActiveHostIp() got = %v, want %v", got, tt.want)
			}
		})
	}
}

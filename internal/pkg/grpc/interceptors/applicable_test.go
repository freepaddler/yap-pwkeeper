package interceptors

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_applicable(t *testing.T) {
	tests := []struct {
		name    string
		apply   map[string]bool
		request string
		want    bool
	}{
		{
			name:    "service match",
			request: "/grpcapi.Wallet/AddCredential",
			apply:   map[string]bool{"Wallet/": true},
			want:    true,
		},
		{
			name:    "method match",
			request: "/grpcapi.Wallet/AddCredential",
			apply:   map[string]bool{"Wallet/AddCredential": true},
			want:    true,
		},
		{
			name:    "no service match",
			request: "/grpcapi.Wallet1/AddCredential",
			apply:   map[string]bool{"Wallet/": true},
			want:    false,
		},
		{
			name:    "no method match",
			request: "/grpcapi.Wallet/AddCredential1",
			apply:   map[string]bool{"Wallet/AddCredential": true},
			want:    false,
		},
		{
			name:    "no match other service same method",
			request: "/grpcapi.Wallet1/AddCredential",
			apply:   map[string]bool{"Wallet/AddCredential": true},
			want:    false,
		},
	}
	for _, tt := range tests {
		got := applicable(tt.apply, tt.request)
		require.Equalf(t, tt.want, got, "expect %t got %t for request %s with map %v", tt.want, got, tt.request, tt.apply)
	}
	//apply := make(map[string]bool)
	//str := "/grpcapi.Wallet/AddCredential"
	//apply["some"] = true
	//fmt.Println(applicable(apply, str))
	//apply["Wallet/AddCred"] = true
	//fmt.Println(applicable(apply, str))
	//apply["Wallet/AddCredential1"] = true
	//fmt.Println(applicable(apply, str))
}

package interceptors

import (
	"fmt"
	"testing"
)

func Test_match(t *testing.T) {
	apply := make(map[string]bool)
	str := "/grpcapi.Wallet/AddCredential"
	apply["some"] = true
	fmt.Println(applicable(apply, str))
	apply["Wallet/AddCred"] = true
	fmt.Println(applicable(apply, str))
	apply["Wallet/AddCredential1"] = true
	fmt.Println(applicable(apply, str))
}

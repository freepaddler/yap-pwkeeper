package jwtToken

import (
	"fmt"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func TestNewToken(t *testing.T) {
	userid := "someuser"
	token, err := NewToken(userid)
	require.NoError(t, err)
	parsed, _, err := jwt.NewParser().ParseUnverified(token, &JWTClaims{})
	fmt.Printf("%+v\n", parsed.Claims)
	claims, err := getClaims(token)
	fmt.Printf("%+v\n", claims)

}

//import (
//	"context"
//	"strconv"
//	"testing"
//	"time"
//
//	"github.com/golang-jwt/jwt/v5"
//	"github.com/stretchr/testify/require"
//
//	"github.com/freepaddler/yap-s56-loyalty/internal/pkg/core"
//)
//
//// get subject from token
//func getSubjectClaim(ss string) string {
//	token, _ := jwt.ParseWithClaims(ss, &jwt.RegisteredClaims{}, nil)
//	s, _ := token.Claims.GetTokenSubject()
//	return s
//}
//
//// returns jwt with ttl
//func getJwt(user *core.User, ttl time.Duration, key []byte) string {
//	oldTTL := jwtTTL
//	jwtTTL = ttl
//	oldKey := JWTKey
//	JWTKey = key
//	jwtTTL = ttl
//	token, _ := GenJWT(context.TODO(), user)
//	jwtTTL = oldTTL
//	JWTKey = oldKey
//	return token
//}
//
//func TestGenJWT(t *testing.T) {
//	userId := 101
//	user := &core.User{}
//	user.Id = userId
//	token, err := GenJWT(context.TODO(), user)
//	require.NoError(t, err)
//	require.Equal(t, strconv.Itoa(userId), getSubjectClaim(token))
//}
//
//func TestValidateJWT(t *testing.T) {
//	id := 101
//	user := &core.User{}
//	user.Id = id
//
//	tests := []struct {
//		name    string
//		token   string
//		want    *core.User
//		wantErr bool
//	}{
//		{
//			name:    "JWT OK",
//			token:   getJwt(user, jwtTTL, JWTKey),
//			want:    user,
//			wantErr: false,
//		},
//		{
//			name:    "JWT Expired",
//			token:   getJwt(user, 0, JWTKey),
//			want:    user,
//			wantErr: true,
//		},
//		{
//			name:    "JWT Bad Signature",
//			token:   getJwt(user, jwtTTL, []byte("qwe")),
//			want:    user,
//			wantErr: true,
//		},
//		{
//			name:    "JWT Malformed",
//			token:   "this is not a token",
//			want:    user,
//			wantErr: true,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := ValidateJWT(context.TODO(), tt.token)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			if err == nil {
//				require.Equal(t, user, got)
//			}
//		})
//	}
//}

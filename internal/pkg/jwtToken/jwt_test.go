package jwtToken

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewToken(t *testing.T) {
	tStart := time.Now()
	ttl := 2 * time.Minute
	SetTTL(ttl)
	subject := "some string"
	token, err := NewToken(subject)
	tEnd := time.Now()
	require.NoError(t, err, "expect no error on token creation")
	claims, err := getClaims(token)
	require.NoError(t, err, "expect no errot on get token claims")
	assert.Equal(t, subject, claims.Subject, "expect subject match")
	assert.NotEqual(t, "", claims.Session, "expect non-empty session")
	assert.WithinRange(t, claims.ExpiresAt.Time, tStart.Add(ttl).Truncate(time.Second), tEnd.Add(ttl).Truncate(time.Second), "expect expire to match ttl")
}

func TestGetTokenExpire(t *testing.T) {
	tStart := time.Now()
	ttl := 2 * time.Minute
	SetTTL(ttl)
	token, err := NewToken("subject")
	tEnd := time.Now()
	require.NoError(t, err, "expect no error on token creation")
	got, err := GetTokenExpire(token)
	require.NoError(t, err, "expect no error on token claims")
	assert.WithinRange(t, got, tStart.Add(ttl).Truncate(time.Second), tEnd.Add(ttl).Truncate(time.Second), "expect expire to match ttl")
}

func TestGetTokenSubject(t *testing.T) {
	subject := "some text"
	token, err := NewToken(subject)
	require.NoError(t, err, "expect no error on token creation")
	got := GetTokenSubject(token)
	require.Equal(t, got, subject, "expect subject to match")
}

func TestGetTokenSession(t *testing.T) {
	token, err := NewToken("subject")
	require.NoError(t, err, "expect no error on token creation")
	got := GetTokenSession(token)
	require.NotEqual(t, "", got, "expect session not empty string")
	token2, err := NewToken("subject")
	got2 := GetTokenSession(token2)
	require.NotEqual(t, "", got, "expect session not empty string")
	require.NotEqual(t, got, got2, "expect session difference")
}

func TestRefreshToken_valid(t *testing.T) {
	subject := "some subject"
	ttl := 2 * time.Minute
	SetTTL(ttl)
	token1, err := NewToken(subject)
	require.NoError(t, err, "expect no error on token creation")
	time.Sleep(time.Second)
	token2, err := RefreshToken(token1)
	require.NoError(t, err, "token should be refreshed without error")
	claims1, err := getClaims(token1)
	require.NoError(t, err, "expect no error on token claims")
	claims2, err := getClaims(token2)
	require.NoError(t, err, "expect no error on token claims")
	assert.Equal(t, claims1.Subject, claims2.Subject, "expect subject to match")
	assert.Equal(t, claims1.Session, claims2.Session, "expect session to match")
	assert.NotEqual(t, claims1.ID, claims2.ID, "expect token id to be changed")
	assert.Less(t, claims1.ExpiresAt.Time, claims2.ExpiresAt.Time, "expect expiredAt to be changed")
}

func TestRefreshToken_invalid(t *testing.T) {
	subject := "some subject"
	SetKey("12345678")
	token1, err := NewToken(subject)
	require.NoError(t, err, "expect no error on token creation")
	SetKey("87654321")
	_, err = RefreshToken(token1)
	require.Error(t, err, "token refresh should fail")
}

// generates test tokens
func generateWithClaims(c JWTClaims, k []byte, s jwt.SigningMethod) (string, error) {
	var signed string
	claims := c
	token := jwt.NewWithClaims(s, claims)
	signed, err := token.SignedString(k)
	return signed, err
}

func TestValid(t *testing.T) {
	tests := []struct {
		name   string
		claims JWTClaims
		key    []byte
		sign   jwt.SigningMethod
		want   bool
	}{
		{
			name: "valid token",
			claims: JWTClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject:   "some subject",
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
					ID:        "someId",
				},
				Session: "someSession",
			},
			key:  jwtKey,
			sign: jwtSign,
			want: true,
		},
		{
			name: "invalid key",
			claims: JWTClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject:   "some subject",
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
					ID:        "someId",
				},
				Session: "someSession",
			},
			key:  []byte("some other key"),
			sign: jwtSign,
			want: false,
		},
		{
			name: "expired",
			claims: JWTClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject:   "some subject",
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Second)),
					ID:        "someId",
				},
				Session: "someSession",
			},
			key:  jwtKey,
			sign: jwtSign,
			want: false,
		},
		{
			name: "invalid sign",
			claims: JWTClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject:   "some subject",
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
					ID:        "someId",
				},
				Session: "someSession",
			},
			key:  jwtKey,
			sign: jwt.SigningMethodHS384,
			want: false,
		},
		{
			name: "no subject",
			claims: JWTClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
					ID:        "someId",
				},
				Session: "someSession",
			},
			key:  jwtKey,
			sign: jwtSign,
			want: false,
		},
		{
			name: "no session",
			claims: JWTClaims{
				RegisteredClaims: jwt.RegisteredClaims{
					Subject:   "some subject",
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute)),
					ID:        "someId",
				},
			},
			key:  jwtKey,
			sign: jwtSign,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := generateWithClaims(tt.claims, tt.key, tt.sign)
			require.NoError(t, err, "expect no error on token generation")
			got := Valid(token)
			require.Equal(t, tt.want, got, "expect token valid to be %t, got %t", tt.want, got)
		})
	}
}

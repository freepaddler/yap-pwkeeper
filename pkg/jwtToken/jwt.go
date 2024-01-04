package jwtToken

import (
	"crypto/rand"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"yap-pwkeeper/internal/pkg/logger"
)

var (
	ErrNoSubject = errors.New("no subject for jwt")
	ErrSign      = errors.New("failed to sign jwt")
	ErrInvalid   = errors.New("invalid jwt")
	jwtKey       []byte
	jwtTTL       = 2 * time.Hour
	jwtSign      = jwt.SigningMethodHS256
)

func init() {
	jwtKey = make([]byte, 64)
	_, _ = rand.Read(jwtKey)
}

func SetKey(key string) {
	if len(key) < 8 {
		logger.Log().Warn("token key is too short, using autogenerated")
		return
	}
	jwtKey = []byte(key)
}

func SetTTL(ttl time.Duration) {
	jwtTTL = ttl
}

type JWTClaims struct {
	jwt.RegisteredClaims
	Session string `json:"session,omitempty"`
}

func genToken(subject string, session string) (string, error) {
	var signed string
	if session == "" {
		session = uuid.NewString()
	}
	if subject == "" {
		return signed, ErrNoSubject
	}
	claims := &JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   subject,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtTTL)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
		Session: session,
	}
	token := jwt.NewWithClaims(jwtSign, claims)
	signed, err := token.SignedString(jwtKey)
	if err != nil {
		return signed, ErrSign
	}
	return signed, nil
}

func NewToken(subject string) (string, error) {
	return genToken(subject, "")
}

func RefreshToken(token string) (string, error) {
	if !ValidateToken(token) {
		return "", ErrInvalid
	}
	return genToken(GetTokenSubject(token), GetTokenSession(token))
}

func ValidateToken(signed string) bool {
	token, err := jwt.ParseWithClaims(signed, &JWTClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
	if err != nil || !token.Valid || token.Method != jwtSign {
		return false
	}
	if GetTokenSession(signed) == "" || GetTokenSubject(signed) == "" {
		return false
	}
	return true
}

func getClaims(token string) (*JWTClaims, error) {
	parsed, _, err := jwt.NewParser().ParseUnverified(token, &JWTClaims{})
	return parsed.Claims.(*JWTClaims), err
}

func GetTokenSubject(token string) string {
	claims, err := getClaims(token)
	if err != nil {
		return ""
	}
	return claims.Subject
}

func GetTokenSession(token string) string {
	claims, err := getClaims(token)
	if err != nil {
		return ""
	}
	return claims.Session
}

//func GenJWT(ctx context.Context, user models.User) (string, error) {
//	//log := logger.Log().WithCtxReqId(ctx)
//	logger.Log().Debug("generating jwt")
//	if user.Id == "" {
//		return "", fmt.Errorf("can't generate jwt: user id is empty")
//	}
//	claims := &jwt.RegisteredClaims{
//		Subject:   user.Id,
//		ExpiresAt: jwt.NewNumericDate(time.Now().Add(jwtTTL)),
//		NotBefore: jwt.NewNumericDate(time.Now()),
//		IssuedAt:  jwt.NewNumericDate(time.Now()),
//		ID:        uuid.New().String(),
//	}
//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	ss, err := token.SignedString(JWTKey)
//	if err != nil {
//		return "", fmt.Errorf("jwt sign failed: %w", err)
//	}
//	return ss, nil
//}

//func ValidateJWT(ctx context.Context, ss string) (*core.User, error) {
//	log := logger.Log.WithCtxReqId(ctx)
//	log.Debug("validating jwt")
//	token, err := jwt.ParseWithClaims(ss, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
//		return JWTKey, nil
//	})
//	if err == nil && token.Valid {
//		user := &core.User{}
//		userId, err := token.Claims.GetTokenSubject()
//		if err != nil {
//			return nil, fmt.Errorf("can't get user.Login from subject claim: %w", err)
//		}
//		user.Id, err = strconv.Atoi(userId)
//		if err != nil {
//			return nil, fmt.Errorf("can't get user.Login from subject claim: %w", err)
//		}
//		log.Debug("jwt is valid")
//		return user, nil
//	}
//	return nil, fmt.Errorf("jwt is not valid %w", err)
//}

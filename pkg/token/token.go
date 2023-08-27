package token

import (
	"github.com/fleezesd/kit-study/internal/pkg/errno"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("jwtSecret")

const JWT_CONTEXT_KEY = "jwt_context_key"

func ParseToken(token string) (jwt.MapClaims, error) {
	jwtToken, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil || jwtToken == nil {
		return nil, errno.ErrTokenInvalid
	}
	claim, ok := jwtToken.Claims.(jwt.MapClaims)
	if ok && jwtToken.Valid {
		return claim, nil
	} else {
		return nil, nil
	}
}

// Sign 使用 jwtSecret 签发 token，token 的 claims 中会存放传入的 subject.
func Sign(secret string) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		JWT_CONTEXT_KEY: secret,                                    // 身份标识
		"nbf":           time.Now().Unix(),                         // 生效时间
		"iat":           time.Now().Unix(),                         // 签发时间
		"exp":           time.Now().Add(100000 * time.Hour).Unix(), // 过期时间
		"sub":           "login",
	})

	// 签发 token
	tokenString, err = token.SignedString(jwtSecret)

	return
}

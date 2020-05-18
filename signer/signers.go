package signer

import (
	"errors"
	"fmt"
	"github.com/azon0320/jwtcli/utils"
	"github.com/pascaldekloe/jwt"
	"strings"
)

var (
	registeredSigners = make(map[string]JwtSigner, 0)
)

type JwtSigner func(claims *jwt.Claims, options *utils.Wrapper) (string, error)

func RegisterSigner(alg string, signer JwtSigner) {
	registeredSigners[alg] = signer
}

func MakeJwt(alg string, claims *jwt.Claims, opts *utils.Wrapper) (string, error) {
	signer, ok := registeredSigners[strings.ToLower(alg)]
	if !ok {
		return "", errors.New(fmt.Sprintf("alg %s not supported", alg))
	}
	return signer(claims, opts)
}

func init() {
	RegisterSigner("hs256", func(claims *jwt.Claims, options *utils.Wrapper) (s string, e error) {
		dat, err := claims.HMACSign(jwt.HS256, []byte(options.GetString("secret")))
		return string(dat), err
	})
	RegisterSigner("hs384", func(claims *jwt.Claims, options *utils.Wrapper) (s string, e error) {
		dat, err := claims.HMACSign(jwt.HS384, []byte(options.GetString("secret")))
		return string(dat), err
	})
	RegisterSigner("hs512", func(claims *jwt.Claims, options *utils.Wrapper) (s string, e error) {
		dat, err := claims.HMACSign(jwt.HS512, []byte(options.GetString("secret")))
		return string(dat), err
	})
}

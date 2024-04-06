package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	//mySigningKey = []byte("test")
	mySigningKey = []byte(os.Getenv("key"))
)

func verifyJWTToken(tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return mySigningKey, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims["username"].(string), nil
	}
	return "", nil
}

func generateJWTToken(username string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	})
	tokenString, _ := token.SignedString(mySigningKey)
	return tokenString
}

func toReqContext(w http.ResponseWriter, r *http.Request) *ReqContext {
	ctx := r.Context().Value(Key{})
	reqCtx, ok := ctx.(*ReqContext)
	if !ok || reqCtx == nil {
		log.Fatalln("toReqContext error")
	}

	reqCtx.Req = r
	reqCtx.Resp = w
	return reqCtx
}

func toHTTPHandlerFunc(h Handler) http.HandlerFunc {
	handle, ok := h.(func(*ReqContext))
	if !ok {
		log.Fatalln("toHTTPHandlerFunc error")
	}

	return func(w http.ResponseWriter, r *http.Request) {
		handle(toReqContext(w, r))
	}
}

func getMWFromHandler(handler Handler) Middleware {
	if mw, ok := handler.(Middleware); ok {
		return mw
	}
	h := toHTTPHandlerFunc(handler)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			h.ServeHTTP(w, r)

			ctx := r.Context().Value(momentumContextKey{}).(*Context)
			next.ServeHTTP(ctx.Resp, ctx.Req)
		})
	}
}

func generateNonce() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

package middleware

import (
	"context"
	"encoding/base64"
	"github.com/suoaao/affordable-copilot/pkg/conf"
	"net/http"
	"strings"
)

func VerifyRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		arr := strings.Split(r.Header.Get("Authorization"), " ")

		authToken, isAuth := "", false
		switch {
		case len(arr) != 2:
			isAuth = false
		case arr[0] == "token":
			authToken = arr[1]
			isAuth = conf.Conf.Auth(authToken)
		case arr[0] == "Basic":
			data, err := base64.StdEncoding.DecodeString(arr[1])
			authToken = strings.Trim(string(data), ":")
			if err != nil {
				isAuth = false
			} else {
				isAuth = conf.Conf.Auth(authToken)
			}
		}
		if !isAuth {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 Unauthorized"))
			return
		}

		ctx := context.WithValue(r.Context(), "auth_token", authToken)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

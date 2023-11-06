package middleware

import (
	"context"
	"github.com/suoaao/affordable-copilot/pkg/conf"
	"net/http"
	"strings"
)

func VerifyRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		arr := strings.Split(r.Header.Get("Authorization"), " ")
		if len(arr) != 2 || !conf.Conf.Auth(arr[1]) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 Unauthorized"))
			return
		}
		ctx := context.WithValue(r.Context(), "auth_token", arr[1])
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

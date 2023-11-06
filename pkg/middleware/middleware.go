package middleware

import (
	"context"
	"github.com/suoaao/affordable-ai/pkg/conf"
	"net/http"
	"path"
	"strings"
)

func VerifyRequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		arr := strings.Split(r.Header.Get("Authorization"), " ")
		if len(arr) != 2 || arr[1] != conf.AuthToken {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("401 Unauthorized"))
			return
		}
		ctx := context.WithValue(r.Context(), "auth_token", arr[1])
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func RemoveFirstPathElement(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) > 1 {
			pathParts = pathParts[2:]
		}
		r.URL.Path = "/" + path.Join(pathParts...)

		next.ServeHTTP(w, r)
	})
}

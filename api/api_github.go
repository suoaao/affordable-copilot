package handler

import (
	"context"
	"encoding/json"
	"github.com/suoaao/affordable-ai/pkg/cache"
	"github.com/suoaao/affordable-ai/pkg/conf"
	"github.com/suoaao/affordable-ai/pkg/extension"
	myHttputil "github.com/suoaao/affordable-ai/pkg/httputil"
	"github.com/suoaao/affordable-ai/pkg/middleware"
	"net/http"
	"net/http/httputil"
	"time"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	handler := middleware.VerifyRequestMiddleware(githubProxyHandler)
	handler.ServeHTTP(w, r)
}

var githubProxyHandler, _ = NewGithubProxyHandler(conf.GhuToken)

type GithubProxyHandler struct {
	ghuToken     string
	proxy        *httputil.ReverseProxy
	coTokenCache *cache.CoTokenCache
}

func NewGithubProxyHandler(ghuToken string) (*GithubProxyHandler, error) {
	var proxy, err = myHttputil.NewReverseProxy("https://api.github.com")
	if err != nil {
		return nil, err
	}
	return &GithubProxyHandler{
		ghuToken:     ghuToken,
		proxy:        proxy,
		coTokenCache: cache.NewCoTokenCache(extension.Redis),
	}, nil
}

func (h *GithubProxyHandler) ProxyRequest(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Authorization", "token "+conf.GhuToken)
	h.proxy.ServeHTTP(w, r)
}

func (h *GithubProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/copilot_internal/v2/token":
		h.cachedCoToken(w, r)
	default:
		h.ProxyRequest(w, r)
	}
}

func (h *GithubProxyHandler) cachedCoToken(w http.ResponseWriter, r *http.Request) {
	authToken, ok := r.Context().Value("auth_token").(string)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 Unauthorized"))
		return
	}
	coToken, err := h.coTokenCache.Get(r.Context(), authToken)
	if err == nil {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(coToken)
		return
	}

	portal := make(chan []byte, 1)
	defer close(portal)
	ctx := context.WithValue(r.Context(), "is_cache", portal)
	r = r.WithContext(ctx)

	h.ProxyRequest(w, r)

	body, ok := <-portal
	if !ok || len(body) == 0 {
		return
	}
	resp := struct {
		ExpiresAt int64 `json:"expires_at"`
	}{}
	json.Unmarshal(body, &resp)
	ttl := time.Duration(resp.ExpiresAt-time.Now().Unix()-1) * time.Second

	h.coTokenCache.Set(r.Context(), authToken, body, ttl)
}

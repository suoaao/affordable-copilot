package handler

import (
	"context"
	"encoding/json"
	"github.com/suoaao/affordable-copilot/pkg/cache"
	"github.com/suoaao/affordable-copilot/pkg/conf"
	"github.com/suoaao/affordable-copilot/pkg/extension"
	myHttputil "github.com/suoaao/affordable-copilot/pkg/httputil"
	"github.com/suoaao/affordable-copilot/pkg/middleware"
	"net/http"
	"net/http/httputil"
	"time"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	handler := middleware.VerifyRequestMiddleware(copilotProxyHandler)
	handler.ServeHTTP(w, r)
}

var copilotProxyHandler, _ = NewCopilotProxyHandler(conf.Conf.GhuToken)

type CopilotProxyHandler struct {
	ghuToken     string
	proxy        *httputil.ReverseProxy
	coTokenCache *cache.CoTokenCache
}

func NewCopilotProxyHandler(ghuToken string) (*CopilotProxyHandler, error) {
	var proxy, err = myHttputil.NewReverseProxy("https://api.github.com")
	if err != nil {
		return nil, err
	}
	return &CopilotProxyHandler{
		ghuToken:     ghuToken,
		proxy:        proxy,
		coTokenCache: cache.NewCoTokenCache(extension.Redis),
	}, nil
}

func (h *CopilotProxyHandler) ProxyRequest(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Authorization", "token "+h.ghuToken)
	h.proxy.ServeHTTP(w, r)
}

func (h *CopilotProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/copilot_internal/v2/token":
		h.cachedCoToken(w, r)
	default:
		h.ProxyRequest(w, r)
	}
}

func (h *CopilotProxyHandler) cachedCoToken(w http.ResponseWriter, r *http.Request) {
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

	var body []byte
	select {
	case body, ok = <-portal:
		if !ok || len(body) == 0 {
			return
		}
	case <-time.After(1 * time.Second):
		return
	}

	resp := struct {
		ExpiresAt int64 `json:"expires_at"`
	}{}
	json.Unmarshal(body, &resp)
	ttl := time.Duration(resp.ExpiresAt-time.Now().Unix()-1) * time.Second

	h.coTokenCache.Set(r.Context(), authToken, body, ttl)
}

package handler

import (
	"github.com/suoaao/affordable-ai/pkg/conf"
	myHttpUtil "github.com/suoaao/affordable-ai/pkg/httputil"
	"github.com/suoaao/affordable-ai/pkg/middleware"
	"net/http"
	"net/http/httputil"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	handler := middleware.RemoveFirstPathElement(middleware.VerifyRequestMiddleware(openaiProxy))
	handler.ServeHTTP(w, r)
}

var openaiProxy, _ = NewOpenaiProxyHandler(conf.Conf.Openai.ApiKey)

type OpenaiProxyHandler struct {
	apiKey string
	proxy  *httputil.ReverseProxy
}

func NewOpenaiProxyHandler(apiKey string) (*OpenaiProxyHandler, error) {
	var proxy, err = myHttpUtil.NewReverseProxy("https://api.openai.com")
	if err != nil {
		return nil, err
	}
	return &OpenaiProxyHandler{
		apiKey: apiKey,
		proxy:  proxy,
	}, nil
}

func (h *OpenaiProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	r.Header.Set("Authorization", "Bearer "+h.apiKey)
	h.proxy.ServeHTTP(w, r)
}

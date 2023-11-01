package handler

import (
	myHttpUtil "github.com/suoaao/affordable-ai/pkg/httputil"
	"github.com/suoaao/affordable-ai/pkg/middleware"
	"net/http"
)

var openaiProxy, _ = myHttpUtil.NewReverseProxy("https://api.openai.com")

func Handler(w http.ResponseWriter, r *http.Request) {
	handler := middleware.RemoveFirstPathElement(openaiProxy)
	handler.ServeHTTP(w, r)
}

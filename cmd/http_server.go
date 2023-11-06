package main

import (
	"fmt"
	github "github.com/suoaao/affordable-ai/api/copilot"
	openai "github.com/suoaao/affordable-ai/api/openai"
	"net/http"
)

func main() {
	http.HandleFunc("/copilot/", github.Handler)
	http.HandleFunc("/openai/", openai.Handler)
	fmt.Println("Starting server at port 8080")
	err := http.ListenAndServe(":8080", nil) // 设置监听的端口
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}

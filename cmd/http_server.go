package main

import (
	"fmt"
	"github.com/suoaao/affordable-copilot/api"
	"net/http"
)

func main() {
	http.HandleFunc("/", handler.Handler)
	fmt.Println("Starting server at port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("ListenAndServe: ", err)
	}
}

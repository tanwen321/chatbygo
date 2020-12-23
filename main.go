package main

import (
	"fmt"
	"net/http"

	"github.com/tanwen321/chatbygo/pkg/setting"
	"github.com/tanwen321/chatbygo/routers"
	"github.com/tanwen321/chatbygo/routers/api/v1"
)

func main() {
	v1.DbCache()
	router := routers.InitRouter()

	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()
}

package main

import (
	"fmt"
	"gympack"
	"net/http"
	"time"
)

func main() {
	time.Sleep(time.Second * 5)
	l := gympack.InitializeLoader()
	c := l.GetConfig()
	r := l.GetRoutes()
	lg := l.GetLogger()

	lg.Info(fmt.Sprintf("starting HTTP server on port %d", c.GetPort()))
	err := http.ListenAndServe(fmt.Sprintf(":%d", c.GetPort()), r)

	if err != nil {
		panic(err)
	}
}

package main

import (
	"engine"
	"net/http"
)

func main() {
	r := engine.New()
	r.GET("/", indexHandler)
	r.GET("/rootca", rootCAHandler)
	r.GET("/subca", subCAHandler)
	r.Start(":443", "./certificate/www.example.com.crt", "./certificate/www.example.com.key.unsecure")
}

func indexHandler(c *engine.Context) {
	c.HTML(http.StatusOK, "<h1>Hello, welcome to goCA</h1>")
}

func rootCAHandler(c *engine.Context) {
	c.HTML(http.StatusOK, "<h1>Hello, welcome to rootCA</h1>")
}

func subCAHandler(c *engine.Context) {
	c.HTML(http.StatusOK, "<h1>Hello, welcome to subCA</h1>")
}

package engine

import (
	"log"
	"net/http"
	"strings"
)

type router struct {
	handlers map[string]HandlerFunc
	//根据HTTP不同类型的请求方法（如 GET、POST等）建立不同的前缀路由树
	roots map[string]*node
}

func newRouter() *router {
	return &router{
		handlers: make(map[string]HandlerFunc),
		roots:    make(map[string]*node),
	}
}

func parsePattern(pattern string) []string {
	pslice := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range pslice {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	log.Printf(" Router : %4s - %s\n", method, pattern)
	key := method + "-" + pattern
	//判断是否存在对应请求方法的前缀路由树，若没有则需要创建
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	//若已经存在相应的前缀路由树，则可以在其中插入路由
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) handle(c *Context) {
	key := c.Req.Method + "-" + c.Req.URL.Path
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", key)
	}
}

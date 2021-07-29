package engine

import (
	"net/http"
)

/*
 * 定义路由映射的处理方法
 */
type HandlerFunc func(*Context)

type Engine struct {
	/*
	 * router 是路由映射表，为 map 类型
	 * key 是用户请求中方法和路径的组合，形如 GET-/
	 * value 是该请求对应的处理方法
	 */
	router *router
}

/*
 * New 是 engine.Engine 的构造函数
 */
func New() *Engine {
	return &Engine{router: newRouter()}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method, pattern, handler)
}

func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

func (engine *Engine) Start(addr string, certFile string, keyFile string) (err error) {
	return http.ListenAndServeTLS(addr, certFile, keyFile, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}

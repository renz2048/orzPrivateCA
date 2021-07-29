package engine

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Path       string
	Method     string
	StatusCode int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
	}
}

func (c *Context) Form(key string) string {
	return c.Req.FormValue(key)
}

// 根据 key 获取 url 中第一个符合的值
func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

// 显示调用 WriteHeader 用于发送错误代码
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

// 设置 http 响应头中字段
//
// 根据键值对
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// 设置响应头 Content-Type 为 text/plain
// 设置响应码为 code
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

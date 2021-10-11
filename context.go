package grape

import (
	"fmt"
	"grape/internal/binding"
	"grape/internal/bytesconv"
	"grape/internal/json"
	"log"
	"math"
	"net/http"
	"sync"
)

// abortIndex表示在abort函数中使用的一个典型值。
const abortIndex int8 = math.MaxInt8 >> 1

// 定义上下文结构体
type Context struct {
	Writer   http.ResponseWriter
	Request  *http.Request
	Method   string
	fullPath string
	handlers HandlersChain
	index    int8
	engine   *Engine
	//这个互斥锁保护key映射
	mu sync.RWMutex
	// key是一个键/值对，专门用于每个请求的上下文。
	Keys map[string]interface{}
}

// 实例化上下文结构体
func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:   w,
		Request:  r,
		Method:   r.Method,
		fullPath: r.URL.Path,
		index:    -1,
	}
}

// 执行下一个handler
func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

// 结束执行handler
func (c *Context) Abort() {
	c.index = abortIndex
}

// Set用于为这个上下文专门存储一个新的键/值对。
//如果c.Keys以前没有被使用过，它也会延迟初始化。
func (c *Context) Set(key string, value interface{}) {
	c.mu.Lock()
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}

	c.Keys[key] = value
	c.mu.Unlock()
}

//返回给定键的值，即:(value, true)。
//如果值不存在，它返回(nil, false)
func (c *Context) Get(key string) (value interface{}, exists bool) {
	c.mu.RLock()
	value, exists = c.Keys[key]
	c.mu.RUnlock()
	return
}

// 写入状态码
func (c *Context) Status(code int) {
	c.Writer.WriteHeader(code)
}

// 跨域处理
func (c *Context) Header(key, value string) {
	if value == "" {
		c.Writer.Header().Del(key)
		return
	}
	c.Writer.Header().Set(key, value)
}

// 设置Header
func (c *Context) SetHeader(key, value string) {
	if key == "" {
		c.Writer.Header().Del(key)
		return
	}
	c.Writer.Header().Set(key, value)
}

// 获取Header
func (c *Context) GetHeader(key string) string {
	return c.Writer.Header().Get(key)
}

// 写入数据
func (c *Context) Write(d []byte) {
	c.Status(http.StatusOK)
	_, _ = c.Writer.Write(d)
}

// 写入字符串
func (c *Context) String(format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain;charset=utf-8")
	if len(values) > 0 {
		_, _ = fmt.Fprintf(c.Writer, format, values...)
		return
	}
	c.Write(bytesconv.StringToBytes(format))
}

// 写入json
func (c *Context) Json(obj interface{}) {
	c.SetHeader("Content-Type", "application/json;charset=utf-8")
	data, err := json.Marshal(obj)
	if err != nil {
		log.Println("序列化失败", err)
		return
	}
	c.Write(data)
}

// 写入html
func (c *Context) HTML(html string) {
	c.SetHeader("Content-Type", "text/html;charset=utf-8")
	c.Write(bytesconv.StringToBytes(html))
}

// 绑定
func (c *Context) ShouldBind(obj interface{}) error {
	//b := binding.Default(c.Request.Method, c.ContentType())
	//return c.ShouldBindWith(obj, binding.Form)
	return binding.Form.Bind(c.Request, obj)
}

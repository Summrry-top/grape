package grape

import (
	"log"
	"net/http"
)

// 定义框架请求处理函数
type HandlerFunc func(c *Context)

// 核心结构体
type Engine struct {
	trees map[string]*node
	*RouterGroup
}

// 实例化核心结构体
func New() *Engine {
	engine := &Engine{}
	engine.RouterGroup = &RouterGroup{engine: engine}
	return engine
}

func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery(), Cors())
	return engine
}

// 添加路由
func (engine *Engine) addRoute(method, path string, handlers HandlersChain) {
	if method == "" {
		panic("method must not be empty")
	}
	if len(path) < 1 || path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}
	if handlers == nil {
		panic("handle must not be nil")
	}
	if engine.trees == nil {
		engine.trees = make(map[string]*node)
	}
	root := engine.trees[method]
	if root == nil {
		root = new(node)
		engine.trees[method] = root
	}
	root.addRoute(path, handlers)
}

// 启动服务
func (engine *Engine) Run(addr string) (err error) {
	log.Println("服务启动", addr)
	return http.ListenAndServe(addr, engine)
}

// Engine 实现ServerHTTP接口(所有的请求都会走到这)
func (engine *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	engine.handleHTTPRequest(c)

}

func (engine *Engine) handleHTTPRequest(c *Context) {
	root := engine.trees[c.Method]
	if root != nil {
		handlers, _ := root.getValue(c.fullPath)
		if handlers != nil {
			c.handlers = handlers
			c.Next()
			return
		}
	}
	c.Json(Err404())
}

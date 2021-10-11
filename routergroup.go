package grape

import (
	"net/http"
	"path"
)

type HandlersChain []HandlerFunc

// 定义路由器组
type RouterGroup struct {
	Handlers     HandlersChain
	absolutePath string
	engine       *Engine
}

// 实例化路由器组
func (group *RouterGroup) Group(relativePath string, handlers ...HandlerFunc) *RouterGroup {
	return &RouterGroup{
		Handlers:     group.combineHandlers(handlers),
		absolutePath: group.calculateAbsolutePath(relativePath),
		engine:       group.engine,
	}
}

// 添加中间件
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.Handlers = append(group.Handlers, middlewares...)
}

// 向路由器组添加路由器
func (group *RouterGroup) Handle(method, part string, handlers HandlersChain) {
	absolutePath := group.calculateAbsolutePath(part)
	handlers = group.combineHandlers(handlers)
	group.engine.addRoute(method, absolutePath, handlers)
}

// 向路由器组添加路由器GET
func (group *RouterGroup) GET(part string, handlers ...HandlerFunc) {
	group.Handle(http.MethodGet, part, handlers)
}

// 向路由器组添加路由器POST
func (group *RouterGroup) POST(part string, handlers ...HandlerFunc) {
	group.Handle(http.MethodPost, part, handlers)
}

func (group *RouterGroup) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	finalSize := len(group.Handlers) + len(handlers)
	mergedHandlers := make([]HandlerFunc, 0, finalSize)
	mergedHandlers = append(mergedHandlers, group.Handlers...)
	return append(mergedHandlers, handlers...)
}

func (group *RouterGroup) calculateAbsolutePath(relativePath string) string {
	if len(relativePath) == 0 {
		return group.absolutePath
	}
	absolutePath := path.Join(group.absolutePath, relativePath)
	appendSlash := lastChar(relativePath) == '/' && lastChar(absolutePath) != '/'
	if appendSlash {
		return absolutePath + "/"
	}
	return absolutePath
}

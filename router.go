package myhttp

import (
	"net/http"
	"reflect"
	"encoding/json"
)

type Handler interface{}

type HandlerChain []Handler

type ErrorResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func ReturnError(status int, err error) *ErrorResponse {
	return &ErrorResponse{
		Status:  status,
		Message: err.Error(),
	}
}

type Router struct {
	root       *Router
	trees      map[string]map[string]HandlerChain
	basePath   string
	middleWare HandlerChain
}

type RouterGroup struct {
	Router
}

func New() *Router {
	return &Router{
		trees: make(map[string]map[string]HandlerChain),
	}
}

// func handler(c *Context) {}
// func handler(c *Context) *ErrorResponse {}
// func handler(in *struct, out *struct) *ErrorResponse {}
// func handler(in *struct, out *struct, c *Context) *ErrorResponse {}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	method := req.Method
	pathHandlers, exist := r.trees[method]
	if !exist {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	handlers, exist := pathHandlers[path]
	if !exist {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	c := newContext(req, w, handlers)
	c.handle()
}

func (c *Context) handle() {
	defer c.response()
	c.processHandler()
}

func (c *Context) processHandler() {
	handlerFunc := c.handlerChain[c.handlerIndex]
	handlerType := reflect.TypeOf(handlerFunc)
	// handler 类型必须是 func
	if handlerType.Kind() != reflect.Func {
		panic("handler type must be func " + handlerType.Name())
	}
	handler := reflect.ValueOf(handlerFunc)

	var res interface{}
	var args []reflect.Value
	var param, response reflect.Value
	var hasResponse bool
	numIn := handlerType.NumIn()
	if numIn == 1 {
		paramType := handlerType.In(0)
		if paramType.Kind() == reflect.TypeOf(&Context{}).Kind() {
			args = append(args, reflect.ValueOf(c))
		} else {
			panic("When handler has one param, it must be ptr to Context" + handlerType.Name())
		}
	}
	if numIn > 1 {
		hasResponse = true
		paramType := handlerType.In(0)
		if paramType.Kind() == reflect.Ptr {
			param = reflect.New(paramType.Elem())
			json.Unmarshal(c.Body(), param.Interface())
			args = append(args, reflect.ValueOf(param.Interface()))
		} else {
			panic("When handler more than one param, first must be ptr to struct" + handlerType.Name())
		}
		responseType := handlerType.In(1)
		if responseType.Kind() == reflect.Ptr {
			response = reflect.New(responseType.Elem())
			args = append(args, reflect.ValueOf(response.Interface()))
		} else {
			panic("When handler more than one param, second must be ptr to struct" + handlerType.Name())
		}
	}
	if numIn == 3 {
		paramType := handlerType.In(2)
		if paramType.Kind() == reflect.TypeOf(&Context{}).Kind() {
			args = append(args, reflect.ValueOf(c))
		} else {
			panic("When handler has three param, third must be ptr to Context" + handlerType.Name())
		}
	}
	errs := handler.Call(args)
	if len(errs) > 0 && !errs[0].IsNil() {
		res = errs[0].Interface()
		c.Json(res)
	} else if hasResponse {
		res = response.Interface()
		c.Json(res)
	} else if !c.hasResponse {
		c.Json(map[string]interface{}{})
	}

	return
}

func (r *Router) Group(path string) *Router {
	router := New()
	router.basePath = r.basePath + path
	copy(router.middleWare, r.middleWare)
	if r.root == nil {
		router.root = r
	} else {
		router.root = r.root
	}
	return router
}

func (r *Router) Use(handler Handler) *Router {
	if r.middleWare == nil {
		r.middleWare = []Handler{}
	}
	r.middleWare = append(r.middleWare, handler)
	return r
}

func (r *Router) GET(path string, handler ...Handler) {
	r.handle(http.MethodGet, path, handler)
}

func (r *Router) POST(path string, handler ...Handler) {
	r.handle(http.MethodPost, path, handler)
}

func (r *Router) handle(method, path string, handlerChain HandlerChain) {
	path = r.basePath + path
	var router *Router

	if r.root != nil {
		router = r.root
	} else {
		router = r
	}

	if router.trees == nil {
		router.trees = make(map[string]map[string]HandlerChain)
	}

	if router.trees[method] == nil {
		router.trees[method] = make(map[string]HandlerChain)
	}

	if router.trees[method][path] == nil {
		router.trees[method][path] = HandlerChain{}
	}

	router.trees[method][path] = append(router.trees[method][path], r.middleWare...)
	router.trees[method][path] = append(router.trees[method][path], handlerChain...)
}

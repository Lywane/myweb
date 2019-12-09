package myhttp

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"reflect"
	"encoding/json"
)

type Handle interface{}

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
	trees map[string]Handle
}

func New() *Router {
	return &Router{
		trees: make(map[string]Handle),
	}
}

// func handler(c *Context) *ErrorResponse {}
// func handler(in *struct, out *struct) *ErrorResponse {}
// func handler(in *struct, out *struct, c *Context) *ErrorResponse {}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	method := req.Method
	var res []byte
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	key := fmt.Sprintf("%s_%s", method, path)
	handleFunc, exist := r.trees[key]
	if !exist {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	handlerType := reflect.TypeOf(handleFunc)
	if handlerType.Kind() != reflect.Func {
		panic("handler type must be func " + handlerType.Name())
	}
	handler := reflect.ValueOf(handleFunc)
	var args []reflect.Value
	var response reflect.Value
	switch handlerType.NumIn() {
	case 1:
		paramType := handlerType.In(0)
		if paramType.Kind() == reflect.TypeOf(&Context{}).Kind() {
			cxt := newContext(req,w)
			args = append(args, reflect.ValueOf(cxt))
		} else {
			panic("When handler has one param, it must be ptr to Context" + handlerType.Name())
		}
		errs := handler.Call(args)
		if !errs[0].IsNil() {
			res, _ = json.Marshal(errs[0].Interface())
			w.Header().Add("Content-Type", "application/json")
			w.Write(res)
		}
	case 2:
		paramType := handlerType.In(0)
		if paramType.Kind() == reflect.Ptr {
			var param reflect.Value
			param = reflect.New(paramType.Elem())
			json.Unmarshal(body, param.Interface())
			args = append(args, reflect.ValueOf(param.Interface()))
		}
		responseType := handlerType.In(1)
		if responseType.Kind() == reflect.Ptr {
			response = reflect.New(responseType.Elem())
			args = append(args, reflect.ValueOf(response.Interface()))
		}
		errs := handler.Call(args)
		if !errs[0].IsNil() {
			res, _ = json.Marshal(errs[0].Interface())
		} else {
			res, _ = json.Marshal(map[string]interface{}{
				"status": 0,
				"data":   response.Interface(),
			})
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(res)

	case 3:
		paramType := handlerType.In(0)
		if paramType.Kind() == reflect.Ptr {
			var param reflect.Value
			param = reflect.New(paramType.Elem())
			json.Unmarshal(body, param.Interface())
			args = append(args, reflect.ValueOf(param.Interface()))
		} else {
			panic("When handler has three param, first must be ptr to struct" + handlerType.Name())
		}
		responseType := handlerType.In(1)
		if responseType.Kind() == reflect.Ptr {
			response = reflect.New(responseType.Elem())
			args = append(args, reflect.ValueOf(response.Interface()))
		} else {
			panic("When handler has three param, second must be ptr to struct" + handlerType.Name())
		}
		if paramType.Kind() == reflect.TypeOf(&Context{}).Kind() {
			cxt := newContext(req,w)
			args = append(args, reflect.ValueOf(cxt))
		} else {
			panic("When handler has three param, third must be ptr to Context" + handlerType.Name())
		}
		errs := handler.Call(args)
		if !errs[0].IsNil() {
			res, _ = json.Marshal(errs[0].Interface())
		} else {
			res, _ = json.Marshal(map[string]interface{}{
				"status": 0,
				"data":   response.Interface(),
			})
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(res)
	}
}

func (r *Router) GET(path string, handle Handle) {
	r.Handle(http.MethodGet, path, handle)
}

func (r *Router) POST(path string, handle Handle) {
	r.Handle(http.MethodPost, path, handle)
}

func (r *Router) DELETE(path string, handle Handle) {
	r.Handle(http.MethodDelete, path, handle)
}

func (r *Router) PUT(path string, handle Handle) {
	r.Handle(http.MethodPut, path, handle)
}

func (r *Router) Handle(method, path string, handle Handle) {
	key := fmt.Sprintf("%s_%s", method, path)
	r.trees[key] = handle
}

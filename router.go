package myrouter

import (
	"net/http"
	"io/ioutil"
	"reflect"
	"encoding/json"
	"fmt"
)

type Handle interface{}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type UrlParam struct {
	data map[string]string
}

func (this *UrlParam) Get(key string) string {
	return this.data[key]
}

func (this *UrlParam) set(key, value string) {
	if this.data == nil {
		this.data = make(map[string]string)
	}
	this.data[key] = value
}

func ReturnError(status int, err error) *Response {
	return &Response{
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

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	method := req.Method
	contentType := req.Header.Get("")
	if contentType != "" && contentType != "application/json" {

	}
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
		panic("handler type must be func but " + handlerType.Name())
	}
	handler := reflect.ValueOf(handleFunc)
	var args []reflect.Value
	var response reflect.Value

	if handlerType.NumIn() != 2 && handlerType.NumIn() != 3 || handlerType.NumOut() != 1 {
		panic("handler must has 2 or 3 param and 1 return")
	}
	paramType := handlerType.In(0)
	var param reflect.Value
	if paramType.Kind() == reflect.Ptr {
		param = reflect.New(paramType.Elem())
		json.Unmarshal(body, param.Interface())
	} else {
		param = reflect.New(paramType)
	}

	responseType := handlerType.In(1)
	if responseType.Kind() == reflect.Ptr {
		response = reflect.New(responseType.Elem())
	} else {
		response = reflect.New(responseType)
	}

	args = []reflect.Value{
		reflect.ValueOf(param.Interface()),
		reflect.ValueOf(response.Interface()),
	}
	if handlerType.NumIn() == 3 {
		fmt.Println(handlerType.In(2).Kind(), reflect.TypeOf(UrlParam{}).Kind())
		if handlerType.In(2).Kind() != reflect.TypeOf(UrlParam{}).Kind() {
			panic("handler the third param must be UrlParam")
		} else {
			urlParam := reflect.New(handlerType.In(2)).Interface().(*UrlParam)
			for k, v := range req.URL.Query() {
				if len(v) > 0 {
					urlParam.set(k, v[0])
				}
			}
			args = append(args, reflect.ValueOf(*urlParam))
		}
	}

	errs := handler.Call(args)
	var res []byte
	if !errs[0].IsNil() {
		res, _ = json.Marshal(errs[0].Interface())
	} else {
		res, _ = json.Marshal(Response{
			Status: 0,
			Data:   response.Interface(),
		})
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(res)
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

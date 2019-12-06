package myhttp

import (
	"net/http"
	"reflect"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Handle interface{}

type ErrorResponse struct {
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

// POST: func handler(in *struct, out *struct)*ErrorResponse{}
// GET:  func handler(urlParam UrlParam, out *struct)*ErrorResponse{}

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
		panic("handler type must be func but " + handlerType.Name())
	}
	handler := reflect.ValueOf(handleFunc)
	var args []reflect.Value
	var response reflect.Value
	switch handlerType.NumIn() {
	case 1:

	case 2:
		paramType := handlerType.In(0)
		if paramType.Kind() == reflect.Ptr {
			var param reflect.Value
			param = reflect.New(paramType.Elem())
			json.Unmarshal(body, param.Interface())
			args = append(args, reflect.ValueOf(param.Interface()))
		} else if paramType.Kind() == reflect.TypeOf(UrlParam{}).Kind() {
			urlParam := &UrlParam{}
			for k, v := range req.URL.Query() {
				if len(v) > 0 {
					urlParam.set(k, v[0])
				}
			}
			args = append(args, reflect.ValueOf(*urlParam))
		} else {
			panic("The first must be a pointer to the structure or UrlParam. " + handlerType.Name())
		}

		responseType := handlerType.In(1)
		if responseType.Kind() == reflect.Ptr {
			response = reflect.New(responseType.Elem())
			args = append(args, reflect.ValueOf(response.Interface()))
		} else {
			panic("The second must be a pointer to the structure. " + handlerType.Name())
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

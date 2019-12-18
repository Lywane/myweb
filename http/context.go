package http

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
)

type Context struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter

	metaData     map[string]interface{}
	handlerIndex int
	handlerChain HandlerChain
	hasReadBody  bool
	body         []byte
	hasResponse  bool

	responseData []byte
	httpStatus   int
	contentType  string
}

func newContext(req *http.Request, w http.ResponseWriter, chain HandlerChain) *Context {
	return &Context{
		Request:        req,
		ResponseWriter: w,
		metaData:       make(map[string]interface{}),
		handlerChain:   chain,
		handlerIndex:   0,
		httpStatus:     http.StatusOK,
	}
}

func (this *Context) Next() {
	this.handlerIndex++
	if this.handlerIndex < len(this.handlerChain) {
		this.processHandler()
	}
}

func (this *Context) GetUrlParam(key string) string {
	values, exist := this.Request.URL.Query()[key]
	if exist && len(values) > 0 {
		return values[0]
	}
	return ""
}

func (this *Context) GetHeader(key string) string {
	return this.Request.Header.Get(key)
}

func (this *Context) GetMetaData(key string) interface{} {
	return this.metaData[key]
}

func (this *Context) SetMetaData(key string, value interface{}) {
	if this.metaData == nil {
		this.metaData = make(map[string]interface{})
	}
	this.metaData[key] = value
}

func (this *Context) Json(data interface{}) {
	res, _ := json.Marshal(map[string]interface{}{
		"status": 0,
		"data":   data,
	})
	this.responseData = res
	this.hasResponse = true
	this.contentType = "application/json;charset=UTF-8"

}

func (this *Context) DieWithHttpStatus(status int) {
	this.httpStatus = status
	this.hasResponse = true
	this.contentType = "text/plain;charset=UTF-8"
}

func (this *Context) response() {
	if this.contentType != "" {
		this.ResponseWriter.Header().Add("Content-Type", this.contentType)
	}
	if this.httpStatus == http.StatusOK {
		this.ResponseWriter.Write(this.responseData)
	} else {
		this.ResponseWriter.WriteHeader(this.httpStatus)
	}

}

func (this *Context) Body() []byte {
	if !this.hasReadBody {
		body, _ := ioutil.ReadAll(this.Request.Body)
		this.hasReadBody = true
		this.body = body
	}
	return this.body
}

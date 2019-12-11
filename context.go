package myhttp

import (
	"net/http"
	"encoding/json"
)

type Context struct {
	Request        *http.Request
	ResponseWriter http.ResponseWriter
	metaData       map[string]interface{}
	next           bool
}

func newContext(req *http.Request, w http.ResponseWriter) *Context {
	return &Context{
		Request:        req,
		ResponseWriter: w,
		metaData:       make(map[string]interface{}),
	}
}

func (this *Context) Next() {
	this.next = true
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
	this.ResponseWriter.Header().Add("Content-Type", "application/json")
	this.ResponseWriter.Write(res)
}

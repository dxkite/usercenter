package handler

import (
	"dxkite.cn/log"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

type callback struct {
	i interface{}
	f ApiCallback
}

type ApiBaseResp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,emitempty"`
}

type HttpContext struct {
	writer  http.ResponseWriter
	request *http.Request
}

const JsonMIMEHeader = "application/json; charset=utf-8"
const JsonSystemError = `{"code":-1, "message":"系统错误，请稍后重试"}`
const JsonContentTypeError = `{"code":-2, "message":"请求内容必须为JSON格式"}`

type ApiCallback func(ctx *HttpContext, input interface{}) (interface{}, int, error)

func NewApiHandler(input interface{}, fun ApiCallback) http.Handler {
	return &callback{input, fun}
}

func (fh *callback) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if t := req.Header.Get("Content-Type"); strings.Index(t, "json") < 0 {
		JsonError(w, JsonContentTypeError, http.StatusBadRequest)
		return
	}
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error(err)
		JsonError(w, JsonSystemError, http.StatusInternalServerError)
		return
	}
	t := reflect.TypeOf(fh.i)
	v := reflect.New(t)
	if err := json.Unmarshal(data, v.Interface()); err != nil {
		log.Error(err)
		JsonError(w, JsonSystemError, http.StatusInternalServerError)
		return
	}
	d, ret, err := fh.f(&HttpContext{
		writer:  w,
		request: req,
	}, v.Elem().Interface())
	WriteData(w, ret, err, d)
}

func WriteData(w http.ResponseWriter, ret int, err error, data interface{}) {
	w.Header().Set("Content-Type", JsonMIMEHeader)
	d := &ApiBaseResp{
		Code:    ret,
		Message: "",
		Data:    data,
	}
	status := http.StatusOK
	if err != nil {
		d.Message = err.Error()
		status = http.StatusBadRequest
	}
	if b, err := json.Marshal(d); err != nil {
		JsonError(w, JsonSystemError, http.StatusInternalServerError)
		log.Error(err)
	} else {
		w.WriteHeader(status)
		_, _ = w.Write(b)
	}
}

func JsonError(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", JsonMIMEHeader)
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	_, _ = fmt.Fprintln(w, error)
}

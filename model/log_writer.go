package middleware

import (
	"bytes"
	"errors"
	"github.com/MadeonZedd/galaxy_log/common"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type OptLog struct {
	User       string `json:"user,omitempty"`
	OpTime     string `json:"op_time，omitempty"`
	OpType     string `json:"op_type,omitempty"`
	Ip         string `json:"ip,omitempty"`
	Company    string `json:"company,omitempty"`
	Permission string `json:"permission,omitempty"`
	Method     string `json:"method,omitempty"`
	Url        string `json:"url,omitempty"`
	Module     string `json:"module,omitempty"`
	Param      string `json:"param,omitempty"`
	Trace      string `json:"trace,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	Rsp        Rsp    `json:"ret,omitempty"`
}

type Option func(*OptLog)

type Rsp struct {
	State     int32       `json:"status,optional"`
	Data      interface{} `json:"data,optional"`
	Message   string      `json:"message,optional"`
	RequestID string      `json:"request_id,optional"`
}

type LogWriter struct {
	http.ResponseWriter
	body *bytes.Buffer
}

func (w LogWriter) Write(p []byte) (int, error) {
	if n, err := w.body.Write(p); err != nil {
		return n, err
	}
	return w.ResponseWriter.Write(p)
}

func (w LogWriter) BaseInfo(r *http.Request) error {
	companyID := r.Header.Get(common.CompanyIDKey)
	userID := r.Header.Get(common.UserIDKey)
	permission := r.Header.Get(common.PermissionKey)

	if companyID == common.EmptyString || userID == common.EmptyString {
		return errors.New("header parameter error")
	}

	op := &OptLog{}
	op.User = userID
	op.Company = companyID
	op.Permission = permission
	op.Method = r.Method
	op.Url = r.RequestURI
	op.Ip = r.RemoteAddr
	op.OpTime = time.Now().Format(common.TimeFormat)
	op.Param = getParam(r)

	return nil
}

func getParam(r *http.Request) string {
	param := common.EmptyString
	if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {
		if r.Body != nil {
			buf := &bytes.Buffer{}
			data := io.TeeReader(r.Body, buf)
			_ = r.Body.Close()
			r.Body = ioutil.NopCloser(buf)
			readAll, err := io.ReadAll(data)
			if err != nil {
				return param
			}
			param = string(readAll)
		}
	}
	if r.Method == http.MethodGet {
		values := r.URL.Query()
		param = common.ToJsonString(values)
	}
	return param
}

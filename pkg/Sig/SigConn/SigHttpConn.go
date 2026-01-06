package SigConn

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
)

// HttpConn	http连接
type HttpConn struct {
	_isHttp bool
	_url    string
	_method string
	_head   map[string]string
	_body   string
	_impl   *http.Client
}

// NewHttpConn 创建http连接
// @param url 连接url
// @param method 连接方法
// @param head http请求首部
// @param body http请求body
// @return *HttpConn http连接
func NewHttpConn(url string, method string, head map[string]string, body string) *HttpConn {
	return &HttpConn{
		_isHttp: true,
		_url:    url,
		_method: method,
		_head:   head,
		_body:   body,
		_impl:   &http.Client{},
	}
}

// NewHttpsConn 创建https连接
// @param url 连接url
// @param method 连接方法
// @param head https请求首部
// @param body https请求body
// @return *HttpConn
func NewHttpsConn(url string, method string, head map[string]string, body string) *HttpConn {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &HttpConn{
		_isHttp: false,
		_url:    url,
		_method: method,
		_head:   head,
		_body:   body,
		_impl:   &http.Client{Transport: tr},
	}
}

// Req http请求
// @receiver c http连接
// @return string http响应
// @return error http请求是否存在错误
func (c *HttpConn) Req() (string, error) {
	// 创建请求
	var req *http.Request = nil
	var err error = nil

	if "" != c._body {
		req, err = http.NewRequest(c._method, c._url, bytes.NewBufferString(c._body))
	} else {
		req, err = http.NewRequest(c._method, c._url, nil)
	}

	if nil != err {
		return "", err
	}

	// 设置请求head
	for k, v := range c._head {
		req.Header.Set(k, v)
	}

	// 发送请求
	resp, err := c._impl.Do(req)
	if nil != err {
		return "", err
	}
	defer resp.Body.Close()

	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if nil != err {
		return "", err
	}
	return string(body), nil
}

// ReqAsync 异步http请求
// @receiver c http连接
// @param cb http请求回调函数
func (c *HttpConn) ReqAsync(cb func(resp string, err error)) {
	go func() {
		resp, err := c.Req()
		cb(resp, err)
	}()
}

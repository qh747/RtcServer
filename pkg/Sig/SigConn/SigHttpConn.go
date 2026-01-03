package SigConn

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
)

func NewSigHttpConn(url string, method string, head map[string]string, body string) *SigHttpConn {
	return &SigHttpConn{
		_url:    url,
		_method: method,
		_head:   head,
		_body:   body,
		_impl:   &http.Client{},
	}
}

func NewSigHttpsConn(url string, method string, head map[string]string, body string) *SigHttpConn {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	return &SigHttpConn{
		_url:    url,
		_method: method,
		_head:   head,
		_body:   body,
		_impl:   &http.Client{Transport: tr},
	}
}

type SigHttpConn struct {
	_url    string
	_method string
	_head   map[string]string
	_body   string
	_impl   *http.Client
}

func (c *SigHttpConn) Req() (string, error) {
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

func (c *SigHttpConn) ReqAsync(cb func(resp string, err error)) {
	go func() {
		resp, err := c.Req()
		cb(resp, err)
	}()
}

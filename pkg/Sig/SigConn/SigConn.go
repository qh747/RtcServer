package SigConn

// const 连接类型
const (
	connRpc   = "rpc"
	connHttp  = "http"
	connHttps = "https"
)

type connType string

// sigConn 媒体服务连接接口
type sigConn interface {
	// 获取连接类型
	GetType() connType

	// 发起请求
	Req() (string, error)
	// 发起异步请求
	ReqAsync(func(string, error))
}

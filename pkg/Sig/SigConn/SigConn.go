package SigConn

// 连接基类

type sigConn interface {
	Req() (string, error)
	ReqAsync(func(string, error))
}

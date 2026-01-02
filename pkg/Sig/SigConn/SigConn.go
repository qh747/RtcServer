package SigConn

type sigConn interface {
	Req() (string, error)
	ReqAsync(func(string, error))
}

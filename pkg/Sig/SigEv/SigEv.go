package SigEv

// dispatcher 事件分发器
var dispatcher *SigEvDispatch

// InitSigEv 		初始化事件分发器
// @return error 	初始化是否存在错误
func InitSigEv() error {
	// 初始化事件分发器
	if nil == dispatcher {
		var err error
		if dispatcher, err = newDispatch("DISPATCH_1"); nil != err {
			return err
		}
	}

	// 注册事件
	dispatcher.Subscribe(EvTopic(EvTopicPush), EvTopicPush+"_1", func(args ...any) {
		// 获取服务连接
		// id := args[0].(string)
		// SigConn.GetSelector().GetAddr(id)

		// reqMsg := args[1].(string)
		// cb := args[2].(func(code int, respMsg string))
	})

	return nil
}

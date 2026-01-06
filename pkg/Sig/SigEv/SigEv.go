package SigEv

// dispatcher 事件分发器
var Dispatcher *EvDispatch

// InitEv 初始化事件分发器
// @return error 初始化是否存在错误
func InitEv() error {
	// 初始化事件分发器
	if nil == Dispatcher {
		var err error
		if Dispatcher, err = newDispatch("DISPATCH_1"); nil != err {
			return err
		}
	}

	// 注册事件
	Dispatcher.Subscribe(EvTopic(EvTopicPush), EvTopicPush+"_1", OnPushEvent)
	return nil
}

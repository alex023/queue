package queue


const (
	_OPENING = 0
	_CLOSED  = 1
)

type ReceiveFunc = func(any interface{})

//异步队列的定义，允许执行消息推送、停止直至消息处理完成、立即停止三种操作
type Queue interface {
	Push(interface{})
	StopGraceful()
	Stop()
}

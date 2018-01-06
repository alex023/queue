package queue

const _BUFFER = 10

//QueueManager 队列管理器
type QueueManager struct {
	queues map[interface{}]Queue
}

//New 创建新的队列管理器实例
func (cm *QueueManager) New() *QueueManager {
	return &QueueManager{
		queues: make(map[interface{}]Queue),
	}
}

//Push 向指定的异步队列推送消息，当队列不存在是，返回 false
func (cm *QueueManager) Push(key interface{}, msg interface{}) (sended bool) {
	queue, founded := cm.queues[key]
	if founded {
		queue.Push(msg)
	}
	return founded
}

//GetOrCreateQueue 创建通道以供使用
func (cm *QueueManager) GetOrCreateQueue(key interface{}, receive ReceiveFunc) (newQueue Queue, newer bool) {
	newQueue, founded := cm.queues[key]
	if !founded {
		newQueue = BoundedQueueMpsc(_BUFFER, receive)
		cm.queues[key] = newQueue
	}
	newer = !founded
	return
}

//Release 释放管理器中指定的队列
func (cm *QueueManager) Release(key interface{}) {
	if queue, founded := cm.queues[key]; founded {
		queue.Stop()
		delete(cm.queues, key)
	}
}

//ReleaseAll 释放管理器中所有的异步队列
func (cm *QueueManager) ReleaseAll() {
	for key, queue := range cm.queues {
		queue.Stop()
		delete(cm.queues, key)
	}
}

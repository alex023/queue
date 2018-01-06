package queue

const _BUFFER = 10

type QueueManager struct {
	queues map[interface{}]Queue
}

func (cm *QueueManager) New() *QueueManager {
	return &QueueManager{
		queues: make(map[interface{}]Queue),
	}
}

func (cm *QueueManager) Push(key interface{}, msg interface{}) (sended bool) {
	queue, founded := cm.queues[key]
	if founded {
		queue.Push(msg)
	}
	return founded
}

//创建通道以供使用
func (cm *QueueManager) GetOrCreateChannel(key interface{},receive ReceiveFunc) (newQueue Queue, newer bool) {
	newQueue, founded := cm.queues[key]
	if !founded {
		newQueue = BoundedQueueCSP(_BUFFER,receive)
		cm.queues[key] = newQueue
	}
	newer = !founded
	return
}

//释放管理器中指定的chan
func (cm *QueueManager) Release(key interface{}) {
	if queue, founded := cm.queues[key]; founded {
		queue.Stop()
		delete(cm.queues, key)
	}
}

func (cm *QueueManager) ReleaseAll() {
	for key, queue := range cm.queues {
		queue.Stop()
		delete(cm.queues, key)
	}
}

package queue

import (
	"log"
	"runtime"
	"runtime/debug"
	"sync/atomic"
)

const (
	_IDLE = iota
	_RUNNING
)

//PosionPill 退出指令定义，用于当前任务队列数据处理完成时退出
type PosionPill struct{}

var poisonPill = &PosionPill{}

type queueMPSC struct {
	closed         int32
	scheduleStatus int32
	userQueue      *MPSC
	invoker        ReceiveFunc
}

//BoundedQueueMpsc 创建一个指定缓冲区大小的队列
func BoundedQueueMpsc(size int, receiver ReceiveFunc) Queue {
	queue := &queueMPSC{
		userQueue: NewMpsc(),
		invoker:   receiver,
	}
	return queue
}

//Push 向消费者推送消息
func (queue *queueMPSC) Push(msg interface{}) {
	if atomic.LoadInt32(&queue.closed) != _CLOSED {
		queue.userQueue.Push(msg)
		go queue.schedule()
	}
}
func (queue *queueMPSC) schedule() {
	if atomic.CompareAndSwapInt32(&queue.scheduleStatus, _IDLE, _RUNNING) {
		queue.run()
		atomic.StoreInt32(&queue.scheduleStatus, _IDLE)
	}
}
func (queue *queueMPSC) run() {
	var msg interface{}
	defer func() {
		if err := recover(); err != nil {
			log.Printf("[queue_mpsc] recovering reason is %+v. More detail:", err)
			log.Println(string(debug.Stack()))
		}
	}()

	i, throughput := 0, 100
	for {
		//长时间消耗，则考虑临时释放系统占用
		if throughput < i {
			i = 0
			runtime.Gosched()
		}
		i++

		if queue.closed == _CLOSED {
			return
		}

		if msg = queue.userQueue.Pop(); msg != nil && msg != poisonPill {
			queue.invoker(msg)
		} else if msg == poisonPill {
			queue.Stop()
			return
		} else {
			return
		}
	}
}

//StopGraceful 当前队列任务执行完毕后，关闭消息队列。该操作后再送入的消息，将会无效。
func (queue *queueMPSC) StopGraceful() {
	if atomic.LoadInt32(&queue.closed) != _CLOSED {
		queue.Push(poisonPill)
	}
}

//Stop 立即关闭任务队列，不论队列中是否还有没被执行的消息
func (queue *queueMPSC) Stop() {
	if atomic.CompareAndSwapInt32(&queue.closed, _OPENING, _CLOSED) {
		queue.userQueue.Empty()
	}
}

package queue

import (
	"context"
	"log"
	"runtime/debug"
	"sync/atomic"
	"runtime"
)

const (
	_IDLE = iota
	_RUNNING
)

type PosionPill struct{}

var posionPill = &PosionPill{}

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

func (queue *queueMPSC) receive(channel chan interface{}, cancelCtx context.Context, receive ReceiveFunc) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("[clock] recovering reason is %+v. More detail:", err)
			log.Println(string(debug.Stack()))
		}
	}()

	for {
		if atomic.LoadInt32(&queue.closed) == _CLOSED {
			break
		}
		data := queue.userQueue.Pop()
		if data != nil && data != posionPill {
			queue.invoker(data)
		} else if data == posionPill {
			break
		}
	}
	//release resource
	queue.userQueue.Empty()
	queue.invoker = nil
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
		defer func() {
			if err := recover(); err != nil {
				log.Printf("[queue_mpsc] recovering reason is %+v. More detail:", err)
				log.Println(string(debug.Stack()))
			}
		}()
	}()

	i,throughput:=0,300
	for  {
		//长时间消耗，则考虑临时释放系统占用
		i++
		if throughput<i{
			i=0
			runtime.Gosched()
		}

		if atomic.LoadInt32(&queue.closed) == _CLOSED {
			return
		}
		if msg = queue.userQueue.Pop(); msg != nil && msg != posionPill {
			queue.invoker(msg)
		} else {
			return
		}
	}
}

//StopGraceful 当队列任务执行完毕后，关闭消息队列
func (queue *queueMPSC) StopGraceful() {
	if atomic.LoadInt32(&queue.closed) != _CLOSED {
		queue.Push(posionPill)
	}
}

//Stop 立即关闭任务队列，不论队列中是否还有没被执行的消息
func (queue *queueMPSC) Stop() {
	if atomic.CompareAndSwapInt32(&queue.closed, _OPENING, _CLOSED) {
		queue.userQueue.Empty()
	}
}

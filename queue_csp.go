package queue

import (
	"context"
	"log"
	"runtime/debug"
	"sync/atomic"
	"runtime"
)

type boundQueue struct {
	closed     int32
	channel    chan interface{}
	cancelFunc context.CancelFunc
}

//创建一个指定缓冲区大小的队列
func BoundedQueueCSP(size int, invoker ReceiveFunc) Queue {
	cancelContext, cancelFunc := context.WithCancel(context.Background())
	queue := &boundQueue{
		channel:    make(chan interface{}, size),
		cancelFunc: cancelFunc,
	}
	go queue.receive(queue.channel, cancelContext, invoker)
	return queue
}

func (queue *boundQueue) receive(channel chan interface{}, cancelCtx context.Context, invoker ReceiveFunc) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("[queue_csp] recovering reason is %+v. More detail:", err)
			log.Println(string(debug.Stack()))
		}
	}()
	var (
		stoppedByChannel  = false
		stoppedByContext = false
	)

	i,throughput:=0,300
	for !stoppedByChannel && !stoppedByContext {
		//长时间消耗，则考虑临时释放系统占用
		i++
		if throughput<i{
			i=0
			runtime.Gosched()
		}
		//Notice:由于 select 的语法特点，即使 cancelCtx.Done的消息已经发出，但只要channel有值，则很可能无法立即退出循环
		select {
		case msg := <-channel:
			if msg != nil {
				invoker(msg)
			} else {
				stoppedByChannel = true
			}
		case <-cancelCtx.Done():
			stoppedByContext = true
		}
	}
	//release resource
	if stoppedByChannel {
		if queue.cancelFunc != nil {
			queue.cancelFunc()
		}
	} else if stoppedByContext {
		for range channel {
			//do nothing for release channel
		}
	}
}

//向消费者推送消息
func (queue *boundQueue) Push(msg interface{}) {
	if atomic.LoadInt32(&queue.closed) != _CLOSED {
		queue.channel <- msg
	}
}

//当队列任务执行完毕后，关闭消息队列
func (queue *boundQueue) StopGraceful() {
	if atomic.CompareAndSwapInt32(&queue.closed, _OPENING, _CLOSED) {
		close(queue.channel)
	}
}

//立即关闭任务队列，不论队列中是否还有没被执行的消息
func (queue *boundQueue) Stop() {
	if atomic.CompareAndSwapInt32(&queue.closed, _OPENING, _CLOSED) {
		if queue.cancelFunc != nil {
			queue.cancelFunc()
		}
		close(queue.channel)
	}
}

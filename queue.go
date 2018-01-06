package queue

import (
	"log"
	"runtime/debug"
	"sync/atomic"
)

const (
	_OPENING = 0
	_CLOSED  = 1
)


type ReceiveFunc=func(any interface{})
type Queue interface{
	Put(interface{})
	Close()
}

type boundQueue struct {
	closed  int32
	channel chan interface{}
}
//创建一个指定缓冲区大小的队列
func BoundedQueue(size int, receiver ReceiveFunc) Queue {
	queue := &boundQueue{
		channel: make(chan interface{}, size),
	}
	go receive(queue.channel, receiver)
	return queue
}

func receive(channel chan interface{}, receive ReceiveFunc) {
	defer func() {
		if err := recover(); err != nil {
			log.Printf("[clock] recovering reason is %+v. More detail:", err)
			log.Println(string(debug.Stack()))
		}
	}()
	for msg := range channel {
		receive(msg)
	}
}

//向消费者推送消息
func (queue *boundQueue) Put(msg interface{}) {
	if atomic.LoadInt32(&queue.closed) != _CLOSED {
		queue.channel <- msg
	}
}

//关闭消息队列
func (queue *boundQueue) Close() {
	if atomic.CompareAndSwapInt32(&queue.closed, _OPENING, _CLOSED) {
		close(queue.channel)
	}
}

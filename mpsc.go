package queue

import (
	"sync/atomic"
	"unsafe"
)

type node struct {
	next *node
	data interface{}
}

// MPSC 基于 go 的multi-produce single-consumer的数据结构
type MPSC struct {
	head, tail *node
}

func NewMpsc() *MPSC {
	q := &MPSC{}
	stub := &node{}
	q.head = stub
	q.tail = stub
	return q
}

// Push 添加一条新的消息到队列的末尾
func (mpsc *MPSC) Push(x interface{}) {
	n := &node{data: x}
	prev := (*node)(atomic.SwapPointer((*unsafe.Pointer)(unsafe.Pointer(&mpsc.head)), unsafe.Pointer(n)))
	prev.next = n
}

// Pop 从队列中提取一条消息交付给消费者
func (mpsc *MPSC) Pop() interface{} {
	tail := mpsc.tail
	next := (*node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&tail.next)))) // acquire
	if next != nil {
		mpsc.tail = next
		v := next.data
		next.data = nil
		return v
	}
	return nil
}

// Empty 清空队列
func (mpsc *MPSC) Empty() bool {
	tail := mpsc.tail
	next := (*node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&tail.next))))
	return next == nil
}

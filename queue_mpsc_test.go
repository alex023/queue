package queue

import (
	"testing"
	"time"
)

func TestBoundedQueueMscp(t *testing.T) {
	var (
		count int
		num=20
	)
	receiver:=func(msg interface{}){
		count++
	}
	queue := BoundedQueueMpsc(num, receiver)

	for i := 0; i < num; i++ {
		queue.Push(struct{}{})
	}
	time.Sleep(time.Second*1)
	if count!=num{
		t.Error("消息接受存在问题")
	}
}

func TestQueueMscp_StopGraceful(t *testing.T) {
	var (
		count int
		num=20
	)
	receiver:=func(msg interface{}){
		count++
		time.Sleep(time.Millisecond*40)
	}
	queue := BoundedQueueMpsc(num, receiver)
	for i := 0; i < num; i++ {
		queue.Push(struct{}{})
	}
	queue.StopGraceful()
	time.Sleep(time.Second*2)
	if num!=count{
		t.Errorf("QueueMpsc_StopGraceful存在问题,应该执行%v,实际执行%2d\n",num,count)
	}
}
func TestQueueMscp_Stop(t *testing.T) {
	var (
		count int
		num=20
	)
	receiver:=func(msg interface{}){
		count++
		time.Sleep(time.Millisecond*50)
	}
	queue := BoundedQueueMpsc(num, receiver)
	for i := 0; i < num; i++ {
		queue.Push(struct{}{})
	}
	queue.Stop()
	time.Sleep(time.Second*2)
	if num==count{
		t.Error("QueueMpsc_StopGraceful存在问题")
	}
}
package queue

import (
	"time"
	"testing"
)


func TestBoundedQueueCSP(t *testing.T) {
	var (
		count int
		num=20
	)
	receiver:=func(msg interface{}){
		count++
	}
	queue := BoundedQueueCSP(num, receiver)

	for i := 0; i < num; i++ {
		queue.Push(struct{}{})
	}
	time.Sleep(time.Second*1)
	if count!=num{
		t.Error("消息接受存在问题")
	}
}

func TestQueueCSP_StopGraceful(t *testing.T) {
	var (
		count int
		num=20
	)
	receiver:=func(msg interface{}){
		count++
		time.Sleep(time.Millisecond*40)
	}
	queue := BoundedQueueCSP(num, receiver)
	for i := 0; i < num; i++ {
		queue.Push(struct{}{})
	}
	queue.StopGraceful()
	time.Sleep(time.Second*2)
	if num!=count{
		t.Errorf("QueueCSP_StopGraceful存在问题,应该执行%v,实际执行%2d\n",num,count)
	}
}
func TestQueueCSP_Stop(t *testing.T) {
	var (
		count int
		num=20
	)
	receiver:=func(msg interface{}){
		count++
		time.Sleep(time.Millisecond*50)
	}
	queue := BoundedQueueCSP(num, receiver)
	for i := 0; i < num; i++ {
		queue.Push(struct{}{})
	}
	queue.Stop()
	queue.Stop()
	time.Sleep(time.Second*2)
	if num==count{
		t.Errorf("QueueCSP_StopGraceful存在问题,应该执行不到%v,实际执行%2d\n",num,count)

	}
}


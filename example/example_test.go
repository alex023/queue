package example

import (
	"fmt"
	"time"
	"github.com/alex023/queue"
)
//消费者
func receiver(msg interface{}) {
	fmt.Println(msg)
	time.Sleep(time.Millisecond * 210)
}

func ExampleBoundQueueCSP_StopGracefull() {
	var que = queue.BoundedQueueCSP(20, receiver)
	for i := 0; i < 10; i++ {
		que.Push(fmt.Sprintf("消息%2d", i))
	}
	time.Sleep(time.Second * 1)
	que.StopGraceful()
	time.Sleep(time.Second * 3)
	//Output:
	//
	//消息 0
	//消息 1
	//消息 2
	//消息 3
	//消息 4
	//消息 5
	//消息 6
	//消息 7
	//消息 8
	//消息 9
}

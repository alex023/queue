package queue

import (
	"fmt"
	"time"
)

//消费者
func receiver(msg interface{}) {
	time.Sleep(time.Millisecond*50)
	fmt.Println(msg)
}

func ExampleBoundedQueue() {
	var que = BoundedQueue(5, receiver)
	for i := 0; i < 10; i++ {
		que.Put(fmt.Sprintf("消息%2d", i))
	}
	que.Close()
	time.Sleep(time.Second * 5)

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

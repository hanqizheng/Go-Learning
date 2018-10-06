# 初识 Go多线程编程

其实感觉自己学node单线程异步仿佛就在昨天，现在摇身一变要来体验一把go的多线程。

起多个线程，他们同时在做不同的事情。

## goroutine
go的线程叫做`goroutine`，是一种非常轻量的。创建一个 goroutine 非常简单，只需要把 go 关键字放在函数调用语句前。

just like this
```go
go doSomething()
```

## channel
试着想象一下，起了多个线程，有可能是为了协作完成某个大的任务。但是现在各个线程如此孤立，并不能达到协作的目的。所以他们之间需要交流，一个人闷声干活是很累的。

channel就诞生了，管道，信道，whaterver道，他就是一个goroutine之间相互通信的道.....

创建一个channel也是非常简单的
```go
myFirstChannel := make(chan string)
```
通过管道发送，接收数据也很简单

```go
myFirstChannel <-"hello" // Send
myVariable := <- myFirstChannel // Receive
```
### buffered and unBuffered channel
我们刚才创建的是每次只能发送一个消息的unBufferd channel，当然在繁忙的任务中这样的传递效率确实太低了，也不太好，所以我们可以创建一下Buffer channel，他可以容纳若干个信息在同时（我的语序为什么变成了英语？？）

创建方法如下
```go
bufferedChan := make(chan string, 3)
```

### Channel Blocking
channel会在几种情况下阻塞当前goroutine

#### Blocking on Send
当一个goroutine 想channel发送消息，这个goroutine就阻塞了。**直到另一个goroutine把这个数据拿走**

just like this
![](https://raw.githubusercontent.com/studygolang/gctt-images/master/Learning-Go-s-Concurrency-Through-Illustrations/blocking-on-send.jpeg)

#### Blocking on Receive
与上述同理，一个goroutine会阻塞接收channel中的数据
just like this

![](https://raw.githubusercontent.com/studygolang/gctt-images/master/Learning-Go-s-Concurrency-Through-Illustrations/blocking-on-receive.jpeg)


## 一个简单的例子
大概接触了一下，我们可以用goroutine实现一个判断两棵二叉查找树是否遍历顺序相等
代码给出，这是go指南的一道[练习题](https://tour.go-zh.org/concurrency/7)
```go
package main
import (
	"fmt"
	"golang.org/x/tour/tree"
)

// Walk to 发送value, 结束后关闭channel
func Walk(t *tree.Tree, ch chan int)  {
	sendValue(t, ch)
	close(ch)
}

func sendValue(t *tree.Tree, ch chan int) {
	if t != nil {
		sendValue(t.Left, ch)
		ch <- t.Value
		sendValue(t.Right, ch)
	}
}

// Same to 判断两个树是否顺序一直
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go Walk(t1, ch1)
	go Walk(t2, ch2)
	
	for i := range ch1 {
		if i != <- ch2 {
			return false
		}
	}

	return true
}
func main() {
	var ch = make(chan int)
	go Walk(tree.New(1), ch)

	for v := range ch {
		fmt.Println(v)
	}

	fmt.Println(Same(tree.New(1), tree.New(1)))
	fmt.Println(Same(tree.New(1), tree.New(2)))
}

```
也许刚看题会被二叉树啥的吓到（什么？你没被吓到？我被吓到了，掀桌.jpg），但他作为一道练习题，已经帮我们把构造树之类的工作都做了，我们只需要关注如何判断两颗树是否相等就好。

这段代码的核心就是如何去判断两树相等顺序，其实就是如何起两个goroutine将两棵树的节点值通过channel发送然后进行比较。

可以看到例子中的`func sendValue(t *tree.Tree, ch chan int)`函数就是发送数据的，发送顺序是中序遍历，所以这道题比较的是两棵二叉树中序遍历结果是否相同

`func Same(t1, t2 *tree.Tree) bool`则是起了两个线程去判断是否相同的函数

本身不难，也许大家对`tree.New(1)`这种语句比较迷(别说你不迷),这个是事先规定好的

```
函数 tree.New(k) 用于构造一个随机结构的已排序二叉查找树，它保存了值 k, 2k, 3k, ..., 10k。

创建一个新的信道 ch 并且对其进行步进：

go Walk(tree.New(1), ch)
然后从信道中读取并打印 10 个值。应当是数字 1, 2, 3, ..., 10。
```

## 参考
[图解 Go 并发编程](https://studygolang.com/articles/13875)
[go 指南](https://tour.go-zh.org/concurrency/1)
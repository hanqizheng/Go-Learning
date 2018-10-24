# goroutine Again

继上次简单的了解了一下Go中的goroutine之后，这次稍微深入一点，来了解一下goroutine的运行机制以及并发编程相关的内容。

## 进程和线程
> 无论是并发并行，还是goroutine，都需要对进程和线程有清晰的认识。

**进程**：当运行一个应用程序(如一个 IDE 或者编辑器)的时候,操作系统会为这个应用程序启动一个进程。可以将这个进程看作一个包含了应用程序在运行中需要用到和维护的各种资源的容器。

可以把进程理解为一个一个的盒子，盒子里放置了某个应用程序运行所需要的资源。这些资源有很多种，比如分配给这个应用程序的内存、操作系统相关的句柄、还有其进程相关的线程等。所有这些东西都放在了进程中。而真正去执行这个应用程序相关代码的是进程中的一部分，线程。

**线程**：一个线程是一个执行空间,这个空间会被操作系统调度来运行函数中所写的代码。

简单的理解就是操作系统要执行某个程序的逻辑（代码），就会调度对应的线程。一个进程至少有一个线程，它叫做`主线程`，主线程结束的时候就意味着这个进程也结束了。当然，进程都已经结束，更别说其他的线程了。所以，主线程运行完毕结束的时候，就算其他线程还在工作，也会随之结束。

![](/sources/images/ss.png)

这张图是进程和线程的关系图，可以根据图片理解一下进程和线程到底是什么。


## 并行和并发

> **并行**是让不同的代码片段同时在不同的物理处理器上执行。并行的关键是同时做很多事情。
> 
> **并发**是指同时管理很多事情,这些事情可能只做了一半就被暂停去做别的事情了。

在很多情况下,并发的效果比并行好,因为操作系统和硬件的总资源一般很少,但能支持系统同时做很多事情。

## goroutine 进一步理解

操作系统会在物理处理器上（比如CPU），调度线程来运行程序。

相对的

Go会在`逻辑处理器`上，调度每个goroutine来运行程序。

但是逻辑处理器也不是能独自运行的，它是绑定到整个操作系统特别的线程`操作系统线程`上去运行。

所以可以可能到其实goroutine并不是等同于线程的概念，有的地方把goroutine翻译成了`协程`

### goroutine从创建到运行

创建一个goroutine 

--->

放到`全局运行队列`中

> 这里的全局运行队列是专门存放goroutine的队列，所有刚创建的goroutine都会放在这个队列中

--->

将`全局运行队列`中的gouroutine分配到某个`逻辑处理器`上，并放到该逻辑处理器对应的`本地运行队列`等待运行

> 本地运行队列是逻辑处理器将要处理的goroutine排成的队列

--->

goroutine会在`本地运行队列`中等待知道被分配的逻辑处理器调用。

给出一张图，就可以搭配描述直观的看出来, 图的左边就是上述描述的过程


![](/sources/images/logicController.png)

图的右边则是描述了 **当一个goroutine执行的任务需要阻塞等待时**，那么调度器会将当前线程分离出去，继续执行需要阻塞等待的这个goroutine。

而因为线程的分离，**会创建一个新的线程**来运行原逻辑处理器，和相关的其他goroutine。

当阻塞任务执行完毕，对应的goroutine会被放回到原来的本地运行队列，而单独执行他的那个线程，则会被保存好，等待将来使用。

### 并行执行
如果希望让 `goroutine 并行`,必须**使用多于一个逻辑处理器**。

当有多个逻辑处理器时,调度器会将 goroutine 平等分配到每个逻辑处理器上。

这会让 goroutine 在不同的线程上运行。不过要想真的实现并行的效果,用户需要让自己的程序运行在有多个物理处理器的机器上。

否则,哪怕 Go 语言运行时使用多个线程,goroutine 依然会在同一个物理处理器上并发运行,达不到并行的效果。

## 竞争状态

当某个程序拥有多个goroutine，但是，**各个goroutine之间却没有实现同步**，就会发生一个共享资源同时被两个甚至更多个的goroutine同时使用，这样，这几个goroutine就会形成**竞争状态**。

为了解决竞争状态的发生，Go提供了几种办法

### sync.WaitGroup

先要明确，**sync.WaitGroup是一个struct**。

这个类型里有三个方法，分别是`Add()`, `Done()`, `Wait()`

声明一个变量
```go
var wg sync.WaitGroup
```

sync.WaitGroup这个结构体内有一个成员变量，用来计数，当声明这个sync.WaitGroup类型的变量的时候，内部这个值将被初始化为0.

Add() 和 Done()方法就是用来改变内部这个值的大小。

Add()用来增加这个值， Done()用来减小这个值。

Wait()方法是用来阻塞当前线程的，他会检查这个计数值的大小，
- 当这个值 == 0 时，Wait()会立刻返回，从而当前被阻塞的线程也就被释放继续执行下去
- 当这个值 > 0 时，Wait()就继续阻塞当前线程

```go
package main
import (
  "fmt"
  "sync"
  "sync/atomic"
  "time"
)
var (
  // shutdown 是通知正在执行的 goroutine 停止工作的标志
  shutdown int64
  // wg 用来等待程序结束
  wg sync.WaitGroup
)
// main 是所有 Go 程序的入口
func main() {
  // 计数加 2,表示要等待两个 goroutine
  wg.Add(2)
  // 创建两个 goroutine
  go doWork("A")
  go doWork("B")
  // 给定 goroutine 执行的时间
  time.Sleep(1 * time.Second)
  // 该停止工作了,安全地设置 shutdown 标志
  fmt.Println("Shutdown Now")
  atomic.StoreInt64(&shutdown, 1)37

  // 等待 goroutine 结束
  wg.Wait()
}
// doWork 用来模拟执行工作的 goroutine,
// 检测之前的 shutdown 标志来决定是否提前终止
func doWork(name string) {
  // 在函数退出时调用 Done 来通知 main 函数工作已经完成
  defer wg.Done()
  for {
    fmt.Printf("Doing %s Work\n", name)
    time.Sleep(250 * time.Millisecond)
    // 要停止工作了吗?
    if atomic.LoadInt64(&shutdown) == 1 {
      fmt.Printf("Shutting %s Down\n", name)
      break
    }
  }
}
```


### 互斥锁

另一种同步访问共享资源的方式是使用互斥锁( mutex )。互斥锁这个名字来自互斥(mutualexclusion)的概念。互斥锁用于在代码上创建一个临界区,保证同一时间只有一个 goroutine 可以执行这个临界区代码。

对应的是Go中的`sync.Mutex`

```go
var mutex sync.Mutex
```
**sync.Mutex也是一个struct**, 他有有两个方法`Lock()`, `UnLock()`

这两个函数的用法就是一头一尾，将一段代码包住，被包住的代码就叫做临界区，临界区被互斥锁锁住之后，资源就只能被一个goroutine使用，这样就避免了资源竞争。


### Channel
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
#### buffered and unBuffered channel
我们刚才创建的是每次只能发送一个消息的unBufferd channel，当然在繁忙的任务中这样的传递效率确实太低了，也不太好，所以我们可以创建一下Buffer channel，他可以容纳若干个信息在同时（我的语序为什么变成了英语？？）

创建方法如下
```go
bufferedChan := make(chan string, 3)
```

#### Channel Blocking
channel会在几种情况下阻塞当前goroutine

##### Blocking on Send
当一个goroutine 想channel发送消息，这个goroutine就阻塞了。**直到另一个goroutine把这个数据拿走**

just like this
![](https://raw.githubusercontent.com/studygolang/gctt-images/master/Learning-Go-s-Concurrency-Through-Illustrations/blocking-on-send.jpeg)

##### Blocking on Receive
与上述同理，一个goroutine会阻塞接收channel中的数据
just like this

![](https://raw.githubusercontent.com/studygolang/gctt-images/master/Learning-Go-s-Concurrency-Through-Illustrations/blocking-on-receive.jpeg)




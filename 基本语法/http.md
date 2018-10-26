# 初步了解Go中http的具体执行流程

在Go中，http包可以迅速的，方便的创建一个http服务器，然后我们就可以做很多很多有趣的事情，但是方便的背后其实有着很多相对复杂的实现，拿到黑色的盒子不光要学会如何使用，这篇文章将一起打开黑盒，看看内部都有些什么。


## 创建一个简单的http服务器

```go
import (
  "fmt"
  "net/http"
)

func sayHello() {
  fmt.Println("hello")
}

func main() {
  http.HandleFunc("/", sayHello) //设置访问的路由
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
```
正如上述一段代码，可以看到Go创建一个http服务器非常的简单，接下来，就拆开这段代码具体是怎么创建的一个服务器的。

## http.ListenAndServer()
可以看到`http.ListenAndServer()`一写，一个服务器就跑起来了，那么`http.ListenAndServer()`它具体做了什么呢？


> http.ListenAndServer(port， handler)有两个参数。
> 
> 会在接下来具体说明两个参数到底是干嘛用的。

- 首先，先初始化了一个server对象，以便后续使用。

- 然后，调用了`net.Listen("tcp", addr)`,其实就是在底层用TCP搭建起来一个服务器，然后监听这我们设置的对应端口，这里的对应端口就是`http.ListenAndServer()`的**第一个参数**。
- 调用`srv.Server(l net.Listener)`，他的作用就是负责接收客户端的请求

- 跑起一个for { }持续接受请求

- 创建一个Conn对象`c ：= srv.NewConn()`，这个Conn就是用户每次请求需要建立的链接，然后把刚才获取到的请求数据都给Conn。Conn保存了这次的请求数据，会在后续将数据传给Handler

- 独自起一个goroutine,`go c.serve()`。这么写可以做到每个用户的每一次请求都是一个新的gorouine，相互不影响。


这里把底层代码给出，可以结合上述解释阅读

```go
func (srv *Server) Serve(l net.Listener) error {
	defer l.Close()
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		rw, e := l.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Printf("http: Accept error: %v; retrying in %v", e, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return e
		}
		tempDelay = 0
		c, err := srv.newConn(rw)
		if err != nil {
			continue
		}
		go c.serve()
	}
}
```

## Conn.server()
上面提到了每一次请求都会创建一个goroutine去单独执行`c.server()`。那么它到底干了什么呢？

- 首先，会解析request()，`c.readRequest()`,获取相应的Handler, `handler := c.serve.Handler`，**这个handler就是http.ListenAndServe()的第二个参数**
- 这个handler就是`http.ListenAndServer()`的第二个参数，第二个参数是代表一个路由，   也就是给出具体的URL可以去执行对应的操作。

  而在最开始我们还写了一句`http.HandleFunc("/", sayHello)`，这就是在注册路由为`"/"`的路由规则，当URL为`/`，就会调用对应的`sayHello`方法。

- 最终因为我们给出的参数是`nil`所以handler会变成默认的`handler = DefaultServerMux`,然后去调用`sayHello`方法。


## ServerMux
上面的过程在最后，调用了默认的`DefaultServerMux`然后最终执行了`sayHello`方法。其实在这里也是有很多步骤的，并非很简单就实现了。

`ServeMux`本质上是一个 HTTP 请求路由器（或者叫多路复用器，Multiplexor）。它把收到的请求与一组预先定义的 URL 路径列表做对比，然后在匹配到路径的时候调用关联的处理器（Handler）。下面给出`ServeMux`的定义

```go
type ServeMux struct {
  mu sync.RWMutex   //锁，由于请求涉及到并发处理，因此这里需要一个锁机制
  m  map[string]muxEntry  // 路由规则，一个string对应一个mux实体，这里的string就是注册的路由表达式
  hosts bool // 是否在任意的规则中带有host信息
}
```

第二个就是用来存放路由规则的一个map，``muxEntry`就是路由对应的具体操作，然后给出`muxEntry`的具体定义

```go
type muxEntry struct {
  explicit bool   // 是否精确匹配
  h        Handler // 这个路由表达式对应哪个handler
  pattern  string  //匹配字符串
}
```

`muxEntry`里面有一个`Handler`类型的变量，就进一步说明了。

```go
type Handler interface {
  ServeHTTP(ResponseWriter, *Request)  // 路由实现器
}
```
可以看到Handler是一个接口类型，所以任何满足了http.Handler接口的对象都可作为一个处理器，通俗的说就是任何实现了`ServeHttp()`方法的都可以作为某个路由的Handler。

到这里也许会有一个问题，就是最开始`sayHello()`并没有实现`ServeHttp()`为什么也可以当Handler呢？

### HandlerFunc
这时候就要引出另一个http包中的类型，`HandlerFunc`。我们定义的函数`sayhelloName`就是这个`HandlerFunc`调用之后的结果，这个类型默认就实现了`ServeHTTP`这个接口，即我们调用了`HandlerFunc(f)`,强制类型转换`f`成为`HandlerFunc`类型，这样`f`就拥有了`ServeHTTP`方法。

```go
type HandlerFunc func(ResponseWriter, *Request)

// ServeHTTP calls f(w, r)
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) {
  f(w, r)
}
```

**简单的说就是**，我们通过`HandlerFunc(f)`将`sayHello`强制类型转换成`HandlerFunc`，这样就实现了`ServeHttp`，就变成了一个`Handler`


**最终会根据专门存放当时注册路由规则的那个map中找到对应的handler执行操作**
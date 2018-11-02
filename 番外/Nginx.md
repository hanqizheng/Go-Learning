# Nginx

Nginx是一款轻量级的Web服务器/反向代理服务器及电子邮件（IMAP/POP3）代理服务器，并在一个BSD-like协议下发行。由俄罗斯的程序设计师IgorSysoev所开发，供俄国大型的入口网站及搜索引擎Rambler（俄文：Рамблер）使用。其特点是占有内存少，并发能力强，事实上nginx的并发能力确实在同类型的网页服务器中表现较好。

## 代理和反向代理

通过两个简单的例子，来理解一下代理和反向代理的概念及区别

### 代理
A同学在大众创业、万众创新的大时代背景下开启他的创业之路，目前他遇到的最大的一个问题就是启动资金，于是他决定去找马云爸爸借钱，可想而知，最后碰一鼻子灰回来了，情急之下，他想到一个办法，找关系开后门，经过一番消息打探，原来A同学的大学老师王老师是马云的同学，于是A同学找到王老师，托王老师帮忙去马云那借500万过来，当然最后事成了。不过马云并不知道这钱是A同学借的，马云是借给王老师的，最后由王老师转交给A同学。这里的王老师在这个过程中扮演了一个非常关键的角色，就是`代理`，也可以说是正向代理，王老师代替A同学办这件事，这个过程中，真正借钱的人是谁，马云是不知道的，这点非常关键。

我们常说的代理也就是只正向代理，正向代理的过程，它隐藏了真实的请求客户端，服务端不知道真实的客户端是谁，客户端请求的服务都被代理服务器代替来请求，某些科学上网工具扮演的就是典型的正向代理角色。用浏览器访问`www.Google.com` 时，被残忍的block，于是你可以在国外搭建一台代理服务器，让代理帮我去请求google.com，代理把请求返回的相应结构再返回给我。

![](https://pic1.zhimg.com/80/v2-07ededff1d415c1fa2db3fd89378eda0_hd.jpg)
### 反向代理
大家都有过这样的经历，拨打10086客服电话，可能一个地区的10086客服有几个或者几十个，你永远都不需要关心在电话那头的是哪一个，叫什么，男的，还是女的，漂亮的还是帅气的，你都不关心，你关心的是你的问题能不能得到专业的解答，你只需要拨通了10086的总机号码，电话那头总会有人会回答你，只是有时慢有时快而已。那么这里的10086总机号码就是我们说的`反向代理`。客户不知道真正提供服务人的是谁。

反向代理隐藏了真实的服务端，当我们请求`www.baidu.com`的时候，就像拨打10086一样，背后可能有成千上万台服务器为我们服务，但具体是哪一台，你不知道，也不需要知道，你只需要知道反向代理服务器是谁就好了，`www.baidu.com` 就是我们的反向代理服务器，反向代理服务器会帮我们把请求转发到真实的服务器那里去。Nginx就是性能非常好的反向代理服务器，用来做负载均衡。

![](https://pic3.zhimg.com/80/v2-816f7595d80b7ef36bf958764a873cba_hd.jpg)


## Nginx简易操作

下载安装有一堆教程，这里就不说了

### 启动

```
nginx
```
### 停止

硬停止， 不在乎当前是否还有进程执行任务，直接推出
```
nginx -s stop
```

软停止， 等当前进程执行的任务完成后再退出，如果期间有新的任务或者请求，会移架到别处运行而不会在当前将要停止的这个nginx服务上运行。

```
nginx -s quit
```

### 重启

重启配置文件
```
nginx -s reload
```

重启日志文件
```
nginx -s reopen
```

## 配置文件

Nginx 配置文件主要分成四部分：
- main（全局设置）
- http（HTTP 的通用设置）
- server（虚拟主机设置）
- location（匹配 URL 路径）
- 还有一些其他的配置段，如 event，upstream 等。

### 通用配置

先给出一段，这是我刚安装的nginx，nginx.conf文件里的通用配置部分，给出这个例子就是为了展示一下通用配置具体长什么样子。
```nginx
user nginx;
worker_processes auto;
error_log /var/log/nginx/error.log;
pid /run/nginx.pid;
```

**通用配置不止这几条，下面具体说下**

- `user nginx`
  
  指定运行 nginx workre 进程的用户和组

- `worker_rlimit_nofile`
  
  指定所有 worker 进程能够打开的最大文件数

- `worker_cpu_affinity`

  设置 worker 进程的 CPU 粘性，以避免进程在 CPU 间切换带来的性能消耗。如 worker_cpu_affinity 0001 0010 0100 1000;（四核）

- `worker_processes auto`

worker 工作进程的个数，这个值可以设置为与 CPU 数量相同，如果开启了 SSL 和 Gzip，那么可以适当增加此数值。这里的`auto`代表让Nginx自己适配当前电脑的CPU数，与当前CPU数目相同

- `worker_connections 1000`

单个 worker 进程能接受的最大并发连接数，**放在 event 段中**

- `error_log logs/error.log info`

错误日志的存放路径和记录级别

- `use epoll`

使用 epoll 事件模型，**放在 event 段中**


```nginx
events {
    worker_connections 1024;
}
```
`event`字段大概就是这个样子


### http 服务器字段

- `server {}`

定义一个虚拟主机，其内部也有许多相关的配置

- `sendfile on`
  
  开启 sendfile 调用来快速的响应客户端

- `keepalive_timeout 65`

  长连接超时时间，单位是秒。

- `send_timeout`

指定响应客户端的超时时间

- `client_max_body_size 10m`

允许客户端请求的实体最大大小

还有很多字段，具体可以看官方文档，这里给出的都是容易见到的。

### server

给出一个server的例子

```nginx
server {
    listen       80 default_server;
    listen       [::]:80 default_server;
    server_name  _;
    #root         /usr/share/nginx/html;

    # Load configuration files for the default server block.
    include /etc/nginx/default.d/*.conf;

    location / {
      proxy_pass http://127.0.0.1:7001;
    }

    error_page 404 /404.html;
        location = /40x.html {
    }

    error_page 500 502 503 504 /50x.html;
        location = /50x.html {
    }
}
```
可以从这个例子上看出，配置一个server需要`listen`监听某个端口，还可以配置`server_ name`等等。最重要的是`location`字段


不同的`server`，监听的端口不一样，当匹配到对应的`server`之后，一个`server`内可以有很多个`location`，就会去匹配对应的`location`，`location`中有`root`或者`proxy`之类的字段就是用来匹配（这里的匹配说的是请求的url）url中的path，找到符合的`location`然后就完成了代理。



**这只是Nginx最简单的一部分，Nginx是一个很复杂的东西。**





## 参考

- [反向代理为何叫反向代理？](https://www.zhihu.com/question/24723688)
- [Nginx教程(一) Nginx入门教程](https://www.cnblogs.com/crazylqy/p/6891929.html)
- [Beginner’s Guide](http://nginx.org/en/docs/beginners_guide.html)
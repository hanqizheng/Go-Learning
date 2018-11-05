# Docker 基本概念

Docker主要就是解决不同application需要特定的开发环境，而这些开发环境配置起来又相当的繁琐麻烦。现在，docker可以帮我们将一个做好的application和它所处的开发环境变成一个整体，然后就可以在别的环境中不需要再配置任何配置运行这款application。


在给出各个概念之前，先给出一张图(无视字体)

![](/sources/images/webwxgetmsgimg&#32;(1).jpeg)

## DockerFile
`DockerFile`定义了`Container`里要做的事情，资源的访问，网络接口的调用，磁盘的虚拟化，都在这个`Docker File`中。这个环境与操作系统的其余部分分离。因为分离，所以需要将`DockerFile`的端口映射到外界，以便外界后续的使用。并且详细说明哪些文件需要“复制”到`DockerFile`里。

## Image(镜像)
是一个 **可以执行的**，包含所有该应用程序运行所需要的代码、 runtime、 库、 环境变量、 配置文件。

## Container(容器)
是一个runtime。其实简单的理解就可以把`Image`看做一个类(Class),而`Container`则是这个类的一个实例化。

## 图片描述

整个`Docker`的运行过程（这个过程是被简易化的，实际Docker还有很多复杂的东西）就是将源码`Project Code`和Docker相关的文件组成了`Docker File`。然后我们可以将`Docker File`进行`build`就能生成一个`Docker Image`， 然后可以在不同的地方根据这个`Docker Image`运行若干个`Container`，这样就实现了不用在每个要运行application的地方再配置一遍环境。Docker帮我们都做好了。


## Service
在我们的程序中， 连接数据库的功能可以算作一个`Service`， 接收用户传来的数据这个功能也可以算作一个`Service`，我们程序中每个单独的功能都能算作一个`Service`。

`Service`可以理解成一个生产中的`Container`，一个`Service`只会run一个`Docker Image`，它定义了Image以哪种方式运行，应该使用哪个端口，这个`Container`有多少个副本需要跑起来。

## Swarm
一个application不可能只在一个机器上跑，他有可能要在好多个机器上一块跑。将多台机器连接到一个`Dockerized`集群中，就称作一个`Swarm`

# Kubernets第一次

Kubernetes是一个开源的，用于管理云平台中多个主机上的容器化的应用，Kubernetes的目标是让部署容器化的应用简单并且高效（powerful）,Kubernetes提供了应用部署，规划，更新，维护的一种机制。

## 几个重要的概念

### Node
Node作为集群中的工作节点，运行真正的应用程序，在Node上Kubernetes管理的最小运行单元是Pod。Node上运行着Kubernetes的Kubelet、kube-proxy服务进程，这些服务进程负责Pod的创建、启动、监控、重启、销毁、以及实现软件模式的负载均衡。

Node包含的信息：

- Node地址：主机的IP地址，或Node ID。
- Node的运行状态：Pending、Running、Terminated三种状态。
- Node Condition：…
- Node系统容量：描述Node可用的系统资源，包括CPU、内存、最大可调度Pod数量等。
- 其他：内核版本号、Kubernetes版本等。

查看Node信息：
```kubernets
kubectl describe node
```

### Pod

Pod是Kubernetes最基本的操作单元，包含一个或多个紧密相关的容器，一个Pod可以被一个容器化的环境看作应用层的“逻辑宿主机”；一个Pod中的多个容器应用通常是紧密耦合的，Pod在Node上被创建、启动或者销毁；每个Pod里运行着一个特殊的被称之为Pause的容器，其他容器则为业务容器，这些业务容器共享Pause容器的网络栈和Volume挂载卷，因此他们之间通信和数据交换更为高效，在设计时我们可以充分利用这一特性将一组密切相关的服务进程放入同一个Pod中。

同一个Pod里的容器之间仅需通过localhost就能互相通信。

一个Pod中的应用容器共享同一组资源：

- PID命名空间：Pod中的不同应用程序可以看到其他应用程序的进程ID；
- 网络命名空间：Pod中的多个容器能够访问同一个IP和端口范围；
- IPC命名空间：Pod中的多个容器能够使用SystemV IPC或POSIX消息队列进行通信；
- UTS命名空间：Pod中的多个容器共享一个主机名；
- Volumes（共享存储卷）：Pod中的各个容器可以访问在Pod级别定义的Volumes；
Pod的生命周期通过Replication Controller来管理；通过模板进行定义，然后分配到一个Node上运行，在Pod所包含容器运行结束后，Pod结束。

Kubernetes为Pod设计了一套独特的网络配置，包括：为每个Pod分配一个IP地址，使用Pod名作为容器间通信的主机名等。

### Service

在Kubernetes的世界里，虽然每个Pod都会被分配一个单独的IP地址，但这个IP地址会随着Pod的销毁而消失，这就引出一个问题：如果有一组Pod组成一个集群来提供服务，那么如何来访问它呢？Service！

一个Service可以看作一组提供相同服务的Pod的对外访问接口，Service作用于哪些Pod是通过Label Selector来定义的。

拥有一个指定的名字（比如my-mysql-server）；
拥有一个虚拟IP（Cluster IP、Service IP或VIP）和端口号，销毁之前不会改变，只能内网访问；
能够提供某种远程服务能力；
被映射到了提供这种服务能力的一组容器应用上；
如果Service要提供外网服务，需指定公共IP和NodePort，或外部负载均衡器；

### NodePort 
系统会在Kubernetes集群中的每个Node上打开一个主机的真实端口，这样，能够访问Node的客户端就能通过这个端口访问到内部的Service了

### Volume
Volume是Pod中能够被多个容器访问的共享目录。

### Label
Label以key/value的形式附加到各种对象上，如Pod、Service、RC、Node等，以识别这些对象，管理关联关系等，如Service和Pod的关联关系。

### RC（Replication Controller）
目标Pod的定义；
目标Pod需要运行的副本数量；
要监控的目标Pod标签（Lable）；
Kubernetes通过RC中定义的Lable筛选出对应的Pod实例，并实时监控其状态和数量，如果实例数量少于定义的副本数量（Replicas），则会根据RC中定义的Pod模板来创建一个新的Pod，然后将此Pod调度到合适的Node上启动运行，直到Pod实例数量达到预定目标。

## 大致整体架构
Kubernetes 的大致整体架构其实可以用下面两张图来代替，然后结合上面给出的基本概念，就可以对Kubernetes的结构有一个相对清晰的认识

![](https://www.kubernetes.org.cn/img/2016/10/20161028141542.jpg)

可以看出上图中给出了一个宏观的大概体系，`API`这部分我们暂时没有提到，就主要看右边的`Node`部分

观看`Node`部分，可以再结合下图一起，可以看到`Node`的大概结构和具体每块里面的内容。

![](/sources/images/kubernetes.jpg)
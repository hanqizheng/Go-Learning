# Go的包

## 程序组成
一个Go程序是由若干个Go pkg组成的，其中 **有且只有**一个`package main`。`main`包中 **有且只有**一个`func main()`。

每个Go包都可以由若干`.go`文件组成,目录可以是这个样子（只是要展示一下什么是一个包有若干个go文件，一个Go程序有若干个包）

```
GOGOGO
|--test
|   |--test1.go
|   |--test2.go
|--some
|   |--some1.go
|--main.go
```
整个Go程序叫做`GOGOGO`,Go程序有3个包构成，分别是`main`、`test`、`some`。可以看到文件目录中的文件夹就是包名，而在每个文件夹内的`.go`文件都属于这个包，实现这个包。这就是Go的基本目录规则。

## 包的导入(import)
Go程序既然是由不同的包组成，那么包的导入就是一个必须的必要的过程。Go的包导入是遵循一种什么规则？有哪几种导入包的方法？？

```go
import (
  "fmt"
  "net/http"
  "github.com/spf13/viper"
  _ "github.com/hqz/amaze"
  myname "myfmt/fmt"
)
```
> 一般有多个包需要导入的时候，写法就类似上面的代码片段，一个`import`关键字 + ()将多个包括住
> 
> **注**：没有用到的包一定不要导入，Go会报错

上面这个代码片段给除了Go程序中的几种包导入方式（不排除有其他的，这些已经很全了)。

### Go标准库的包
Go标准有许多已经写好，并且且的很好的包。这些包就相当与Go内置的，比如`net`包、`fmt`包。

- 导入这类包时，Go的编译器会先去**安装Go的地方**($GOROOT)中去找。
  > 如果不知道什么是GOROOT请点击[这里](https://github.com/hanqizheng/Go-Learning/blob/master/%E5%9F%BA%E6%9C%AC%E8%AF%AD%E6%B3%95/GOENV.md)学习
- 如果找不到会使用Go的环境变量配置的路径($GOPATH)下的对应位置去找。
- 如果找不到就会报错，提示说找不到

Go编译器在找**到合适的包之后就不会再继续往下找**了。

下面给出一个例子，比如我要导入标准库的`net/http`包,假设我安装Go的目录(GOROOT)为`/usr/local/go`,我定义的Go环境变量(GOPATH)为`/home/hqz/go`
```
/usr/local/go/pkg/net/http  //这里就找到了，不会继续下面

/home/hqz/go/src/net/http
```
可以看到Go的查找先后顺序。这里要说一下如果自己测试有可能发现路径会和上面描述的有些许不同，是因为Go在创建文件夹的时候会带上当前操作系统的名称，所以有可能会变成下述
```
/usr/local/go/pkg/linux_amd64/net/http
```

### 远程导入
Go的社区很活跃，大家会互相分享自己写好的包，这些包会放在代码托管的平台上供大家下载，such as `github`。

所以自己的程序难免会用到别人写好的包，这时候就需要远程导入包。

```go
import "github.com/hqz/gin"
```
就像这种，给出了一个类似URL的导入路径，其实Go的编译器还是会现在本地找，**找的方法同上面标准库导入是一样的**，当然不可能在`GOROOT`找到，所以要去`GOPATH`去找，还是没有找到。

这个时候就会报错，说没有这个包，可以尝试用`go get`命令获取这个包

> 不知道`go get`可以点击[这里](https://github.com/hanqizheng/Go-Learning/blob/master/%E5%9F%BA%E6%9C%AC%E8%AF%AD%E6%B3%95/GOCMD.md)

然后就能从远程仓库获取到这个包，下载好的包会存在`GOPATH`对应的目录下面，类似这样

```
/home/hqz/go/src/github.com/gin
```
### 命名导入
有的时候，自己实现的包的名字会和标准库或者远程导入的包名相同，这样会造成冲突
```go
import (
  "fmt"
  "my/fmt"
)
```
这个时候就可以自定义包的名字，后续使用就直接使用自定义的名字即可

```go
import (
  "fmt"
  myfmt "my/fmt"
)
```
有的时候，我们只需要使用一个包的`init()`函数，而不需要使用包内的东西，这个时候可以使用空白命名

```go
import _ "fmt"
```
起名字起成 `_` 即可


## 与Node包机制比较
- Node使用`require()`进行导入，多个包要写多个`require()`
  
  Go使用`import`进行导入，多个包建议只写一个`import`然后用小括号包裹

- Node查找顺序是在当前工程内由内向外在`node_modules`中查找

  Go则是在环境变量`GOROOT`与`GOPATH`对应的目录去查询

- Node使用其他开发者自定义的包可以通过`npm`进行`npm install`，这一点Node做的很好

  Go则要使用`go get`进行安装，而且Go至今没有一个统一的平台来管理这些包，多数都是放在`github`上的

- Node的`node_module`太过沉重，一个程序根本用不了那么多的包，但是方便

  Go不允许多导入任何一个多余的没有使用的包


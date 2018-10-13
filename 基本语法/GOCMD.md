# GO的编译过程及相关命令

## go run

这个命令是用来**专门运行命令源码文件**的， 他不是用来运行所有源码文件的。

记得在上一节说到go源码文件分为三种，而`go run`只能运行`属于main包，包含main函数的那种，和属于main包的库源码文件`

如果参数中有多个或者没有命令源码文件，那么 go run 命令就只会打印错误提示信息并退出，而不会继续执行。

可以使用`go run -n`查看具体这条命令在执行的时候干了些什么

```
HQZ ~/LeetCode_Go/helloworld/src/me $  go run -n helloworld.go

#
# command-line-arguments
#

mkdir -p $WORK/command-line-arguments/_obj/
mkdir -p $WORK/command-line-arguments/_obj/exe/
cd /Users/YDZ/LeetCode_Go/helloworld/src/me
/usr/local/Cellar/go/1.8.3/libexec/pkg/tool/darwin_amd64/compile -o $WORK/command-line-arguments.a -trimpath $WORK -p main -complete -buildid 2841ae50ca62b7a3671974e64d76e198a2155ee7 -D _/Users/YDZ/LeetCode_Go/helloworld/src/me -I $WORK -pack ./helloworld.go
cd .
/usr/local/Cellar/go/1.8.3/libexec/pkg/tool/darwin_amd64/link -o $WORK/command-line-arguments/_obj/exe/helloworld -L $WORK -w -extld=clang -buildmode=exe -buildid=2841ae50ca62b7a3671974e64d76e198a2155ee7 $WORK/command-line-arguments.a
$WORK/command-line-arguments/_obj/exe/helloworld
```

可以看到先是创建了两个临时文件夹`_obj和exe`,然后编译(compile)，链接(link)，生成了归档文件.a 和 最终可执行文件，最终的可执行文件放在 exe 文件夹里面。命令的最后一步就是执行了可执行文件。

可以看到，最终go run命令是生成了2个文件，一个是归档文件，一个是可执行文件。command-line-arguments 这个归档文件是 Go 语言为命令源码文件临时指定的一个代码包。在接下来的几个命令中，生成的临时代码包都叫这个名字。

## go build

当代码包中有且仅有一个命令源码文件的时候，在文件夹所在目录中执行 go build 命令，会在该目录下生成一个与目录同名的可执行文件。

go build 用于编译我们指定的源码文件或代码包以及它们的依赖包。，但是注意如果用来编译非命令源码文件，即库源码文件，go build 执行完是不会产生任何结果的。这种情况下，go build 命令只是检查库源码文件的有效性，只会做检查性的编译，而不会输出任何结果文件。

go build 编译命令源码文件，则会在该命令的执行目录中生成一个可执行文件，上面的例子也印证了这个过程。

go build 后面不追加目录路径的话，它就把当前目录作为代码包并进行编译。go build 命令后面如果跟了代码包导入路径作为参数，那么该代码包及其依赖都会被编译。

接下来可看看`go build`都做了些什么

```
#
# command-line-arguments
#

mkdir -p $WORK/command-line-arguments/_obj/
mkdir -p $WORK/command-line-arguments/_obj/exe/
cd /Users/YDZ/MyGitHub/LeetCode_Go/helloworld/src/me
/usr/local/Cellar/go/1.8.3/libexec/pkg/tool/darwin_amd64/compile -o $WORK/command-line-arguments.a -trimpath $WORK -p main -complete -buildid 2841ae50ca62b7a3671974e64d76e198a2155ee7 -D _/Users/YDZ/MyGitHub/LeetCode_Go/helloworld/src/me -I $WORK -pack ./helloworld.go
cd .
/usr/local/Cellar/go/1.8.3/libexec/pkg/tool/darwin_amd64/link -o $WORK/command-line-arguments/_obj/exe/a.out -L $WORK -extld=clang -buildmode=exe -buildid=2841ae50ca62b7a3671974e64d76e198a2155ee7 $WORK/command-line-arguments.a
mv $WORK/command-line-arguments/_obj/exe/a.out helloworld
```

可以看到，执行过程和 go run 大体相同，唯一不同的就是在最后一步，go run 是执行了可执行文件，但是 go build 命令是把可执行文件移动到了当前目录的文件夹中。

如果`go build`了一个库源码文件的话，又是什么样子呢
```
#
# _/Users/YDZ/Downloads/goc2p-master/src/pkgtool
#

mkdir -p $WORK/_/Users/YDZ/Downloads/goc2p-master/src/pkgtool/_obj/
mkdir -p $WORK/_/Users/YDZ/Downloads/goc2p-master/src/
cd /Users/YDZ/Downloads/goc2p-master/src/pkgtool
/usr/local/Cellar/go/1.8.3/libexec/pkg/tool/darwin_amd64/compile -o $WORK/_/Users/YDZ/Downloads/goc2p-master/src/pkgtool.a -trimpath $WORK -p _/Users/YDZ/Downloads/goc2p-master/src/pkgtool -complete -buildid cef542c3da6d3126cdae561b5f6e1470aff363ba -D _/Users/YDZ/Downloads/goc2p-master/src/pkgtool -I $WORK -pack ./envir.go ./fpath.go ./ipath.go ./pnode.go ./util.go
```

这里可以看到 go build 命令只是把库源码文件编译了一遍，其他什么事情都没有干。

## go install
go install 命令是用来编译并安装代码包或者源码文件的。

与 go build 命令一样，传给 go install 命令的代码包参数应该以导入路径的形式提供。并且，go build 命令的绝大多数标记也都可以用于
go install 命令。实际上，go install 命令只比 go build 命令多做了一件事，即：安装编译后的结果文件到指定目录。

```
#
# command-line-arguments
#

mkdir -p $WORK/command-line-arguments/_obj/
mkdir -p $WORK/command-line-arguments/_obj/exe/
cd /Users/YDZ/MyGitHub/LeetCode_Go/helloworld/src/me
/usr/local/Cellar/go/1.8.3/libexec/pkg/tool/darwin_amd64/compile -o $WORK/command-line-arguments.a -trimpath $WORK -p main -complete -buildid 2841ae50ca62b7a3671974e64d76e198a2155ee7 -D _/Users/YDZ/MyGitHub/LeetCode_Go/helloworld/src/me -I $WORK -pack ./helloworld.go
cd .
/usr/local/Cellar/go/1.8.3/libexec/pkg/tool/darwin_amd64/link -o $WORK/command-line-arguments/_obj/exe/a.out -L $WORK -extld=clang -buildmode=exe -buildid=2841ae50ca62b7a3671974e64d76e198a2155ee7 $WORK/command-line-arguments.a
mkdir -p /Users/YDZ/Ele_Project/clairstormeye/bin/
mv $WORK/command-line-arguments/_obj/exe/a.out /Users/YDZ/Ele_Project/clairstormeye/bin/helloworld
```

前面几步依旧和 go run 、go build 完全一致，只是最后一步的差别，go install 会把命令源码文件安装到当前工作区的 bin 目录（如果 GOPATH 下有多个工作区，就会放在 GOBIN 目录下）。如果是库源码文件，就会被安装到当前工作区的 pkg 的平台相关目录下。

## go get
go get 命令用于从远程代码仓库（比如 Github ）上下载并安装代码包。

**注意，go get 命令会把当前的代码包下载到 $GOPATH 中的第一个工作区的 src 目录中，并安装。**

go get 命令究竟做了些什么呢？我们还是来打印一下每一步的执行过程。

```
cd .
git clone https://github.com/go-errors/errors /Users/YDZ/Ele_Project/clairstormeye/src/github.com/go-errors/errors
cd /Users/YDZ/Ele_Project/clairstormeye/src/github.com/go-errors/errors
git submodule update --init --recursive
cd /Users/YDZ/Ele_Project/clairstormeye/src/github.com/go-errors/errors
git show-ref
cd /Users/YDZ/Ele_Project/clairstormeye/src/github.com/go-errors/errors
git submodule update --init --recursive
WORK=/var/folders/66/dcf61ty92rgd_xftrsxgx5yr0000gn/T/go-build124856678
mkdir -p $WORK/github.com/go-errors/errors/_obj/
mkdir -p $WORK/github.com/go-errors/
cd /Users/YDZ/Ele_Project/clairstormeye/src/github.com/go-errors/errors
/usr/local/Cellar/go/1.8.3/libexec/pkg/tool/darwin_amd64/compile -o $WORK/github.com/go-errors/errors.a -trimpath $WORK -p github.com/go-errors/errors -complete -buildid bb3526a8c1c21853f852838637d531b9fcd57d30 -D _/Users/YDZ/Ele_Project/clairstormeye/src/github.com/go-errors/errors -I $WORK -pack ./error.go ./parse_panic.go ./stackframe.go
mkdir -p /Users/YDZ/Ele_Project/clairstormeye/pkg/darwin_amd64/github.com/go-errors/
mv $WORK/github.com/go-errors/errors.a /Users/YDZ/Ele_Project/clairstormeye/pkg/darwin_amd64/github.com/go-errors/errors.a

```

这里可以很明显的看到，执行完 go get 命令以后，会调用 git clone 方法下载源码，并编译，最终会把库源码文件编译成归档文件安装到 pkg 对应的相关平台目录下。


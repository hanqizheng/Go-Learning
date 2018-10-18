# Go 的 interface

在没有学习Go的 interface 之前，我不知道Go没有类这个概念（感觉怎么是一句废话？？）

没有类之后，最大的问题就是Go怎么实现面向对象的三个特性？

今天的接口(interface)就和其中一个特性`多态`有关系。

## 基本语法
```go
// 代码段1.1

type myInterface interface {
    // 里面有若干方法，但是都没有具体实现
    doSomething(a int) int
    doOtherThing(b string) (string, int)
}
```

基本语法不是重点看看就明白，所以继续。

## 该如何理解interface

**首先，也是最重要的，interface是一种类型，就像int,string一样，它是一种类型。**

go中的接口和C++中的纯虚 基类 （只有虚函数，没有数据成员）类似。在Go语言中，任何实现了接口的函数的类型，都可以 看作是接口的一个实现。 类型在实现某个接口的时候，不需要显式关联该接口的信息。接口的实现 和接口的定义完全分离了。

接上部分代码段1.1
```go
//代码段1.2

type myType struct {
    age int
}

func (m *myType) doSomething(a int) int {
    var result int
    m.age = 50
    result = m.age - a
    return result
}

func main() {
    var test myType
    test.doSomething(5)
}
```
结合代码段1.1 、 1.2不难看出，定义的接口`myInterface`内部只有函数的名字，参数，返回值之类的信息。具体的实现却没有给出来。

代码段1.2中我们定义了一个自己的结构体`myType`，将具体实现的`doSomething()`挂载(或者说依附或者说别的，道理你懂得)到`myType`上面,这样`myType`这个结构体就有了自己的一个 **方法**(具体方法是什么我会在别的章节学习)。

```go
//代码段1.3

func (m *myType) doOtherThing(b string) (string, int) {
    //.....
}
```
加入代码段1.3已经类似1.2一样，将`doOtherThing()`也实现并且挂在`myType`上了，那么这么一来，**`myType`自此已经实现了`myInterface`中的所有抽象方法，那么，就说`myType`实现了`myInterface`**

必须注意，必须把接口内 **所有的**方法都实现，才能说实现了这个接口。

### 值得一提的

#### 第一个值得一提的是

不光我们自己定义的struct才能挂在方法，Go自带的类型也是可以的，只不过，需要给他起一个别名，just like this
```go
// 代码段1.4

type myInt int

func (i myInt) doSomething(a int) int {
    .....
} 
```

#### 第二个值得一提的
当一个类型（可以是自定义的也可以是Go自带的）实现了一个接口，那么，所有方法中含有该接口类型的参数都可以用`那个类型（实现了那个接口的那个类型）`的变量来代替。

晕了，掀桌.jpg

```go
//代码片段1.5

   type MyInterface interface{
       Print()
   }
   
   func TestFunc(x MyInterface) {}

   type MyStruct struct {}

   func (me MyStruct) Print() {}
   
   func main() {
       var me MyStruct
       TestFunc(me)
   }
```
片段1.5里所示，其中类型`MyStruct`实现了接口`MyInterface`。在一个普通函数`TestFunc()`中有一个参数是`MyInterface`类型的，但是在`main()`中，却是将`MyStruct`类型的变量当作参数传递给了`TestFucn()`这就是我要表达的意思。

## interface值
概念上讲一个接口的值，接口值，由两个部分组成，一个具体的类型和那个类型的值。

在Go语言中，变量总是被一个定义明确的值初始化，即使接口类型也不例外。对于一个接口的零值就是它的类型和值的部分都是nil

![](https://yar999.gitbooks.io/gopl-zh/content/images/ch7-01.png)

正如上图，这是一个空接口的接口值的具体表示

> 空接口可以写成 `interface {}`

它们被称为接口的动态类型和动态值。对于像Go语言这种静态类型的语言，类型是编译期的概念；因此一个类型不是一个值。在我们的概念模型中，一些提供每个类型信息的值被称为类型描述符，比如类型的名称和方法。在一个接口值中，类型部分代表与之相关类型的描述符。

其实这段话`简单但不准确但好理解的说`可以比如一下，比如不同的类型实现了同一个接口，那么接口值其实是改变的，不论是`type`还是`value`。

下面就用Go圣经的经典例子，来理解一下接口值

```go
// 代码片段2.1

var w io.Writer
w = os.Stdout
```
看一下代码片段2.1，首先是定义了一个接口类型的变量`w`，这里的`io.Writer`就是一个接口，所以他是一个类型，所以可以定义一个以他为类型的值。

```go
w = os.Stdout
```
这句话是将`os.Stdout`类型的值复制给了`w`，其实是存在一个隐式转换的

> os.Stdout --> io.Writer

其实这里的`os.Stdout`是`os.File`类型，具体为什么是os包内的具体内容，暂时不做解释，暂时先直接理解成`os.File`，所以其实是

> os.File --> io.Writer

这一个变化就会改变接口值对应的`type`和`value`

![](https://yar999.gitbooks.io/gopl-zh/content/images/ch7-02.png)
可看到`type`和`value`都已经改变了

### 值得一提又来了

#### 第一 接口值可以持有任意大动态值
一个接口值可以持有任意大的动态值。

例如，表示时间实例的time.Time类型

```go
var x interface{} = time.Now()
```

从概念上讲，不论接口值多大，动态值总是可以容下它。（这只是一个概念上的模型；具体的实现可能会非常不同）
![](https://yar999.gitbooks.io/gopl-zh/content/images/ch7-04.png)

#### 第二 接口值的比较

如果想比较两个接口值，那么就有两种情况，接口值相等

- 两个接口都是空接口 nil == nil
- 两个接口值的动态类型`type`相等，且在当前相等`type`所对应的`value`也想等

#### 第三 其实也是第二的一部分  nil 和 空接口

```go
// 代码片段2.2
type my interface {}

func main() {
		var null my
		if null == nil {
				fmt.Println("yes")
		} else {
				fmt.Println("no")
		}
}

// yes
```
如代码片段2.2 答案是yes，因为空接口的`type`和`value`都是nil，所以和nil相等，其他情况不等

## 参考

- [The Go Programming Language](https://yar999.gitbooks.io/gopl-zh/content/index.html)
- [Golang interface接口深入理解](https://juejin.im/post/5a6873fd518825734501b3c5)
- [Go 中文文档](https://wizardforcel.gitbooks.io/golang-doc/content/index.html)
- [Go语言interface详解](https://studygolang.com/articles/9099)
- [Go interface 详解 (三) ：interface 的值](https://sanyuesha.com/2017/10/18/go-interface-3/)
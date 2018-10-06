# fmt的Reader接口实现io读

今天阅读到了官方文档的接口部分，文档列举出了几个常用的接口 Reader，Stringer、Error等。这篇咱们就来说一下用io.Reader内置的Read就可以实现字符串的读。

首先Reader是在Go内置的io包中

先来看一个例子，**注意看Read的返回值和最终输出的结果有什么特点**
```go
package main
import (
	"fmt"
	"io"
	"strings"
)

func main() {
	r := strings.NewReader("Hello, Reader!")

	b := make([]byte, 8)
	for {
		n, err := r.Read(b)
		fmt.Printf("n = %v err = %v b = %v\n", n, err, b)
		fmt.Printf("b[:n] = %q\n", b[:n])
		if err == io.EOF {
			break
		}
	}
}
```
先看一下`n, err := r.Read(b)`，这句话可以看到Read返回两个值，这里我来说明一下，返回的两个值中:

**第一个值**代表读取的数据个数或者长度，这样说太过模糊，不知道有没有注意到这句话`b := make([]byte, 8)`，创建的切片是`byte`类型的。我最开始看这个例子，把这个byte变成int，float之类的都会直接报错，报错如下
```
# command-line-arguments
./hello.go:13:19: cannot use b (type []rune) as type []byte in argument to r.Read
```
也就是说，这个Read他要去i必须是`byte`型,所以他每次最多读`8位`，这一点也可以在输出结果清晰的看到
```
n = 8 err = <nil> b = [72 101 108 108 111 44 32 82]
b[:n] = "Hello, R"
n = 6 err = <nil> b = [101 97 100 101 114 33 32 82]
b[:n] = "eader!"
n = 0 err = EOF b = [101 97 100 101 114 33 32 82]
b[:n] = ""
```
可以看到`hello， Reader!`这个字符串长度超过了8, 所以：
第一次读取之后，返回了读取数据的个数`8`,

然后读取了字符串的前8个字符，然后第二次读取返回了`6`,因为只剩下6个字符。

最后一次返回`0`,因为已经读取完毕

**其实真真正正为什么必须是byte，他每次只读1 byte的内容然后填充在对应的切片里，就是这么简单**

**第二个值**代表读取中有可能遇到的错误`error`，如果没有遇到错误，就返回`nil`


## 另一个简单的例子

使用io.Reader来破解一串rot13密码

```
ROT13  是过去在古罗马开发的**凯撒加密**的一种变体。

 套用ROT13到一段文字上仅仅只需要检查字元字母顺序并取代它在13位之后的对应字母 ，
 有需要超过时则重新绕回26英文字母开头即可[2]。A换成N、B换成O、依此类推到M换成Z，
 然后序列反转：N换成A、O换成B、最后Z换成M。只有这些出现在英文字母里头的字元受影响；
 数字、符号、空白字元以及所有其他字元都不变。 
 因为只有在英文字母表里头只有26个，并且26=2×13，ROT13函数是它自己的逆反：
```
我们给除了一段密码`Lbh penpxrq gur pbqr!`，根据上面的原理讲解利用Reader进行破解。下面代码已经给出

```go
package main

import (
    "io"
    "os"
    "strings"
)

type rot13Reader struct {
    r io.Reader
}

func rot13(b byte) byte {
    switch {
        case 'A' <= b && 'M' >= b:
            b = b + 13
        case 'N' <= b && 'Z' >= b:
            b = b -13
        case 'a' <= b && 'm' >= b:
            b = b + 13
        case 'n' <= b && 'z' >= b:
           b = b -13
    }
    return b
}

func (rr rot13Reader) Read(b []byte) (int, error) {
    n, e := rr.r.Read(b)
    for i := 0; i < n; i++ {
	b[i] = rot13(b[i])	
    }
    return n, e
}

func main() {
    s := strings.NewReader("Lbh penpxrq gur pbqr!")
    r := rot13Reader{s}
    io.Copy(os.Stdout, &r)
}
```

包装现有的Reader形成一个自己的全新的特定功能的Reader是这个例子的核心思想。

这个例子中的Reader需要完成字符之间的关系转换（其实就是解密），每读一次，就把这次读的内容解密一次，这样就把一串密码解密了

好奇最后的答案是什么- -

```go
'You cracked the code!'
```
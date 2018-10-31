# Go 与 JSON

JSON是平时最常用的一种数据格式，那么Go怎么操作JSON呢？

说到JSON对它最多的操作就是解析某个JSON或者将数据组合生成一个JSON。

## Go与JSON的数据类型对应关系

```go
bool,                   for        JSON booleans
float64,                for        JSON numbers
string,                 for        JSON strings
[]interface{},          for        JSON arrays
map[string]interface{}, for        JSON objects
nil                     for        JSON null
```

## Unmarshal
先说解析一个拿到的JSON

我们用到的包是Go内置的
```go
import "encoding/json"
```

现在我们有一个JSON需要解析一下
```json
{
  "species": "pigeon",
  "decription": "likes to perch on rocks"
}
```
在Go中，如果想要解析某个JSON，就必须要创建一个对应的结构体。具体的说就是根据JSON中的`key & value`，创建一个拥有对应类型成员变量的`struct`，just like this
```go
type Bird struct {
  Species string
  Description string
}
```
上面`json`中有`species & decription`两个key，那么结构体也要有两个成员变量，**而且注意类型都要一致**

> **important**
> 
> 有注意到结构体中每个变量命名的特点吗？
> 
> **都是大写字母开头**，这一点是有原因的:
> 
> 大写开头是说明将要使用的，接下来所关注的一些json中的字段，**只有以大写字母开头的成员变量，才能解析与json中对应的key**，如果一个成员变量以小写开头定义或者不写，都是不会被解析的。

然后就是解析了
```go
birdJson := `{"species": "pigeon","description": "likes to perch on rocks"}`
var bird Bird	
json.Unmarshal([]byte(birdJson), &bird)
fmt.Printf("Species: %s, Description: %s", bird.Species, bird.Description)
```
可以看到有一个JSON结构的变量`birdJson`，使用`json.Unmarshal([]byte(birdJson), &bird)`解析，这个函数第一个参数是要解析的对象，**必须是[]byte类型**，这里将一整串数据8位8位读取到byte类型切片中。第二个参数是目标变量，也就是存放解析好的json的，所以可以看到这里使用了先前定义好的`Bird`类型，也就是与json的key想对应的`struct`。

### Unmarshal a JSON Array

当然在实际的开发中，很多情况并不是单纯的只接受一个JSON，而是多个JSON组成的数组，所以知道如何解析JSON数组也是一个非常重要的。

假设现在有一个这样的JSON需要解析
```json
[
  {
    "species": "pigeon",
    "decription": "likes to perch on rocks"
  },
  {
    "species":"eagle",
    "description":"bird of prey"
  }
]
```
只需要将上面`var bird Bird`变为`var bird []Bird`就可以了，对就是将一个Bird类型的变量变成Bird类型的数组即可

```go
birdJson := `[{"species":"pigeon","decription":"likes to perch on rocks"},{"species":"eagle","description":"bird of prey"}]`
var birds []Bird
json.Unmarshal([]byte(birdJson), &birds)
fmt.Printf("Birds : %+v", birds)
//Birds : [{Species:pigeon Description:} {Species:eagle Description:bird of prey}]
```

### JSON中包含一个对象
然而简单的JSON往往也无法满足日常实际中的开发，一个复杂的JSON对象可能包含一个`Object`类型的`key`，就像这个

```json
{
  "species": "pigeon",
  "decription": "likes to perch on rocks",
  "dimensions": {
    "height": 24,
    "width": 10
  }
}
```
解析这类的JSON也是同理，只需要在定义一个对应的`struct`即可，like this
```go
type Dimensions struct {
  Height int
  Width int
}

type Bird struct {
  Species string
  Description string
  Dimensions Dimensions
}
```
然后解析的方法就完全一样了

```go
birdJson := `{"species":"pigeon","description":"likes to perch on rocks", "dimensions":{"height":24,"width":10}}`
var birds Bird
json.Unmarshal([]byte(birdJson), &birds)
fmt.Printf(bird)
// {pigeon likes to perch on rocks {24 10}}
```

### 自定义属性名称

有时候并不像严格按照JSON中每个`key`的名称来定义`struct`中每个属性的名称，这个时候就可以用到`filed tags`来解决这个问题

```json
{
  "birdType": "pigeon",
  "what it does": "likes to perch on rocks"
}
```
别如上面这个JSON，如果要按照`key`来定义结构体中的属性名称，那么可能就要定义一个`WhatItDose`这样的属性...不是很好，所以可以这样写，来完成自定义的属性名
```go
type Bird struct {
  Species string `json:"birdType"`
  Description string `json:"what it does"`
}
```

## Marshal
说完解析JSON，拼凑一个JSON则可以使用完全相反的过程来完成，想对应的要使用 `json.Marshal()`这个函数

同样，先定义好对应的结构体
```go
type Animal struct {
    Name  string `json:"name"`
    Order string `json:"order"`
}
```

然后，将想要转成JSON的数据进行拼装，然后Marshl

```go
var animals []Animal
animals = append(animals, Animal{Name: "Platypus", Order: "Monotremata"})
animals = append(animals, Animal{Name: "Quoll", Order: "Dasyuromorphia"})

jsonStr, err := json.Marshal(animals)
if err != nil {
    fmt.Println("error:", err)
}

fmt.Println(string(jsonStr))
```

其实还有其他方法来组装，比如Gin框架中的`c.JSON`我会在后续的例子中给出

## Gin & JSON
Gin框架中使用JSON也是相当的方便

先给出我要POST的一段JSON
```json
{
	"name": "hqz",
	"age" : "12",
	"habit": "play tennis"
}
```

然后在代码中定义好对应的结构体，这些都和上述方法一致
```go
type Info struct {
	Name string `json:"name"`
	Age string	`json:"age"`
	Habit string `json:"habit"`
}
```
然后接受`request body`中POST过来的JSON，然后将得到的数据Bind，就可以拿到新鲜的JSON了，具体这个Bind之类的是什么，会在后续学习Gin框架中给出。

```go
func testPage(c *gin.Context)  {
	i := &Info{}
	c.Bind(i)
	fmt.Println(i.Name)
	c.JSON(http.StatusOK, i)
}

func main()  {
	router := gin.Default()
	router.POST("/test", testPage)
	router.Run()
}
```
 可以看到`c.JSON(http.StatusOK, i)`这个方法，接受两个参数，第一个参数是一个`int`类型，就是状态码，第二个参数则是要返回的具体内容，是一个`interface{}`具体为什么是interface我现在还不知道。

 然后这样就把拿到的JSON返回了，这个过程适用于写API的时候。

## 参考

[Parsing JSON in Golang](https://www.sohamkamani.com/blog/2017/10/18/parsing-json-in-golang/)

[golang中的json处理](https://segmentfault.com/a/1190000009820634)
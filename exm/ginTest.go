package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"github.com/gin-gonic/gin"
)

// Info to Unmarshal JSON
type Info struct {
	Name string `json:"name"`
	Age string	`json:"age"`
	Habit string `json:"habit"`
}

// HomePage is 首页
func HomePage(c *gin.Context)  {
	c.JSON(200, gin.H{
		"message": "hello,world!",
	})
}

// PostPage is 尝试POST text然后显示
func PostPage(c *gin.Context)  {
	body := c.Request.Body
	value, err := ioutil.ReadAll(body)

	if err != nil {
		fmt.Println(err.Error())
	}

	c.JSON(200, gin.H{
		"message": string(value),
	})
}

// QueryTest is 测试Gin框架怎么接受query
func QueryTest(c *gin.Context)  {
	name := c.Query("name")
	age := c.Query("age")

	c.JSON(200, gin.H{
		"name": name,
		"age": age,
	})
}

// QueryTest is 测试Gin框架怎么接受params
func ParamsTest(c *gin.Context)  {
	name := c.Param("name")
	age := c.Param("age")

	c.JSON(200, gin.H{
		"name": name,
		"age": age,
	})
}

// testPage is 解析JSON的Handler
func testPage(c *gin.Context)  {
	i := &Info{}
	c.Bind(i)
	fmt.Println(i.Name)
	c.JSON(http.StatusOK, i)

}

func main()  {
	router := gin.Default()

	router.GET("/", HomePage)
	router.GET("/query", QueryTest)
	router.GET("/params/:name/:age", ParamsTest)
	router.POST("/", PostPage)
	router.POST("/test", testPage)

	router.Run()
}
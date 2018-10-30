package user

import (
	"net/http"

	userService "../../service/user"

	"github.com/gin-gonic/gin"
)

//User is a struct 对应解析JSON
type User struct {
	UserName string `json: "username"`
	Password string `json: "password"`
}

// GetList is 获取用户列表
func GetList(c *gin.Context) {
	users := userService.FindAll()
	c.JSON(http.StatusOK, users)
}

// AddUser is 添加一个新用户
func AddUser(c *gin.Context) {
	user := &User{}
	c.Bind(user)
	userService.AddOne(user.UserName, user.Password)
	c.JSON(201, gin.H{
		"message": "create",
	})
}

// DeleteUser is 删除一个用户
func DeleteUser(c *gin.Context) {
	username := c.Param("username")
	userService.DelOne(username)
	c.JSON(204, gin.H{
		"message": "delete",
	})
}

// UpdateUser is 更新一个用户信息
func UpdateUser(c *gin.Context) {
	username := c.Param("username")
	info := &User{}
	c.Bind(info)
	userService.UpdateOne(info)
	c.JSON(200, gin.H{
		"message": "update",
	})
}

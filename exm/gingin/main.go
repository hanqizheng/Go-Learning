package main

import (
	user "./controller/user"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	router := gin.Default()

	router.GET("/user", user.GetList)
	router.POST("/user", user.AddUser)
	router.DELETE("/user/:username", user.DeleteUser)
	router.PUT("/user/:username", user.UpdateUser)

	router.Run()
}

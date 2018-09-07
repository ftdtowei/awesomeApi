package main

import (
	"awesomeProject/service/project"
	"awesomeProject/service/user"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
)

func main() {
	//日志模块
	gin.DisableConsoleColor()

	// Logging to a file.
	f, _ := os.Create("gin.log")
	gin.DefaultWriter = io.MultiWriter(f)

	//路由模块
	router := gin.Default()
	//用户权限管理
	user := router.Group("/user")
	{
		//登录
		user.POST("/login", userService.UserLogin)
		//注册
		user.POST("/register", userService.UserRegister)
	}

	//项目接口管理
	project := router.Group("/project")
	{
		project.Use(Validate())
		project.POST("/qry", projectService.GetSelfPro)
	}
	router.Run(":8080")
}

func Validate() gin.HandlerFunc{
	return func(c *gin.Context) {
		if cookie, err := c.Request.Cookie("session_id"); err == nil {
			value := cookie.Value
			fmt.Println(value)
			if value == "aweSomeApi" {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		c.Abort()
		return
	}
}


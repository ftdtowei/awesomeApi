package main

import (
	"awesomeApi/service"
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
		user.POST("/login", service.UserLogin)
		//注册
		user.POST("/register", service.UserRegister)
	}


	//项目接口管理
	project := router.Group("/project")
	{
		project.Use(Validate())
		//查询自己关联的项目
		project.POST("/qryMyPro", service.QryMyPro)

		//增删项目成员
		project.POST("manageMember", service.ManageMember)

		//查询接口模块
		project.POST("/qryModule", service.QryModule)

		//管理接口模块  增删改
		project.POST("/manageModule", service.ManageModule)

		//管理分组
		project.POST("/managePage", service.ManagePage)

		//查询模块下的分组及接口列表
		project.POST("/qryActionList", service.QryActionList)

		//查询接口详情
		project.POST("/qryActionDetail", service.QryActionDetail)

		//管理接口详情  增删改
		project.POST("/manageActionDetail", service.ManageActionDetail)
	}

	router.Run(":8080")
}

func Validate() gin.HandlerFunc {
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

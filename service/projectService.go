package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

//查询自己关联的项目
func QryMyPro(c *gin.Context) {
	c.JSON(http.StatusOK,gin.H{"message":"your pro"})
}

//查询接口模块
func QryModule(c *gin.Context) {
	c.JSON(http.StatusOK,gin.H{"message":"your pro"})
}
//增删成员
func ManageMember(c *gin.Context) {
	c.JSON(http.StatusOK,gin.H{"message":"your pro"})
}

//管理接口模块  增删改
func ManageModule(c *gin.Context) {
	c.JSON(http.StatusOK,gin.H{"message":"your pro"})
}

//管理分组
func ManagePage(c *gin.Context) {
	c.JSON(http.StatusOK,gin.H{"message":"your pro"})
}

//查询模块下的分组及接口列表
func QryActionList(c *gin.Context) {
	c.JSON(http.StatusOK,gin.H{"message":"your pro"})
}

//查询接口详情
func QryActionDetail(c *gin.Context) {
	c.JSON(http.StatusOK,gin.H{"message":"your pro"})
}

//管理接口详情  增删改
func ManageActionDetail(c *gin.Context) {
	c.JSON(http.StatusOK,gin.H{"message":"your pro"})
}

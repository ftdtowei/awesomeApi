package service

import (
	"awesomeApi/utils"
	. "encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

//将body转化成map
func bodytomap(c *gin.Context) (map[string]interface{}, error) {
	//将 body 流转换成 map
	data, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Println("----params----")
	fmt.Println("ctx.Request.body: %v", string(data))
	var m map[string]interface{}
	err := Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

//查询自己关联的项目 userId
func QryMyPro(c *gin.Context) {
	//转换body
	params, err := bodytomap(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "json error"})
		return
	}
	userId := params["userId"]
	//查询自己创建的project和自己加入得project
	rows, err := db.Query("select id,project_name,description from tb_project where create_user =?", userId)
	if err != nil {
		//查询报错
		fmt.Println(err.Error())
	}
	selfjson := utils.Rowtojson(rows)

	row1, err := db.Query("select p.id,p.project_name,p.description from tb_project p,project_and_user pu where p.id = pu.project_id and pu.user_id = ?", userId)
	if err != nil {
		//查询报错
		fmt.Println(err.Error())
	}
	joinjson := utils.Rowtojson(row1)

	c.JSON(http.StatusOK, gin.H{"self": selfjson, "join": joinjson})
}

//增删成员
func ManageMember(c *gin.Context) {
	//转换body
	params, err := bodytomap(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "json error"})
		return
	}
	projectId := params["projectId"]
	userId := params["userId"]
	opt := params["opt"]
	if opt == "add" {
		println("add")

		id := utils.UniqueId()
		stmt, err := db.Prepare("insert into project_and_user (id,project_id,user_id) values(?,?,?)")
		if err != nil {
			log.Fatalln(err)
		}
		res, err := stmt.Exec(id, projectId, userId)
		if err != nil {
			log.Fatalln(err)
		}
		println(res)
	} else if opt == "del" {
		println("del")

		id := utils.UniqueId()
		stmt, err := db.Prepare("delete project_and_user where ")
		if err != nil {
			log.Fatalln(err)
		}
		res, err := stmt.Exec(id, projectId, userId)
		if err != nil {
			log.Fatalln(err)
		}
		println(res)

	}

	c.JSON(http.StatusOK, gin.H{"message": "your pro"})
}

//查询项目模块
func QryModule(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "your pro"})
}

//锁定项目模块
func LockModule(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "your pro"})
}

//管理项目模块  增删改
func ManageModule(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "your pro"})
}

//管理分组
func ManagePage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "your pro"})
}

//查询模块下的分组及接口列表
func QryActionList(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "your pro"})
}

//查询接口详情
func QryActionDetail(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "your pro"})
}

//管理接口详情  增删改
func ManageActionDetail(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "your pro"})
}

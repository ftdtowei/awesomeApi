package service

import (
	"awesomeApi/utils"
	. "encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//将body转化成map
func bodytomap(c *gin.Context) (map[string]interface{}, error) {
	//将 body 流转换成 map
	data, _ := ioutil.ReadAll(c.Request.Body)

	fmt.Println(time.Now())
	fmt.Println("url: ===>", c.Request.RequestURI)
	fmt.Println("ctx.Request.body: ===>", string(data))
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
		c.JSON(http.StatusOK, gin.H{"return": "json error"})
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

	row1, err := db.Query("select p.id,p.project_name,p.description "+
		"from tb_project p,project_and_user pu "+
		"where p.id = pu.project_id and pu.user_id = ?", userId)
	if err != nil {
		//查询报错
		fmt.Println(err.Error())
	}
	joinjson := utils.Rowtojson(row1)

	c.JSON(http.StatusOK, gin.H{"self": selfjson, "join": joinjson, "return": "success"})
}

//查询项目详情 projectId
func QryProjectDetail(c *gin.Context) {
	//转换body
	params, err := bodytomap(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"return": "json error"})
		return
	}
	projectId := params["projectId"]

	rows, err := db.Query("select id,project_name,description ,create_user,create_time from tb_project where id =?", projectId)
	if err != nil {
		//查询报错
		fmt.Println(err.Error())
	}
	projectDetail := utils.Rowtojson(rows)

	userRow, err := db.Query("select u.id,u.name "+
		"from tb_project p , tb_user u, project_and_user pu "+
		"where p.id = ? and pu.project_id = p.id and u.id = pu.user_id", projectId)
	if err != nil {
		//查询报错
		fmt.Println(err.Error())
	}
	memberList := utils.Rowtojson(userRow)

	c.JSON(http.StatusOK, gin.H{"project": projectDetail, "members": memberList, "return": "success"})
}

//增删成员 projectId  userId（数组） opt（add，del）
func ManageMember(c *gin.Context) {
	//转换body
	params, err := bodytomap(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"return": "json error"})
		return
	}
	projectId := params["projectId"]
	userId := params["userId"].([]interface{})
	opt := params["opt"]

	//批量增加
	if opt == "add" {
		println("--------------add------------")
		stmt, err := db.Prepare("insert into project_and_user (id,project_id,user_id) values(?,?,?)")
		if err != nil {
			log.Fatalln(err)
		}
		//循环添加user
		for _, value := range userId {

			id := utils.UniqueId()
			res, err := stmt.Exec(id, projectId, value)
			if err != nil {
				log.Fatalln(err)
			}
			println(res)
		}
	} else if opt == "del" {
		println("------------del---------------")

		stmt, err := db.Prepare("delete from project_and_user where project_id= ? and user_id=?")
		if err != nil {
			log.Fatalln(err)
		}
		res, err := stmt.Exec(projectId, userId)
		if err != nil {
			log.Fatalln(err)
		}
		println(res)

	}
	c.JSON(http.StatusOK, gin.H{"return": "success"})
}

//查询项目模块  以及模块下的分组 接口
func QryModule(c *gin.Context) {
	//转换body
	params, err := bodytomap(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"return": "json error"})
		return
	}
	projectId := params["projectId"]
	//
	moduleRow, err := db.Query("select m.module_name,m.id from tb_module m, project_and_module pm where pm.project_id = ? and pm.module_id = m.id", projectId)
	if err != nil {
		//查询报错
		fmt.Println(err.Error())
	}
	moduleList := utils.Rowtojson(moduleRow)

	//循环模块
	for _, module := range moduleList {
		pageRow, err := db.Query("select p.id,p.page_name from tb_page p, module_and_page mp where mp.module_id = ? and mp.page_id = p.id", module["id"])
		if err != nil {
			//查询报错
			fmt.Println(err.Error())
		}
		pageList := utils.Rowtojson(pageRow)
		//循环页面
		for _, page := range pageList {
			actionRow, err := db.Query("select a.id,a.action_name from tb_action a, page_and_action pa where  pa.page_id = ? and pa.action_id = a.id", page["id"])
			if err != nil {
				//查询报错
				fmt.Println(err.Error())
			}
			actionList := utils.Rowtojson(actionRow)
			//给页面增加接口信息
			page["actionList"] = actionList
		}
		//给模块增加页面信息
		module["pageList"] = pageList
	}

	c.JSON(http.StatusOK, gin.H{"return": "success", "moduleList": moduleList})
}

//锁定项目模块
func LockModule(c *gin.Context) {

	//转换body
	params, err := bodytomap(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"return": "json error"})
		return
	}
	projectId := params["projectId"].(string)
	moduleId := params["moduleId"].(string)
	userId := params["userId"].(string)
	utils.LockKey(projectId, moduleId, userId)
	c.JSON(http.StatusOK, gin.H{"return": "success"})
}

//解锁项目模块
func UnlockModule(c *gin.Context) {
	//转换body
	params, err := bodytomap(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"return": "json error"})
		return
	}
	projectId := params["projectId"].(string)
	moduleId := params["moduleId"].(string)
	utils.UnlockKey(projectId, moduleId)
	c.JSON(http.StatusOK, gin.H{"return": "success"})

}

//管理项目模块  增删改
func ManageModule(c *gin.Context) {
	//转换body
	params, err := bodytomap(c)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"return": "json error"})
		return
	}
	projectId := params["projectId"].(string)
	opt := params["opt"].(string)

	if opt == "add" {
		id := utils.UniqueId()
		moduleName := params["moduleName"].(string)
		//新增模块
		stmt, err := db.Prepare("insert into tb_module (id,module_name,create_time) values (?,?,?)")
		if err != nil {
			log.Fatalln(err)
		}
		res, err := stmt.Exec(id, moduleName, time.Now())
		if err != nil {
			log.Fatalln(err)
		}
		println(res)
		//新增项目与模块映射
		idExt := utils.UniqueId()
		stmtExt, err := db.Prepare("insert into project_and_module (id,project_id,module_id) values (?,?,?)")
		if err != nil {
			log.Fatalln(err)
		}
		resExt, err := stmtExt.Exec(idExt, projectId, id)
		if err != nil {
			log.Fatalln(err)
		}
		println(resExt)

		c.JSON(http.StatusOK, gin.H{"return": "success"})
		return
	} else if opt == "upt" {
		moduleId := params["moduleId"].(string)
		moduleName := params["moduleName"].(string)

		stmtExt, err := db.Prepare("update tb_module set module_name = ? where id = ?")
		if err != nil {
			log.Fatalln(err)
		}
		resExt, err := stmtExt.Exec(moduleName, moduleId)
		if err != nil {
			log.Fatalln(err)
		}
		println(resExt)
		c.JSON(http.StatusOK, gin.H{"return": "success"})
		return
	} else if opt == "del" {
		moduleId := params["moduleId"].(string)
		//判断是否被锁定
		isLock, val := utils.IsLockKey(projectId, moduleId)
		//被锁返回 锁定人名字
		if isLock {
			c.JSON(http.StatusOK, gin.H{"return": "success", "user_name": val})
			return
		}

		stmtExt, err := db.Prepare("delete from tb_module  where id = ?")
		if err != nil {
			log.Fatalln(err)
		}
		resExt, err := stmtExt.Exec(moduleId)
		if err != nil {
			log.Fatalln(err)
		}
		println(resExt)
		c.JSON(http.StatusOK, gin.H{"return": "success"})
		return

	}

}

//管理分组
func ManagePage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"return": "success"})
}

////查询模块下的分组及接口列表
//func QryActionList(c *gin.Context) {
//	c.JSON(http.StatusOK, gin.H{"return": "success"})
//}

//查询接口详情
func QryActionDetail(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"return": "success"})
}

//管理接口详情  增删改
func ManageActionDetail(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"return": "success"})
}

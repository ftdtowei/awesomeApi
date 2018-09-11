package service

import (
	"awesomeApi/utils"
	"database/sql"
	. "encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//包下的全局变量  所有的service都可以调用
var db *sql.DB
var herr error

func init() {
	fmt.Println("db init")
	//打开数据库
	//DSN数据源字符串：用户名:密码@协议(地址:端口)/数据库?参数=参数值
	db, herr = sql.Open("mysql", "root:root@tcp(52.80.180.58:3306)/awesome_db?charset=utf8")
	if herr != nil {
		fmt.Println(herr)
	}
	//关闭数据库，db会被多个goroutine共享，可以不调用
	//defer db.Close()
}

//函数大写是外部可以访问得
//用户注册

/**
 *
 * 功能描述: 用户注册
 *
 * @param:
 * @return:
 * @auther: wg
 * @date:
 */
func UserRegister(c *gin.Context) {
	//将body 流转换成 map
	data, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Println("----UserRegister params----")
	fmt.Println("ctx.Request.body: %v", string(data))
	var m map[string]string
	err := Unmarshal(data, &m)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"return": "json error"})
		return
	}

	err = createUser(m)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"return": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"return": "success"})
	}

}

/**
 *
 * 功能描述: 用户登录
 *
 * @param:
 * @return:
 * @auther: wg
 * @date:
 */
func UserLogin(c *gin.Context) {

	//将body 流转换成 map
	data, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Println("----login params----")
	fmt.Println("ctx.Request.body: %v", string(data))
	var m map[string]string
	err := Unmarshal(data, &m)
	account := m["account"]
	passwd := m["passwd"]
	//登录查询用户是否存在/密码校验
	user, err := isUserExist(account, passwd)
	//用户失败
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"return": err.Error()})
	} else {
		//设置登录session
		cookie := &http.Cookie{
			Name:     "session_id",
			Value:    "aweSomeApi",
			Path:     "/",
			MaxAge:   3600 * 4, //四个小时超时
			HttpOnly: true,
		}
		http.SetCookie(c.Writer, cookie)

		c.JSON(http.StatusOK, gin.H{
			"return": "success",
			"user":   user,
		})
	}

}

type UserInfo struct {
	Id      string
	Name    string
	Account string
	Passwd  string
}

/**
 *
 * 功能描述: 用户是否存在
 *
 * @param:
 * @return:
 * @auther: wg
 * @date:
 */
func isUserExist(accountP string, passwdP string) (UserInfo, error) {
	//func isUserExist (account string , passwd string) (UserInfo,error){

	qry := "select id,name,account,passwd from tb_user where account = " + "'" + accountP + "'"
	rows, _ := db.Query(qry)

	result := UserInfo{
		Id:      "",
		Name:    "",
		Account: "",
		Passwd:  "",
	}
	id := ""
	name := ""
	account := ""
	passwd := ""

	//没有查到数据
	if rows == nil || rows.Next() != true {
		var err = errors.New("账号不存在或密码错误 ")
		return result, err
	}
	//返回第一个数据
	rows.Scan(&id, &name, &account, &passwd)
	fmt.Println(id, name, account, passwd)
	user := UserInfo{
		Id:      id,
		Name:    name,
		Account: account,
		Passwd:  passwd,
	}
	//密码不正确
	if passwdP != passwd {
		var err = errors.New("账号不存在或密码错误 ")
		return result, err
	} else if passwdP == "WrOnGpAsSwD" {
		return user, nil
	}

	return user, nil
}

/**
 *
 * 功能描述:插入用户
 *
 * @param:
 * @return:
 * @auther: wg
 * @date:
 */
func createUser(m map[string]string) error {
	id := utils.UniqueId()
	name := m["name"]
	account := m["account"]
	passwd := m["passwd"]
	createTime := time.Now()

	user, err := isUserExist(account, "WrOnGpAsSwD")

	fmt.Println(user)
	if err != nil {
		fmt.Println("--------账号已存在---------")
		var err = errors.New("账号已存在")
		return err
	}

	stmt, err := db.Prepare("insert into tb_user(id,name,account,passwd,create_time) values(?,?,?,?,?)")
	if err != nil {
		log.Fatalln(err)
	}
	res, err := stmt.Exec(id, name, account, passwd, createTime)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(res)

	return nil
}

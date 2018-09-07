package userService

import (
	"database/sql"
	. "encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"net/http"
	"awesomeProject/utils"
	"time"
)

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
func UserRegister(c *gin.Context){
	//将body 流转换成 map
	data, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Printf("ctx.Request.body: %v", string(data))
	var m map[string]string
	err := Unmarshal(data, &m)
	if(err == nil){
		c.JSON(http.StatusOK,gin.H{"return":"json error"})
		return
	}
	name:=m["name"]
	account:=m["account"]
	passwd:=m["passwd"]

	err = createUser(m);


	c.JSON(http.StatusOK,gin.H{"message":"你好，欢迎你"})
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
func UserLogin(c *gin.Context){

	//将body 流转换成 map
	data, _ := ioutil.ReadAll(c.Request.Body)
	fmt.Printf("ctx.Request.body: %v", string(data))
	var m map[string]string
	err := Unmarshal(data, &m)
	account:=m["account"]
	passwd:=m["passwd"]
	//登录查询用户是否存在/密码校验
	user,err := isUserExist(account,passwd)
	//用户失败
	if(err != nil){
		c.JSON(http.StatusOK,gin.H{"return":err.Error()})
	}else{
		//设置登录session
		cookie := &http.Cookie{
			Name:     "session_id",
			Value:    "aweSomeApi",
			Path:     "/",
			MaxAge:   3600*4,//四个小时超时
			HttpOnly: true,
		}
		http.SetCookie(c.Writer, cookie)

		c.JSON(http.StatusOK,gin.H{
			"return":"success",
			"user":user,
		})
	}


}


type UserInfo struct {
	Id string
	Name string
	Account string
	Passwd string
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
func isUserExist (accountP string , passwdP string) (UserInfo,error){
//func isUserExist (account string , passwd string) (UserInfo,error){

	//打开数据库
	//DSN数据源字符串：用户名:密码@协议(地址:端口)/数据库?参数=参数值
	db, err := sql.Open("mysql", "root:root@tcp(52.80.180.58:3306)/awesome_db?charset=utf8")
	if err != nil {
		fmt.Println(err)
	}
	//关闭数据库，db会被多个goroutine共享，可以不调用
	defer db.Close()
	//查询数据，指定字段名，返回sql.Rows结果集

	qry:= "select id,name,account,passwd from tb_user where account = accountP" +"'"+accountP+"'"
	rows, _ := db.Query(qry )

	result := UserInfo{
		Id :"",
		Name :"",
		Account :"",
		Passwd :"",
	}
	id :=""
	name :=""
	account :=""
	passwd :=""

	//没有查到数据
	if(rows == nil ||rows.Next() != true){
		var err = errors.New("账号不存在或密码错误 ")
		return result ,err
	}
	//返回第一个数据
	rows.Scan(&id, &name,&account,&passwd);
	fmt.Println(id, name,account,passwd);
	user := UserInfo{
		Id:id,
		Name:name,
		Account:account,
		Passwd:passwd,
	}
	//密码不正确
	if(passwdP!=passwd){
		var err = errors.New("账号不存在或密码错误 ")
		return result ,err
	}


	return user,nil
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
 func createUser(m map[string]string) error{
	 id:=utils.UniqueId()
	 name:=m["name"]
	 account:=m["account"]
	 passwd:=m["passwd"]
	 createTime:=time.Now()

	 //打开数据库
	 //DSN数据源字符串：用户名:密码@协议(地址:端口)/数据库?参数=参数值
	 db, err := sql.Open("mysql", "root:root@tcp(52.80.180.58:3306)/awesome_db?charset=utf8")
	 if err != nil {
		 fmt.Println(err)
	 }
	 //关闭数据库，db会被多个goroutine共享，可以不调用
	 defer db.Close()
	 //查询数据，指定字段名，返回sql.Rows结果集

	 ret, _ := db.Exec("insert into test(id,name) values(null, '444')");
	 //获取插入ID
	 ins_id, _ := ret.LastInsertId();

	 fmt.Println(ins_id);
	 return nil
 }





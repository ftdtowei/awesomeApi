package utils

import (
	"github.com/go-redis/redis"
	"time"
)

//redis 全局变量
var redisConn *redis.Client
var prefix string = "lock_"
var midfix string = "_"

//初始化redis链接
func init() {
	redisConn = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", //默认空密码
		DB:       0,  //使用默认数据库
	})
	defer redisConn.Close()

}

//加操作锁
func LockKey(project_id string, module_id string) {

	redisConn.Set(prefix+project_id+midfix+module_id, time.Now(), 0) //锁得时间确定一下
}

//判断是否被锁
func IsLockKey(project_id string, module_id string) bool {
	StringCmd, err := redisConn.Get(prefix + project_id + midfix + module_id).Result()
	println(StringCmd)
	if err != nil { //被锁
		return true
	}
	return false
}

//解锁删除key
func UnlockKey(project_id string, module_id string) {

	redisConn.Del(prefix + project_id + midfix + module_id)

}

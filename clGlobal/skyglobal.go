package clGlobal

import (
	"errors"
	"github.com/xiaolan580230/clUtil/clLog"
	"github.com/xiaolan580230/clUtil/clMysql"
	"github.com/xiaolan580230/clUtil/clRedis"
	"github.com/xiaolan580230/clhttp-framework/core/skyconfig"

)

var ServerVersion = `v1.0.0`

type SkyConfig struct {
	MgoUrl       string
	MgoDBName    string
	MgoUser      string
	MgoPass      string

	MysqlHost string
	MysqlName string
	MysqlUser string
	MysqlPass string

	MysqlMaxConnections uint32		// 数据库最大连接数
	MysqlIdleConnections uint32		// 数据库最小连接数
	MysqlIdleLifeTime uint32		// 空闲连接的存活时间

	RedisHost    string
	RedisPrefix  string
	RedisPass    string

	LogType  uint32
	LogLevel uint32

	IsCluster bool
	DebugRouter bool
}

var SkyConf SkyConfig
var mRedis *clRedis.RedisObject
var mMysql *clMysql.DBPointer
var conf *skyconfig.Config

func Init(_filename string) {

	conf = skyconfig.New(_filename, 0)

	conf.GetStr("mongodb", "mgo_url", "", &SkyConf.MgoUrl)
	conf.GetStr("mongodb", "mgo_dbname", "", &SkyConf.MgoDBName)
	conf.GetStr("mongodb", "mgo_user", "", &SkyConf.MgoUser)
	conf.GetStr("mongodb", "mgo_pass", "", &SkyConf.MgoPass)

	conf.GetStr("mysql", "mysql_host", "", &SkyConf.MysqlHost)
	conf.GetStr("mysql", "mysql_name", "", &SkyConf.MysqlName)
	conf.GetStr("mysql", "mysql_user", "", &SkyConf.MysqlUser)
	conf.GetStr("mysql", "mysql_pass", "", &SkyConf.MysqlPass)
	conf.GetUint32("mysql", "max_connections", 30, &SkyConf.MysqlMaxConnections)
	conf.GetUint32("mysql", "idle_connections", 10, &SkyConf.MysqlIdleConnections)
	conf.GetUint32("mysql", "max_life_sec", 3600*4, &SkyConf.MysqlIdleLifeTime)


	conf.GetStr("redis", "redis_host", "", &SkyConf.RedisHost)
	conf.GetStr("redis", "redis_prefix", "", &SkyConf.RedisPrefix)
	conf.GetStr("redis", "redis_password", "", &SkyConf.RedisPass)

	conf.GetStr("system", "version", "", &ServerVersion)
	conf.GetBool("system", "is_cluster", false, &SkyConf.IsCluster)
	conf.GetBool("system", "debug_router", false, &SkyConf.DebugRouter)

	if SkyConf.DebugRouter {
		clLog.Debug("%+v", SkyConf)
	}
}


// 获取redis连线
func GetRedis() *clRedis.RedisObject {
	if mRedis != nil && mRedis.Ping() {
		return mRedis
	}
	newRedis, err := clRedis.New(SkyConf.RedisHost, SkyConf.RedisPass, SkyConf.RedisPrefix)
	if err != nil {
		clLog.Error("连接redis [%v] [%v] 失败! %v", SkyConf.RedisHost, SkyConf.RedisPass, err)
		return nil
	}
	mRedis = newRedis
	return mRedis
}


// 获取mysql连线
func GetMysql() *clMysql.DBPointer {
	if mMysql != nil && mMysql.IsUsefull() {
		return mMysql
	}

	db, err := clMysql.NewWithOpt(SkyConf.MysqlHost, SkyConf.MysqlUser, SkyConf.MysqlPass, SkyConf.MysqlName, clMysql.DBOpitions{
		MaxConnection:  SkyConf.MysqlMaxConnections,
		IdleConnection: SkyConf.MysqlIdleConnections,
		MaxLifeTime:    SkyConf.MysqlIdleLifeTime,
	})
	if err != nil {
		clLog.Error("连接数据库错误: %v", err)
		clLog.Error("Host: %v User: %v Pass: %v", SkyConf.MysqlHost, SkyConf.MysqlUser, SkyConf.MysqlPass)
		return nil
	}
	mMysql = db
	return mMysql
}



// 获取事务连线
func GetMysqlTx() (*clMysql.ClTranslate, error) {
	db := GetMysql()
	if db == nil {
		return nil, errors.New("数据库连接失败")
	}
	return db.StartTrans()
}



// 获取mongodb连线
func GetMongo() {

}
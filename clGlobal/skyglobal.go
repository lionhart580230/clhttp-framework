package clGlobal

import (
	"errors"
	"github.com/lionhart580230/clUtil/clCrypt"
	"github.com/lionhart580230/clUtil/clLog"
	"github.com/lionhart580230/clUtil/clMysql"
	"github.com/lionhart580230/clUtil/clRedis"
	"github.com/lionhart580230/clhttp-framework/core/skyconfig"
	"strings"
)

var ServerVersion = `v1.0.0`

type MysqlConf struct {
	DBHost string
	DBUser string
	DBPass string
	DBName string
}

type SkyConfig struct {
	MgoUrl    string
	MgoDBName string
	MgoUser   string
	MgoPass   string

	MysqlHost string
	MysqlName string
	MysqlUser string
	MysqlPass string

	MysqlList  []MysqlConf // 数据库配置列表
	MysqlCount uint32      // 数据库数量

	MysqlMaxConnections  uint32 // 数据库最大连接数
	MysqlIdleConnections uint32 // 数据库最小连接数
	MysqlIdleLifeTime    uint32 // 空闲连接的存活时间

	RedisHost   string
	RedisPrefix string
	RedisPass   string

	LogType  uint32
	LogLevel uint32

	IsCluster   bool
	DebugRouter bool
}

var SkyConf SkyConfig
var mRedis *clRedis.RedisObject
var mMysqlList []*clMysql.DBPointer
var conf *skyconfig.Config

func Init(_filename string) {

	conf = skyconfig.New(_filename, 0)

	//conf.GetStr("mongodb", "mgo_url", "", &SkyConf.MgoUrl)
	//conf.GetStr("mongodb", "mgo_dbname", "", &SkyConf.MgoDBName)
	//conf.GetStr("mongodb", "mgo_user", "", &SkyConf.MgoUser)
	//conf.GetStr("mongodb", "mgo_pass", "", &SkyConf.MgoPass)

	conf.GetUint32("mysql", "max_connections", 30, &SkyConf.MysqlMaxConnections)
	conf.GetUint32("mysql", "idle_connections", 10, &SkyConf.MysqlIdleConnections)
	conf.GetUint32("mysql", "max_life_sec", 3600*4, &SkyConf.MysqlIdleLifeTime)

	conf.GetStr("redis", "redis_host", "", &SkyConf.RedisHost)
	conf.GetStr("redis", "redis_prefix", "", &SkyConf.RedisPrefix)
	conf.GetStr("redis", "redis_password", "", &SkyConf.RedisPass)

	conf.GetStr("system", "version", "", &ServerVersion)
	conf.GetBool("system", "is_cluster", false, &SkyConf.IsCluster)
	conf.GetBool("system", "debug_router", false, &SkyConf.DebugRouter)

	var mysqlEncryptStr string
	conf.GetStr("mysql", "connection_str", "", &mysqlEncryptStr)

	// 新版数据库连线
	//conf.GetUint32("mysql", "count", 0, &SkyConf.MysqlCount)
	//if SkyConf.MysqlCount > 0 {
	if mysqlEncryptStr != "" {
		SkyConf.MysqlList = make([]MysqlConf, 0)
		clLog.Debug("获取到mysql_connection_str: %v", mysqlEncryptStr)
		mysqlConnDecode := DecryptMysql(mysqlEncryptStr)
		clLog.Debug("解密后: %v", mysqlConnDecode)
		mysqlConnItems := strings.Split(mysqlConnDecode, "$$")
		mMysqlList = make([]*clMysql.DBPointer, len(mysqlConnItems))
		for _, mysqlConnItem := range mysqlConnItems {
			var MysqlConnectsItem = strings.Split(mysqlConnItem, "|")
			if len(MysqlConnectsItem) >= 4 {
				SkyConf.MysqlList = append(SkyConf.MysqlList, MysqlConf{
					DBHost: MysqlConnectsItem[0],
					DBUser: MysqlConnectsItem[1],
					DBPass: MysqlConnectsItem[2],
					DBName: MysqlConnectsItem[3],
				})
			}
		}
	} else {
		// 走旧的读取方式
		var MysqlHost, MysqlName, MysqlUser, MysqlPass string
		conf.GetStr("mysql", "mysql_host", "", &MysqlHost)
		conf.GetStr("mysql", "mysql_name", "", &MysqlName)
		conf.GetStr("mysql", "mysql_user", "", &MysqlUser)
		conf.GetStr("mysql", "mysql_pass", "", &MysqlPass)
		SkyConf.MysqlList = make([]MysqlConf, 1)
		mMysqlList = make([]*clMysql.DBPointer, 1)
		SkyConf.MysqlList[0] = MysqlConf{
			DBHost: MysqlHost,
			DBUser: MysqlUser,
			DBPass: MysqlPass,
			DBName: MysqlName,
		}
	}

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
	return GetMysqlById(0)
}

// 获取mysql连线
func GetMysqlById(_id int) *clMysql.DBPointer {
	if _id >= len(mMysqlList) {
		_id = 0
	}
	if mMysqlList[_id] != nil && mMysqlList[_id].IsUsefull() {
		return mMysqlList[_id]
	}
	DBHost := SkyConf.MysqlList[_id].DBHost
	DBUser := SkyConf.MysqlList[_id].DBUser
	DBPass := SkyConf.MysqlList[_id].DBPass
	DBName := SkyConf.MysqlList[_id].DBName
	db, err := clMysql.NewWithOpt(DBHost, DBUser, DBPass, DBName, clMysql.DBOpitions{
		MaxConnection:  SkyConf.MysqlMaxConnections,
		IdleConnection: SkyConf.MysqlIdleConnections,
		MaxLifeTime:    SkyConf.MysqlIdleLifeTime,
	})
	if err != nil {
		clLog.Error("连接数据库错误: %v", err)
		clLog.Error("Host: %v User: %v Pass: %v", DBHost, DBUser, DBPass)
		return nil
	}
	mMysqlList[_id] = db
	return mMysqlList[_id]
}

// 获取事务连线
func GetMysqlTx() (*clMysql.ClTranslate, error) {
	db := GetMysql()
	if db == nil {
		return nil, errors.New("数据库连接失败")
	}
	return db.StartTrans()
}

// 加密Mysql
const MysqlEncryptKey = "DFRvxOaSPr6Btwc0"
const MysqlEncryptIV = "5qhdOgVyBSlpeQ5a"

func EncryptMysql(_p string) string {
	return clCrypt.AesCBCEncode(_p, MysqlEncryptKey, MysqlEncryptIV)
}

func DecryptMysql(_p string) string {
	return clCrypt.AesCBCDecode([]byte(_p), []byte(MysqlEncryptKey), []byte(MysqlEncryptIV))
}

package skyredis

import (
	"github.com/xiaolan580230/clhttp-framework/clCommon"
	"github.com/xiaolan580230/clhttp-framework/core/skylog"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"strings"
	"time"
)

type RedisObject struct {
	myredis   *redis.Client
	prefix    string
	isCluster bool
}

type RedisCacheInfo struct {
	Data   string `json:"data"`
	Expire uint32 `json:"expire"`
	Sign   string `json:"sign"`
}

//@author xiaolan
//@lastUpdate 2019-08-05
//@comment 创建一个新的redis对象
func New(_addr string, _password string, _prefix string) (*RedisObject, error) {

	client := redis.NewClient(&redis.Options{
		Addr:        _addr,
		PoolSize:    10,
		PoolTimeout: 30 * time.Second,
		Password:    _password,
	})

	_, err := client.Ping().Result()
	if err != nil {
		return nil, err
	}

	clrd := &RedisObject{
		myredis:   client,
		prefix:    _prefix,
		isCluster: false,
	}

	return clrd, nil
}

//@author xiaolan
//@lastUpdate 2019-08-05
//@comment 关闭redis连线
func (this *RedisObject) Close() {

	if this.myredis != nil {
		this.myredis.Close()
	}
}

//@author xiaolan
//@lastUpdate 2019-08-05
//@comment 删除一个key
//@param _key key的名称
func (this *RedisObject) Del(_key string) error {

	keys := _key
	if this.prefix != "" {
		keys = this.prefix + "_" + _key
	}

	i := this.myredis.Del(keys)
	return i.Err()
}

//@author xiaolan
//@lastUpdate 2019-08-05
//@comment 设置一个key
//@param _key key的名字
//@param _val key对应的值
//@param _expire 有效期
func (this *RedisObject) Set(_key string, _val interface{}, _expire int32) error {

	keys := _key
	if this.prefix != "" {
		keys = this.prefix + "_" + _key
	}
	err := this.myredis.Set(keys, buildRedisValue(keys, uint32(_expire), _val),
		time.Duration(time.Second*time.Duration(_expire))).Err()
	return err
}

//@author xiaolan
//@lastUpdate 2019-08-05
//@comment 获取一个key的值
//@param _key 需要获取的值
func (this *RedisObject) Get(_key string) string {

	keys := _key
	if this.prefix != "" {
		keys = this.prefix + "_" + _key
	}
	resp := this.myredis.Get(keys)
	result := checkRedisValid(keys, resp)
	if result == "" {
		this.myredis.Del(keys)
	}
	return result
}

// 设置的key不会过期
func (this *RedisObject) SetNoExpire(_key string, _val interface{}) error {

	keys := _key
	if this.prefix != "" {
		keys = this.prefix + "_" + _key
	}
	cmd := this.myredis.Set(keys, _val, 0)
	return cmd.Err()
}

// 获取永不过期的key
func (this *RedisObject) GetNoExpire(_key string) (result string) {

	keys := _key
	if this.prefix != "" {
		keys = this.prefix + "_" + _key
	}
	result = this.myredis.Get(keys).Val()
	if result == "" {
		this.myredis.Del(keys)
	}
	return
}

//@author xiaolan
//@lastUpdate 2019-08-05
//@comment 获取没有前缀的key
//@param _key 获取的key名称
func (this *RedisObject) GetNoPrefix(_key string) string {

	resp := this.myredis.Get(_key)
	result := checkRedisValid(_key, resp)
	if result == "" {
		this.myredis.Del(_key)
	}
	return result
}

//@author xiaolan
//@lastUpdate 2019-08-05
//@comment 设置一个hash结构
//@param _key 设置的key
//@param _field 设置的key的字段
//@param _value 设置的值
//@param _expire 有效期
func (this *RedisObject) HSet(_key string, _field string, _value interface{}, _expire uint32) bool {

	keys := _key
	if this.prefix != "" {
		keys = this.prefix + "_" + _key
	}
	value := buildRedisValue(keys+_field, _expire, _value)
	rest := this.myredis.HSet(keys, _field, value)
	if rest == nil {
		return false
	}

	if _, err := rest.Result(); err != nil {
		return false
	}

	return true
}

//@author xiaolan
//@lastUpdate 2019-08-05
//@comment 获取一个hash结构的值
//@param _key 获取hash结构的key
//@param _field 获取hash结构的key的字段
func (this *RedisObject) HGet(_key string, _field string) string {

	keys := _key
	if this.prefix != "" {
		keys = this.prefix + "_" + _key
	}

	resp := this.myredis.HGet(keys, _field)
	result := checkRedisValid(keys+_field, resp)
	if result == "" {
		this.myredis.HDel(keys, _field)
	}
	return result
}

//@author xiaolan
//@lastUpdate 2019-08-05
//@comment 获取删除一个hash结构
//@param _key 获取hash结构的key
//@param _field 获取hash结构的字段
func (this *RedisObject) HDel(_key string, _field string) bool {

	keys := _key
	if this.prefix != "" {
		keys = this.prefix + "_" + _key
	}
	resp := this.myredis.HDel(keys, _field)
	return resp.Val() > 0
}

//@author xiaolan
//@lastUpdate 2019-08-05
//@comment 获取全部的key
//@param _key 获取的key名称
//@param _prefix 获取的前缀
func (this *RedisObject) HGetKeys(_key string, _prefix string) []string {

	keys := _key
	if this.prefix != "" {
		keys = this.prefix + "_" + _key
	}

	val := this.myredis.HKeys(keys)
	if val == nil {
		return []string{}
	}
	resp := make([]string, 0)
	for _, val := range val.Val() {
		if strings.HasPrefix(val, _prefix) {
			resp = append(resp, val)
		}
	}
	return resp
}

//@author xiaolan
//@lastUpdate 2019-08-05
//@comment 删除指定开头的keys
//@param _key 删除的key
//@param _prefix 需要删除的key的前缀
func (this *RedisObject) HDelKeys(_key string, _prefix string) {

	keys := _key
	if this.prefix != "" {
		keys = this.prefix + "_" + _key
	}
	keylist := this.HGetKeys(keys, _prefix)
	if len(keylist) > 0 {
		this.myredis.HDel(keys, keylist...)
	}
}

//@author xiaolan
//@lastUpdate 2019-08-05
//@comment 获取全部的hash字段
//@param _key 获取所有字段的key
func (this *RedisObject) HGetAll(_key string) map[string]string {

	keys := _key
	if this.prefix != "" {
		keys = this.prefix + "_" + _key
	}

	val := this.myredis.HGetAll(keys)
	return checkRedisValidMap(keys, val)
}

//@author xiaolan
//@lastUpdate 2019-08-05
//@comment 设置一个分布式锁
//@param _key 锁的名称
//@param _value 锁的内容
//@param _expire 锁的持续时间
func (this *RedisObject) SetNx(_key string, _value interface{}, _expire uint32) bool {

	keys := _key
	if this.prefix != "" {
		keys = this.prefix + "_" + _key
	}
	value := buildRedisValue(keys, _expire, _value)
	rest := this.myredis.SetNX(keys, value, time.Duration(_expire)*time.Second)
	if rest == nil {
		return false
	}

	if _, err := rest.Result(); err != nil {
		return false
	}

	return rest.Val()
}

//@author xiaolan
//@lastUpdate 2019-08-05
//@comment 检验redis缓存是否有效
//@param _keys string redis缓存的键名
//@param _targetData *StringCmd 目标数据
func checkRedisValidMap(_keys string, _targetData *redis.StringStringMapCmd) map[string]string {
	if _targetData == nil || len(_targetData.Val()) == 0 {
		return nil
	}

	resp := make(map[string]string)
	for key, val := range _targetData.Val() {

		var redisInfo RedisCacheInfo
		err := json.Unmarshal([]byte(val), redisInfo)
		if err != nil {
			return nil
		}

		// 缓存到期
		if redisInfo.Expire > 0 && redisInfo.Expire < uint32(time.Now().Unix()) {
			continue
		}

		// 签名不匹配
		sign := clCommon.Md5([]byte("Cache:__" + _keys + key))
		if redisInfo.Sign != sign {
			continue
		}

		resp[key] = redisInfo.Data
	}
	return resp
}

// 检验redis缓存是否有效
// @param keys string redis缓存的键名
// @param targetData *StringCmd 目标数据
func checkRedisValid(_keys string, targetData *redis.StringCmd) string {
	if targetData == nil || targetData.Val() == "" {
		return ""
	}

	var redisInfo RedisCacheInfo
	err := json.Unmarshal([]byte(targetData.Val()), &redisInfo)
	if err != nil {
		skylog.LogErr( "json.Unmarshl error: %v", err)
		return ""
	}

	// 缓存到期
	if redisInfo.Expire > 0 && redisInfo.Expire < uint32(time.Now().Unix()) {
		skylog.LogErr( "expire: %v < %v", redisInfo.Expire, uint32(time.Now().Unix()))

		return ""
	}

	// 签名不匹配
	sign := clCommon.Md5([]byte("Cache:__" + _keys))
	if redisInfo.Sign != sign {
		skylog.LogErr( "check sign error: %v != %v", redisInfo.Sign, sign)
		return ""
	}

	return redisInfo.Data
}

// 组装缓存的值
func buildRedisValue(_keys string, expire uint32, data interface{}) string {

	cache_data, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	resp, err := json.Marshal(RedisCacheInfo{
		Data:   string(cache_data),
		Expire: uint32(time.Now().Unix()) + expire,
		Sign:   clCommon.Md5([]byte("Cache:__" + _keys)),
	})

	if err != nil {
		return ""
	}
	return string(resp)
}

// 删除指定用户缓存
func (this *RedisObject) DelUserCache(uid uint32) {

	this.Del("USER" + "_" + string(uid))
}

//@author xiaolan
//@lastUpdate 2019-08-05
//@comment 删除指定接口缓存
//@param apiname string 接口名称
//@param uid uint32 用户的uid，如果为0则删除全部
func (this *RedisObject) DelApiCache(_apiname string, _uid uint32) {

	if _uid == 0 {
		this.Del(_apiname)
		return
	}
	this.HDelKeys(_apiname, fmt.Sprintf("U%v_", _uid))
}

func (this *RedisObject) Ping() bool {
	if this.myredis == nil {
		return false
	}
	return this.myredis.Ping().Err() == nil
}

#### 如何获取clhttp-framework?

```
go get github.com/lionhart580230/clhttp-framework
```


#### v1.4.14 更新内容
1. 修复了接口配置 `respContent` 的一些错误逻辑
2. 移除部分与clUtil包重复冗余的代码
3. 为了支持更加自定义化的参数配置，添加如下参数类型:
```
rule.AddRule(rule.Rule{
		Request: "request",
		Name: "api_example",
		Params: []rule.ParamInfo{
			// 参数名为id，它必须是整数，并且值范围必须在1到10之间
			rule.NewIntParamRange("id", true, "1", 1, 10),
			// 参数名为name, 它必须是字符串，并且这个字符串的长度必须为2到5之间
			rule.NewStrParamRange("name", true, "", 2, 5),
		},
		Login: true,
		CallBack: example.ApiExample,
		CacheExpire: 180,
	})
```
通过这种方式，可以更加严格的对参数进行控制

4. 鉴于可能存在前后台使用同一个redis，并且前后台用户的 `redis key` 冲突的可能性, 增加一个配置用于规避这个问题, 可以在程序入口处添加:
```
    clAuth.SetAuthPrefix("U_INFO")
```
这样即可将当前服务的用户在REDIS中以 `U_INFO` 前缀，这个前缀可以任意设置，建议不要超过128个字符

5. 由于之前版本默认开启了上传功能，这对原本不需要上传功能的服务来说，可能存在安全隐患。所以增加了这个功能的开关控制.
```
    // 关闭上传功能
	httpserver.SetEnableUploadFile(false)
	// 关闭上传调试页
	httpserver.SetEnableUploadTest(false)
```
一样，你需要在程序入口处添加这个代码进行关闭

6. 为了应对一些比较复杂的参数，减少开发代码量，特提供了如下三个获取参数的方式:
```
    // 获取字符串列表数组
	strArr := _param.GetStrSplit("str_list", ",")

	// 获取整数列表数组
	numArr := _param.GetUint32Split("id_list")

	// 获取浮点数列表数组
	posArr := _param.GetFloatSplit("pos_list")
```
这只是一个简单的语法糖，有了它就不用自己获取，手动分割，再手动转化了

### 可用配置列表
```yaml
environments:
  - MYSQL_MAX_CONNECTIONS=30  // MYSQL的最大连接数
  - MYSQL_IDLE_CONNECTIONS=10 // MYSQL最大空闲连接数
  - MYSQL_MAX_LIFE_SEC=12400  // MYSQL连线最大存活时间（秒）
  - MYSQL_COUNT=1       // 数据库连线数量（必须跟配置统一）
  - MYSQL_CONNECTIONS=  // 数据库配置加密文本
```
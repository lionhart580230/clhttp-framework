### 网络请求回应包

原本通过`clCommon`包的`JCode`方法进行对客户端请求的回应
但是`clCommon`包应该是处理一些通用方法的，把回应放在里面有点混淆，于是将其从clCommon包中抽离独立一个`clResponse`包
通过封装一些常用的回应类型来简化编程

之前的写法：
```
return clCommon.JCode(0, "ok", nil)
```

现在的写法:
```
return clResponse.JCode(0, "ok", nil)
```

> 以上各位可以使用批量替换一次性替换过来

除了这个改变还添加了很多通用的回应:
```

// 标示成功, data为要返回给前端的数据
clResponse.Success( data )

// 标示失败, 此方法与JCode类似，但是从名字上可以更直观的知道发生了错误
// code 代表错误代码 > 0
// param 代表错误信息
// data 代表返回的一些扩展信息
clResponse.Failed(code, param, data)

// 一些系统错误简单返回
clResponse.SystemError()

// 一些请求错误简单返回
clResponse.ServerError()

// 需要前台登录
clResponse.NotLogin()

// 前台操作过于频繁
clResponse.TooQuickly()

```

同时提供一个多语言支持的内置方法 `JCodeByLang` 此方法可以支持多语言切换
用它之前，你需要先给他注入一个多语言字典(`clResponse.ImportI18N`)，以便他能够根据langType和code找到对应的文本并进行输出, 首先你需要一个配置文件来配置code和字符串之间的关系，比如:
```
1,系统错误
2,你好,%v
3,测试测试
```

然后将这个文本保存为xxx.csv 放置在项目目录中，假设我放于`./language/zhcn.csv`
然后我可以通过如下代码进行字典导入:
```

const LangTypeChinese = 0
buffer, err := ioutil.ReadFile("./language/zhcn.csv")
if err != nil {
    skylog.LogErr("读取文件失败! 错误:%v", err)
    return
}

clResponse.ImportI18N(LangTypeChinese, strings.Split(string(buffer), "\n"))
```

然后我就可以很方便的进行语言包管理了
```
return clResponse.JCodeByLang(LangTypeChinese, 1, nil)
```
以上代码就可以得到:
```
{
   "code": 1,
   "msg": "系统错误",
   "data": nil
}
```

如果文本中带有参数的我们可以如此:
```
return clResponse.JCodeByLang(LangTypeChinese, 2, nil, "张三")
```
以上代码就可以得到:
```
{
   "code": 2,
   "msg": "你好,张三",
   "data": nil
}
```

package httpserver

import (
	"fmt"
	"github.com/lionhart580230/clUtil/clCrypt"
	"github.com/lionhart580230/clUtil/clJson"
	"github.com/lionhart580230/clUtil/clLog"
	"github.com/lionhart580230/clhttp-framework/clCommon"
	"github.com/lionhart580230/clhttp-framework/core/rule"
	"log"
	"net/http"
	"os"
	"strings"
)

//var TempDirPath string

// 是否启用上传测试页面
var mEnableUploadTest = false

// 是否启用上传文件
var mEnableUploadFile = false

// 设置AES 密钥
var mAesKey = ""

// 启用上传测试页面的访问
// 访问url为 http://your_domain/upload_test
func SetEnableUploadTest(_enable bool) {
	mEnableUploadTest = _enable
}

// 启用上传文件功能，启用后将允许前台上传文件到服务器上
// 访问url为 http://your_domain/upload
func SetEnableUploadFile(_enable bool) {
	mEnableUploadFile = _enable
}

// 设置解密密钥
func SetAESKey(_aesKey string) {
	mAesKey = _aesKey
}

// @author xiaolan
// @lastUpdate 2019-08-09
// @comment 启动HTTP服务
func StartServer(_listenPort uint32) {

	http.HandleFunc("/", rootHandler)
	if mEnableUploadTest {
		http.HandleFunc("/upload_test", uploadFileHtml)
	}
	if mEnableUploadFile {
		http.HandleFunc("/uploadFile", uploadFile)
	}
	http.HandleFunc("/mysql_conf_encrypt", mysqlConfEncrypt)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", _listenPort), nil))
}

var uploadFileSizeLimit int64 = 1024 * 1024 * 300

func SetUploadFileSizeLimit(_limit int64) {
	uploadFileSizeLimit = _limit
}

// @author xiaolan
// @lastUpdate 2019-08-09
// @comment http请求主入口回调
func rootHandler(rw http.ResponseWriter, rq *http.Request) {

	// 跨域支持
	rw.Header().Set("Access-Control-Allow-Origin", "*")  //允许访问所有域
	rw.Header().Add("Access-Control-Allow-Headers", "*") //header的类型
	rw.Header().Set("Access-Control-Expose-Headers", "Encrypt-Type,Encrypt-iv")

	// 需要过滤的请求文件类型列表
	filter_file_ext := []string{
		".ico", ".png", ".jpg", ".gif", ".js", ".css", ".html",
	}

	if strings.ToUpper(rq.Method) == "OPTIONS" {
		rw.WriteHeader(200)
		return
	}

	// 判断路由
	requestURI := strings.Split(rq.RequestURI, "?")
	requestArr := strings.Split(requestURI[0], "/")
	if len(requestArr) < 2 || requestArr[1] == "" {
		rw.WriteHeader(404)
		return
	}
	requestName := requestArr[1]
	actionName := ""
	if len(requestArr) > 2 {
		actionName = requestArr[2]
	}
	// 过滤请求
	for _, filterExt := range filter_file_ext {
		if strings.HasSuffix(requestArr[1], filterExt) {
			rw.WriteHeader(404)
			return
		}
	}

	var isEncrypt = rq.Header.Get("Encrypt-Type") == "AES"
	var iv = rq.Header.Get("Encrypt-iv")

	var contentType = strings.ToLower(rq.Header.Get("Content-Type"))
	var values = make(map[string]string)

	getData := rq.URL.Query()
	var uriValues = make(map[string]string)
	for key, val := range getData {
		uriValues[key] = val[0]
		values[key] = val[0]
	}

	var rawData = ""
	if strings.Contains(contentType, "text/json") || strings.Contains(contentType, "application/json") {
		var jsonBytes = make([]byte, 10*1024)

		n, err := rq.Body.Read(jsonBytes)
		if err != nil && err.Error() != "EOF" {
			clLog.Error("读取json参数失败! 错误:%v", err)
			rw.WriteHeader(502)
			return
		}
		var jsonStr = jsonBytes[:n]
		if isEncrypt {
			jsonStr = []byte(clCrypt.AesCBCDecode(jsonStr, []byte(mAesKey), []byte(iv)))
		}
		if jsonStr == nil || len(jsonStr) == 0 {
			clLog.Error("数据: %v 结构化失败! 加密:%v 长度:%v", string(jsonStr), isEncrypt, n)
			rw.WriteHeader(502)
			return
		}
		jsonObj := clJson.New(jsonStr)
		if jsonObj != nil {
			{
				jsonMap := jsonObj.ToMap().ToCustom()
				for k, v := range jsonMap {
					values[k] = v
				}
			}
		}
		rawData = string(jsonBytes[:n])
	} else if strings.Contains(contentType, "multipart/form-data") || strings.Contains(contentType, "application/x-www-form-urlencoded") {
		rq.ParseForm()
		if len(rq.Form) > 0 {
			for key, val := range rq.Form {
				if len(val) == 1 {
					values[key] = val[0]
				}
			}
		}
		if len(rq.PostForm) > 0 {
			for key, val := range rq.PostForm {
				if len(val) == 1 {
					values[key] = val[0]
				}
			}
		}
	} else if strings.Contains(contentType, "text/xml") {
		var xmlBytes = make([]byte, 4096)
		n, err := rq.Body.Read(xmlBytes)
		if err != nil && err.Error() != "EOF" {
			clLog.Error("读取json参数失败! 错误:%v", err)
			rw.WriteHeader(502)
			return
		}
		rawData = string(xmlBytes[:n])
	} else if rq.Method == "GET" && len(requestURI) > 1 {
		GetParams := strings.Split(requestURI[1], "&")
		for _, paramItem := range GetParams {
			params := strings.Split(paramItem, "=")
			if len(params) < 2 {
				continue
			}
			values[params[0]] = params[1]
		}
	} else {
		rq.ParseMultipartForm(2 << 32)
		if nil != rq.MultipartForm && nil != rq.MultipartForm.Value {
			for key, val := range rq.MultipartForm.Value {
				if len(val) == 1 {
					values[key] = val[0]
				}
			}
		}
	}
	var rqObj = rule.NewHttpParam(values)

	remoteip := rq.Header.Get("X-Forwarded-For")
	if remoteip == "" {
		remoteip = rq.Header.Get("X-Real-Ip")
		if remoteip == "" {
			remoteip = rq.RemoteAddr
		}
	} else {
		remotes := strings.Split(remoteip, ",")
		remoteip = remotes[0]
	}
	request_url := ""
	proctol := rq.Header.Get("Proxy-X-Forwarded-Proto")
	if proctol == "" {
		proctol = rq.Header.Get("X-Forwarded-Proto")
	}

	request_url = "https://" + rq.Host + rq.RequestURI

	myUA := rq.Header.Get("Platform")
	if myUA == "" {
		myUA = "web"
	}

	myLang := rq.Header.Get("Lang-Type")
	if myLang == "" {
		myLang = "1"
	}
	if strings.HasPrefix(remoteip, "[::1]") {
		remoteip = strings.ReplaceAll(remoteip, "[::1]", "127.0.0.1")
	}
	RemoteIpArr := strings.Split(remoteip, ":")
	remoteIp := RemoteIpArr[0]
	if len(RemoteIpArr) > 2 {
		remoteIp = remoteip
	}

	// 获取header
	headerAuthorization := rq.Header.Get("Authorization")
	if headerAuthorization != "" {
		authorArr := strings.Split(headerAuthorization, ":")
		if len(authorArr) == 2 {
			rqObj.Add("uid", authorArr[0])
			rqObj.Add("token", authorArr[1])
		}
	}

	var serObj = rule.ServerParam{
		RemoteIP:    remoteIp,
		RequestURI:  rq.RequestURI,
		UriData:     rule.NewHttpParam(uriValues),
		Host:        rq.Host,
		Method:      rq.Method,
		Header:      rq.Header,
		RequestURL:  request_url,
		UA:          myUA,
		UAType:      UAToInt(myUA),
		Proctol:     proctol,
		Port:        "",
		Language:    myLang,
		LangType:    clCommon.Uint32(myLang),
		ContentType: rq.Header.Get("Content-Type"),
		RawData:     rawData,
		RawParam:    rqObj,     // 原始参数
		Encrypt:     isEncrypt, // 是否加密
		AesKey:      mAesKey,
		Iv:          iv,
		JWT:         rq.Header.Get("Jwt"),
		AcName:      actionName,
	}
	content, contentType := CallHandler(rq, &rw, requestName, rqObj, &serObj)
	if contentType == "" {
		contentType = "text/json"
	}

	if contentType == "error" {
		rw.WriteHeader(502)
		return
	}

	if isEncrypt {
		rw.Header().Set("Encrypt-Type", "AES")
		rw.Header().Set("Encrypt-iv", iv)
	}
	rw.Header().Set("Content-Type", contentType)
	rw.Header().Set("Charset", "UTF-8")
	rw.Write([]byte(content))
}

// 测试上传页面
func uploadFileHtml(rw http.ResponseWriter, rq *http.Request) {

	rw.Write([]byte(`
	<html>  
  <head>  
    <title>选择文件</title>
  </head>  
  <body>  
    <form enctype="multipart/form-data" action="/uploadFile" method="post">  
      <input type="file" name="uploadfile" />  
      <input type="submit" value="上传文件" />  
    </form>  
  </body>  
</html>
`))
}

// 数据库配置加密页面
func mysqlConfEncrypt(rw http.ResponseWriter, rq *http.Request) {

	rw.Write([]byte(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>数据库配置文件加密</title>
    <style>
      *{margin: 0px;padding:0px}
      .main_box { padding: 6px; }
      .tip_title{ font-size:14px; color:#ff3333; font-weight: bold}
      .code_format {font-size: 14px; font-style: italic;}
      #code{ width: 100%; height: 200px; resize: none; }
      #encrypt_code{ width: 100%; height: 100px; resize: none; }
      #btn_encrypt{ width: 100%; height: 34px; font-weight: bold; line-height: 34px; color:#ffffff; background-color: #0e79eb; border-radius:4px; text-align: center;}
      #btn_encrypt:hover{ background-color: #2b6dde; user-select: none;}
    </style>
</head>
<body>
  <div class="main_box">
    <h2>请将数据库配置填入以下文本中并点击加密</h2>
    <p class="tip_title">* 数据库连线配置格式:</p>
    <p class="code_format">数据库地址:端口|数据库账号|数据库密码|数据库名称\n数据库地址:端口|数据库账号|数据库密码|数据库名称</p>
    <textarea id="code">127.0.0.1:3306|root|asdasd001|dbName
127.0.0.2:3306|root|asdasd001|dbName</textarea>
    <p>加密结果:</p>
    <textarea id="encrypt_code"></textarea>
    <div id="btn_encrypt">加密</div>
  </div>
  <script>
    var textBox = document.getElementById("code")
    document.getElementById("btn_encrypt").addEventListener("click", ()=>{
        var xhr = new XMLHttpRequest()
        xhr.onreadystatechange = function () {
            if (xhr.readyState == 4 && xhr.status == 200) {
                var jsobj = JSON.parse(xhr.responseText);
                if (jsobj.code > 0) {
                    alert(jsobj.msg)
                } else {
                    document.getElementById("encrypt_code").value = jsobj.data;
                }
            }
        }
        xhr.open("POST","/sys", true);
        xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded")
        xhr.send("ac=mysql_encrypt&p=" + textBox.value.split("\n").join("$$"));
    })
  </script>
</body>
</html>`))
}

// 上传文件
func uploadFile(rw http.ResponseWriter, rq *http.Request) {

	// 跨域支持
	rw.Header().Set("Access-Control-Allow-Origin", "*")  //允许访问所有域
	rw.Header().Add("Access-Control-Allow-Headers", "*") //header的类型

	if strings.ToUpper(rq.Method) == "OPTIONS" {
		rw.WriteHeader(200)
		return
	}

	// 判断路由
	requestURI := strings.Split(rq.RequestURI, "?")
	requestArr := strings.Split(requestURI[0], "/")
	if len(requestArr) < 2 || requestArr[1] == "" {
		rw.WriteHeader(404)
		return
	}

	// 最大内存限制为10MB
	rq.ParseMultipartForm(uploadFileSizeLimit)
	clientfd, handler, err := rq.FormFile("uploadfile")
	if err != nil {
		clLog.Error("加载文件失败: %v", err)
		rw.WriteHeader(502)
		return
	}
	defer clientfd.Close()

	fileNameArr := strings.Split(handler.Filename, ".")
	fileExt := fileNameArr[len(fileNameArr)-1]

	request_url := ""
	proctol := rq.Header.Get("Proxy-X-Forwarded-Proto")
	if proctol == "" {
		proctol = rq.Header.Get("X-Forwarded-Proto")
	}
	port := rq.Header.Get("X-Forwarded-Port")
	if proctol == "https" {
		request_url = "https://" + rq.Host
		if port != "443" {
			request_url += ":" + rq.Header.Get("X-Forwarded-Port")
		}
		request_url += rq.RequestURI
		proctol = "https"
	} else {
		request_url = "http://" + rq.Host
		if port != "80" {
			request_url += ":" + rq.Header.Get("X-Forwarded-Port")
		}
		request_url += rq.RequestURI
	}

	myUA := rq.Header.Get("Bxvip-Ua")
	if myUA == "" {
		myUA = "web"
	}

	remoteip := rq.Header.Get("X-Forwarded-For")
	if remoteip == "" {
		remoteip = rq.Header.Get("X-Real-Ip")
		if remoteip == "" {
			remoteip = rq.RemoteAddr
		}
	} else {
		remotes := strings.Split(remoteip, ",")
		remoteip = remotes[0]
	}

	buffers := make([]byte, uploadFileSizeLimit)
	lenOfBuffer, err := clientfd.Read(buffers)
	if err != nil {
		clLog.Error("读取文件内容失败: %v", err)
		rw.WriteHeader(502)
		return
	}
	fileName := fmt.Sprintf("%v.%v", clCommon.Md5(buffers[:lenOfBuffer]), fileExt)

	localfd, err := os.OpenFile(os.TempDir()+fileName, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		clLog.Error("打开文件失败! %v", err)
		rw.WriteHeader(502)
		return
	}
	localfd.Write(buffers[:lenOfBuffer])
	localfd.Close()

	var rqObj = rule.NewHttpParam(map[string]string{
		"ac":        "UploadFile",            // 执行的动作
		"filename":  fileName,                // 文件名
		"fileExt":   fileExt,                 // 文件扩展名
		"localPath": os.TempDir() + fileName, // 本地路径
	})

	// 获取header
	headerAuthorization := rq.Header.Get("Authorization")
	if headerAuthorization != "" {
		authorArr := strings.Split(headerAuthorization, ":")
		if len(authorArr) == 2 {
			rqObj.Add("uid", authorArr[0])
			rqObj.Add("token", authorArr[1])
		}
	}

	var serObj = rule.ServerParam{
		RemoteIP:   strings.Split(remoteip, ":")[0],
		RequestURI: rq.RequestURI,
		Host:       rq.Host,
		Method:     rq.Method,
		Header:     rq.Header,
		RequestURL: request_url,
		UA:         myUA,
		UAType:     UAToInt(myUA),
		Proctol:    proctol,
		Language:   "zh-cn",
	}

	content, contentType := CallHandler(rq, &rw, "upload", rqObj, &serObj)
	if contentType == "" {
		contentType = "text/json"
	}

	rw.Header().Set("Content-Type", contentType)
	rw.Write([]byte(content))
	return
}

// @author xiaolan
// @lastUpdate 2019-08-11
// @comment UA转整数型
func UAToInt(_ua string) uint32 {
	switch _ua {
	case "web":
		return 1
	case "android":
		return 2
	case "ios":
		return 3
	default:
		return 0
	}
}

// @author xiaolan
// @lastUpdate 2019-08-10
// @comment 调用回调函数, 成功返回文档类型(默认为text/json)和回传结果,失败返回空字符串
// @param _name 回调函数名称
// @param _param 参数列表
// @param _serconf 服务器信息
func CallHandler(rq *http.Request, rw *http.ResponseWriter, _name string, _param *rule.HttpParam, _serconf *rule.ServerParam) (string, string) {
	return rule.CallRule(rq, rw, _name, _param, _serconf)
}

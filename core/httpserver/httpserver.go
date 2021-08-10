package httpserver

import (
	"fmt"
	"github.com/xiaolan580230/clhttp-framework/clCommon"
	"github.com/xiaolan580230/clhttp-framework/core/cljson"
	"github.com/xiaolan580230/clhttp-framework/core/rule"
	"github.com/xiaolan580230/clhttp-framework/core/skylog"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)


var TempDirPath string

//@author xiaolan
//@lastUpdate 2019-08-09
//@comment 启动HTTP服务
func StartServer(_listenPort uint32) {

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/upload_test", uploadFileHtml)
	http.HandleFunc("/uploadFile", uploadFile)

	tempDirPath, _ := ioutil.TempDir("__clhttp_tempfile__", "")
	TempDirPath = tempDirPath
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", _listenPort), nil))
}

//@author xiaolan
//@lastUpdate 2019-08-09
//@comment http请求主入口回调
func rootHandler(rw http.ResponseWriter, rq *http.Request) {

	// 跨域支持
	rw.Header().Set("Access-Control-Allow-Origin", "*")             //允许访问所有域
	rw.Header().Add("Access-Control-Allow-Headers", "Content-Type") //header的类型

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

	// 过滤请求
	for _, filterExt := range filter_file_ext {
		if strings.HasSuffix(requestArr[1], filterExt) {
			rw.WriteHeader(404)
			return
		}
	}
	var contentType = rq.Header.Get("Content-Type")

	var values = make(map[string]string)
	var rawData = ""
	if contentType == "text/json" || contentType == "application/json" {
		var jsonBytes = make([]byte, 4096)
		n, err := rq.Body.Read(jsonBytes)
		if err != nil && err.Error() != "EOF"{
			skylog.LogErr("读取json参数失败! 错误:%v", err)
			rw.WriteHeader(502)
			return
		}
		jsonObj := cljson.New(jsonBytes[:n])
		if jsonObj != nil {
			values = jsonObj.ToMap().ToCustom()
		}
		rawData = string(jsonBytes[:n])
	} else {
		rq.ParseMultipartForm(2 << 32)
		if nil != rq.MultipartForm && nil != rq.MultipartForm.Value {
			for key, val := range rq.MultipartForm.Value {
				if len(val) == 1 {
					values[key] = val[0]
				}
			}
		} else {
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

	myUA := rq.Header.Get("platform")
	if myUA == "" {
		myUA = "web"
	}

	myLang := rq.Header.Get("Sky-Server-Lang")
	if myLang == "" {
		myLang = "zhcn_simple"
	}
	if strings.HasPrefix(remoteip, "[::1]") {
		remoteip = strings.ReplaceAll(remoteip, "[::1]", "127.0.0.1")
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
		Port:       "",
		Language:   myLang,
		RawData:    rawData,
	}
	content, contentType := CallHandler(requestName, rqObj, &serObj)
	if contentType == "" {
		contentType = "text/json"
	}

	if contentType == "error" {
		rw.WriteHeader(502)
		return
	}

	rw.Header().Set("Content-Type", contentType)
	rw.Header().Set("Charset", "UTF-8")
	rw.Write([]byte(content))
}



// 测试上传页面
func uploadFileHtml (rw http.ResponseWriter, rq *http.Request) {

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


// 上传文件
func uploadFile (rw http.ResponseWriter, rq *http.Request) {
	if rq.Method != "POST" {
		// 不是使用post的一律拒绝
		rw.WriteHeader(502)
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
	rq.ParseMultipartForm(10 << 20)
	clientfd, handler, err := rq.FormFile("uploadfile")
	if err != nil {
		skylog.LogErr("加载文件失败: %v", err)
		rw.WriteHeader(502)
		return
	}
	defer clientfd.Close()

	fileNameArr := strings.Split(handler.Filename, ".")
	fileExt := fileNameArr[len(fileNameArr)-1]

	//// 需要过滤的请求文件类型列表
	//filter_file_ext := []string {
	//	"png", "jpg", "gif", "jpeg",
	//}
	//
	//// 过滤请求
	//isPass := false
	//for _, filter_ext := range filter_file_ext {
	//	if filter_ext == strings.ToLower(fileExt) {
	//		isPass = true
	//	}
	//}
	//
	//if !isPass {
	//	rw.WriteHeader(501)
	//	return
	//}

	request_url := ""
	proctol := rq.Header.Get("Proxy-X-Forwarded-Proto")
	if proctol == "" {
		proctol = rq.Header.Get("X-Forwarded-Proto")
	}
	port := rq.Header.Get("X-Forwarded-Port")
	if proctol == "https" {
		request_url = "https://"+rq.Host
		if port != "443" {
			request_url += ":"+rq.Header.Get("X-Forwarded-Port")
		}
		request_url += rq.RequestURI
		proctol = "https"
	} else {
		request_url = "http://"+rq.Host
		if port != "80" {
			request_url += ":"+rq.Header.Get("X-Forwarded-Port")
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


	buffers := make([]byte, 10 << 20)
	lenOfBuffer, err := clientfd.Read(buffers)
	if err != nil {
		skylog.LogErr("读取文件内容失败: %v", err)
		rw.WriteHeader(502)
		return
	}
	fileName := fmt.Sprintf("%v.%v", clCommon.Md5(buffers[:lenOfBuffer]), fileExt)

	localfd, err := os.OpenFile(os.TempDir() + fileName, os.O_CREATE | os.O_TRUNC | os.O_WRONLY, 0666)
	if err != nil {
		skylog.LogErr("打开文件失败! %v", err)
		rw.WriteHeader(502)
		return
	}
	localfd.Write(buffers[:lenOfBuffer])
	localfd.Close()


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

	var rqObj = rule.NewHttpParam(map[string] string{
		"ac": "UploadFile",			// 执行的动作
		"filename": fileName,		// 文件名
		"fileExt": fileExt,			// 文件扩展名
		"localPath": os.TempDir() + fileName,  // 本地路径
	})

	content, contentType := CallHandler("upload", rqObj, &serObj)
	if contentType == "" {
		contentType = "text/json"
	}

	rw.Header().Set("Content-Type", contentType)
	rw.Write([]byte(content))
	return
}


//@author xiaolan
//@lastUpdate 2019-08-11
//@comment UA转整数型
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

//@author xiaolan
//@lastUpdate 2019-08-10
//@comment 调用回调函数, 成功返回文档类型(默认为text/json)和回传结果,失败返回空字符串
//@param _name 回调函数名称
//@param _param 参数列表
//@param _serconf 服务器信息
func CallHandler(_name string, _param *rule.HttpParam, _serconf *rule.ServerParam) (string, string) {

	return Request(_name, _param, _serconf)
}



//@author xiaolan
//@lastUpdate 2019-08-10
//@comment 请求信息
func Request(_name string, _param *rule.HttpParam, _server *rule.ServerParam) (string, string) {

	return rule.CallRule(_name, _param, _server), "text/json"
}
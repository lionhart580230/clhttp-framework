package httpserver

import (
	"clhttp-framework/core/rule"
	"fmt"
	"log"
	"net/http"
	"strings"
)


//@author xiaolan
//@lastUpdate 2019-08-09
//@comment 启动HTTP服务
func StartServer(_listenPort uint32) {

	http.HandleFunc("/", rootHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", _listenPort), nil))
}

//@author xiaolan
//@lastUpdate 2019-08-09
//@comment http请求主入口回调
func rootHandler(rw http.ResponseWriter, rq *http.Request) {
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

	request_name := requestArr[1]

	// 过滤请求
	for _, filter_ext := range filter_file_ext {
		if strings.HasSuffix(requestArr[1], filter_ext) {
			rw.WriteHeader(404)
			return
		}
	}

	rq.ParseMultipartForm(2 << 32)

	var values = make(map[string]string)

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

	myUA := rq.Header.Get("platform")
	if myUA == "" {
		myUA = "web"
	}

	myLang := rq.Header.Get("Sky-Server-Lang")
	if myLang == "" {
		myLang = "zhcn_simple"
	}
	var serObj = rule.ServerParam{
		RemoteIP:   strings.Split(remoteip, ":")[0],
		RequestURI: rq.RequestURI,
		Host:       rq.Host,
		Method:     rq.Method,
		RequestURL: request_url,
		UA:         myUA,
		UAType:     UAToInt(myUA),
		Proctol:    proctol,
		Port:       port,
		Language:   myLang,
	}
	content, contentType := CallHandler(request_name, rqObj, &serObj)
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
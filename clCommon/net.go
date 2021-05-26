package clCommon

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"
)



//@author ciaolan
//@lastUpdate 2019-08-31
//@comment 模拟HTTP POST请求
func HttpPost(url string, param string) string {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				conn, err := net.DialTimeout(netw, addr, time.Second*10)
				if err != nil {
					return nil, err
				}
				conn.SetDeadline(time.Now().Add(time.Second * 10))
				return conn, nil
			},
			ResponseHeaderTimeout: time.Second * 10,
		},
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(param))
	if err != nil {
		return fmt.Sprintf(`{"msg":400,"param":"%v"}`, err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Mobile Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf(`{"msg":400,"param":"%v"}`, err)
	}
	jsonStr, err2 := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err2 != nil {
		return fmt.Sprintf(`{"msg":400,"param":"%v"}`, err)
	}

	return string(jsonStr)
}
package modules

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

// 定义header结构体，请求头与数据
type header struct {
	key   string
	value string
}

func Init(URL string, proxy string, AddHeader []string, bypassip []string) {
	//判断字符串最后是否有反斜杠
	if URL[len(URL)-1:] != "/" {
		URL = URL + "/"
	}
	//如果url没有http或者https,则添加http
	if !strings.HasPrefix(URL, "http://") && !strings.HasPrefix(URL, "https://") {
		URL = "http://" + URL
	}

	//判断是否有代理
	if proxy != "" {
		fmt.Println("\033[34m[+] Using proxy: ", proxy, "\033[0m")
	}
	var headers []header
	//判断是否有添加的请求头
	if len(AddHeader) != 0 {
		fmt.Println("\033[34m[+] Using headers: ", AddHeader, "\033[0m")
		//单独输出每一个请求头
		for _, v := range AddHeader {
			//将每一个请求头根据:分割，并存储结构体
			split := strings.Split(v, ":")
			headers = append(headers, header{split[0], split[1]})
		}
	}
	//判断是否有添加的ip
	if len(bypassip) != 0 {
		fmt.Println("\033[34m[+] Using bypass ip: ", bypassip, "\033[0m")
	}
	TestMethods(URL, proxy, headers)
	Testheaders("GET", URL, proxy, headers, bypassip)
	Testheaders("POST", URL, proxy, headers, bypassip)
	TestendPath("GET", URL, proxy, headers)
	TestendPath("POST", URL, proxy, headers)
	TestmidPath("GET", URL, proxy, headers)
	TestmidPath("POST", URL, proxy, headers)

	//fmt.Println(Request("GET", URL, proxy, headers))

}

// 封装请求
func Request(method string, URL string, proxy string, headers []header) (statusCode int, response []byte, err error) {
	//判断方法是否为空
	if method == "" {
		method = "GET"
	}
	//处理代理,http或socks5
	PROXY := func(_ *http.Request) (*url.URL, error) {
		if proxy == "" {
			return nil, nil
		}
		return url.Parse(proxy)
	}
	client := &http.Client{Transport: &http.Transport{
		Proxy:           PROXY,
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		DialContext:     (&net.Dialer{Timeout: 1 * time.Second}).DialContext,
	}}
	req, err := http.NewRequest(method, URL, nil)
	if err != nil {
		return 0, nil, err
	}
	//处理添加的请求头
	for _, header := range headers {
		req.Header.Add(header.key, header.value)
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	}
	res, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	resp, err := httputil.DumpResponse(res, true)
	if err != nil {
		return 0, nil, err
	}
	return res.StatusCode, resp, nil
}

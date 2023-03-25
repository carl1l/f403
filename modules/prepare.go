package modules

import (
	"crypto/tls"
	"f403/util"
	"fmt"
	"github.com/go-resty/resty/v2"
	"strings"
)

type header struct {
	key   string
	value string
}

func Init(URL string, proxy string, AddHeader []string, bypassip []string, postdata string) {
	if !strings.HasPrefix(URL, "http://") && !strings.HasPrefix(URL, "https://") {
		URL = "http://" + URL
	}

	if proxy != "" {
		fmt.Println("\033[34m[+] Using proxy: ", proxy, "\033[0m")
	}
	var headers []header
	//判断是否有添加的请求头
	if len(AddHeader) != 0 {
		headerstr := util.ArrayToString(AddHeader, ",")
		util.Blue("[+] Using headers: " + headerstr)

		for _, v := range AddHeader {
			split := strings.SplitN(v, ":", 2)
			headers = append(headers, header{strings.TrimSpace(split[0]), strings.TrimSpace(split[1])})
		}
	}
	if len(bypassip) != 0 {
		bypassipstr := util.ArrayToString(bypassip, ",")
		util.Blue("[+] Using bypassip: " + bypassipstr)
	}

	TestMethods(URL, proxy, headers, postdata)
	Testheaders("GET", URL, proxy, headers, bypassip, "")
	Testheaders("POST", URL, proxy, headers, bypassip, postdata)
	TestendPath("GET", URL, proxy, headers, "")
	TestendPath("POST", URL, proxy, headers, postdata)
	TestmidPath("GET", URL, proxy, headers, "")
	TestmidPath("POST", URL, proxy, headers, postdata)
	if URL[len(URL)-1:] == "/" {
		URL = URL[:len(URL)-1]
	}
	TestpathCase("GET", URL, proxy, headers, "")
	TestpathCase("POST", URL, proxy, headers, postdata)
	URL = URL + "/"
	TestpathCase("GET", URL, proxy, headers, "")
	TestpathCase("POST", URL, proxy, headers, postdata)

}

func Request(method string, url string, proxy string, headers []header, postdata string) (int, string, error) {
	// create a resty client
	var resp *resty.Response = nil
	var err error = nil
	if method == "" {
		method = "GET"
	}
	client := resty.New()
	if proxy != "" {
		client.SetProxy(proxy)
	}
	client.SetHeader("User-Agent", " Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.71 Safari/537.36")
	if headers != nil {
		for _, header := range headers {
			client.SetHeader(header.key, header.value)
		}
	}
	// disable TLS verification
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	if method == "POST" {
		if util.IsJSON(postdata) {
			client.SetHeader("Content-Type", "application/json")
		} else {
			client.SetHeader("Content-Type", "application/x-www-form-urlencoded")
		}

		resp, err = client.R().SetBody(postdata).Execute(method, url)
		if err != nil {
			fmt.Println(err)
			return 0, "", err
		}
	} else {
		resp, err = client.R().Execute(method, url)
		if err != nil {
			fmt.Println(err)
			return 0, "", err
		}
	}
	return resp.StatusCode(), resp.String(), nil
}

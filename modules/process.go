package modules

import (
	"f403/util"
	"fmt"
	"github.com/spf13/viper"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

func TestMethods(URL string, proxy string, headers []header, postdata string) {
	fmt.Println()
	util.Blue("[*]Testing methods: \n")
	methods := viper.Get("http.allow_methods")
	var wg sync.WaitGroup
	wg.Add(len(methods.([]interface{})))
	for _, method := range methods.([]interface{}) {
		go func(method string) {
			defer wg.Done()
			statusCode, respone, err := Request(method, URL, proxy, headers, postdata)
			if err != nil {
				util.Red("[-] " + method + string(len(respone)) + err.Error())
			} else {
				if statusCode/100 == 2 {
					util.Green("[+] " + strconv.Itoa(statusCode) + " " + strconv.Itoa(len(respone)) + " " + method)
				} else {
					util.Red("[-] " + strconv.Itoa(statusCode) + " " + strconv.Itoa(len(respone)) + " " + method)
				}
			}
		}(method.(string))
	}
	wg.Wait()
}

func Testheaders(method string, URL string, proxy string, addheaders []header, bypassip []string, postdata string) {
	fmt.Println()
	util.Blue("[*] Testing headers by " + method + "-method")
	fmt.Println()
	TestHeaders := viper.Get("http.headers")
	if len(bypassip) == 0 {
		bypassip = []string{"127.0.0.1", "localhost"}
	}
	var wg sync.WaitGroup
	wg.Add(len(TestHeaders.([]interface{})) * len(bypassip))
	for _, ip := range bypassip {
		for _, every := range TestHeaders.([]interface{}) {
			go func(every string, ip string) {
				defer wg.Done()
				headers := append(addheaders, header{every, ip})
				statusCode, respone, err := Request(method, URL, proxy, headers, postdata)
				if err != nil {
					util.Red("[-] " + strconv.Itoa(len(respone)) + " " + ip + " " + every + " " + err.Error())
				} else {
					if statusCode/100 == 2 {
						util.Green("[+] " + strconv.Itoa(statusCode) + " " + strconv.Itoa(len(respone)) + " " + every + ": " + ip)
					} else {
						util.Red("[-] " + strconv.Itoa(statusCode) + " " + strconv.Itoa(len(respone)) + " " + every + ": " + ip)
					}
				}
			}(every.(string), ip)
		}
	}
	wg.Wait()
}

func TestendPath(method string, URL string, proxy string, addheaders []header, postdata string) {
	fmt.Println()
	util.Blue("\033[32m[*] Testing endpath by " + method + "-method")
	fmt.Println()
	Testendpath := viper.Get("http.end_path")
	var wg sync.WaitGroup
	wg.Add(len(Testendpath.([]interface{})))
	for _, every := range Testendpath.([]interface{}) {
		go func(every string) {
			defer wg.Done()
			statusCode, respone, err := Request(method, URL+every, proxy, addheaders, postdata)
			if err != nil {
				util.Red("[-] " + strconv.Itoa(len(respone)) + " " + URL + every + " " + err.Error())
			} else {
				if statusCode/100 == 2 {
					util.Green("[+] " + strconv.Itoa(statusCode) + " " + strconv.Itoa(len(respone)) + " " + URL + every)
				} else {
					util.Red("[-] " + strconv.Itoa(statusCode) + " " + strconv.Itoa(len(respone)) + " " + URL + every)
				}
			}
		}(every.(string))
	}
	wg.Wait()
}

func TestmidPath(method string, URL string, proxy string, addheaders []header, postdata string) {
	fmt.Println()
	util.Blue("[*] Testing midpath by " + method + "-method")
	fmt.Println()
	Testmidpath := viper.Get("http.mid_path")
	var wg sync.WaitGroup
	wg.Add(len(Testmidpath.([]interface{})))
	for _, every := range Testmidpath.([]interface{}) {
		go func(every string) {
			defer wg.Done()
			u, _ := url.Parse(URL)
			urlend := every + u.Path
			Url := u.Scheme + "://" + u.Host + "/" + urlend
			statusCode, respone, err := Request(method, Url, proxy, addheaders, postdata)
			if err != nil {
				util.Red("[-] " + strconv.Itoa(len(respone)) + " " + Url + " " + err.Error())
			} else {
				if statusCode/100 == 2 {
					util.Red("[+] " + strconv.Itoa(statusCode) + " " + strconv.Itoa(len(respone)) + " " + Url)
				} else {
					util.Red("[-] " + strconv.Itoa(statusCode) + " " + strconv.Itoa(len(respone)) + " " + Url)
				}
			}
			if u.Path != "" {
				Url1 := Url + "/"
				statusCode, respone, err := Request(method, Url1, proxy, addheaders, postdata)
				if err != nil {
					util.Red("[-] " + strconv.Itoa(len(respone)) + " " + Url1 + " " + err.Error())
				} else {
					if statusCode/100 == 2 {
						util.Red("[+] " + strconv.Itoa(statusCode) + " " + strconv.Itoa(len(respone)) + " " + Url1)
					} else {
						util.Red("[-] " + strconv.Itoa(statusCode) + " " + strconv.Itoa(len(respone)) + " " + Url1)
					}
				}
			}
		}(every.(string))
	}
	wg.Wait()
}

// 将url后的path路径进行大小写转换
func TestpathCase(method string, URL string, proxy string, addheaders []header, postdata string) {
	fmt.Println()
	util.Blue("[*] Test case conversion in the path by " + method + "-method")
	fmt.Println()
	var wg sync.WaitGroup
	u, _ := url.Parse(URL)
	wg.Add(len(u.Path))
	for i := 0; i < len(u.Path); i++ {
		go func(i int) {
			defer wg.Done()
			if u.Path[i] >= 'a' && u.Path[i] <= 'z' {
				u.Path = u.Path[:i] + strings.ToUpper(string(u.Path[i])) + u.Path[i+1:]
				Url := u.Scheme + "://" + u.Host + u.Path
				statusCode, respone, err := Request(method, Url, proxy, addheaders, postdata)
				if err != nil {
					util.Red("[-] " + strconv.Itoa(len(respone)) + " " + Url + " " + err.Error())
				} else {
					if statusCode/100 == 2 {
						util.Green("[+] " + strconv.Itoa(statusCode) + " " + strconv.Itoa(len(respone)) + " " + Url)
					} else {
						util.Red("[-] " + strconv.Itoa(statusCode) + " " + strconv.Itoa(len(respone)) + " " + Url)
					}
				}
				u.Path = u.Path[:i] + strings.ToLower(string(u.Path[i])) + u.Path[i+1:]
			}
			if u.Path[i] >= 'A' && u.Path[i] <= 'Z' {
				u.Path = u.Path[:i] + strings.ToLower(string(u.Path[i])) + u.Path[i+1:]
				Url := u.Scheme + "://" + u.Host + u.Path
				statusCode, respone, err := Request(method, Url, proxy, addheaders, postdata)
				if err != nil {
					util.Red("[-] " + strconv.Itoa(len(respone)) + " " + Url + " " + err.Error())
				} else {
					if statusCode/100 == 2 {
						util.Green("[+] " + strconv.Itoa(statusCode) + " " + strconv.Itoa(len(respone)) + " " + Url)
					} else {
						util.Red("[-] " + strconv.Itoa(statusCode) + " " + strconv.Itoa(len(respone)) + " " + Url)
					}
				}
				u.Path = u.Path[:i] + strings.ToUpper(string(u.Path[i])) + u.Path[i+1:]
			}
		}(i)
	}
	wg.Wait()
}

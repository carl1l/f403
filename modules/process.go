package modules

import (
	"fmt"
	"github.com/spf13/viper"
	"net/url"
	"sync"
)

func TestMethods(URL string, proxy string, headers []header) {
	fmt.Println()
	fmt.Println("\033[32m[*] Testing methods: \033[0m")
	fmt.Println()
	methods := viper.Get("http.allow_methods")
	var wg sync.WaitGroup
	wg.Add(len(methods.([]interface{})))
	for _, method := range methods.([]interface{}) {
		go func(method string) {
			defer wg.Done()
			statusCode, respone, err := Request("GET", URL, proxy, headers)
			if err != nil {
				fmt.Println("\033[31m[-] ", method, len(respone), err, "\033[0m")
			} else {
				if statusCode/100 == 2 {
					fmt.Println("\033[32m[+] ", statusCode, len(respone), method, "\033[0m")
				} else {
					fmt.Println("\033[31m[-] ", statusCode, len(respone), method, "\033[0m")
				}
			}
		}(method.(string))
	}
	wg.Wait()
}

func Testheaders(method string, URL string, proxy string, addheaders []header, bypassip []string) {
	fmt.Println()
	fmt.Println("\033[32m[*] Testing headers by "+method+"-method", "\033[0m")
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
				statusCode, respone, err := Request(method, URL, proxy, headers)
				if err != nil {
					fmt.Println("\033[31m[-] ", len(respone), ip, every, err, "\033[0m")
				} else {
					if statusCode/100 == 2 {
						fmt.Println("\033[32m[+] ", statusCode, len(respone), every+": "+ip, "\033[0m")
					} else {
						fmt.Println("\033[31m[-] ", statusCode, len(respone), every+": "+ip, "\033[0m")
					}
				}
			}(every.(string), ip)
		}
	}
	wg.Wait()
}

func TestendPath(method string, URL string, proxy string, addheaders []header) {
	fmt.Println()
	fmt.Println("\033[32m[*] Testing endpath by "+method+"-method", "\033[0m")
	fmt.Println()
	Testendpath := viper.Get("http.end_path")
	var wg sync.WaitGroup
	wg.Add(len(Testendpath.([]interface{})))
	for _, every := range Testendpath.([]interface{}) {
		go func(every string) {
			defer wg.Done()
			statusCode, respone, err := Request(method, URL+every, proxy, addheaders)
			if err != nil {
				fmt.Println("\033[31m[-] ", len(respone), every, err, "\033[0m")
			} else {
				if statusCode/100 == 2 {
					fmt.Println("\033[32m[+] ", statusCode, len(respone), URL+every, "\033[0m")
				} else {
					fmt.Println("\033[31m[-] ", statusCode, len(respone), URL+every, "\033[0m")
				}
			}
		}(every.(string))
	}
	wg.Wait()
}

func TestmidPath(method string, URL string, proxy string, addheaders []header) {
	fmt.Println()
	fmt.Println()
	fmt.Println("\033[32m[*] Testing midpath by "+method+"-method", "\033[0m")
	fmt.Println()
	Testmidpath := viper.Get("http.mid_path")
	var wg sync.WaitGroup
	wg.Add(len(Testmidpath.([]interface{})))
	for _, every := range Testmidpath.([]interface{}) {
		go func(every string) {
			defer wg.Done()
			u, _ := url.Parse(URL)
			u.Path = every + u.Path
			Url := u.Scheme + "://" + u.Host + "/" + u.Path
			statusCode, respone, err := Request(method, Url, proxy, addheaders)
			if err != nil {
				fmt.Println("\033[31m[-] ", len(respone), Url, err, "\033[0m")
			} else {
				if statusCode/100 == 2 {
					fmt.Println("\033[32m[+] ", statusCode, len(respone), Url, "\033[0m")
				} else {
					fmt.Println("\033[31m[-] ", statusCode, len(respone), Url, "\033[0m")
				}
			}
		}(every.(string))
	}
	wg.Wait()

}

package util

import (
	"encoding/json"
	"fmt"
	"strings"
)

func Red(str string) {
	fmt.Println("\033[31m" + str + "\033[0m")
}

func Green(str string) {
	fmt.Println("\033[32m" + str + "\033[0m")
}

func Blue(str string) {
	fmt.Println("\033[34m" + str + "\033[0m")
}

func Yellow(str string) {
	fmt.Println("\033[33m" + str + "\033[0m")
}

func ArrayToString(a []string, sep string) string {
	return strings.Join(a, sep)
}

func IsJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

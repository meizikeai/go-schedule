package tool

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"regexp"
	"strconv"
)

func Contain(arr []string, element string) bool {
	for _, v := range arr {
		if v == element {
			return true
		}
	}
	return false
}

func MarshalJson(date interface{}) []byte {
	res, err := json.Marshal(date)

	if err != nil {
		fmt.Println(err)
	}

	return res
}

func UnmarshalJson(date string) map[string]interface{} {
	var res map[string]interface{}

	_ = json.Unmarshal([]byte(date), &res)

	return res
}

func GetRandmod(length int) int64 {
	result := int64(0)
	res, err := rand.Int(rand.Reader, big.NewInt(int64(length)))

	if err != nil {
		return result
	}

	return res.Int64()
}

func IntToString(value int64) string {
	v := strconv.FormatInt(value, 10)

	return v
}

func StringToInt(value string) int64 {
	res, err := strconv.ParseInt(value, 10, 64)

	if err != nil {
		res = 0
	}

	return res
}

func CheckPassword(password string, min, max int) int {
	level := 0

	if len(password) < min {
		return -1
	}

	if len(password) > max {
		return 5
	}

	patternList := []string{`[0-9]+`, `[a-z]+`, `[A-Z]+`, `[~!@#$%^&amp;*?_-]+`}

	for _, pattern := range patternList {
		match, _ := regexp.MatchString(pattern, password)

		if match == true {
			level++
		}
	}

	return level
}

var compileRegex = regexp.MustCompile(`\D`)

func ClearNotaNumber(str string) string {
	return compileRegex.ReplaceAllString(str, "")
}

func HandleEscape(source string) string {
	var j int = 0

	if len(source) == 0 {
		return ""
	}

	tempStr := source[:]
	desc := make([]byte, len(tempStr)*2)

	for i := 0; i < len(tempStr); i++ {
		flag := false
		var escape byte

		switch tempStr[i] {
		case '\r':
			flag = true
			escape = '\r'
			break
		case '\n':
			flag = true
			escape = '\n'
			break
		case '\\':
			flag = true
			escape = '\\'
			break
		case '\'':
			flag = true
			escape = '\''
			break
		case '"':
			flag = true
			escape = '"'
			break
		case '\032':
			flag = true
			escape = 'Z'
			break
		default:
		}

		if flag {
			desc[j] = '\\'
			desc[j+1] = escape
			j = j + 2
		} else {
			desc[j] = tempStr[i]
			j = j + 1
		}
	}

	return string(desc[0:j])
}

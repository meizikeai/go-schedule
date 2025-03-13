package tool

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Units struct{}

func NewUnits() *Units {
	return &Units{}
}

func (u *Units) Contain(arr []string, element string) bool {
	for _, v := range arr {
		if v == element {
			return true
		}
	}
	return false
}

func (u *Units) MarshalJson(date any) []byte {
	res, err := json.Marshal(date)

	if err != nil {
		fmt.Println(err)
	}

	return res
}

func (u *Units) UnmarshalJson(date string) map[string]any {
	var res map[string]any

	_ = json.Unmarshal([]byte(date), &res)

	return res
}

func (u *Units) IntToString(value int64) string {
	v := strconv.FormatInt(value, 10)

	return v
}

func (u *Units) StringToInt(value string) int64 {
	res, err := strconv.ParseInt(value, 10, 64)

	if err != nil {
		res = 0
	}

	return res
}

func (u *Units) CheckPassword(password string, min, max int) int {
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

func (u *Units) HandleEscape(source string) string {
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

func (u *Units) GenerateRandomNumber(start, end, count int) ([]int, error) {
	var result []int

	for i := 0; i < count; i++ {
		rangeBig := big.NewInt(int64(end - start + 1))
		n, err := rand.Int(rand.Reader, rangeBig)

		if err != nil {
			return nil, err
		}

		num := int(n.Int64()) + start

		result = append(result, num)
	}

	return result, nil
}
func (u *Units) RemoveDuplicateElement(strs []string) []string {
	result := make([]string, 0, len(strs))
	temp := map[string]struct{}{}

	for _, item := range strs {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}

func (u *Units) IsSlice(v any) bool {
	kind := reflect.ValueOf(v).Kind()

	if kind == reflect.Slice || kind == reflect.Array {
		return true
	}

	return false
}

func (u *Units) StringToArray(data string) []string {
	result := []string{}

	if len(data) > 0 {
		result = strings.Split(data, ",")
	}

	return result
}

func (u *Units) FloatToString(data float64) string {
	result := strconv.FormatFloat(data, 'f', 2, 64)
	return result
}

package tool

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"regexp"
	"strings"
)

type Units struct{}

func NewUnits() *Units {
	return &Units{}
}

// 字符串转换
// https://github.com/spf13/cast

func (u *Units) MarshalJson(date any) string {
	res, err := json.Marshal(date)

	if err != nil {
		return ""
	}

	return string(res)
}

func (u *Units) UnmarshalJson(date string) map[string]any {
	res := make(map[string]any, 0)

	_ = json.Unmarshal([]byte(date), &res)

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
	if len(source) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.Grow(len(source) * 2)

	for _, c := range source {
		switch c {
		case '\r', '\n', '\\', '\'', '"', '\032', '\x00', '\b', '\t':
			builder.WriteByte('\\')
			if c == '\032' {
				builder.WriteByte('Z') // 处理 MySQL 的 Ctrl+Z
			} else {
				builder.WriteByte(byte(c))
			}
		default:
			builder.WriteByte(byte(c))
		}
	}
	return builder.String()
}

func (u *Units) GenerateRandomNumber(start, end, count int) ([]int, error) {
	var result []int

	for range count {
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

func (u *Units) ArrayIntToString(array []int, delim string) string {
	return strings.Trim(strings.Replace(fmt.Sprint(array), " ", delim, -1), "[]")
}

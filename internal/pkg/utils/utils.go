// internal/pkg/utils/utils.go
package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"
	"unicode"
)

const (
	letters      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	numbers      = "0123456789"
	alphanumeric = letters + numbers
)

func RandomString(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	for i := range b {
		b[i] = alphanumeric[b[i]%byte(len(alphanumeric))]
	}
	return string(b)
}

func RandomNumber(min, max int) (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	if err != nil {
		return 0, err
	}
	return int(n.Int64()) + min, nil
}

func RemoveDuplicates[T comparable](data []T) []T {
	seen := make(map[T]struct{}, len(data))
	result := make([]T, 0, len(data))
	for _, v := range data {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			result = append(result, v)
		}
	}
	return result
}

func ArrayToString[T any](array []T, delim string) string {
	if len(array) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.Grow(len(array) * 8)
	for i, v := range array {
		if i > 0 {
			sb.WriteString(delim)
		}
		fmt.Fprint(&sb, v)
	}
	return sb.String()
}

func Marshal(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(b)
}

func Unmarshal(str string, v any) error {
	return json.Unmarshal([]byte(str), v)
}

type PasswordRule struct {
	MinLength     int
	MaxLength     int
	RequireUpper  bool
	RequireLower  bool
	RequireNumber bool
	RequireSymbol bool
}

func CheckPasswordStrength(pwd string, rule PasswordRule) (int, error) {
	if len(pwd) < rule.MinLength {
		return 0, fmt.Errorf("password too short")
	}
	if len(pwd) > rule.MaxLength {
		return 0, fmt.Errorf("password too long")
	}

	score := 1
	if rule.RequireUpper && hasUpper(pwd) {
		score++
	}
	if rule.RequireLower && hasLower(pwd) {
		score++
	}
	if rule.RequireNumber && hasNumber(pwd) {
		score++
	}
	if rule.RequireSymbol && hasSymbol(pwd) {
		score++
	}
	return score, nil
}

func hasUpper(s string) bool  { return strings.IndexFunc(s, unicode.IsUpper) >= 0 }
func hasLower(s string) bool  { return strings.IndexFunc(s, unicode.IsLower) >= 0 }
func hasNumber(s string) bool { return strings.IndexFunc(s, unicode.IsDigit) >= 0 }
func hasSymbol(s string) bool {
	return strings.IndexFunc(s, func(r rune) bool { return unicode.IsPunct(r) || unicode.IsSymbol(r) }) >= 0
}

func NowStr() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

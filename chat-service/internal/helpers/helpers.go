package helpers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
)

func PrintJSON(v interface{}) {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling to JSON:", err)
		return
	}
	fmt.Println(string(b))
}

func RandomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, n)
	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return ""
		}
		result[i] = letters[num.Int64()]
	}
	return string(result)
}

// intToStr converts an int64 to a string
func IntToStr(n int64) string {
	return strconv.FormatInt(n, 10)
}

// strToInt converts a string to int64, ignoring errors
func StrToInt(s string) int64 {
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}

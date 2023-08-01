package util

import (
	"bytes"
	"math/big"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
)

func StandardizationText(value string) string {
	var buffer bytes.Buffer
	value = strings.ToLower(value)
	for _, char := range value {
		if strings.Contains(charset, string(char)) {
			buffer.WriteByte(byte(char))
		}
	}
	text := buffer.String()
	padCount := maxPageContentLength - len(text)
	padString := strings.Repeat(" ", padCount)
	return text + padString
}

func AppendRandomString(value string) string {
	var result string
	valueLen := maxPageContentLength - len(value)
	index := rand.Intn((valueLen - 1) + 1)
	// rand 1 -> valueLen: 3194
	for i := 1; i < index; i++ {
		result += string(charset[rand.Intn((39-0)+1)])
	}
	result += value
	for i := len(result); i < maxPageContentLength; i++ {
		result += string(charset[rand.Intn((39-0)+1)])
	}
	return result
}

func IsAlpha(c byte) bool {
	return regexp.MustCompile(`^[a-z]`).MatchString(string(c))
}

func ConvertToIntBase(s string, base int) (hexInt *big.Int) {
	hexInt, _ = new(big.Int).SetString(s, base)
	return hexInt
}
func CalculateLibraryCoordinate() int {
	libraryCoordinate, _ := strconv.Atoi(page + volume + shelf + wall)
	return libraryCoordinate
}

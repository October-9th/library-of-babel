package util

import (
	"math/big"
	"math/rand"
	"strconv"
	"strings"
)

const (
	charset              = "abcdefghijklmnopqrstuvwxyz, ./0123456789"
	alphabet             = "abcdefghijklmnopqrstuvwxyz"
	digits               = "0123456789"
	maxPageContentLength = 3200
	maxWalls             = 4
	maxShelves           = 5
	maxVolumes           = 32
	maxPages             = 410
	hexagonBase          = 36
)

var (
	charsetLen = big.NewInt(int64(len(charset)))
	wall       = strconv.Itoa(rand.Intn(maxWalls-1) + 1)
	shelf      = strconv.Itoa(rand.Intn(maxShelves-1) + 1)
	volume     = strconv.Itoa(rand.Intn(maxVolumes-1) + 1)
	page       = strconv.Itoa(rand.Intn(maxPages-1) + 1)
)

// Algorithm to encode text into hexagon address for the library of babel
func GetAddress(text string, libraryCoordinate int64) string {
	standardText := StandardizationText(text)

	sum := big.NewInt(0)
	for i := 0; i < len(standardText); i++ {
		c := standardText[i]
		charVal := big.NewInt(0)
		if IsAlpha(c) {
			charVal.SetInt64(int64(c - 'a'))
		} else if c == '.' {
			charVal.SetInt64(28)
		} else if c == ' ' {
			charVal.SetInt64(27)
		} else {
			charVal.SetInt64(26)
		}
		switch c {
		case '/':
			charVal.SetInt64(29)
		case '0':
			charVal.SetInt64(30)
		case '1':
			charVal.SetInt64(31)
		case '2':
			charVal.SetInt64(32)
		case '3':
			charVal.SetInt64(33)
		case '4':
			charVal.SetInt64(34)
		case '5':
			charVal.SetInt64(35)
		case '6':
			charVal.SetInt64(36)
		case '7':
			charVal.SetInt64(37)
		case '8':
			charVal.SetInt64(38)
		case '9':
			charVal.SetInt64(39)
		}
		powerVal := big.NewInt(0)
		powerVal.Exp(charsetLen, big.NewInt(int64(i)), nil)
		// charsetLen ** i

		mulVal := big.NewInt(0)
		mulVal.Mul(charVal, powerVal)
		// charVal * powerVal

		sum.Add(sum, mulVal)
	}

	libCoord := big.NewInt(0)
	libCoord.SetInt64(libraryCoordinate)

	powerVal := big.NewInt(0)
	powerVal.Exp(charsetLen, big.NewInt(int64(maxPageContentLength)), nil)
	// charsetLen ** maxpageContentLength
	mulVal := big.NewInt(0)
	mulVal.Mul(libCoord, powerVal)
	// libCoord * powerVal

	result := big.NewInt(0)
	result.Add(mulVal, sum)

	finalRes := convertToBase(result, hexagonBase)
	return finalRes + ":w" + wall + ":s" + shelf + ":v" + volume + ":p" + page
}

// Algorithm to reverse the process, decode hexagon address into text
func GetContent(address string) string {
	parts := strings.Split(address, ":")
	hexagonAddress := parts[0]
	wall := parts[1][2:]
	shelf := parts[2][2:]
	volume := parts[3][2:]
	page := parts[4][2:]

	libraryCoordinate, _ := strconv.Atoi(page + volume + shelf + wall)

	hexInt := ConvertToIntBase(hexagonAddress, hexagonBase)
	powerVal := big.NewInt(0)
	powerVal.Exp(charsetLen, big.NewInt(maxPageContentLength), nil)
	powerVal.Mul(powerVal, big.NewInt(int64(libraryCoordinate)))

	seed := big.NewInt(0)
	seed.Sub(hexInt, powerVal)

	hexagonBaseResult := convertToBase(seed, hexagonBase)
	hexagonBaseResultInt := ConvertToIntBase(hexagonBaseResult, hexagonBase)

	result := convertToBase(hexagonBaseResultInt, int(charsetLen.Int64()))

	if len(result) < maxPageContentLength {
		resultInt, _ := strconv.Atoi(result)
		rand.NewSource(int64(resultInt))
		for len(result) < maxPageContentLength {
			index := rand.Intn(len(charset))
			result += string(charset[index])
		}
	} else if len(result) > maxPageContentLength {
		result = result[len(result)-maxPageContentLength:]
	}
	resultBytes := []byte(result)
	for i, j := 0, len(resultBytes)-1; i < j; i, j = i+1, j-1 {
		resultBytes[i], resultBytes[j] = resultBytes[j], resultBytes[i]
	}
	finalString := string(resultBytes)
	// finalString = strings.TrimRight(finalString, " ")
	return finalString

}

// Convert x to base
func convertToBase(n *big.Int, base int) string {
	var digs string
	var sign int64
	var chars []string
	if base == 36 {
		digs = digits + alphabet
	} else if base == 10 {
		digs = digits
	} else if base == 60 {
		digs = digits + "ABCDEFGHIJKLMNOPQRSTUVWXYZ" + alphabet
	} else {
		digs = charset
	}
	if n.Cmp(big.NewInt(0)) < 0 {
		sign = -1
	} else if n.Cmp(big.NewInt(0)) == 0 {
		return string(digs[0])
	} else {
		sign = 1
	}
	n.Mul(n, big.NewInt(sign))
	for n.Cmp(big.NewInt(0)) > 0 {
		index := big.NewInt(0)
		n, index = n.DivMod(n, big.NewInt(int64(base)), index)
		chars = append(chars, string(digs[index.Int64()]))
	}
	if sign < 0 {
		chars = append(chars, "-")
	}

	for i, j := 0, len(chars)-1; i < j; i, j = i+1, j-1 {
		chars[i], chars[j] = chars[j], chars[i]
	}

	return strings.Join(chars, "")
}

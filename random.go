package awsutils

import (
	"crypto/rand"
	"errors"
	"math/big"
)

const (
	exceedsTotalLength          = "numbers plus symbols must be less than total length"
	allowedNumbersMustBeDefined = "the list of allowed numbers must be specified"
	allowedSymbolsMustBeDefined = "the list of allowed symbols must be specified"
)

func GeneratePassword(length, numbersLength, symbolsLength int, letters, numbers, symbols *string) (string, error) {

	// default permitted letters.
	defaultLetters := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// default permitted digits.
	defaultNumbers := "0123456789"
	// default permitted symbols.
	defaultSymbols := "~!@#$%^&*()_+`-={}|[]\\:\"<>?,./"

	if letters == nil {
		letters = &defaultLetters
	}
	if numbers == nil {
		numbers = &defaultNumbers
	}
	if symbols == nil {
		symbols = &defaultSymbols
	}

	if numbersLength > 0 && len(*numbers) == 0 {
		return "", errors.New(allowedNumbersMustBeDefined)
	}

	if symbolsLength > 0 && len(*symbols) == 0 {
		return "", errors.New(allowedSymbolsMustBeDefined)
	}

	lettersLength := length - numbersLength - symbolsLength
	if lettersLength < 0 {
		return "", errors.New(exceedsTotalLength)
	}

	password := scramble("", *letters, lettersLength)
	password = scramble(password, *numbers, numbersLength)
	password = scramble(password, *symbols, symbolsLength)
	return password, nil
}
func scramble(password, allowed string, n int) string {
	if len(allowed) == 0 {
		return password
	}
	for i := 0; i < n; i++ {
		selection := randomSelect(allowed)
		password = randomInsert(password, selection)
	}
	return password
}
func randomInsert(in, val string) string {
	if in == "" {
		return val
	}
	n := len(in)
	i := randomInt(n + 1)
	return in[0:i] + val + in[i:n]
}
func randomSelect(in string) string {
	i := randomInt(len(in))
	return string(in[i])
}
func randomInt(max int) int64 {
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0
	}
	i := n.Int64()
	return i
}

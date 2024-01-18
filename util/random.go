package util

import (
	"math/rand"
	"strings"
)

const alphabet = "abcdefgijklmnopqrstuvwxyz"

// Generates a random int64 between max and min
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// Generates a random string of length n
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		char := alphabet[rand.Intn(k)]
		sb.WriteByte(char)
	}

	return sb.String()
}

// Generates a random owner name of length 6
func RandomOwner() string {
	return RandomString(6)
}

// Generates a random money between 0 and 1000
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

func RandomCurrency() string {

	currencies := make([]string, 0, len(SupportedCurrencies))
	for k := range SupportedCurrencies {
		currencies = append(currencies, k)
	}

	n := len(currencies)

	return currencies[rand.Intn(n)]
}

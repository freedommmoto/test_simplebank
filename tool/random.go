package tool

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "asdfhsdkfjhweurgbroqwhsdhdgtyphjfgv"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(6)
}

func RandomMoney() int64 {
	return RandomInt(1, 999)
}

func RandomCurrency() string {
	currency := []string{"THB", "EUR", "USD"}
	n := len(currency)
	return currency[rand.Intn(n)]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@%s.com", RandomString(4), RandomString(7))
}

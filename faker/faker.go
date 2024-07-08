package faker

import (
	"fmt"
	"math/rand"
	"strings"
)

//add password variable and put hashed password

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}

func RandInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}

func Username() string {
	return RandStringRunes(RandInt(2, 10))
}
func ID() string {
	return fmt.Sprintf("%s-%s-%s-%s", RandStringRunes(4), RandStringRunes(4), RandStringRunes(4), RandStringRunes(4))
}
func Email() string {
	return fmt.Sprintf("%s@example.com", strings.ToLower(RandStringRunes(RandInt(5, 10))))
}

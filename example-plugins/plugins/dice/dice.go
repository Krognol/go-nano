package dice

import (
	"math/rand"
	"strconv"
)

func D6() int {
	return rand.Intn(6)
}

func Roll(n string) int {
	i, _ := strconv.Atoi(n)
	return rand.Intn(i)
}

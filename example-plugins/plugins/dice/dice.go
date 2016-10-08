package dice

import (
	"math/rand"
	"strconv"
	"time"
)

func D6() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(6)
}

func Roll(n string) int {
	rand.Seed(time.Now().UnixNano())
	i, _ := strconv.Atoi(n)
	return rand.Intn(i)
}

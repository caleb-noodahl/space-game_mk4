package utils

import (
	"crypto/rand"
	"math/big"
)

func Roll(sides, amount int) []int {
	out := []int{}
	for i := 0; i < amount; i++ {
		r, _ := rand.Int(rand.Reader, big.NewInt(int64(sides)))
		out = append(out, int(r.Int64())+1)
	}
	return out
}

func Rand(lowest, max int) int {
	out, _ := rand.Int(rand.Reader, big.NewInt(int64(max)-int64(lowest)))
	out = big.NewInt(out.Int64() + int64(lowest))
	return int(out.Int64())
}

func RandArray(input []int) int {
	r, _ := rand.Int(rand.Reader, big.NewInt(int64(len(input))))
	index := r.Int64()
	return input[index]
}

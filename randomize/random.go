package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	fmt.Println("Hello, playground")

	//yy := []int{1, 2, 3, 9, 10}
	//fmt.Println(Shuffle(yy))

	xx := []int{1, 2, 3, 9, 10}
	fmt.Println(xx)
	xx = Reverse(xx)
	fmt.Println(xx)
	fmt.Println(Reverse(xx))

}

func Shuffle(vals []int) []int {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	ret := make([]int, len(vals))
	n := len(vals)
	for i := 0; i < n; i++ {
		randIndex := r.Intn(len(vals))
		ret[i] = vals[randIndex]
		vals = append(vals[:randIndex], vals[randIndex+1:]...)
	}
	return ret
}
func Reverse(vals []int) []int {
	ret := make([]int, len(vals))
	for i, j := 0, len(vals)-1; i < j; i, j = i+1, j-1 {
		ret[i], ret[j] = vals[j], vals[i]
	}
	return ret
}

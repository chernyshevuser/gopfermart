package lunh

import (
	"math/rand"
	"strconv"
)

func Validate(number string) bool {
	var sum int
	double := false
	for i := len(number) - 1; i >= 0; i-- {
		n, err := strconv.Atoi(string(number[i]))
		if err != nil {
			return false
		}
		if double {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		double = !double
	}
	return sum%10 == 0
}

func Generate(len int) string {
	number := make([]int, len-1)
	for i := range number {
		number[i] = rand.Intn(10)
	}
	sum := 0
	double := true
	for i := len - 2; i >= 0; i-- {
		n := number[i]
		if double {
			n *= 2
			if n > 9 {
				n -= 9
			}
		}
		sum += n
		double = !double
	}
	checkDigit := (10 - (sum % 10)) % 10
	number = append(number, checkDigit)
	result := ""
	for _, n := range number {
		result += strconv.Itoa(n)
	}
	return result
}

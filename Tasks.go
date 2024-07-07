package awesomeProject

import "fmt"
import "math"
import "strings"
import "golang.org/x/exp/slices"

func one() {
	fmt.Println("Hello World")
}

func two(a int, b int) int {
	return a + b
}

func three(a int) bool {
	if a%2 == 0 {
		return true
	} else {
		return false
	}
}

func four(a, b, c float64) float64 {
	var max_ab = math.Max(a, b)
	var max_abc = math.Max(max_ab, c)
	return max_abc
}

func five(n int) int {
	if n == 0 {
		return 1
	} else {
		return n * five(n-1)
	}
}

func six(s string) bool {
	runes := []rune(s)
	vowels := "уеёыаоэяиюУЕËЫАОЭЯИЮ"
	return strings.ContainsRune(vowels, runes[0])
}

func is_prime(n int) bool {
	if n <= 1 {
		return false
	}
	for i := 2; i*i <= n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}

func seven(limit int) int {
	var arr []int
	for i := 2; i <= limit; i++ {
		if is_prime(i) {
			arr = append(arr, i)
		}
	}
	return arr[]
}

func eight(s string) string {
	runes := []rune(s)
	slices.Reverse(runes)
	return string(runes)
}

func nine(arr []int) int {
	sum := 0
	for _, num := range arr {
		sum += num
	}
	return sum
}

type Rectangle struct {
	Width, Height float64
}

func ten(r Rectangle) float64 {
	return r.Width * r.Height
}

func main() {
	one()
	two(1, 3)
	three(3)
	four(1, 2, 3)
	five(3)
	six("Есть гласная")
	seven(40)
	eight("привет")
	nine([]int{1, 2, 3, 4, 5})
	ten(Rectangle{Width: 5, Height: 5})
}

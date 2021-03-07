package sum

func Sum(values []int) int {
	var s int
	for _, i := range values {
		s += i
	}
	return s
}

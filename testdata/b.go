package testdata

func magic(a, b int) (num int) {
	switch {
	case a == 1:
		num = 2 * b
	case b > 3:
		num = 7
	default:
		num = 0
	}
	return
}

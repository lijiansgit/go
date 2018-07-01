package libs

func IntBoolCompare(i int, b bool) bool {
	if i > 1 {
		i = 1
	}
	n := 0
	if b {
		n = 1
	}
	if n == i {
		return true
	}

	return false
}

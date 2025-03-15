package mathutils

func Adjacent(x1, y1, x2, y2 int) bool {
	return Abs(x1-x2) == 1 && y1 == y2 || Abs(y1-y2) == 1 && x1 == x2
}

func Abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

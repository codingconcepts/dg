package generator

// CartesianProduct returns the Cartesian product of a variable number of arrays.
func CartesianProduct(a ...[]string) (c [][]string) {
	if len(a) == 0 {
		return [][]string{nil}
	}

	last := len(a) - 1
	l := CartesianProduct(a[:last]...)
	for _, e := range a[last] {
		for _, p := range l {
			c = append(c, append(p, e))
		}
	}
	return
}

// Transpose a multi-dimensional array.
func Transpose(m [][]string) [][]string {
	r := make([][]string, len(m[0]))

	for x := range r {
		r[x] = make([]string, len(m))
	}

	for y, s := range m {
		for x, e := range s {
			r[x][y] = e
		}
	}
	return r
}

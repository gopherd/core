package fuzzy

// Dist calculates the distance between string s and t
func Dist(s, t string) float32 {
	dist, _, _ := editDistance(s, t)
	return dist
}

// Match returns whether s matches t with similarity
func Match(s, t string, similarity float32) bool {
	const base = 1
	dist, sizes, sizet := editDistance(s, t)
	return (float32(dist)+base)/float32(max(sizes, sizet)+base) <= (1 - similarity)
}

func editDistance(s, t string) (float32, int, int) {
	var (
		m int
		n int
		d []float32
	)
	for range s {
		m++
	}
	for range t {
		n++
	}
	d = make([]float32, (m+1)*(n+1))
	for i := 0; i < m+1; i++ {
		d[i] = float32(i)
	}
	for j := 0; j < n+1; j++ {
		d[j*(m+1)] = float32(j)
	}

	for j, x := range t {
		var jm = j * m
		for i, y := range s {
			var off = jm + i + j
			if x == y {
				d[off+m+2] = d[off]
			} else {
				d[off+m+2] = min(d[off+m+1], min(d[off+1], d[off])) + 1
			}
		}
	}

	return d[m+n*(m+1)], m, n
}

func min(x, y float32) float32 {
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}

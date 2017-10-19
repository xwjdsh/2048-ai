package utils

type Point struct {
	X, Y int
}

func Diff(a1, a2 [4][4]int) bool {
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			if a1[i][j] != a2[i][j] {
				return true
			}
		}
	}
	return false
}

package ai

import (
	"time"

	"github.com/xwjdsh/2048-ai/grid"
	"github.com/xwjdsh/2048-ai/utils"
)

type AI struct {
	Grid   *grid.Grid
	Active bool
}

type badPoint struct {
	point utils.Point
	fill  int // 2 or 4
}

var directions = []grid.Direction{
	grid.UP,
	grid.LEFT,
	grid.DOWN,
	grid.RIGHT,
}

func (a *AI) GetBest() grid.Direction {
	start := time.Now().UnixNano() / 1000000
	bestDire := grid.NONE
	dept := 0
	for ; ; dept++ {
		dire, _ := a.search(dept, 0, 0)
		if dire == grid.NONE {
			break
		}
		bestDire = dire
		if time.Now().UnixNano()/1000000-start > 100 {
			break
		}
	}
	return bestDire
}

func (a *AI) search(dept, alpha, beta int) (grid.Direction, int) {
	bestDire, bestScore := grid.NONE, 0
	tempScore := 0

	if a.Active {
		bestScore = alpha
		for _, dire := range directions {
			newGrid := a.Grid.Clone()
			if newGrid.Move(dire) {
				newAI := &AI{Grid: newGrid, Active: false}
				if dept == 0 {
					_, tempScore = bestDire, newAI.score()
				} else {
					_, tempScore = newAI.search(dept-1, bestScore, beta)
				}
				if tempScore > bestScore {
					bestScore = tempScore
					bestDire = dire
				}
				if bestScore > beta {
					return bestDire, beta
				}
			}
		}
	} else {
		bestScore = beta
		badScore := 0
		badPoints := []badPoint{}
		for _, p := range a.Grid.VacantPoints() {
			for _, f := range []int{2, 4} {
				a.Grid.Data[p.X][p.Y] = f
				score := -a.score()
				bp := badPoint{point: p, fill: f}
				if score < badScore {
					badScore = score
					badPoints = []badPoint{bp}
				} else if score == badScore {
					badPoints = append(badPoints, bp)
				}
				a.Grid.Data[p.X][p.Y] = 0
			}
		}

		for _, bp := range badPoints {
			newGrid := a.Grid.Clone()
			newGrid.Data[bp.point.X][bp.point.Y] = bp.fill
			newAI := &AI{Grid: newGrid, Active: true}
			_, tempScore = newAI.search(dept, alpha, bestScore)
			if tempScore < bestScore {
				bestScore = tempScore
			}
			if bestScore < alpha {
				return grid.NONE, alpha
			}
		}
	}
	return bestDire, bestScore
}

var (
	model1 = [][]int{
		{16, 15, 14, 13},
		{9, 10, 11, 12},
		{8, 7, 6, 5},
		{1, 2, 3, 4},
	}
	model2 = [][]int{
		{16, 15, 12, 4},
		{14, 13, 11, 3},
		{10, 9, 8, 2},
		{7, 6, 5, 1},
	}
	model3 = [][]int{
		{16, 15, 14, 4},
		{13, 12, 11, 3},
		{10, 9, 8, 2},
		{7, 6, 5, 1},
	}
)

func (a *AI) score() int {
	result := make([]int, 24)
	for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {
			modelScore(0, x, y, a.Grid.Data[x][y], model1, &result)
			modelScore(1, x, y, a.Grid.Data[x][y], model2, &result)
			modelScore(2, x, y, a.Grid.Data[x][y], model3, &result)
		}
	}
	var max int
	for _, v := range result {
		if v > max {
			max = v
		}
	}
	return max
}

func modelScore(index, x, y, value int, model [][]int, result *[]int) {
	start := index * 8
	r := *result
	r[start] += value * model[x][y]
	r[start+1] += value * model[x][3-y]

	r[start+2] += value * model[y][x]
	r[start+3] += value * model[3-y][x]

	r[start+4] += value * model[3-x][3-y]
	r[start+5] += value * model[3-x][y]

	r[start+6] += value * model[y][3-x]
	r[start+7] += value * model[3-y][3-x]
}

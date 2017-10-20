package ai

import (
	"github.com/xwjdsh/2048-ai/grid"
	"github.com/xwjdsh/2048-ai/utils"
)

type AI struct {
	Grid *grid.Grid
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

func (a *AI) Search(dept, alpha, beta int) (grid.Direction, int) {
	bestDire, bestScore := grid.NONE, 0
	tempScore := 0

	if a.Grid.Active {
		bestScore = alpha
		for _, dire := range directions {
			newGrid := a.Grid.Clone()
			if newGrid.Move(dire) {
				newGrid.Active = false
				newAI := &AI{Grid: newGrid}
				if dept == 0 {
					_, tempScore = bestDire, newAI.modelScore()
				} else {
					_, tempScore = newAI.Search(dept-1, bestScore, beta)
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
				score := -a.modelScore()
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
			newGrid.Active = true
			newAI := &AI{Grid: newGrid}
			_, tempScore = newAI.Search(dept, alpha, bestScore)
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

func (a *AI) modelScore() int {
	for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {

		}
	}
	return 0
}

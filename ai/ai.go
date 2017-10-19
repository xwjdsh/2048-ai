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
	tempDire, tempScore := grid.NONE, 0

	if a.Grid.Active {
		bestScore = alpha
		for _, dire := range directions {
			newGrid := a.Grid.Clone()
			if newGrid.Move(dire) {
				newGrid.Active = false
				newAI := &AI{Grid: newGrid}
				if dept == 0 {
					tempDire, tempScore = bestDire, newAI.modelScore()
				} else {
					tempDire, tempScore = newAI.Search(dept-1, bestScore, beta)
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
			newGrid.Active = true
			newAI := &AI{Grid: newGrid}
			tempDire, tempScore = newAI.Search(dept, alpha, bestScore)
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

func (a *AI) modelScore() int {
	return 0
}

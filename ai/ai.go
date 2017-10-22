package ai

import (
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

var expectMap = map[int]float64{
	2: 0.9,
	4: 0.1,
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

func (a *AI) Search() grid.Direction {
	var (
		bestDire  = grid.NONE
		bestScore float64
	)
	dept := a.deptSelect()
	for _, dire := range directions {
		newGrid := a.Grid.Clone()
		if newGrid.Move(dire) {
			newAI := &AI{Grid: newGrid, Active: false}
			if newScore := newAI.expectSearch(dept); newScore > bestScore {
				bestDire = dire
				bestScore = newScore
			}
		}
	}
	return bestDire
}

func (a *AI) expectSearch(dept int) float64 {
	if dept == 0 {
		return float64(a.score())
	}
	var score float64
	if a.Active {
		for _, d := range directions {
			newGrid := a.Grid.Clone()
			if newGrid.Move(d) {
				newAI := &AI{Grid: newGrid, Active: false}
				if newScore := newAI.expectSearch(dept - 1); newScore > score {
					score = newScore
				}
			}
		}
	} else {
		points := a.Grid.VacantPoints()
		for k, v := range expectMap {
			for _, point := range points {
				newGrid := a.Grid.Clone()
				newGrid.Data[point.X][point.Y] = k
				newAI := &AI{Grid: newGrid, Active: true}
				newScore := newAI.expectSearch(dept - 1)
				score += float64(newScore) * v
			}
		}
		score /= float64(len(points))
	}
	return score
}

func (a *AI) score() int {
	result := make([]int, 24)
	for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {
			if value := a.Grid.Data[x][y]; value != 0 {
				modelScore(0, x, y, value, model1, &result)
				modelScore(1, x, y, value, model2, &result)
				modelScore(2, x, y, value, model3, &result)
			}
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

func (a *AI) deptSelect() int {
	dept := 4
	max := a.Grid.Max()
	if max >= 2048 {
		dept = 6
	} else if max >= 1024 {
		dept = 5
	}
	return dept
}

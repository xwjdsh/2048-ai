package ai

import (
	"github.com/xwjdsh/2048-ai/grid"
)

type AI struct {
	// Grid is 4x4 grid.
	Grid *grid.Grid
	// Active is true represent need to select a direction to move, else represent computer need fill a number("2" or "4") into grid.
	Active bool
}

var directions = []grid.Direction{
	grid.UP,
	grid.LEFT,
	grid.DOWN,
	grid.RIGHT,
}

// The chance is 10% about fill "4" into grid and 90% fill "2" in the 2048 game.
var expectMap = map[int]float64{
	2: 0.9,
	4: 0.1,
}

var (
	// There are three model weight matrix, represents three formation for 2048 game, it from internet.
	// The evaluate function is simple and crude, so actually it's not stable.
	// If you feel interesting in evaluation function, you can read https://github.com/ovolve/2048-AI project source code.
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

// Search method compute each could move direction score result by expect search algorithm
func (a *AI) Search() grid.Direction {
	var (
		bestDire          = grid.NONE
		bestScore float64 = -1
	)
	// depth value depending on grid's max value.
	dept := a.deptSelect()
	for _, dire := range directions {
		newGrid := a.Grid.Clone()
		if newGrid.Move(dire) {
			// Could move.
			// Active is false represent computer should fill number to grid now.
			newAI := &AI{Grid: newGrid, Active: false}
			if newScore := newAI.expectSearch(dept); newScore > bestScore {
				bestDire = dire
				bestScore = newScore
			}
		}
	}
	return bestDire
}

// expect search implements
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
		// computer fill a number to grid now, it will try each vacant point with "2" or "4"
		points := a.Grid.VacantPoints()
		for k, v := range expectMap {
			for _, point := range points {
				newGrid := a.Grid.Clone()
				newGrid.Data[point.X][point.Y] = k
				// Change active, select a direction to move now.
				newAI := &AI{Grid: newGrid, Active: true}
				newScore := newAI.expectSearch(dept - 1)
				score += float64(newScore) * v
			}
		}
		score /= float64(len(points))
	}
	return score
}

// score method evaluate a grid
func (a *AI) score() int {
	result := make([]int, 24)
	for x := 0; x < 4; x++ {
		for y := 0; y < 4; y++ {
			if value := a.Grid.Data[x][y]; value != 0 {
				// get eight result(rotate and flip grid) for each model,
				modelScore(0, x, y, value, model1, &result)
				modelScore(1, x, y, value, model2, &result)
				modelScore(2, x, y, value, model3, &result)
			}
		}
	}
	// get max score in above 24 result, apply best formation
	var max int
	for _, v := range result {
		if v > max {
			max = v
		}
	}
	return max
}

// get eight result(rotate and flip grid) for each model
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

// the return value is search depth, it depending on grid's max value
// the max value larger and depth larger, this will takes more calculations and make move became slowly but maybe have a better score result.
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

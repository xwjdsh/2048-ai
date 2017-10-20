package grid

import "github.com/xwjdsh/2048-ai/utils"

type Direction int

const (
	UP Direction = iota
	RIGHT
	DOWN
	LEFT

	NONE
)

type Grid struct {
	Active bool      `json:"-"`
	Data   [4][4]int `json:"data"`
}

func (g *Grid) Clone() *Grid {
	gridClone := &Grid{}
	*gridClone = *g
	return gridClone
}

func (g *Grid) VacantPoints() []utils.Point {
	points := []utils.Point{}
	for x, row := range g.Data {
		for y, value := range row {
			if value == 0 {
				points = append(points, utils.Point{X: x, Y: y})
			}
		}
	}
	return points
}

func (g *Grid) Move(d Direction) bool {
	originData := g.Data
	data := &g.Data
	switch d {
	case UP:
		for y := 0; y < 4; y++ {
			for x := 0; x < 3; x++ {
				for nx := x + 1; nx <= 3; nx++ {
					if data[nx][y] > 0 {
						if data[x][y] <= 0 {
							data[x][y] = data[nx][y]
							data[nx][y] = 0
							x -= 1
						} else if data[x][y] == data[nx][y] {
							data[x][y] += data[nx][y]
							data[nx][y] = 0
						}
						break
					}
				}
			}
		}
	case DOWN:
		for y := 0; y < 4; y++ {
			for x := 3; x > 0; x-- {
				for nx := x - 1; nx >= 0; nx-- {
					if data[nx][y] > 0 {
						if data[x][y] <= 0 {
							data[x][y] = data[nx][y]
							data[nx][y] = 0
							x += 1
						} else if data[x][y] == data[nx][y] {
							data[x][y] += data[nx][y]
							data[nx][y] = 0
						}
						break
					}
				}
			}
		}
	case LEFT:
		for x := 0; x < 4; x++ {
			for y := 0; y < 3; y++ {
				for ny := y + 1; ny <= 3; ny++ {
					if data[x][ny] > 0 {
						if data[x][y] <= 0 {
							data[x][y] = data[x][ny]
							data[x][ny] = 0
							y -= 1
						} else if data[x][y] == data[x][ny] {
							data[x][y] += data[x][ny]
							data[x][ny] = 0
						}
						break
					}
				}
			}
		}
	case RIGHT:
		for x := 0; x < 4; x++ {
			for y := 3; y > 0; y-- {
				for ny := y - 1; ny >= 0; ny-- {
					if data[x][ny] > 0 {
						if data[x][y] <= 0 {
							data[x][y] = data[x][ny]
							data[x][ny] = 0
							y += 1
						} else if data[x][y] == data[x][ny] {
							data[x][y] += data[x][ny]
							data[x][ny] = 0
						}
						break
					}
				}
			}
		}
	}
	return utils.Diff(*data, originData)
}

package ai

import (
	"testing"

	"github.com/xwjdsh/2048-ai/grid"
)

func TestSearch(t *testing.T) {
	g := &grid.Grid{
		Data: [4][4]int{
			{2, 128, 32, 4},
			{4, 512, 128, 4},
			{2048, 32, 2, 16},
			{4096, 4, 16, 4},
		},
	}
	a := &AI{Grid: g}
	dire := a.Search()
	if dire != grid.UP && dire != grid.DOWN {
		t.Errorf("direction error, should be UP or DOWN, but got %v", dire)
	}

}

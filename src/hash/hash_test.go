package hash

import (
	"testing"
)

func TestDist(t *testing.T) {
	var x, y uint64
	x, y = 1, 0
	if Dist(x, y) != 1 {
		t.Error()
	}
}

// TODO ...

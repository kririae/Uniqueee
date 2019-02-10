package bar

import "testing"

func TestBar(t *testing.T) {
	b := New("Test", 100)
	for i := 1; i <= 100; i++ {
		b.Update(float64(i))
	}
}

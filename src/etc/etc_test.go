package etc

import (
	"testing"
)

func TestSplit(t *testing.T) {
	s := "C:\\qwq\\qwq.go"
	path, name := SplitPath(s)
	if path != "C:\\qwq\\" {
		t.Error()
	}
	if name != "qwq.go" {
		t.Error()
	}
}

func TestImage(t *testing.T) {
	s := []string{"qwq.go", "qwq.png", "qwq.jpg"}
	res := []bool{false, true, true}
	for i, v := range s {
		if IsImage(v) != res[i] {
			t.Error()
		}
	}
}
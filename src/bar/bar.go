package bar

import (
	"fmt"
	"strings"
)

// Bar ...
type Bar struct {
	prefix string
	len    float64
}

// New ...
func New(_prefix string, _len float64) *Bar {
	return &Bar{
		prefix: _prefix,
		len:    _len,
	}
}

// Update ...
func (b *Bar) Update(x float64) {
	proc := int(100 * (x / b.len))
	fmt.Printf("%s %.3f%% [%s]\r", b.prefix, x/b.len*100, completeStr(proc, "=")+completeStr(100-proc, " "))
	if proc == 100 {
		fmt.Println("")
	}
}

func completeStr(x int, s string) string {
	lst := make([]string, x)
	if x == 0 {
		return ""
	}
	return strings.Join(lst, s) + s
}

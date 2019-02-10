package hash

import (
	"image"
	"github.com/anthonynsimon/bild/transform"
)

func Ahash(img image.Image) uint64 {
	rst := transform.Resize(img, 8, 8, transform.Linear)
	px := rgb2Gray6b(rst)
	return flat(px, getAve(px))
}

func rgb2Gray6b(img image.Image) (px [][]int64) {
	h, w := img.Bounds().Size().Y, img.Bounds().Size().X
	px = make([][]int64, h)

	for i := range px {
		px[i] = make([]int64, w)
		for j := range px[i] {
			r, g, b, _ := img.At(i, j).RGBA()
			r, g, b = r/257, g/257, b/257
			// fmt.Println(r, g, b)
			gr := int64((r*19 + g*37 + b*7) >> 6)
			px[i][j] = gr
		}
	}
	return px
}

func getAve(x [][]int64) float64 {
	var sum, sz int64
	for i := range x {
		for j := range x[i] {
			sum += x[i][j]
			sz++
		}
	}
	return float64(sum) / float64(sz)
}

func flat(x [][]int64, ave float64) (ret uint64) {
	var sz uint64
	for i := range x {
		for j := range x[i] {
			if float64(x[i][j]) >= ave {
				ret |= (uint64(1) << sz)
			}
			sz++
		}
	}
	return ret
}

// TODO
func Dist(x, y uint64) int32 {
	v, res := x^y, 0
	for v > 0 {
		if v&1 == 1 {
			res++
		}
		v >>= 1
	}
	return int32(res)
}
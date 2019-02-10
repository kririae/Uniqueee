// by kririae
package main

import (
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

// SIMILARITY ...
var (
	SIMILARITY int32 = 3
	fa         []int
)

// Bar ...
type Bar struct {
	prefix string
	len    int
	proc   int
}

// Img ...
type Img struct {
	width, height      int
	filePath, fileName string
	pixels             int
	hash               uint64
}

// NewBar ...
func NewBar(_prefix string, _len int) *Bar {
	return &Bar{
		prefix: _prefix,
		len:    _len,
	}
}

// Update ...
func (b *Bar) Update(x int) {
	b.proc = x
	fmt.Printf("%s [%s]\r", b.prefix, completeStr(x, "#")+completeStr(b.len-x, " "))
	if x == b.len {
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

func main() {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	path, _ = splitPath(path)
	os.Mkdir(path+"tmp", os.ModePerm)
	_files, _ := ioutil.ReadDir(path) // 可能包含文件夹

	fmt.Println("Working on path:", path)
	b := NewBar("IO LOADING: ", 100) // 加载进度条

	var totalSize, currSize int64 // 文件总大小
	for _, f := range _files {
		if f.IsDir() || !isImage(f.Name()) {
			continue
		}
		totalSize += f.Size()
	}

	files := make([]Img, 0) // 去重文件夹之后的
	for _, f := range _files {
		if f.IsDir() || !isImage(f.Name()) {
			continue
		}
		currSize += f.Size()
		im := New(path + f.Name())
		files = append(files, *im)
		b.Update(int((float64(currSize) / float64(totalSize)) * float64(100)))
	}
	fa = make([]int, len(files))
	for i := range fa {
		fa[i] = i
	}

	for i := range files {
		for j := i + 1; j < len(files); j++ {
			if dhashDist(files[i].hash, files[j].hash) <= SIMILARITY {
				union(i, j)
			}
		}
	}

	combList := make(map[int][]Img, 0)
	for i := 0; i < len(files); i++ {
		rt := find(i)
		combList[rt] = append(combList[rt], files[i])
	}

	sub := 0
	for _, value := range combList {
		if len(value) == 1 {
			File := value[0]
			srcPath := File.filePath + File.fileName
			dstPath := File.filePath + "tmp\\" + File.fileName
			copyFile(srcPath, dstPath)
			continue
		}

		sub += len(value) - 1
		maxPixels := -1
		File := Img{}

		for _, x := range value {
			if x.pixels > maxPixels {
				maxPixels = x.pixels
				File = x
			}
		}

		srcPath := File.filePath + File.fileName
		dstPath := File.filePath + "tmp\\" + File.fileName
		copyFile(srcPath, dstPath)

		fmt.Println("Source:", File.fileName)
		for _, x := range value {
			if x.filePath+x.fileName != File.filePath+File.fileName {
				fmt.Println("---- ", x.fileName)
			}
		}
		fmt.Println("")
	}
	fmt.Println("Total uniques:", sub)
	fmt.Scanln()
}

func dhash(img *image.Image) uint64 {
	rst := transform.Resize(*img, 8, 8, transform.Linear)
	px := rgb2Gray6b(rst)
	return flat(px, getAve(px))
}

// New ...
func New(path string) *Img {
	ret := Img{}

	img, err := imgio.Open(path)

	if err != nil {
		panic(err)
	}

	ret.filePath, ret.fileName = splitPath(path)
	ret.width = img.Bounds().Size().X
	ret.height = img.Bounds().Size().Y
	ret.pixels = ret.width * ret.height

	ret.hash = dhash(&img)
	return &ret
}

func isImage(s string) bool {
	lst := strings.Split(s, ".")
	switch lst[len(lst)-1] {
	case "png":
		return true
	case "jpg":
		return true
	default:
		return false
	}
}

func copyFile(srcPath, dstPath string) {
	src, err := os.Open(srcPath)
	if err != nil {
		panic(err)
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		panic(err)
	}
	defer dst.Close()

	io.Copy(dst, src)
}

func find(x int) int {
	if fa[x] == x {
		return x
	}
	fa[x] = find(fa[x])
	return fa[x]
}

func union(x, y int) {
	x, y = find(x), find(y)
	fa[x] = y
}

func splitPath(file string) (filePath string, fileName string) {
	tmp := strings.Split(file, "\\")

	fileName = tmp[len(tmp)-1]
	filePath = strings.Join(tmp[0:len(tmp)-1], "\\") + "\\"

	return filePath, fileName
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
func dhashDist(x, y uint64) int32 {
	v, res := x^y, 0
	for v > 0 {
		if v&1 == 1 {
			res++
		}
		v >>= 1
	}
	return int32(res)
}

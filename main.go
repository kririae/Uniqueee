// by kririae
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/anthonynsimon/bild/imgio"

	"./src/bar"
	"./src/etc"
	"./src/hash"
)

// SIMILARITY ...
var (
	SIMILARITY int32 = 3
	fa         []int
)

// Img ...
type Img struct {
	width, height      int
	filePath, fileName string
	pixels             int
	hash               uint64
}

func main() {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	path, _ = etc.SplitPath(path)
	os.Mkdir(path+"tmp", os.ModePerm)
	_files, _ := ioutil.ReadDir(path) // 可能包含文件夹

	fmt.Println("Working on path:", path)
	b := bar.New("IO LOADING:", float64(1)) // 加载进度条

	var totalSize, currSize int64 // 文件总大小
	for _, f := range _files {
		if f.IsDir() || !etc.IsImage(f.Name()) {
			continue
		}
		totalSize += f.Size()
	}

	files := make([]Img, 0) // 去重文件夹之后的
	for _, f := range _files {
		if f.IsDir() || !etc.IsImage(f.Name()) {
			continue
		}
		currSize += f.Size()
		im := New(path + f.Name())
		files = append(files, *im)
		b.Update(float64(currSize) / float64(totalSize))
	}
	fa = make([]int, len(files))
	for i := range fa {
		fa[i] = i
	}

	for i := range files {
		for j := i + 1; j < len(files); j++ {
			if hash.Dist(files[i].hash, files[j].hash) <= SIMILARITY {
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
			etc.CopyFile(srcPath, dstPath)
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
		etc.CopyFile(srcPath, dstPath)

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

// New ...
func New(path string) *Img {
	ret := Img{}
	img, err := imgio.Open(path)
	if err != nil {
		panic(err)
	}

	ret.filePath, ret.fileName = etc.SplitPath(path)
	ret.width = img.Bounds().Size().X
	ret.height = img.Bounds().Size().Y
	ret.pixels = ret.width * ret.height

	ret.hash = hash.Ahash(img)
	return &ret
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

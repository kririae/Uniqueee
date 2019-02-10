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
  SIMILARITY int32 = 4
  fa         []int
)

// Bar ...
type Bar struct {
  prefix string
  len    int
  proc   int
}

func main() {
  file, _ := exec.LookPath(os.Args[0])
  path, _ := filepath.Abs(file)
  path, _ = splitPath(path)
  os.Mkdir(path+"tmp", os.ModePerm)
  _files, _ := ioutil.ReadDir(path) // 可能包含文件夹

  fmt.Println("Working on path:", path)
  b := NewBar("IO LOADING: ", 100) // 加载进度条

  var tb int64                    // 文件总大小
  files := make([]os.FileInfo, 0) // 去重文件夹之后的
  for _, f := range _files {
    if f.IsDir() || !isImage(f.Name()) {
      continue
    }
    tb += f.Size()
    files = append(files, f)
  }
  hshLst := make([]uint64, len(files))

  var cb int64
  for i, f := range files {
    hshLst[i] = calc(path + f.Name())
    cb += f.Size()
    b.Update(int((float64(cb) / float64(tb)) * float64(100)))
    // fmt.Printf("%b\n", hshLst[i])
  }
  fa = make([]int, len(hshLst))
  for i := range fa {
    fa[i] = i
  }

  for i := range hshLst {
    for j := i + 1; j < len(hshLst); j++ {
      if dist(hshLst[i], hshLst[j]) <= SIMILARITY {
        union(i, j)
      }
    }
  }

  simLst := make(map[string][]string, 0)
  for i := 0; i < len(hshLst); i++ {
    rt := find(i)
    if rt == i {
      srcPath := path + files[i].Name()
      dstPath := path + "tmp\\" + files[rt].Name()
      cpFile(srcPath, dstPath)
    } else {
      rtPath := files[rt].Name()
      simLst[rtPath] = append(simLst[rtPath], files[i].Name())
    }
  }

  for key, value := range simLst {
    fmt.Println("Source:", key)
    for _, x := range value {
      fmt.Println("---- ", x)
    }
    fmt.Println("")
  }
  fmt.Scanln()
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

func cpFile(srcPath, dstPath string) {
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
  fmt.Printf("%s [%s]\r", b.prefix, calcS(x, "#")+calcS(b.len-x, " "))
  if x == b.len {
    fmt.Println("")
  }
}

func calcS(x int, s string) string {
  lst := make([]string, x)
  if x == 0 {
    return ""
  }
  return strings.Join(lst, s) + s
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

func dist(x, y uint64) int32 {
  v, res := x^y, 0
  for v > 0 {
    v -= v & -v
    res++
  }
  return int32(res)
}

func calc(fl string) uint64 {
  img, err := imgio.Open(fl)

  if err != nil {
    panic(err)
  }

  rst := transform.Resize(img, 8, 8, transform.Linear)
  px := rgb2Gray6b(rst)
  return flat(px, getAve(px))
}

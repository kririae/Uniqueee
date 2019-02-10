package etc

import (
	"strings"
	"os"
	"io"
)

func IsImage(s string) bool {
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

func SplitPath(file string) (filePath string, fileName string) {
	tmp := strings.Split(file, "\\")

	fileName = tmp[len(tmp)-1]
	filePath = strings.Join(tmp[0:len(tmp)-1], "\\") + "\\"

	return filePath, fileName
}

func CopyFile(srcPath, dstPath string) {
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
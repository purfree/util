package file

import (
	"bytes"
	"os"
	"unicode/utf8"
)

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

// 输出的内容排版，使每列数据宽度一致。（不完全准确）
func FillBlank(l ...int) func(cells ...string) string {
	return func(cells ...string) string {
		var bb bytes.Buffer
		for i, cell := range cells {
			bb.WriteString(cell)
			var sum int
			for _, r := range cell {
				if utf8.RuneLen(r) > 2 {
					sum += 2
				} else {
					sum += 1
				}
			}
			for ; sum < l[i]; sum++ {
				bb.WriteRune(' ')
			}
		}
		return bb.String()
	}
}

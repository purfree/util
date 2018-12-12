package file

import (
	"bytes"
	"os"
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

// 预设每列数据的长度
func FillBlank(l ...int) func(cells ...string) string {
	return func(cells ...string) string {
		var bb bytes.Buffer
		for i, cell := range cells {
			var b []byte
			if i+1 > len(l) {
				b = make([]byte, len(cell))
			} else {
				b = make([]byte, l[i])
			}
			if len(cell) > len(b) {
				// 预设宽度不够，重新设置
				b = make([]byte, len(cell))
				l[i] = len(cell)
			}
			for j, _ := range b {
				if j+1 > len(cell) {
					b[j] = ' '
				} else {
					b[j] = cell[j]
				}
			}
			bb.Write(b)
		}
		return bb.String()
	}
}

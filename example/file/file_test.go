package file

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func Test_File(t *testing.T) {
	path := "test.txt"
	fpa, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	fmt.Printf("文件绝对路径：%s\n", fpa)
	fpd := filepath.Dir(fpa)
	fmt.Printf("文件所在目录：%s(传入绝对路径)", fpd)
}

/*
 读取行数据
*/
func Test_R(t *testing.T) {
	path := "test.txt"
	fpa, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	f, err := os.OpenFile(fpa, os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
	// bufio.NewReader默认缓存4096，如果一行数据大于4096，则只返回4096个字节。
	// 可以使用bufio.NewReaderSize(rd io.Reader, size int)调整缓存大小
	bi := bufio.NewReader(f)
	for {
		l, p, err := bi.ReadLine()
		if err == io.EOF {
			// 已读取到文件末尾
			break
		}
		if p {
			// 缓存已满，未能读取完整的一行数据
			// TODO 需要合并数据
			continue
		}
		// 行数据
		fmt.Println(l)
	}
}

/*
 读取数据
*/
func Test_R_2(t *testing.T) {
	path := "test.txt"
	fpa, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	f, err := os.OpenFile(fpa, os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}

	buf := make([]byte, 1024)
	for {
		n, err := f.Read(buf)
		if err == io.EOF {
			// 已读取到文件末尾
			break
		}
		fmt.Println(string(buf[:n]))
	}
}

/*
 读取所有数据
*/
func Test_R_ALL_1(t *testing.T) {
	path := "test.txt"
	fpa, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	f, err := os.Open(fpa)
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	fmt.Println(b)
}

/*
 读取所有数据
*/
func Test_R_ALL_2(t *testing.T) {
	path := "test.txt"
	fpa, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadFile(fpa)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}

func Test_W(t *testing.T) {
	path := "test.txt"
	fpa, err := filepath.Abs(path)
	if err != nil {
		panic(err)
	}
	f, err := os.OpenFile(fpa, os.O_RDWR, 0777)
	if err != nil {
		panic(err)
	}
	// 写入数据
	f.WriteString("hello world")
	// 写入数据到指定位置，该位置向后的数据会覆盖，而不是后移。
	f.WriteAt([]byte("hello world"), 2)
}

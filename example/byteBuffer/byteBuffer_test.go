package byteBuffer

import (
	"bytes"
	"fmt"
	"testing"
)

func Test_Buffer(t *testing.T) {
	buf := bytes.Buffer{}
	buf.Write([]byte("hello1"))
	buf.Write([]byte("hello2"))
	buf.Write([]byte("hello3"))
	buf.Write([]byte("hello4"))
	buf.Write([]byte("hello5"))

	b := make([]byte, 6)
	buf.Read(b)
	fmt.Println(string(b))
	buf.Read(b)
	fmt.Println(string(b))
	buf.Read(b)
	fmt.Println(string(b))
	//buf.Reset()
	//buf.Read(b)
	//fmt.Println(string(b))
	fmt.Println(buf.Len(), buf.Cap())
}

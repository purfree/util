package byteBuffer

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/purfree/util"
)

func Test_Buffer(t *testing.T) {
	buf := bytes.Buffer{}

	for i := 0; i < 1000000; i++ {
		buf.Write([]byte("hello"))
		fmt.Println(i)
	}

	//runtime.GC()
	//time.Sleep(time.Second * 2)
	//
	//b := make([]byte, 5)
	//for i := 0; i < 1000000; i++ {
	//	buf.Read(b)
	//	fmt.Println(i)
	//}
	//runtime.GC()

	buf.Reset()
	util.Wait(nil)

	//buf.Write([]byte("hello2"))
	//buf.Write([]byte("hello3"))
	//buf.Write([]byte("hello4"))
	//buf.Write([]byte("hello5"))
	//
	//b := make([]byte, 6)
	//buf.Read(b)
	//fmt.Println(string(b))
	//buf.Read(b)
	//fmt.Println(string(b))
	//buf.Read(b)
	//fmt.Println(string(b))
	//buf.Reset()
	////buf.Read(b)
	////fmt.Println(string(b))
	//fmt.Println(buf.Len(), buf.Cap())

}

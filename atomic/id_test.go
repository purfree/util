package atomic

import (
	"fmt"
	"testing"
)

func TestID_Add(t *testing.T) {
	i := &ID{}
	j := &ID{}
	i.Add()
	i.Add()
	j.Add()
	fmt.Println(i.id, j.id)
}

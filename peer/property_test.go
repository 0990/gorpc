package peer

import (
	"fmt"
	"testing"
)

type Test struct {
	name string
	age  int32
}

func Test_Property(t *testing.T) {
	c := CoreContextSet{}
	c.SetContext(3, 4)
	c.SetContext(true, 2)
	key := &Test{"xu", 30}
	c.SetContext(key, Test{"jia", 100})

	fmt.Println(c.GetContext(key))
	fmt.Println(c.GetContext(3))

	a := 1
	valuePtr := &a

	ok := c.FetchContext(3, valuePtr)
	fmt.Println(ok, *valuePtr)
}

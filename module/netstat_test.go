package module

import (
	"fmt"
	"testing"
)

func TestSocketAll(t *testing.T) {
	s := SocketAll()
	for _, line := range s {
		fmt.Println(line.string())
	}
}

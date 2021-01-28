package module

import (
	"fmt"
	"testing"
	"time"
)

func TestNewTcpProbe(t *testing.T) {
	n := NewTcpProbe("10.11.209.102", 8080, 3*time.Second)
	status, err := n.RetryProbe()
	fmt.Println(status, err)
}

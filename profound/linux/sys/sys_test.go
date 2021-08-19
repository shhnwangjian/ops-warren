package sys

import (
	"fmt"
	"testing"
)

func TestParseNodeMemInfo(t *testing.T) {
	var data = `
Node 0 MemTotal:       66881860 kB
Node 0 MemFree:         9153560 kB
Node 0 MemUsed:        57728300 kB
Node 0 Active:         44689872 kB
Node 0 Inactive:        3424948 kB
Node 0 Active(anon):   39474064 kB
Node 0 Inactive(anon):   411776 kB
Node 0 Active(file):    5215808 kB
Node 0 Inactive(file):  3013172 kB
Node 0 Unevictable:           0 kB
Node 0 Mlocked:               0 kB
Node 0 Dirty:               776 kB
Node 0 Writeback:             0 kB
Node 0 FilePages:       8923100 kB
Node 0 Mapped:           165608 kB
Node 0 AnonPages:      39179536 kB
Node 0 Shmem:            694100 kB
Node 0 KernelStack:       72512 kB
Node 0 PageTables:        78028 kB
Node 0 NFS_Unstable:          0 kB
Node 0 Bounce:                0 kB
Node 0 WritebackTmp:          0 kB
Node 0 Slab:            7233624 kB
Node 0 SReclaimable:    4147064 kB
Node 0 SUnreclaim:      3086560 kB
Node 0 AnonHugePages:  33001472 kB
Node 0 HugePages_Total:     0
Node 0 HugePages_Free:      0
Node 0 HugePages_Surp:      0
`
	fmt.Println(ParseNodeMemInfo(data))
}

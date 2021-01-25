package module

import (
	"fmt"
	"testing"
	"time"
)

func TestNewPingProbe(t *testing.T) {
	pingProbe, err := NewPingProbe("10.10.144.11")
	if err != nil {
		panic(err)
	}
	pingProbe.Count = 3
	pingProbe.SetPrivileged(true)
	pingProbe.Timeout = 10 * time.Second

	pingProbe.OnRecv = func(pkt *Packet) {
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
			pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}

	pingProbe.OnFinish = func(stats *Statistics) {
		fmt.Printf("\n--- %s ping statistics ---\n", stats.Addr)
		fmt.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
			stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		fmt.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
			stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}

	err = pingProbe.Run() // Blocks until finished.
	if err != nil {
		panic(err)
	}
	// pingProbe.Stop()
}

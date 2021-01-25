package module

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"net"
	"runtime"
	"sync"
	"syscall"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

const (
	timeSliceLength  = 8
	trackerLength    = 8
	protocolICMP     = 1
	protocolIPv6ICMP = 58
)

var (
	ipv4Proto = map[string]string{"icmp": "ip4:icmp", "udp": "udp4"}
	ipv6Proto = map[string]string{"icmp": "ip6:ipv6-icmp", "udp": "udp6"}
)

// New returns a new PingProbe struct pointer.
func New(addr string) *PingProbe {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return &PingProbe{
		Count:      -1,
		Interval:   time.Second,
		RecordRtts: true,
		Size:       timeSliceLength,
		Timeout:    time.Second * 100000,
		Tracker:    r.Int63n(math.MaxInt64),
		addr:       addr,
		done:       make(chan bool),
		id:         r.Intn(math.MaxInt16),
		ipaddr:     nil,
		isIpv4:     false,
		network:    "ip",
		protocol:   "udp",
	}
}

// NewPingProbe returns a new PingProbe and resolves the address.
func NewPingProbe(addr string) (*PingProbe, error) {
	p := New(addr)
	return p, p.Resolve()
}

// PingProbe represents a packet sender/receiver.
type PingProbe struct {
	// Interval is the wait time between each packet send. Default is 1s.
	Interval time.Duration

	// Timeout specifies a timeout before ping exits, regardless of how many
	// packets have been received.
	Timeout time.Duration

	// Count tells PingProbe to stop after sending (and receiving) Count echo
	// packets. If this option is not specified, PingProbe will operate until
	// interrupted.
	Count int

	// Debug runs in debug mode
	Debug bool

	// Number of packets sent
	PacketsSent int

	// Number of packets received
	PacketsRecv int

	// If true, keep a record of rtts of all received packets.
	// Set to false to avoid memory bloat for long running pings.
	RecordRtts bool

	// rtts is all of the Rtts
	rtts []time.Duration

	// OnSend is called when PingProbe sends a packet
	OnSend func(*Packet)

	// OnRecv is called when PingProbe receives and processes a packet
	OnRecv func(*Packet)

	// OnFinish is called when PingProbe exits
	OnFinish func(*Statistics)

	// Size of packet being sent
	Size int

	// Tracker: Used to uniquely identify packet when non-priviledged
	Tracker int64

	// Source is the source IP address
	Source string

	// stop chan bool
	done chan bool

	ipaddr *net.IPAddr
	addr   string

	isIpv4   bool
	id       int
	sequence int
	// network is one of "ip", "ip4", or "ip6".
	network string
	// protocol is "icmp" or "udp".
	protocol string
}

type packet struct {
	bytes  []byte
	nBytes int
	ttl    int
}

// Packet represents a received and processed ICMP echo packet.
type Packet struct {
	// Rtt is the round-trip time it took to ping.
	Rtt time.Duration

	// IPAddr is the address of the host being pinged.
	IPAddr *net.IPAddr

	// Addr is the string address of the host being pinged.
	Addr string

	// NBytes is the number of bytes in the message.
	Nbytes int

	// Seq is the ICMP sequence number.
	Seq int

	// TTL is the Time To Live on the packet.
	Ttl int
}

// Statistics represent the stats of a currently running or finished
// PingProbe operation.
type Statistics struct {
	// PacketsRecv is the number of packets received.
	PacketsRecv int

	// PacketsSent is the number of packets sent.
	PacketsSent int

	// PacketLoss is the percentage of packets lost.
	PacketLoss float64

	// IPAddr is the address of the host being pinged.
	IPAddr *net.IPAddr

	// Addr is the string address of the host being pinged.
	Addr string

	// Rtts is all of the round-trip times sent via this pinger.
	Rtts []time.Duration

	// MinRtt is the minimum round-trip time sent via this pinger.
	MinRtt time.Duration

	// MaxRtt is the maximum round-trip time sent via this pinger.
	MaxRtt time.Duration

	// AvgRtt is the average round-trip time sent via this pinger.
	AvgRtt time.Duration

	// StdDevRtt is the standard deviation of the round-trip times sent via
	// this pinger.
	StdDevRtt time.Duration
}

// SetIPAddr sets the ip address of the target host.
func (p *PingProbe) SetIPAddr(ipaddr *net.IPAddr) {
	p.isIpv4 = isIP4(ipaddr.IP)
	p.ipaddr = ipaddr
	p.addr = ipaddr.String()
}

// IPAddr returns the ip address of the target host.
func (p *PingProbe) GetIPAddr() *net.IPAddr {
	return p.ipaddr
}

// Resolve does the DNS lookup for the PingProbe address and sets IP protocol.
func (p *PingProbe) Resolve() error {
	if len(p.addr) == 0 {
		return errors.New("addr cannot be empty")
	}
	ipaddr, err := net.ResolveIPAddr(p.network, p.addr)
	if err != nil {
		return err
	}
	p.isIpv4 = isIP4(ipaddr.IP)
	p.ipaddr = ipaddr
	return nil
}

// SetAddr resolves and sets the ip address of the target host, addr can be a
// DNS name like "www.google.com" or IP like "127.0.0.1".
func (p *PingProbe) SetAddr(addr string) error {
	oldAddr := p.addr
	p.addr = addr
	err := p.Resolve()
	if err != nil {
		p.addr = oldAddr
		return err
	}
	return nil
}

// Addr returns the string ip address of the target host.
func (p *PingProbe) GetAddr() string {
	return p.addr
}

// SetNetwork allows configuration of DNS resolution.
// * "ip" will automatically select IPv4 or IPv6.
// * "ip4" will select IPv4.
// * "ip6" will select IPv6.
func (p *PingProbe) SetNetwork(n string) {
	switch n {
	case "ip4":
		p.network = "ip4"
	case "ip6":
		p.network = "ip6"
	default:
		p.network = "ip"
	}
}

// SetPrivileged sets the type of ping PingProbe will send.
// false means PingProbe will send an "unprivileged" UDP ping.
// true means PingProbe will send a "privileged" raw ICMP ping.
// NOTE: setting to true requires that it be run with super-user privileges.
func (p *PingProbe) SetPrivileged(privileged bool) {
	if privileged {
		p.protocol = "icmp"
	} else {
		p.protocol = "udp"
	}
}

// Run runs the PingProbe. This is a blocking function that will exit when it's
// done. If Count or Interval are not specified, it will run continuously until
// it is interrupted.
func (p *PingProbe) Run() error {
	var (
		conn *icmp.PacketConn
		err  error
	)

	if p.ipaddr == nil {
		err = p.Resolve()
	}
	if err != nil {
		return err
	}
	if p.isIpv4 {
		if conn, err = p.listen(ipv4Proto[p.protocol]); err != nil {
			return err
		}
		if err = conn.IPv4PacketConn().SetControlMessage(ipv4.FlagTTL, true); runtime.GOOS != "windows" && err != nil {
			return err
		}
	} else {
		if conn, err = p.listen(ipv6Proto[p.protocol]); err != nil {
			return err
		}
		if err = conn.IPv6PacketConn().SetControlMessage(ipv6.FlagHopLimit, true); runtime.GOOS != "windows" && err != nil {
			return err
		}
	}
	defer conn.Close()
	defer p.finish()

	var wg sync.WaitGroup
	recv := make(chan *packet, 5)
	defer close(recv)
	wg.Add(1)
	// nolint: err check
	go func() {
		err = p.recvICMP(conn, recv, &wg)
		if err != nil {
			fmt.Println("FATAL: ", err.Error())
		}
	}()

	err = p.sendICMP(conn)
	if err != nil {
		return err
	}

	timeout := time.NewTicker(p.Timeout)
	defer timeout.Stop()
	interval := time.NewTicker(p.Interval)
	defer interval.Stop()

	for {
		select {
		case <-p.done:
			wg.Wait()
			return nil
		case <-timeout.C:
			close(p.done)
			wg.Wait()
			return nil
		case <-interval.C:
			if p.Count > 0 && p.PacketsSent >= p.Count {
				continue
			}
			err = p.sendICMP(conn)
			if err != nil {
				close(p.done)
				wg.Wait()
				return err
			}
		case r := <-recv:
			err = p.processPacket(r)
			if err != nil {
				close(p.done)
				wg.Wait()
				return err
			}
		}
		if p.Count > 0 && p.PacketsRecv >= p.Count {
			close(p.done)
			wg.Wait()
			return nil
		}
	}
}

func (p *PingProbe) Stop() {
	close(p.done)
}

func (p *PingProbe) finish() {
	handler := p.OnFinish
	if handler != nil {
		s := p.Statistics()
		handler(s)
	}
}

// Statistics returns the statistics of the PingProbe. This can be run while the
// PingProbe is running or after it is finished. OnFinish calls this function to
// get it's finished statistics.
func (p *PingProbe) Statistics() *Statistics {
	loss := float64(p.PacketsSent-p.PacketsRecv) / float64(p.PacketsSent) * 100
	var min, max, total time.Duration
	if len(p.rtts) > 0 {
		min = p.rtts[0]
		max = p.rtts[0]
	}
	for _, rtt := range p.rtts {
		if rtt < min {
			min = rtt
		}
		if rtt > max {
			max = rtt
		}
		total += rtt
	}
	s := Statistics{
		PacketsSent: p.PacketsSent,
		PacketsRecv: p.PacketsRecv,
		PacketLoss:  loss,
		Rtts:        p.rtts,
		Addr:        p.addr,
		IPAddr:      p.ipaddr,
		MaxRtt:      max,
		MinRtt:      min,
	}
	if len(p.rtts) > 0 {
		s.AvgRtt = total / time.Duration(len(p.rtts))
		var sumSquares time.Duration
		for _, rtt := range p.rtts {
			sumSquares += (rtt - s.AvgRtt) * (rtt - s.AvgRtt)
		}
		s.StdDevRtt = time.Duration(math.Sqrt(
			float64(sumSquares / time.Duration(len(p.rtts)))))
	}
	return &s
}

func (p *PingProbe) recvICMP(conn *icmp.PacketConn, recv chan<- *packet, wg *sync.WaitGroup) error {
	defer wg.Done()
	for {
		select {
		case <-p.done:
			return nil
		default:
			rb := make([]byte, 512)
			if err := conn.SetReadDeadline(time.Now().Add(time.Millisecond * 100)); err != nil {
				return err
			}
			var n, ttl int
			var err error
			if p.isIpv4 {
				var cm *ipv4.ControlMessage
				n, cm, _, err = conn.IPv4PacketConn().ReadFrom(rb)
				if cm != nil {
					ttl = cm.TTL
				}
			} else {
				var cm *ipv6.ControlMessage
				n, cm, _, err = conn.IPv6PacketConn().ReadFrom(rb)
				if cm != nil {
					ttl = cm.HopLimit
				}
			}
			if err != nil {
				if netErr, ok := err.(*net.OpError); ok {
					if netErr.Timeout() {
						// Read timeout
						continue
					} else {
						close(p.done)
						return err
					}
				}
			}

			select {
			case <-p.done:
				return nil
			case recv <- &packet{bytes: rb, nBytes: n, ttl: ttl}:
			}
		}
	}
}

func (p *PingProbe) processPacket(recv *packet) error {
	receivedAt := time.Now()
	var proto int
	if p.isIpv4 {
		proto = protocolICMP
	} else {
		proto = protocolIPv6ICMP
	}

	var m *icmp.Message
	var err error
	if m, err = icmp.ParseMessage(proto, recv.bytes); err != nil {
		return fmt.Errorf("error parsing icmp message: %s", err.Error())
	}

	if m.Type != ipv4.ICMPTypeEchoReply && m.Type != ipv6.ICMPTypeEchoReply {
		// Not an echo reply, ignore it
		return nil
	}

	outPkt := &Packet{
		Nbytes: recv.nBytes,
		IPAddr: p.ipaddr,
		Addr:   p.addr,
		Ttl:    recv.ttl,
	}

	switch pkt := m.Body.(type) {
	case *icmp.Echo:
		// If we are priviledged, we can match icmp.ID
		if p.protocol == "icmp" {
			// Check if reply from same ID
			if pkt.ID != p.id {
				return nil
			}
		}

		if len(pkt.Data) < timeSliceLength+trackerLength {
			return fmt.Errorf("insufficient data received; got: %d %v",
				len(pkt.Data), pkt.Data)
		}

		tracker := bytesToInt(pkt.Data[timeSliceLength:])
		timestamp := bytesToTime(pkt.Data[:timeSliceLength])

		if tracker != p.Tracker {
			return nil
		}

		outPkt.Rtt = receivedAt.Sub(timestamp)
		outPkt.Seq = pkt.Seq
		p.PacketsRecv++
	default:
		// Very bad, not sure how this can happen
		return fmt.Errorf("invalid ICMP echo reply; type: '%T', '%v'", pkt, pkt)
	}

	if p.RecordRtts {
		p.rtts = append(p.rtts, outPkt.Rtt)
	}
	handler := p.OnRecv
	if handler != nil {
		handler(outPkt)
	}

	return nil
}

func (p *PingProbe) sendICMP(conn *icmp.PacketConn) error {
	var typ icmp.Type
	if p.isIpv4 {
		typ = ipv4.ICMPTypeEcho
	} else {
		typ = ipv6.ICMPTypeEchoRequest
	}

	var dst net.Addr = p.ipaddr
	if p.protocol == "udp" {
		dst = &net.UDPAddr{IP: p.ipaddr.IP, Zone: p.ipaddr.Zone}
	}

	t := append(timeToBytes(time.Now()), intToBytes(p.Tracker)...)
	if remainSize := p.Size - timeSliceLength - trackerLength; remainSize > 0 {
		t = append(t, bytes.Repeat([]byte{1}, remainSize)...)
	}

	body := &icmp.Echo{
		ID:   p.id,
		Seq:  p.sequence,
		Data: t,
	}

	msg := &icmp.Message{
		Type: typ,
		Code: 0,
		Body: body,
	}

	msgBytes, err := msg.Marshal(nil)
	if err != nil {
		return err
	}

	for {
		if _, err := conn.WriteTo(msgBytes, dst); err != nil {
			if netErr, ok := err.(*net.OpError); ok {
				if netErr.Err == syscall.ENOBUFS {
					continue
				}
			}
		}
		handler := p.OnSend
		if handler != nil {
			outPkt := &Packet{
				Nbytes: len(msgBytes),
				IPAddr: p.ipaddr,
				Addr:   p.addr,
				Seq:    p.sequence,
			}
			handler(outPkt)
		}

		p.PacketsSent++
		p.sequence++
		break
	}

	return nil
}

// listen icmp.ListenPacket
func (p *PingProbe) listen(netProto string) (*icmp.PacketConn, error) {
	conn, err := icmp.ListenPacket(netProto, p.Source)
	if err != nil {
		close(p.done)
		return nil, err
	}
	return conn, nil
}

func bytesToTime(b []byte) time.Time {
	var nSec int64
	for i := uint8(0); i < 8; i++ {
		nSec += int64(b[i]) << ((7 - i) * 8)
	}
	return time.Unix(nSec/1000000000, nSec%1000000000)
}

func isIP4(ip net.IP) bool {
	return len(ip.To4()) == net.IPv4len
}

func timeToBytes(t time.Time) []byte {
	nSec := t.UnixNano()
	b := make([]byte, 8)
	for i := uint8(0); i < 8; i++ {
		b[i] = byte((nSec >> ((7 - i) * 8)) & 0xff)
	}
	return b
}

func bytesToInt(b []byte) int64 {
	return int64(binary.BigEndian.Uint64(b))
}

func intToBytes(tracker int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(tracker))
	return b
}

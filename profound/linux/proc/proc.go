package proc

import (
	"encoding/binary"
	"fmt"
	"net"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	// DefaultProcMountPoint is the common mount point of the proc filesystem.
	DefaultProcMountPoint = "/proc"
)

// tcpSocketINode socket文件描述符
// sample socket:[539211079] 获取中间的数字
func tcpSocketINode(fdPath string) (iNode string) {
	if !strings.HasPrefix(fdPath, "socket:") {
		return
	}
	reg := regexp.MustCompile(`\d+`)
	iNode = reg.FindString(fdPath)
	return
}

// ParseIPV4One 转化ip方法1
// sample 3301700A 16进制的字符串
func ParseIPV4One(hexIP string) string {
	ip1, _ := strconv.ParseUint(strings.TrimSpace(hexIP[6:]), 16, 64)
	ip2, _ := strconv.ParseUint(strings.TrimSpace(hexIP[4:6]), 16, 64)
	ip3, _ := strconv.ParseUint(strings.TrimSpace(hexIP[2:4]), 16, 64)
	ip4, _ := strconv.ParseUint(strings.TrimSpace(hexIP[0:2]), 16, 64)
	return fmt.Sprintf("%d.%d.%d.%d", ip1, ip2, ip3, ip4)
}

// ParseIPv6
func ParseIPv6(s string) (string, error) {
	ip := make(net.IP, net.IPv6len)
	const grpLen = 4
	i, j := 0, 4
	for len(s) != 0 {
		grp := s[0:8]
		u, err := strconv.ParseUint(grp, 16, 32)
		binary.LittleEndian.PutUint32(ip[i:j], uint32(u))
		if err != nil {
			return "", err
		}
		i, j = i+grpLen, j+grpLen
		s = s[8:]
	}
	return ip.String(), nil
}

//  ParseIPV4Two 转化ip方法2
// sample 3301700A 16进制的字符串
func ParseIPV4Two(hexIP string) net.IP {
	b := []byte(hexIP)
	for i, j := 1, len(b)-2; i < j; i, j = i+2, j-2 { // 反转字节，转换成小端法
		b[i], b[i-1], b[j], b[j+1] = b[j+1], b[j], b[i-1], b[i]
	}
	l, _ := strconv.ParseInt(string(b), 16, 64)
	fmt.Println(l)
	return net.IPv4(byte(l>>24), byte(l>>16), byte(l>>8), byte(l))
}

// tcpConnectionStateTransform tcp链接状态对应关系
func TcpConnectionStateTransform(s string) string {
	/*
		00  "ERROR_STATUS",
		01  "TCP_ESTABLISHED",
		02  "TCP_SYN_SENT",
		03  "TCP_SYN_RECV",
		04  "TCP_FIN_WAIT1",
		05  "TCP_FIN_WAIT2",
		06  "TCP_TIME_WAIT",
		07  "TCP_CLOSE",
		08  "TCP_CLOSE_WAIT",
		09  "TCP_LAST_ACK",
		0A  "TCP_LISTEN",
		0B  "TCP_CLOSING",
	*/
	switch s {
	case "00":
		return "ERROR_STATUS"
	case "01":
		return "TCP_ESTABLISHED"
	case "02":
		return "TCP_SYN_SENT"
	case "03":
		return "TCP_SYN_RECV"
	case "04":
		return "TCP_FIN_WAIT1"
	case "05":
		return "TCP_FIN_WAIT2"
	case "06":
		return "TCP_TIME_WAIT"
	case "07":
		return "TCP_CLOSE"
	case "08":
		return "TCP_CLOSE_WAIT"
	case "09":
		return "TCP_LAST_ACK"
	case "0A":
		return "TCP_LISTEN"
	case "0B":
		return "TCP_CLOSING"
	default:
		return "unknown"
	}
}

// procFilePath /proc/{name}
func procFilePath(name string) string {
	return filepath.Join(DefaultProcMountPoint, name)
}

/*
NetTcpFile /proc/net/tcp file
https://github.com/torvalds/linux/blob/v5.9/Documentation/networking/proc_net_tcp.rst
https://www.kernel.org/doc/html/latest/networking/proc_net_tcp.html
46: 010310AC:9C4C 030310AC:1770 01
|      |      |      |      |   |--> connection state
|      |      |      |      |------> remote TCP port number
|      |      |      |-------------> remote IPv4 address
|      |      |--------------------> local TCP port number
|      |---------------------------> local IPv4 address
|----------------------------------> number of entry

00000150:00000000 01:00000019 00000000
   |        |     |     |       |--> number of unrecovered RTO timeouts （超时重传次数）
   |        |     |     |----------> number of jiffies until timer expires (超时时间，单位是jiffies)
   |        |     |----------------> timer_active (see below)
   |        |----------------------> receive-queue 当状态是ESTABLISHED，表示接收队列中数据长度；状态是LISTEN，表示已经完成连接队列的长度
   |-------------------------------> transmit-queue 发送队列中数据长度

1000        0 54165785 4 cd1e6040 25 4 27 3 -1
 |          |    |     |    |     |  | |  | |--> slow start size threshold,
 |          |    |     |    |     |  | |  |      or -1 if the threshold
 |          |    |     |    |     |  | |  |      is >= 0xFFFF
 |          |    |     |    |     |  | |  |----> sending congestion window （当前拥塞窗口大小）
 |          |    |     |    |     |  | |-------> (ack.quick<<1)|ack.pingpong
 |          |    |     |    |     |  |---------> Predicted tick of soft clock
 |          |    |     |    |     |              (delayed ACK control data)
 |          |    |     |    |     |------------> retransmit timeout （RTO，单位是clock_t）
 |          |    |     |    |------------------> location of socket in memory （socket实例的地址）
 |          |    |     |-----------------------> socket reference count （socket结构体的引用数）
 |          |    |-----------------------------> inode （套接字对应的inode）
 |          |----------------------------------> unanswered 0-window probes
 |---------------------------------------------> uid (用户id)

timer_active:
0	no timer is pending // 没有启动定时器
1	retransmit-timer is pending // 重传定时器
2	another timer (e.g. delayed ack or keepalive) is pending // 连接定时器、FIN_WAIT_2定时器或TCP保活定时器
3	this is a socket in TIME_WAIT state. Not all fields will contain data (or even exist) // TIME_WAIT定时器
4	zero window probe timer is pending // 持续定时器
*/
func NetTcpFile() string {
	return filepath.Join(procFilePath("net"), "tcp")
}

// NetTcp6File /proc/net/tcp6 file
func NetTcp6File() string {
	return filepath.Join(procFilePath("net"), "tcp6")
}

// NetUdpFile /proc/net/udp
func NetUdpFile() string {
	return filepath.Join(procFilePath("net"), "udp")
}

// NetUdp6File /proc/net/udp6
func NetUdp6File() string {
	return filepath.Join(procFilePath("net"), "udp6")
}

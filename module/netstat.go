package module

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os/user"
	"strconv"
	"strings"
	"time"

	"github.com/shhnwangjian/ops-warren/profound/linux/proc"
)

type socketEntry struct {
	proto         string // 协议
	id            int64  // sl
	localIP       string // 源地址
	localPort     int    // 源端口
	remoteIP      string // 目标地址
	remotePort    int    // 目标端口
	state         string // 状态
	transmitQueue int64
	receiveQueue  int64
	timerActive   int8
	timerDuration time.Duration
	rto           time.Duration
	uname         string // 用户
	inode         int64
}

func (s *socketEntry) string() string {
	return fmt.Sprintf("%s:%d	%s:%d	%s	%s	%d",
		s.localIP, s.localPort, s.remoteIP, s.remotePort, s.state, s.uname, s.inode)
}

// SocketAll
func SocketAll() (s []*socketEntry) {
	s = append(s, addSocket("tcp", proc.NetTcpFile())...)
	s = append(s, addSocket("tcp6", proc.NetTcp6File())...)
	s = append(s, addSocket("udp", proc.NetUdpFile())...)
	s = append(s, addSocket("udp6", proc.NetUdp6File())...)
	return
}

// addSocket
func addSocket(proto, path string) (sockets []*socketEntry) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return make([]*socketEntry, 0)
	}
	s, err := parseSocketEntry(proto, string(content))
	if err != nil {
		return make([]*socketEntry, 0)
	}
	return s
}

// parseSocketEntry
func parseSocketEntry(proto, entry string) (sockets []*socketEntry, err error) {
	reader := strings.NewReader(entry)
	scanner := bufio.NewScanner(reader)
	scanner.Split(bufio.ScanLines)
	if b := scanner.Scan(); !b {
		return nil, scanner.Err()
	}
	for scanner.Scan() {
		line := scanner.Text()
		lineList := strings.Fields(strings.TrimSpace(line))
		if len(lineList) == 0 {
			continue
		}
		if lineList[0] == "sl" {
			continue
		}
		se := &socketEntry{
			proto: proto,
		}
		for i := 0; i < len(lineList); i++ {
			switch i {
			case 0:
				id, err := strconv.ParseInt(lineList[i][:len(lineList[0])-1], 10, 64)
				if err != nil {
					continue
				}
				se.id = id
			case 1:
				ipPort := strings.Split(lineList[i], ":")
				if len(ipPort) != 2 {
					continue
				}
				if strings.HasSuffix(proto, "6") {
					se.localIP, _ = proc.ParseIPv6(ipPort[0])
				} else {
					se.localIP = proc.ParseIPV4One(ipPort[0])
				}
				port, err := strconv.ParseInt(ipPort[1], 16, 64)
				if err != nil {
					continue
				}
				se.localPort = int(port)
			case 2:
				ipPort := strings.Split(lineList[i], ":")
				if len(ipPort) != 2 {
					continue
				}
				if strings.HasSuffix(proto, "6") {
					se.remoteIP, _ = proc.ParseIPv6(ipPort[0])
				} else {
					se.remoteIP = proc.ParseIPV4One(ipPort[0])
				}
				port, err := strconv.ParseInt(ipPort[1], 16, 64)
				if err != nil {
					continue
				}
				se.remotePort = int(port)
			case 3:
				if se.proto == "tcp" || se.proto == "tcp6" {
					se.state = proc.TcpConnectionStateTransform(lineList[i])
				}
			case 4:
				queueList := strings.Split(lineList[i], ":")
				if len(queueList) != 2 {
					continue
				}
				se.transmitQueue, _ = strconv.ParseInt(queueList[0], 16, 64)
				se.receiveQueue, _ = strconv.ParseInt(queueList[1], 16, 64)
			case 7:
				u, err := strconv.ParseInt(lineList[7], 10, 64)
				if err != nil {
					continue
				}
				au, err := user.LookupId(strconv.Itoa(int(u)))
				if err != nil {
					continue
				}
				se.uname = au.Username
			case 9:
				inode, err := strconv.ParseInt(lineList[9], 10, 64)
				if err != nil {
					continue
				}
				se.inode = inode
			default:
				continue
			}
		}
		sockets = append(sockets, se)
	}
	return sockets, nil
}

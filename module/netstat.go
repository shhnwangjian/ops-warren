package module

import (
	"bufio"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/shhnwangjian/ops-warren/profound/linux"
)

type socketEntry struct {
	proto         string
	id            int64
	localIP       string
	localPort     int
	remoteIP      string
	remotePort    int
	state         string
	transmitQueue int64
	receiveQueue  int64
	timerActive   int8
	timerDuration time.Duration
	rto           time.Duration
	uid           int
	uname         string
	inode         string
}

// SocketAll
func SocketAll() (s []*socketEntry) {
	s = append(s, addSocket("tcp", linux.NetTcpFile())...)
	s = append(s, addSocket("tcp6", linux.NetTcp6File())...)
	s = append(s, addSocket("udp", linux.NetUdpFile())...)
	s = append(s, addSocket("udp6", linux.NetUdp6File())...)
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
				se.localIP = linux.TcpParseIPV4One(ipPort[0])
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
				se.remoteIP = linux.TcpParseIPV4One(ipPort[0])
				port, err := strconv.ParseInt(ipPort[1], 16, 64)
				if err != nil {
					continue
				}
				se.remotePort = int(port)
			case 3:
				se.state = linux.TcpConnectionStateTransform(lineList[i])
			case 4:
				queueList := strings.Split(lineList[i], ":")
				if len(queueList) != 2 {
					continue
				}
				se.transmitQueue, _ = strconv.ParseInt(queueList[0], 16, 64)
				se.receiveQueue, _ = strconv.ParseInt(queueList[1], 16, 64)
			case 5:

			default:
				continue
			}
		}
		sockets = append(sockets, se)
	}
	return sockets, nil
}

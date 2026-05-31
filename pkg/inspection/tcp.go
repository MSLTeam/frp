package inspection

import (
	"net"
	"sync"

	"github.com/apernet/OpenGFW/analyzer"
)


type activeTCPStream struct {
	stream analyzer.TCPStream
	name   string
}

type InspectTCPConn struct {
	net.Conn
	streams   []activeTCPStream
	totalRead int
	proxyName string
	serverToken string
	knownSafe bool
	mu        sync.Mutex // 并发安全锁
}

func WrapTCP(c net.Conn, proxyName string, serverToken string) net.Conn {
	var srcIP, dstIP net.IP
	var srcPort, dstPort uint16

	if tcpAddr, ok := c.RemoteAddr().(*net.TCPAddr); ok {
		srcIP, srcPort = tcpAddr.IP, uint16(tcpAddr.Port)
	}
	if tcpAddr, ok := c.LocalAddr().(*net.TCPAddr); ok {
		dstIP, dstPort = tcpAddr.IP, uint16(tcpAddr.Port)
	}

	info := analyzer.TCPInfo{SrcIP: srcIP, DstIP: dstIP, SrcPort: srcPort, DstPort: dstPort}

	streams := make([]activeTCPStream, 0, len(tcpAnalyzers))
	for _, a := range tcpAnalyzers {
		streams = append(streams, activeTCPStream{
			stream: a.NewTCP(info, logger),
			name:   a.Name(),
		})
	}

	return &InspectTCPConn{Conn: c, streams: streams, proxyName: proxyName,serverToken: serverToken,}
}

func (c *InspectTCPConn) Close() error {
	c.mu.Lock()
	if len(c.streams) > 0 {
		for _, s := range c.streams {
			s.stream.Close(false) // 主动释放底层资源
		}
		c.streams = nil
	}
	c.mu.Unlock()

	// 执行底层的真正关闭
	return c.Conn.Close()
}

func (c *InspectTCPConn) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)

	c.mu.Lock()
	defer c.mu.Unlock()

	// 如果客户端强制断开或发生网络错误，立刻销毁引擎资源
	if err != nil {
		if len(c.streams) > 0 {
			for _, s := range c.streams {
				s.stream.Close(false)
			}
			c.streams = nil
		}
		return n, err
	}

	if n > 0 && len(c.streams) > 0 {
		payload := b[:n]
		isStart := (c.totalRead == 0)
		c.totalRead += n

		nextActive := make([]activeTCPStream, 0, len(c.streams))
		hasFiredThreat := false

		for _, s := range c.streams {
			update, done := s.stream.Feed(false, isStart, false, 0, payload)

			if update != nil {
				detectedProto := s.name
				if update.M != nil {
					if p, ok := update.M.Get("protocol").(string); ok {
						detectedProto = p
					} else if app, ok := update.M.Get("app").(string); ok {
						detectedProto = app
					}
				}

				// 豁免检查 (让 TLS, HTTP 等认领后，解决掉 FET 误报)
				if detectedProto == "fet" && c.knownSafe {
					done = true
					continue
				}

				policy := getThreatPolicy(detectedProto)
				go reportToSystem(c.serverToken, c.RemoteAddr().String(), detectedProto, "TCP", c.proxyName)

				if policy != PolicyIgnore {
					hasFiredThreat = true // 存在威胁
					break 
				} else {
					c.knownSafe = true
					done = true // 正常下班收工
				}
			}

			if !done {
				nextActive = append(nextActive, s)
			} else {
				// 回收资源
				s.stream.Close(false)
			}
		}

		if hasFiredThreat {
			for _, s := range nextActive {
				s.stream.Close(false)
			}
			c.streams = nil // 彻底下班收工喵
		} else {
			c.streams = nextActive
		}
	}

	return n, err
}
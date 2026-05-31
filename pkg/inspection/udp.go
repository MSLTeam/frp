package inspection

import (
	"sync"
	"time"

	"github.com/apernet/OpenGFW/analyzer"
)


type activeUDPStream struct {
	stream analyzer.UDPStream
	name   string
}

type udpSession struct {
	streams  []activeUDPStream
	passed   bool
	lastSeen int64      // 记录最后的活跃时间
	mu       sync.Mutex // 会话级并发锁，保护切片和 DPI 状态机
}

var udpSessionMap sync.Map

// 周期性清扫超时的 UDP 会话
func cleanupUDPSessions() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		now := time.Now().Unix()
		udpSessionMap.Range(func(key, value interface{}) bool {
			session := value.(*udpSession)
			
			session.mu.Lock()
			// 如果 5 分钟没有新的数据包，视为会话过期，彻底销毁
			if now-session.lastSeen > 300 {
				if len(session.streams) > 0 {
					for _, s := range session.streams {
						s.stream.Close(false)
					}
					session.streams = nil
				}
				session.mu.Unlock() // 解锁后再删除
				udpSessionMap.Delete(key)
			} else {
				session.mu.Unlock() // 还没过期，放开锁
			}
			return true
		})
	}
}

func CheckUDP(srcAddr string, payload []byte, proxyName string, serverToken string) bool {
	val, ok := udpSessionMap.Load(srcAddr)
	var session *udpSession

	if !ok {
		info := analyzer.UDPInfo{}
		// 预分配容量
		streams := make([]activeUDPStream, 0, len(udpAnalyzers))
		for _, a := range udpAnalyzers {
			streams = append(streams, activeUDPStream{
				stream: a.NewUDP(info, logger),
				name:   a.Name(),
			})
		}
		session = &udpSession{
			streams:  streams,
			lastSeen: time.Now().Unix(),
		}
		actualVal, _ := udpSessionMap.LoadOrStore(srcAddr, session)
		session = actualVal.(*udpSession)
	} else {
		session = val.(*udpSession)
	}

	session.mu.Lock()
	defer session.mu.Unlock()

	session.lastSeen = time.Now().Unix() // 刷新存活时间

	if session.passed {
		return false // 已审计过，直接放行
	}

	nextActive := make([]activeUDPStream, 0, len(session.streams))
	hasFiredLog := false

	for _, s := range session.streams {
		update, done := s.stream.Feed(false, payload)

		if update != nil {
			detectedProto := s.name
			if update.M != nil {
				if p, ok := update.M.Get("protocol").(string); ok {
					detectedProto = p
				} else if app, ok := update.M.Get("app").(string); ok {
					detectedProto = app
				}
			}

			go reportToSystem(serverToken, srcAddr, detectedProto, "UDP", proxyName)
			
			hasFiredLog = true
			break
		}

		if !done {
			nextActive = append(nextActive, s)
		} else {
			// 正常识别结束的分析器，主动释放
			s.stream.Close(false)
		}
	}

	if hasFiredLog || len(nextActive) == 0 {
		// 安全释放剩余资源
		for _, s := range nextActive {
			s.stream.Close(false)
		}
		session.passed = true
		session.streams = nil
	} else {
		session.streams = nextActive
	}

	return false
}
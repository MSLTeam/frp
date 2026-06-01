package inspection

import (
	"bytes"
	"github.com/apernet/OpenGFW/analyzer"
)

// MC协议识别

type MCAnalyzer struct{}

func (a *MCAnalyzer) Name() string { return "minecraft" }
func (a *MCAnalyzer) Limit() int { return 32 }

func (a *MCAnalyzer) NewTCP(info analyzer.TCPInfo, logger analyzer.Logger) analyzer.TCPStream {
	return &mcStream{}
}

type mcStream struct{}

func (s *mcStream) Feed(rev, start, end bool, skip int, data []byte) (*analyzer.PropUpdate, bool) {
	if start && len(data) >= 2 {
		if data[0] == 0xFE {
			return &analyzer.PropUpdate{M: map[string]interface{}{"protocol": "minecraft"}}, true
		}
		// 排除 RDP 特征
		if len(data) >= 3 && data[0] == 0x03 && data[1] == 0x00 {
			return nil, true
		}
		if len(data) >= 5 && data[1] == 0x00 && data[0] > 0x05 {
			return &analyzer.PropUpdate{M: map[string]interface{}{"protocol": "minecraft"}}, true
		}
	}
	return nil, true
}
func (s *mcStream) Close(limited bool) *analyzer.PropUpdate { return nil }


// RDP 协议识别

type RDPAnalyzer struct{}

func (a *RDPAnalyzer) Name() string { return "rdp" }
func (a *RDPAnalyzer) Limit() int { return 32 }

func (a *RDPAnalyzer) NewTCP(info analyzer.TCPInfo, logger analyzer.Logger) analyzer.TCPStream {
	return &rdpStream{}
}

type rdpStream struct{}

func (s *rdpStream) Feed(rev, start, end bool, skip int, data []byte) (*analyzer.PropUpdate, bool) {
	if start && len(data) >= 6 {
		if data[0] == 0x03 && data[1] == 0x00 && data[5] == 0xE0 {
			return &analyzer.PropUpdate{M: map[string]interface{}{"protocol": "rdp"}}, true
		}
	}
	return nil, true
}
func (s *rdpStream) Close(limited bool) *analyzer.PropUpdate { return nil }


// BT 协议识别

type BTAnalyzer struct{}

func (a *BTAnalyzer) Name() string { return "bittorrent" }
// BT 握手包定长 68 字节，看前 20 字节就行
func (a *BTAnalyzer) Limit() int { return 32 }

func (a *BTAnalyzer) NewTCP(info analyzer.TCPInfo, logger analyzer.Logger) analyzer.TCPStream {
	return &btStream{}
}

type btStream struct{}

func (s *btStream) Feed(rev, start, end bool, skip int, data []byte) (*analyzer.PropUpdate, bool) {
	if start && len(data) >= 20 {
		// BT 起手式：1 字节长度 (0x13 = 19) + 19 字节协议标识
		if data[0] == 0x13 && string(data[1:20]) == "BitTorrent protocol" {
			return &analyzer.PropUpdate{M: map[string]interface{}{"protocol": "bittorrent"}}, true
		}
	}
	return nil, true
}
func (s *btStream) Close(limited bool) *analyzer.PropUpdate { return nil }


// VNC 协议

type VNCAnalyzer struct{}

func (a *VNCAnalyzer) Name() string { return "vnc" }
func (a *VNCAnalyzer) Limit() int { return 16 }

func (a *VNCAnalyzer) NewTCP(info analyzer.TCPInfo, logger analyzer.Logger) analyzer.TCPStream {
	return &vncStream{}
}

type vncStream struct{}

func (s *vncStream) Feed(rev, start, end bool, skip int, data []byte) (*analyzer.PropUpdate, bool) {
	if start && len(data) >= 12 {
		// VNC 起手式通常是：RFB 003.008\n
		if string(data[:4]) == "RFB " && data[11] == '\n' {
			return &analyzer.PropUpdate{M: map[string]interface{}{"protocol": "vnc"}}, true
		}
	}
	return nil, true
}
func (s *vncStream) Close(limited bool) *analyzer.PropUpdate { return nil }


// mysql协议

type MySQLAnalyzer struct{}

func (a *MySQLAnalyzer) Name() string { return "mysql" }
func (a *MySQLAnalyzer) Limit() int { return 256 }

func (a *MySQLAnalyzer) NewTCP(info analyzer.TCPInfo, logger analyzer.Logger) analyzer.TCPStream {
	return &mysqlStream{}
}

type mysqlStream struct{}

func (s *mysqlStream) Feed(rev, start, end bool, skip int, data []byte) (*analyzer.PropUpdate, bool) {
	if start && len(data) > 36 {
		// MySQL 客户端回复服务端的登录包，Sequence ID 通常是 1 (data[3] == 0x01)
		// 并且通常带有特定的认证插件名称特征
		if data[3] == 0x01 {
			if bytes.Contains(data, []byte("mysql_native_password")) || bytes.Contains(data, []byte("caching_sha2_password")) {
				return &analyzer.PropUpdate{M: map[string]interface{}{"protocol": "mysql"}}, true
			}
		}
	}
	return nil, true
}
func (s *mysqlStream) Close(limited bool) *analyzer.PropUpdate { return nil }


// redis协议
type RedisAnalyzer struct{}

func (a *RedisAnalyzer) Name() string { return "redis" }
func (a *RedisAnalyzer) Limit() int { return 64 }

func (a *RedisAnalyzer) NewTCP(info analyzer.TCPInfo, logger analyzer.Logger) analyzer.TCPStream {
	return &redisStream{}
}

type redisStream struct{}

func (s *redisStream) Feed(rev, start, end bool, skip int, data []byte) (*analyzer.PropUpdate, bool) {
	if start && len(data) >= 14 {
		// Redis 客户端的 RESP 协议起手通常是一个数组 (*数字\r\n)
		// 大多数客户端连上后会发 PING 或者 AUTH 命令
		if data[0] == '*' && bytes.Contains(data[:14], []byte("\r\n")) {
			// 二次确认，必须包含指令长度符号 $
			if bytes.Contains(data[:14], []byte("$4\r\nPING")) || bytes.Contains(data, []byte("AUTH")) {
				return &analyzer.PropUpdate{M: map[string]interface{}{"protocol": "redis"}}, true
			}
		}
	}
	return nil, true
}
func (s *redisStream) Close(limited bool) *analyzer.PropUpdate { return nil }

// RTMP 协议识别
type RTMPAnalyzer struct{}

func (a *RTMPAnalyzer) Name() string { return "rtmp" }
// RTMP 的 C0 和 C1 握手包头通常只需要看前 32 字节就能锁定特征
func (a *RTMPAnalyzer) Limit() int { return 32 }

func (a *RTMPAnalyzer) NewTCP(info analyzer.TCPInfo, logger analyzer.Logger) analyzer.TCPStream {
	return &rtmpStream{}
}

type rtmpStream struct{}

func (s *rtmpStream) Feed(rev, start, end bool, skip int, data []byte) (*analyzer.PropUpdate, bool) {
	if start && len(data) >= 9 {
		// RTMP 客户端发起的 C0 握手永远是 1 个字节：0x03 (代表 RTMP 版本 3)
		if data[0] == 0x03 {
			// 防误判：排除 RDP 协议 (RDP 的 TPDU 也是 0x03 开头，但它的第2和第6字节是固定的网络层特征)
			if len(data) >= 6 && data[1] == 0x00 && data[5] == 0xE0 {
				return nil, true
			}

			// 接下来是 C1 块 (1536 字节)。
			// C1 的前 4 字节是时间戳，紧接着的 4 字节是版本号。
			// 严格的标准 RTMP 推流中，这 4 个字节要求是 0x00 0x00 0x00 0x00
			if data[5] == 0x00 && data[6] == 0x00 && data[7] == 0x00 && data[8] == 0x00 {
				return &analyzer.PropUpdate{M: map[string]interface{}{"protocol": "rtmp"}}, true
			}

			// 兼容现代推流工具（如 OBS、FFmpeg）。它们经常在这里塞入自己的 Flash/客户端版本号。
			// 常见的推流工具版本号首字节通常是 0x0A, 0x09, 0x80, 0x04 等。
			if data[5] == 0x0A || data[5] == 0x09 || data[5] == 0x80 || data[5] == 0x04 {
				return &analyzer.PropUpdate{M: map[string]interface{}{"protocol": "rtmp"}}, true
			}
		}
	}
	return nil, true
}
func (s *rtmpStream) Close(limited bool) *analyzer.PropUpdate { return nil }


// RTSP 协议识别 (安防摄像头、局域网串流)
type RTSPAnalyzer struct{}

func (a *RTSPAnalyzer) Name() string { return "rtsp" }
func (a *RTSPAnalyzer) Limit() int { return 64 }

func (a *RTSPAnalyzer) NewTCP(info analyzer.TCPInfo, logger analyzer.Logger) analyzer.TCPStream {
	return &rtspStream{}
}

type rtspStream struct{}

func (s *rtspStream) Feed(rev, start, end bool, skip int, data []byte) (*analyzer.PropUpdate, bool) {
	if start && len(data) >= 10 {
		// RTSP 协议是明文控制流，推流或拉流的起手式通常是 OPTIONS 或 SETUP
		if bytes.HasPrefix(data, []byte("OPTIONS ")) || 
		   bytes.HasPrefix(data, []byte("SETUP ")) || 
		   bytes.HasPrefix(data, []byte("DESCRIBE ")) ||
		   bytes.HasPrefix(data, []byte("PLAY ")) {
			
			// 二次校验：必须包含 rtsp 标识
			if bytes.Contains(data, []byte("rtsp://")) || bytes.Contains(data, []byte("RTSP/1.0\r\n")) {
				return &analyzer.PropUpdate{M: map[string]interface{}{"protocol": "rtsp"}}, true
			}
		}
	}
	return nil, true
}
func (s *rtspStream) Close(limited bool) *analyzer.PropUpdate { return nil }
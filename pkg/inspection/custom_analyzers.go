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
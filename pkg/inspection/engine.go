package inspection

import (
	"github.com/apernet/OpenGFW/analyzer"
	"github.com/apernet/OpenGFW/analyzer/tcp"
	"github.com/apernet/OpenGFW/analyzer/udp"
)

var (
	tcpAnalyzers []analyzer.TCPAnalyzer
	udpAnalyzers []analyzer.UDPAnalyzer
	logger       = &gfwLogger{}
)

func init() {
	tcpAnalyzers = []analyzer.TCPAnalyzer{
		&tcp.TLSAnalyzer{},
		&tcp.HTTPAnalyzer{},
		&tcp.SSHAnalyzer{},
		&MCAnalyzer{},
		&RDPAnalyzer{}, 
		&BTAnalyzer{},    // BT 检测
		&VNCAnalyzer{},   // VNC 检测
		&MySQLAnalyzer{}, // MySQL 检测
		&RedisAnalyzer{}, // Redis 检测
		&tcp.SocksAnalyzer{},
		&tcp.TrojanAnalyzer{},
		&udp.OpenVPNAnalyzer{},
		&tcp.FETAnalyzer{},
	}

	udpAnalyzers = []analyzer.UDPAnalyzer{
		&udp.WireGuardAnalyzer{},
		&udp.OpenVPNAnalyzer{},
		&udp.QUICAnalyzer{},
	}

	// 启动 UDP 会话清理协程
	go cleanupUDPSessions()
}
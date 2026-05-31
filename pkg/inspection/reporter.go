package inspection

import (
	"sync"
	"time"

	"github.com/fatedier/frp/pkg/mslapi"
	"github.com/fatedier/frp/pkg/util/log"
)

const (
	// 【配置项】
	AuditWindowSeconds = 120 // 统计时间窗口 (秒)
	AuditThreshold     = 20  // 触发上报系统的阈值次数

	// ANSI 颜色转义码 (在部分纯文本 log 文件中可能会显示为转义字符，如果不喜欢可以删掉)
	colorReset  = "\033[0m"
	colorYellow = "\033[33m"
)

// 集中定义哪些协议属于“需要触发风控报警的威胁”
func isThreat(proto string) bool {
	switch proto {
	case "fet", "trojan", "socks", "openvpn", "wireguard", "bittorrent", "rdp":
		return true
	default:
		// mc, vnc, mysql, redis 等视为安全业务协议
		return false
	}
}

type alertCounter struct {
	count     int
	startTime int64
	mu        sync.Mutex
}

var alertCache sync.Map

func reportToSystem(serverToken string, ip string, proto string, netType string, proxyName string, isThreatProto bool) {
	log.Infof("[GFW_AUDIT] 特征命中 | 隧道: [%s] | 网络: %s | 来源IP: %s | 识别协议: %s",
		proxyName, netType, ip, proto)

	if !isThreatProto {
		return
	}

	// 聚合统计逻辑 (以 "隧道名_协议名" 为粒度进行防抖)
	cacheKey := proxyName + "_" + proto
	val, ok := alertCache.Load(cacheKey)
	var counter *alertCounter

	if !ok {
		counter = &alertCounter{startTime: time.Now().Unix(), count: 0}
		val, _ = alertCache.LoadOrStore(cacheKey, counter)
		counter = val.(*alertCounter)
	} else {
		counter = val.(*alertCounter)
	}

	// ========== 临界区开始 ==========
	counter.mu.Lock()
	now := time.Now().Unix()

	if now-counter.startTime > AuditWindowSeconds {
		counter.startTime = now
		counter.count = 0
	}

	counter.count++
	currentCount := counter.count
	shouldReport := currentCount >= AuditThreshold

	if shouldReport {
		// 达到阈值，重置计数器防止重复上报
		counter.count = 0
		counter.startTime = now
	}
	counter.mu.Unlock()
	// ========== 临界区结束 ==========

	// 异步触发上报
	if shouldReport {
		log.Warnf("%s[GFW_WARNING] 流量审计异常 | 隧道 [%s] 在 %d 秒内受到 %d 次 [%s] 协议探测！正在调用 API...%s",
			colorYellow, proxyName, AuditWindowSeconds, currentCount, proto, colorReset)

		go sendToBackend(serverToken, proxyName, proto, netType, ip, currentCount)
	}
}

func sendToBackend(serverToken, proxyName, protocol, netType, srcIp string, count int) {
	apiService, err := mslapi.MyAPIService()
	if err != nil {
		log.Errorf("[GFW_ERROR] 无法获取 API Service: %v", err)
		return
	}

	payload := mslapi.ThreatPayload{
		ServerToken:  serverToken,
		ProxyName:    proxyName,
		Protocol:     protocol,
		Type:         netType,
		SrcIp:        srcIp,
		TriggerCount: count,
	}

	retMsg, err := apiService.SubmitThreatLog(payload)
	if err != nil {
		log.Errorf("[GFW_ERROR] 流量风控上报失败 | 错误: %v | 后端信息: %s", err, retMsg)
	} else {
		log.Infof("%s[GFW_SUCCESS] 流量风控上报成功 | 隧道: %s | 结果: %s%s",
			colorYellow, proxyName, retMsg, colorReset)
	}
}

type gfwLogger struct{}

func (l *gfwLogger) Debugf(format string, args ...interface{}) {
	log.Debugf(format, args...)
}
func (l *gfwLogger) Infof(format string, args ...interface{}) {
	log.Infof(format, args...)
}
func (l *gfwLogger) Errorf(format string, args ...interface{}) {
	log.Errorf(format, args...)
}

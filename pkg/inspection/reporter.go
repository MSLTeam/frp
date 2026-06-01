package inspection

import (
	"sync"
	"time"

	"github.com/fatedier/frp/pkg/mslapi"
	"github.com/fatedier/frp/pkg/util/log"
)

// 计数聚合配置
const (
	AuditWindowSeconds = 10 
	AuditThreshold     = 5  
)

type Policy int

const (
	PolicyIgnore  Policy = iota // 0: 忽略 (安全业务)
	PolicyCount                 // 1: 计数聚合 (易误报协议，走高频并发检测)
	PolicyInstant               // 2: 立即上报 (实锤违规协议，一发入魂)
)

// getThreatPolicy 判定协议的危险级别和审计策略
func getThreatPolicy(proto string) Policy {
	switch proto {
	// 实锤违规：抓到一次握手特征直接击毙
	case "trojan", "socks", "openvpn", "wireguard", "bittorrent":
		return PolicyInstant

	// 嫌疑观察：FET(全加密流量)
	case "fet":
		return PolicyCount

	default:
		// mc, rdp, vnc, mysql, redis, tls, http 等全部视为安全
		return PolicyIgnore
	}
}

type alertCounter struct {
	count     int
	startTime int64
	mu        sync.Mutex
}

var alertCache sync.Map

func reportToSystem(serverToken string, ip string, proto string, netType string, proxyName string) {
	log.Infof("[GFW_AUDIT] 特征命中 | 隧道: [%s] | 网络: %s | 来源IP: %s | 识别协议: %s",
		proxyName, netType, ip, proto)

	policy := getThreatPolicy(proto)
	if policy == PolicyIgnore {
		return
	}

	currentCount := 1
	shouldReport := false

	if policy == PolicyInstant {
		shouldReport = true
	} else if policy == PolicyCount {
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

		counter.mu.Lock()
		now := time.Now().Unix()

		if now-counter.startTime > AuditWindowSeconds {
			counter.startTime = now
			counter.count = 0
		}

		counter.count++
		currentCount = counter.count
		shouldReport = currentCount >= AuditThreshold

		if shouldReport {
			counter.count = 0
			counter.startTime = now
		}
		counter.mu.Unlock()
	}

	if shouldReport {
		if policy == PolicyInstant {
			log.Warnf("[GFW_WARNING] 触发现行违规 | 隧道 [%s] 捕获到实锤高危协议 [%s]！正在即时同步云端...", proxyName, proto)
		} else {
			log.Warnf("[GFW_WARNING] 流量并发异常 | 隧道 [%s] 在 %d 秒内受到 %d 次 [%s] 协议请求！疑似代理穿透，正在上报...",
				proxyName, AuditWindowSeconds, currentCount, proto)
		}

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
		log.Errorf("[GFW_ERROR] 风控上报失败 | 错误: %v | 后端信息: %s", err, retMsg)
	} else {
		log.Infof("[GFW_SUCCESS] 风控同步成功 | 隧道: %s | 结果: %s", proxyName, retMsg)
	}
}

type gfwLogger struct{}

func (l *gfwLogger) Debugf(format string, args ...interface{}) { log.Debugf(format, args...) }
func (l *gfwLogger) Infof(format string, args ...interface{})  { log.Infof(format, args...) }
func (l *gfwLogger) Errorf(format string, args ...interface{}) { log.Errorf(format, args...) }
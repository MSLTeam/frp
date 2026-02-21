package inspect

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	frpint "github.com/apernet/OpenGFW/integration/frp"
	"github.com/samber/lo"

	v1 "github.com/fatedier/frp/pkg/config/v1"
	"github.com/fatedier/frp/pkg/util/log"
)

var ErrBlockedByOpenGFW = errors.New("blocked by openGFW policy")

type OpenGFW struct {
	enabled  bool
	sessions *frpint.SessionManager

	sessionSeq atomic.Uint64
	sessionMu  sync.Map

	stateMu sync.Mutex
	states  map[string]*sessionState
}

type sessionState struct {
	warned  bool
	blocked bool
}

func NewOpenGFW(cfg v1.OpenGFWServerConfig) (*OpenGFW, error) {
	cfg.Complete()
	og := &OpenGFW{
		enabled: lo.FromPtr(cfg.Enable),
		states:  make(map[string]*sessionState),
	}
	if !og.enabled {
		log.Infof("openGFW inspection disabled")
		return og, nil
	}

	engineCfg := frpint.DefaultConfig()
	engineCfg.ProxyFeatureFile = strings.TrimSpace(cfg.ProxyFeatureFile)
	engineCfg.TrafficFeatureFile = strings.TrimSpace(cfg.TrafficFeatureFile)
	engineCfg.ProxyPolicy = frpint.ProxyPolicy{
		Enabled:         lo.FromPtr(cfg.ProxyPolicy.Enable),
		HighThreshold:   cfg.ProxyPolicy.HighThreshold,
		MediumThreshold: cfg.ProxyPolicy.MediumThreshold,
	}
	engineCfg.TrafficPolicy = frpint.TrafficPolicy{
		Enabled:     lo.FromPtr(cfg.TrafficPolicy.Enable),
		BlockWebTCP: lo.FromPtr(cfg.TrafficPolicy.BlockWebTCP),
	}

	engine, err := frpint.NewEngine(engineCfg)
	if err != nil {
		return nil, err
	}
	og.sessions = frpint.NewSessionManager(engine)

	log.Infof(
		"openGFW inspection enabled, proxyPolicy(enabled=%v high=%d medium=%d), trafficPolicy(enabled=%v blockWebTCP=%v)",
		engineCfg.ProxyPolicy.Enabled,
		engineCfg.ProxyPolicy.HighThreshold,
		engineCfg.ProxyPolicy.MediumThreshold,
		engineCfg.TrafficPolicy.Enabled,
		engineCfg.TrafficPolicy.BlockWebTCP,
	)
	return og, nil
}

func (og *OpenGFW) Enabled() bool {
	return og != nil && og.enabled && og.sessions != nil
}

func (og *OpenGFW) NewTCPSessionID(proxyName string) string {
	if proxyName == "" {
		proxyName = "proxy"
	}
	seq := og.sessionSeq.Add(1)
	return fmt.Sprintf("%s-tcp-%d", proxyName, seq)
}

func (og *OpenGFW) NewUDPSessionID(proxyName string, srcAddr, dstAddr *net.UDPAddr) string {
	if proxyName == "" {
		proxyName = "proxy"
	}
	src := "unknown-src"
	dst := "unknown-dst"
	if srcAddr != nil {
		src = srcAddr.String()
	}
	if dstAddr != nil {
		dst = dstAddr.String()
	}
	return fmt.Sprintf("%s-udp-%s-%s", proxyName, src, dst)
}

func (og *OpenGFW) WrapTCPConn(conn net.Conn, sessionID string, srcAddr, dstAddr net.Addr, readIsRev bool) net.Conn {
	if !og.Enabled() || conn == nil {
		return conn
	}
	meta, ok := streamMetaFromAddrs(srcAddr, dstAddr)
	if !ok {
		log.Warnf("openGFW skipped tcp session[%s], invalid addr src=%v dst=%v", sessionID, srcAddr, dstAddr)
		return conn
	}
	return &inspectedTCPConn{
		Conn:      conn,
		inspector: og,
		sessionID: sessionID,
		meta:      meta,
		readIsRev: readIsRev,
	}
}

func (og *OpenGFW) InspectUDP(sessionID string, srcAddr, dstAddr *net.UDPAddr, rev bool, payload []byte) bool {
	if !og.Enabled() || len(payload) == 0 || srcAddr == nil || dstAddr == nil {
		return false
	}
	unlock := og.lockSession(sessionID)
	defer unlock()

	meta := frpint.StreamMeta{
		SrcIP:   append(net.IP(nil), srcAddr.IP...),
		DstIP:   append(net.IP(nil), dstAddr.IP...),
		SrcPort: uint16(srcAddr.Port),
		DstPort: uint16(dstAddr.Port),
	}
	rep, err := og.sessions.FeedUDP(sessionID, meta, rev, payload)
	if err != nil {
		log.Warnf("openGFW inspect udp session[%s] failed: %v", sessionID, err)
		return false
	}
	og.maybeLogDecision(sessionID, rep)
	if rep.Action == frpint.ActionBlock {
		if closeRep, ok := og.sessions.CloseUDP(sessionID, false); ok {
			og.maybeLogDecision(sessionID, closeRep)
		}
		og.clearSessionState(sessionID)
		return true
	}
	if rep.Done {
		og.clearSessionState(sessionID)
	}
	return false
}

func (og *OpenGFW) CloseTCPSession(sessionID string) {
	if !og.Enabled() || sessionID == "" {
		return
	}
	unlock := og.lockSession(sessionID)
	defer unlock()

	if rep, ok := og.sessions.CloseTCP(sessionID, false); ok {
		og.maybeLogDecision(sessionID, rep)
	}
	og.clearSessionState(sessionID)
}

func (og *OpenGFW) CloseUDPSession(sessionID string) {
	if !og.Enabled() || sessionID == "" {
		return
	}
	unlock := og.lockSession(sessionID)
	defer unlock()

	if rep, ok := og.sessions.CloseUDP(sessionID, false); ok {
		og.maybeLogDecision(sessionID, rep)
	}
	og.clearSessionState(sessionID)
}

func (og *OpenGFW) inspectTCPPayload(sessionID string, meta frpint.StreamMeta, rev, start bool, payload []byte) bool {
	if !og.Enabled() || len(payload) == 0 {
		return false
	}
	unlock := og.lockSession(sessionID)
	defer unlock()

	rep, err := og.sessions.FeedTCP(sessionID, meta, rev, start, false, 0, payload)
	if err != nil {
		log.Warnf("openGFW inspect tcp session[%s] failed: %v", sessionID, err)
		return false
	}
	og.maybeLogDecision(sessionID, rep)
	return rep.Action == frpint.ActionBlock
}

func (og *OpenGFW) maybeLogDecision(sessionID string, rep frpint.InspectionReport) {
	if !og.Enabled() {
		return
	}
	switch rep.Action {
	case frpint.ActionWarn:
		if !og.markWarned(sessionID) {
			return
		}
		log.Warnf("openGFW warn session[%s], reason=%s%s", sessionID, rep.Reason, formatReport(rep))
	case frpint.ActionBlock:
		if !og.markBlocked(sessionID) {
			return
		}
		log.Warnf("openGFW block session[%s], reason=%s%s", sessionID, rep.Reason, formatReport(rep))
	}
}

func (og *OpenGFW) markWarned(sessionID string) bool {
	og.stateMu.Lock()
	defer og.stateMu.Unlock()

	st := og.getOrCreateStateLocked(sessionID)
	if st.warned {
		return false
	}
	st.warned = true
	return true
}

func (og *OpenGFW) markBlocked(sessionID string) bool {
	og.stateMu.Lock()
	defer og.stateMu.Unlock()

	st := og.getOrCreateStateLocked(sessionID)
	if st.blocked {
		return false
	}
	st.blocked = true
	return true
}

func (og *OpenGFW) getOrCreateStateLocked(sessionID string) *sessionState {
	if st, ok := og.states[sessionID]; ok {
		return st
	}
	st := &sessionState{}
	og.states[sessionID] = st
	return st
}

func (og *OpenGFW) clearSessionState(sessionID string) {
	og.stateMu.Lock()
	delete(og.states, sessionID)
	og.stateMu.Unlock()
	og.sessionMu.Delete(sessionID)
}

func (og *OpenGFW) lockSession(sessionID string) func() {
	if sessionID == "" {
		return func() {}
	}
	raw, _ := og.sessionMu.LoadOrStore(sessionID, &sync.Mutex{})
	mu := raw.(*sync.Mutex)
	mu.Lock()
	return mu.Unlock
}

func formatReport(rep frpint.InspectionReport) string {
	out := ""
	if rep.Proxy != nil {
		out += fmt.Sprintf(", proxy(protocol=%s type=%s class=%s score=%d level=%s)",
			rep.Proxy.Protocol,
			rep.Proxy.Type,
			rep.Proxy.Class,
			rep.Proxy.SensitivityScore,
			rep.Proxy.SensitivityLevel,
		)
	}
	if rep.Traffic != nil {
		out += fmt.Sprintf(", traffic(label=%s family=%s role=%s method=%s transport=%s confidence=%d)",
			rep.Traffic.Label,
			rep.Traffic.Family,
			rep.Traffic.Role,
			rep.Traffic.Method,
			rep.Traffic.Transport,
			rep.Traffic.Confidence,
		)
	}
	return out
}

type inspectedTCPConn struct {
	net.Conn

	inspector *OpenGFW
	sessionID string
	meta      frpint.StreamMeta
	readIsRev bool

	readStarted  atomic.Bool
	writeStarted atomic.Bool
	blocked      atomic.Bool
	closeOnce    sync.Once
}

func (c *inspectedTCPConn) Read(p []byte) (int, error) {
	if c.blocked.Load() {
		return 0, ErrBlockedByOpenGFW
	}
	n, err := c.Conn.Read(p)
	if n > 0 {
		start := c.readStarted.CompareAndSwap(false, true)
		if c.inspector.inspectTCPPayload(c.sessionID, c.meta, c.readIsRev, start, p[:n]) {
			c.block()
			return 0, ErrBlockedByOpenGFW
		}
	}
	return n, err
}

func (c *inspectedTCPConn) Write(p []byte) (int, error) {
	if c.blocked.Load() {
		return 0, ErrBlockedByOpenGFW
	}
	if len(p) > 0 {
		start := c.writeStarted.CompareAndSwap(false, true)
		if c.inspector.inspectTCPPayload(c.sessionID, c.meta, !c.readIsRev, start, p) {
			c.block()
			return 0, ErrBlockedByOpenGFW
		}
	}
	return c.Conn.Write(p)
}

func (c *inspectedTCPConn) Close() error {
	c.closeOnce.Do(func() {
		if c.inspector != nil {
			c.inspector.CloseTCPSession(c.sessionID)
		}
	})
	return c.Conn.Close()
}

func (c *inspectedTCPConn) block() {
	if !c.blocked.CompareAndSwap(false, true) {
		return
	}
	if c.inspector != nil {
		c.inspector.CloseTCPSession(c.sessionID)
	}
	_ = c.Conn.Close()
}

func streamMetaFromAddrs(srcAddr, dstAddr net.Addr) (frpint.StreamMeta, bool) {
	srcIP, srcPort, ok := ipPortFromAddr(srcAddr)
	if !ok {
		return frpint.StreamMeta{}, false
	}
	dstIP, dstPort, ok := ipPortFromAddr(dstAddr)
	if !ok {
		return frpint.StreamMeta{}, false
	}
	return frpint.StreamMeta{
		SrcIP:   srcIP,
		DstIP:   dstIP,
		SrcPort: srcPort,
		DstPort: dstPort,
	}, true
}

func ipPortFromAddr(addr net.Addr) (net.IP, uint16, bool) {
	switch a := addr.(type) {
	case *net.TCPAddr:
		if a == nil || a.IP == nil || a.Port < 0 || a.Port > 65535 {
			return nil, 0, false
		}
		return append(net.IP(nil), a.IP...), uint16(a.Port), true
	case *net.UDPAddr:
		if a == nil || a.IP == nil || a.Port < 0 || a.Port > 65535 {
			return nil, 0, false
		}
		return append(net.IP(nil), a.IP...), uint16(a.Port), true
	}
	if addr == nil {
		return nil, 0, false
	}
	host, portStr, err := net.SplitHostPort(addr.String())
	if err != nil {
		return nil, 0, false
	}
	if idx := strings.LastIndex(host, "%"); idx >= 0 {
		host = host[:idx]
	}
	host = strings.Trim(host, "[]")
	ip := net.ParseIP(host)
	if ip == nil {
		return nil, 0, false
	}
	port, err := strconv.Atoi(portStr)
	if err != nil || port < 0 || port > 65535 {
		return nil, 0, false
	}
	return append(net.IP(nil), ip...), uint16(port), true
}

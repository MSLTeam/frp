package inspect

import (
	"errors"
	"net"
	"testing"

	frpint "github.com/apernet/OpenGFW/integration/frp"

	v1 "github.com/fatedier/frp/pkg/config/v1"
)

func newTestInspector(t *testing.T) *OpenGFW {
	t.Helper()
	og, err := NewOpenGFW(v1.OpenGFWServerConfig{})
	if err != nil {
		t.Fatalf("new inspector: %v", err)
	}
	if !og.Enabled() {
		t.Fatal("inspector should be enabled by default")
	}
	return og
}

func tcpMeta(srcIP, dstIP net.IP, srcPort, dstPort int) frpint.StreamMeta {
	return frpint.StreamMeta{
		SrcIP:   srcIP,
		DstIP:   dstIP,
		SrcPort: uint16(srcPort),
		DstPort: uint16(dstPort),
	}
}

func TestInspectClassifyAndProxySeverity(t *testing.T) {
	og := newTestInspector(t)

	blockSession := og.NewTCPSessionID("proxy")
	blockRep, err := og.sessions.FeedTCP(
		blockSession,
		tcpMeta(net.IPv4(10, 0, 0, 1), net.IPv4(10, 0, 0, 2), 50000, 1080),
		false, true, false, 0,
		[]byte{0x05, 0x01, 0x00},
	)
	if err != nil {
		t.Fatalf("feed tcp block session: %v", err)
	}
	if blockRep.Action != frpint.ActionBlock {
		t.Fatalf("action=%s, want block", blockRep.Action)
	}
	if blockRep.Proxy == nil || blockRep.Proxy.Protocol != "socks5" {
		t.Fatalf("proxy report=%+v, want protocol socks5", blockRep.Proxy)
	}
	if blockRep.Traffic == nil || blockRep.Traffic.Family != "proxy" {
		t.Fatalf("traffic report=%+v, want family proxy", blockRep.Traffic)
	}

	warnSession := og.NewTCPSessionID("proxy")
	warnRep, err := og.sessions.FeedTCP(
		warnSession,
		tcpMeta(net.IPv4(10, 0, 0, 3), net.IPv4(10, 0, 0, 4), 50001, 6022),
		false, true, false, 0,
		[]byte("SSH-2.0-OpenSSH_9.7\r\n"),
	)
	if err != nil {
		t.Fatalf("feed tcp warn session: %v", err)
	}
	if warnRep.Action != frpint.ActionWarn {
		t.Fatalf("action=%s, want warn", warnRep.Action)
	}
	if warnRep.Proxy == nil || warnRep.Proxy.Protocol != "ssh_tunnel" {
		t.Fatalf("proxy report=%+v, want protocol ssh_tunnel", warnRep.Proxy)
	}
}

func TestInspectBlockWebTCP(t *testing.T) {
	og := newTestInspector(t)

	httpSession := og.NewTCPSessionID("web")
	httpRep, err := og.sessions.FeedTCP(
		httpSession,
		tcpMeta(net.IPv4(10, 0, 1, 1), net.IPv4(10, 0, 1, 2), 50010, 80),
		false, true, false, 0,
		[]byte("GET / HTTP/1.1\r\nHost: example.com\r\n\r\n"),
	)
	if err != nil {
		t.Fatalf("feed http session: %v", err)
	}
	if httpRep.Action != frpint.ActionBlock {
		t.Fatalf("action=%s, want block", httpRep.Action)
	}
	if httpRep.Reason != "traffic policy block web tcp" {
		t.Fatalf("reason=%q, want traffic policy block web tcp", httpRep.Reason)
	}
	if httpRep.Traffic == nil || httpRep.Traffic.Label != "Web Service (HTTP)" {
		t.Fatalf("traffic report=%+v, want Web Service (HTTP)", httpRep.Traffic)
	}

	tlsSession := og.NewTCPSessionID("web")
	tlsRep, err := og.sessions.FeedTCP(
		tlsSession,
		tcpMeta(net.IPv4(10, 0, 1, 3), net.IPv4(10, 0, 1, 4), 50011, 443),
		false, true, false, 0,
		[]byte{0x16, 0x03, 0x03, 0x00, 0x04, 0x01, 0x00, 0x00, 0x00},
	)
	if err != nil {
		t.Fatalf("feed tls session: %v", err)
	}
	if tlsRep.Action != frpint.ActionBlock {
		t.Fatalf("action=%s, want block", tlsRep.Action)
	}
	if tlsRep.Traffic == nil || tlsRep.Traffic.Label != "Web Service (HTTPS/TLS)" {
		t.Fatalf("traffic report=%+v, want Web Service (HTTPS/TLS)", tlsRep.Traffic)
	}
}

func TestWrappedTCPConnCanBlockHTTPPayload(t *testing.T) {
	og := newTestInspector(t)

	serverConn, clientConn := net.Pipe()
	t.Cleanup(func() {
		_ = clientConn.Close()
		_ = serverConn.Close()
	})

	sessionID := og.NewTCPSessionID("wrapped")
	wrapped := og.WrapTCPConn(
		serverConn,
		sessionID,
		&net.TCPAddr{IP: net.IPv4(10, 1, 0, 1), Port: 52000},
		&net.TCPAddr{IP: net.IPv4(10, 1, 0, 2), Port: 80},
		false,
	)
	t.Cleanup(func() { _ = wrapped.Close() })

	go func() {
		_, _ = clientConn.Write([]byte("GET / HTTP/1.1\r\nHost: example.com\r\n\r\n"))
	}()

	buf := make([]byte, 1024)
	n, err := wrapped.Read(buf)
	if !errors.Is(err, ErrBlockedByOpenGFW) {
		t.Fatalf("err=%v, want %v", err, ErrBlockedByOpenGFW)
	}
	if n != 0 {
		t.Fatalf("n=%d, want 0 on block", n)
	}
}

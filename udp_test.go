package utils

import (
	"net"
	"testing"
	"time"
)

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String() // 本机 IP 也就是内网IP
			}
		}
	}

	return ""
}

func TestMakeUDPAddr(t *testing.T) {
	_, err := MakeUDPAddr("127.0.0.1", 9000)
	if err != nil {
		t.Errorf("MakeUDPAddr Failed: %v ", err)
	}

	_, err = MakeUDPAddr("aa.0.0.1", 9000)
	if err == nil {
		t.Errorf("MakeUDPAddr failed: %v", err)
	}
}

func TestGetIPAddr(t *testing.T) {
	ip := getLocalIP()
	if ip == "" {
		t.Errorf("Get local ip failed\n")
	}
}

func TestUDPDaemon(t *testing.T) {
	ip := getLocalIP()
	if ip == "" {
		t.Errorf("Get local ip failed\n")
	}

	recvCh := make(chan struct{})
	sendNum := 1
	count := 0
	fn := func(b []byte, from *net.UDPAddr, ts time.Time) {
		count++
		if count == sendNum {
			close(recvCh)
		}
	}

	d1, err := NewUDPListener(ip, 3001, fn)
	if err != nil {
		t.Errorf("NewUDPListener failed: %v\n", err)
	}

	d2, err := NewUDPListener(ip, 4001, fn)
	if err != nil {
		t.Errorf("NewUDPListener failed: %v\n", err)
	}

	for i := 0; i < sendNum; i++ {
		d1.SendTo([]byte("hello"), d2.conn.LocalAddr())
	}

	timeout := 10 * time.Second
	select {
	case <-time.NewTimer(timeout).C:
		t.Error("not recv packet")
	case <-recvCh:
	}
}

package utils

import (
	"fmt"
	"net"
	"time"
)

const (
	//UDPPacketSize: max udp buffer
	UDPPacketSize = 65536
)

//MakeUDPAddr return udp address
func MakeUDPAddr(ipaddr string, port int) (*net.UDPAddr, error) {
	ip := net.ParseIP(ipaddr)
	if ip == nil {
		return nil, fmt.Errorf("%v is not a valid ip address", ipaddr)
	}

	return &net.UDPAddr{ip, port, ""}, nil
}

//UDPDaemon struct
type UDPDaemon struct {
	conn    *net.UDPConn
	handler func([]byte, *net.UDPAddr, time.Time)
	running bool
}

//NewUDPListener create listener
func NewUDPListener(ipaddr string, port int, fn func([]byte, *net.UDPAddr, time.Time)) (*UDPDaemon, error) {
	udpaddr, err := MakeUDPAddr(ipaddr, port)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", udpaddr)
	if err != nil {
		return nil, err
	}

	daemon := &UDPDaemon{
		conn:    conn,
		handler: fn,
		running: true,
	}

	go daemon.Listen()
	return daemon, nil
}

func (d *UDPDaemon) Stop() {
	d.running = false
}
func (d *UDPDaemon) Listen() {
	for {
		buf := make([]byte, UDPPacketSize)
		n, clientAddr, err := d.conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}

		if !d.running {
			break
		}

		//add debug info
		d.handler(buf[:n], clientAddr, time.Now())
	}

	fmt.Printf("%v listen exit", d.conn)
}

//Send msg
func (d *UDPDaemon) SendUDP(b []byte, to *net.UDPAddr) {
	d.conn.WriteToUDP(b, to)
}

func (d *UDPDaemon) SendTo(b []byte, to net.Addr) {
	d.conn.WriteTo(b, to)
}

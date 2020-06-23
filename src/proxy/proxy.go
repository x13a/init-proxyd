package proxy

import (
	"errors"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	DefaultTimeout    = 1 << 3 * time.Second
	DefaultBufferSize = 1 << 9
)

func StartStream(fd int, dest string, timeout time.Duration) (*Stream, error) {
	file, err := newFile(fd)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	ln, err := net.FileListener(file)
	if err != nil {
		return nil, err
	}
	if timeout < 1 {
		timeout = DefaultTimeout
	}
	addr := ln.Addr()
	stream := NewStream(ln, dest, timeout)
	log.Printf("[%s] %s -> %q\n", addr.Network(), addr, dest)
	go stream.Start()
	return stream, nil
}

func StartPacket(
	fd int,
	dest string,
	timeout time.Duration,
	bufSize int,
) (*Packet, error) {
	file, err := newFile(fd)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	pc, err := net.FilePacketConn(file)
	if err != nil {
		return nil, err
	}
	if timeout < 1 {
		timeout = DefaultTimeout
	}
	if bufSize < 1 {
		bufSize = DefaultBufferSize
	}
	addr := pc.LocalAddr()
	packet := NewPacket(pc, dest, timeout, bufSize)
	log.Printf("[%s] %s -> %q\n", addr.Network(), addr, dest)
	go packet.Start()
	return packet, nil
}

func newFile(fd int) (*os.File, error) {
	file := os.NewFile(uintptr(fd), "")
	if file == nil {
		return nil, errors.New("Invalid file descriptor: " + strconv.Itoa(fd))
	}
	return file, nil
}

func splitDest(dest string) (network, address string) {
	substr := "://"
	idx := strings.Index(dest, substr)
	if idx == -1 {
		address = dest
	} else {
		network = dest[:idx]
		address = dest[idx+len(substr):]
	}
	return
}

func NewStream(ln net.Listener, dest string, timeout time.Duration) *Stream {
	network, address := splitDest(dest)
	if network == "" {
		network = ln.Addr().Network()
	}
	return &Stream{
		ln:      ln,
		network: network,
		address: address,
		timeout: timeout,
	}
}

type Stream struct {
	ln      net.Listener
	network string
	address string
	timeout time.Duration
}

func (p *Stream) Start() {
	defer p.ln.Close()
	for {
		conn, err := p.ln.Accept()
		if err != nil {
			log.Println(err)
			break
		}
		go p.proxy(conn)
	}
	addr := p.ln.Addr()
	log.Printf("[%s] %s -> exit\n", addr.Network(), addr)
}

func (p *Stream) proxy(in net.Conn) {
	defer in.Close()
	out, err := net.DialTimeout(p.network, p.address, p.timeout)
	if err != nil {
		log.Println(err)
		return
	}
	defer out.Close()
	errchan := make(chan error, 2)
	go func() {
		_, err := io.Copy(out, in)
		errchan <- err
	}()
	go func() {
		_, err := io.Copy(in, out)
		errchan <- err
	}()
	for i, n := 0, cap(errchan); i < n; i++ {
		if err = <-errchan; err != nil {
			log.Println(err)
		}
	}
}

func makePacketMap() packetMap {
	return packetMap{m: make(map[string]net.Conn)}
}

type packetMap struct {
	rw sync.RWMutex
	m  map[string]net.Conn
}

func (m *packetMap) Load(key string) (value net.Conn, ok bool) {
	m.rw.RLock()
	value, ok = m.m[key]
	m.rw.RUnlock()
	return
}

func (m *packetMap) Store(key string, value net.Conn) {
	m.rw.Lock()
	m.m[key] = value
	m.rw.Unlock()
}

func (m *packetMap) Delete(key string) {
	m.rw.Lock()
	delete(m.m, key)
	m.rw.Unlock()
}

func NewPacket(
	pc net.PacketConn,
	dest string,
	timeout time.Duration,
	bufSize int,
) *Packet {
	network, address := splitDest(dest)
	if network == "" {
		network = pc.LocalAddr().Network()
	}
	return &Packet{
		pc:      pc,
		storage: makePacketMap(),
		network: network,
		address: address,
		timeout: timeout,
		bufSize: bufSize,
	}
}

type Packet struct {
	pc      net.PacketConn
	storage packetMap
	network string
	address string
	timeout time.Duration
	bufSize int
}

func (p *Packet) Start() {
	defer p.pc.Close()
	for {
		buf := make([]byte, p.bufSize)
		n, addr, err := p.pc.ReadFrom(buf)
		if n > 0 {
			go p.handle(buf[:n], addr)
		}
		if err != nil {
			log.Println(err)
			break
		}
	}
	addr := p.pc.LocalAddr()
	log.Printf("[%s] %s -> exit\n", addr.Network(), addr)
}

func (p *Packet) nextDeadline() time.Time {
	return time.Now().Add(p.timeout)
}

func (p *Packet) handle(data []byte, addr net.Addr) {
	var err error
	addrstr := addr.String()
	out, ok := p.storage.Load(addrstr)
	if !ok {
		out, err = p.dial()
		if err != nil {
			log.Println(err)
			return
		}
		p.storage.Store(addrstr, out)
		go p.proxy(out, addr)
	}
	if _, err = out.Write(data); err != nil {
		log.Println(err)
		return
	}
	if err = out.SetReadDeadline(p.nextDeadline()); err != nil {
		log.Println(err)
	}
}

func (p *Packet) dial() (net.Conn, error) {
	out, err := net.DialTimeout(p.network, p.address, p.timeout)
	if err != nil {
		return nil, err
	}
	if err = out.SetReadDeadline(p.nextDeadline()); err != nil {
		out.Close()
		return nil, err
	}
	return out, nil
}

func (p *Packet) proxy(out net.Conn, addr net.Addr) {
	defer func() {
		p.storage.Delete(addr.String())
		out.Close()
	}()
	buf := make([]byte, p.bufSize)
	for {
		n, err := out.Read(buf)
		if n > 0 {
			if _, err1 := p.pc.WriteTo(buf[:n], addr); err == nil {
				err = err1
			}
		}
		if err != nil {
			if !os.IsTimeout(err) {
				log.Println(err)
			}
			break
		}
		if err = out.SetReadDeadline(p.nextDeadline()); err != nil {
			log.Println(err)
		}
	}
}

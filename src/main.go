package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"syscall"
	"time"

	"github.com/x13a/go-proxy"

	"github.com/x13a/init-proxyd/launchd"
	"github.com/x13a/init-proxyd/systemd"
)

const (
	Version = "0.2.0"

	FlagConfig      = "c"
	FlagDestination = "d"
	FlagTimeout     = "t"
	FlagBufferSize  = "b"

	ExitSuccess = 0
	ExitUsage   = 2
)

func start(fd int, opts *Opts) error {
	syscall.CloseOnExec(fd)
	socketType, err := syscall.GetsockoptInt(
		fd,
		syscall.SOL_SOCKET,
		syscall.SO_TYPE,
	)
	if err != nil {
		return err
	}
	var prx proxy.Proxy
	switch socketType {
	case syscall.SOCK_STREAM:
		prx, err = proxy.NewFileStreamProxy(fd, opts.dest, opts.timeout)
	case syscall.SOCK_DGRAM:
		prx, err = proxy.NewFilePacketProxy(
			fd,
			opts.dest,
			opts.timeout,
			opts.bufSize,
		)
	default:
		return errors.New("Unsupported socket type: " +
			strconv.Itoa(socketType))
	}
	if err != nil {
		return err
	}
	go func(prx proxy.Proxy) {
		prx.Start()
		<-prx.WaitChan()
		log.Fatalln("Proxy exited")
	}(prx)
	return nil
}

type Opts struct {
	config  string
	dest    string
	timeout time.Duration
	bufSize int
}

func getOpts() *Opts {
	opts := &Opts{}
	isHelp := flag.Bool("h", false, "Print help and exit")
	isVersion := flag.Bool("V", false, "Print version and exit")
	if launchd.Is() {
		flag.StringVar(
			&opts.config,
			FlagConfig,
			launchd.DefaultConfig,
			"Path to config file",
		)
	}
	flag.StringVar(&opts.dest, FlagDestination, "", "Destination address")
	flag.DurationVar(
		&opts.timeout,
		FlagTimeout,
		proxy.DefaultTimeout,
		"Timeout",
	)
	flag.IntVar(
		&opts.bufSize,
		FlagBufferSize,
		proxy.DefaultBufferSize,
		"UDP buffer size",
	)
	flag.Parse()
	if *isHelp {
		flag.Usage()
		os.Exit(ExitSuccess)
	}
	if *isVersion {
		fmt.Println(Version)
		os.Exit(ExitSuccess)
	}
	if opts.dest == "" {
		os.Stderr.WriteString("Empty destination\n")
		os.Exit(ExitUsage)
	}
	return opts
}

func main() {
	opts := getOpts()
	var fds []int
	var err error
	switch {
	case launchd.Is():
		fds, err = launchd.Sockets(opts.config)
	case systemd.Is():
		fds, err = systemd.Sockets()
	default:
		log.Fatalln("Unsupported init system")
	}
	if err != nil {
		log.Fatalln(err)
	}
	if len(fds) == 0 {
		log.Fatalln("Sockets not found")
	}
	for _, fd := range fds {
		if err = start(fd, opts); err != nil {
			log.Fatalln(err)
		}
	}
	select {}
}

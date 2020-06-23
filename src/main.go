package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"./launch"
	"./proxy"
)

const (
	Version = "0.0.1"

	FlagNames       = "n"
	FlagDestination = "d"
	FlagTimeout     = "t"
	FlagBufferSize  = "b"

	TCP  = "tcp"
	UDP  = "udp"
	TCP6 = TCP + "6"
	UDP6 = UDP + "6"

	ExitSuccess = 0
	ExitUsage   = 2
)

func start(name string, noHost bool, opts *Opts) error {
	fds, err := launch.ActivateSocket(name)
	if err != nil {
		return err
	}
	dest := opts.dest
	if noHost {
		switch name {
		case TCP6, UDP6:
			dest = "[::1]" + dest
		default:
			dest = "127.0.0.1" + dest
		}
	}
	name = name[:min(3, len(name))]
	for _, fd := range fds {
		switch name {
		case TCP:
			_, err = proxy.StartStream(fd, dest, opts.timeout)
		case UDP:
			_, err = proxy.StartPacket(fd, dest, opts.timeout, opts.bufSize)
		default:
			return fmt.Errorf("Invalid name: %q", name)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func min(v ...int) int {
	res := v[0]
	for _, i := range v[1:] {
		if i < res {
			res = i
		}
	}
	return res
}

type Opts struct {
	names   []string
	dest    string
	timeout time.Duration
	bufSize int
}

func getOpts() *Opts {
	opts := &Opts{}
	isHelp := flag.Bool("h", false, "Print help and exit")
	isVersion := flag.Bool("V", false, "Print version and exit")
	names := flag.String(
		FlagNames,
		TCP+","+UDP,
		fmt.Sprintf("Comma separated socket names { %s | %s }", TCP, UDP),
	)
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
	opts.names = strings.Split(*names, ",")
	if opts.dest == "" {
		os.Stderr.WriteString("Empty destination\n")
		os.Exit(ExitUsage)
	}
	return opts
}

func main() {
	opts := getOpts()
	noHost := false
	if host, _, err := net.SplitHostPort(opts.dest); err == nil && host == "" {
		noHost = true
	}
	for _, name := range opts.names {
		if err := start(strings.TrimSpace(name), noHost, opts); err != nil {
			log.Fatalln(err)
		}
	}
	select {}
}

package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"./launch"
	"./proxy"
)

const (
	Version = "0.1.2"

	FlagPlist       = "p"
	FlagDestination = "d"
	FlagTimeout     = "t"
	FlagBufferSize  = "b"

	ExitSuccess = 0
	ExitUsage   = 2

	DefaultPlist = "/Library/LaunchDaemons/me.lucky.launch-proxy.plist"
)

func start(name string, opts *Opts) error {
	fds, err := launch.ActivateSocket(name)
	if err != nil {
		return err
	}
	name = strings.ToLower(name)
	isStream := strings.HasPrefix(name, proxy.TCP)
	isPacket := strings.HasPrefix(name, proxy.UDP)
	var prx proxy.Proxy
	for _, fd := range fds {
		switch {
		case isStream:
			prx, err = proxy.NewFileStreamProxy(fd, opts.dest, opts.timeout)
		case isPacket:
			prx, err = proxy.NewFilePacketProxy(
				fd,
				opts.dest,
				opts.timeout,
				opts.bufSize,
			)
		default:
			return fmt.Errorf("Invalid name: %q", name)
		}
		if err != nil {
			return err
		}
		go func(prx proxy.Proxy) {
			prx.Start()
			<-prx.WaitChan()
			log.Fatalln("Proxy exited")
		}(prx)
	}
	return nil
}

func parsePlist(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	dec := xml.NewDecoder(file)
	const (
		Key   = "key"
		Array = "array"
		Dict  = "dict"
	)
	depth := -1
	found := false
	res := []string{}
Loop:
	for {
		token, err := dec.Token()
		if err != nil {
			return nil, err
		}
		switch element := token.(type) {
		case xml.StartElement:
			switch element.Name.Local {
			case Key:
				if (!found && depth != 0) || (found && depth != 1) {
					continue
				}
				var key string
				if err = dec.DecodeElement(&key, &element); err != nil {
					log.Println(err)
					continue
				}
				if !found {
					if key == "Sockets" {
						found = true
					}
				} else {
					res = append(res, key)
				}
			case Array, Dict:
				depth++
			}
		case xml.EndElement:
			switch element.Name.Local {
			case Array:
				depth--
			case Dict:
				depth--
				if found && depth == 0 {
					break Loop
				}
			}
		}
	}
	return res, nil
}

type Opts struct {
	plist   string
	dest    string
	timeout time.Duration
	bufSize int
}

func getOpts() *Opts {
	opts := &Opts{}
	isHelp := flag.Bool("h", false, "Print help and exit")
	isVersion := flag.Bool("V", false, "Print version and exit")
	flag.StringVar(&opts.plist, FlagPlist, DefaultPlist, "Path to plist file")
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
	names, err := parsePlist(opts.plist)
	if err != nil {
		log.Fatalln(err)
	}
	if len(names) == 0 {
		log.Fatalln("Socket names not found in", opts.plist)
	}
	log.Println("Found socket names:", names)
	for _, name := range names {
		if err = start(name, opts); err != nil {
			log.Fatalln(err)
		}
	}
	select {}
}

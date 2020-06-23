# launch-proxy

Starting from version 2.0.43, `dnscrypt-proxy` can't drop privileges on macOS.
To bind to 127.0.0.1:53 you need to use sudo. Where are two options: use, 
for example, dnsmasq or use launch socket activation. It's last.

Issue: [1371](https://github.com/DNSCrypt/dnscrypt-proxy/issues/1371)  
Issue: [1367](https://github.com/DNSCrypt/dnscrypt-proxy/issues/1367)  

History: [launch_socket_server](https://github.com/sstephenson/launch_socket_server)

## Installation
```sh
$ make
$ sudo make install
```
or
```sh
$ brew tap x31a/tap https://bitbucket.org/x31a/homebrew-tap.git
$ brew install x31a/tap/launch-proxy
```

## Usage
```text
Usage of launch-proxy:
  -V	Print version and exit
  -b int
    	UDP buffer size (default 512)
  -d string
    	Destination address
  -h	Print help and exit
  -n string
    	Comma separated socket names { tcp | udp } (default "tcp,udp")
  -t duration
    	Timeout (default 8s)
```

## Example

To load:
```sh
$ sudo launchctl load /Library/LaunchDaemons/me.lucky.launch-proxy.plist
```

To unload:
```sh
$ sudo launchctl unload /Library/LaunchDaemons/me.lucky.launch-proxy.plist
```

The same using Makefile (will copy and remove plist for you):
```sh
$ sudo make load
$ sudo make unload
```

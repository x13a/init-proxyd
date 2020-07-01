# launch-proxy

Launch daemon sockets proxy for macOS.

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
  -p string
    	Path to plist file (default "/Library/LaunchDaemons/me.lucky.launch-proxy.plist")
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

## Caveats

Proxy may not load under `nobody:nogroup`, then you need to change:
```xml
<key>UserName</key>
<string>nobody</string>
<key>GroupName</key>
<string>nogroup</string>
```

## Friends
- [launch_socket_server](https://github.com/sstephenson/launch_socket_server)

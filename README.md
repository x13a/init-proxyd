# init-proxyd

Launch daemon sockets proxy for macOS.

## Installation
```sh
$ make
$ sudo make install
```
or
```sh
$ brew tap x31a/tap https://bitbucket.org/x31a/homebrew-tap.git
$ brew install x31a/tap/init-proxyd
```

## Usage
```text
Usage of init-proxyd:
  -V	Print version and exit
  -b int
    	UDP buffer size (default 512)
  -c string (launchd only)
    	Path to config file (default "/Library/LaunchDaemons/me.lucky.init-proxyd.plist")
  -d string
    	Destination address
  -h	Print help and exit
  -t duration
    	Timeout (default 8s)
```

## Example

To load:
```sh
$ sudo launchctl load /Library/LaunchDaemons/me.lucky.init-proxyd.plist
```

To unload:
```sh
$ sudo launchctl unload /Library/LaunchDaemons/me.lucky.init-proxyd.plist
```

The same using Makefile (will copy and remove config for you):
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

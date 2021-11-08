# init-proxyd

Init daemons socket activation proxy.

Supported init systems:
- launchd
- systemd

## Installation
```sh
$ make
$ sudo make install
```
or
```sh
$ brew tap x13a/tap
$ brew install x13a/tap/init-proxyd
```

## Usage
```text
Usage of init-proxyd:
  -V	Print version and exit
  -b int
    	UDP buffer size (default 512)
  -c string (darwin only)
    	Path to config file (default "/Library/LaunchDaemons/me.lucky.init-proxyd.plist")
  -d string
    	Destination address
  -h	Print help and exit
  -t duration
    	Timeout (default 8s)
```

## Example

To load (macOS):
```sh
$ sudo install -m 0644 ./config/launchd/me.lucky.init-proxyd.plist /Library/LaunchDaemons/
$ sudo launchctl load /Library/LaunchDaemons/me.lucky.init-proxyd.plist
```

To unload (macOS):
```sh
$ sudo launchctl unload /Library/LaunchDaemons/me.lucky.init-proxyd.plist
$ sudo rm -f /Library/LaunchDaemons/me.lucky.init-proxyd.plist
```

The same using Makefile (macOS only):
```sh
$ sudo make load
$ sudo make unload
```

## Caveats

Proxy may not load under `nobody:nogroup`, then you should change:
```xml
<key>UserName</key>
<string>nobody</string>
<key>GroupName</key>
<string>nogroup</string>
```

## Friends
- [launch_socket_server](https://github.com/sstephenson/launch_socket_server)

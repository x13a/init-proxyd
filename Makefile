NAME        := launch-proxy
PLISTNAME   := me.lucky.launch-proxy.plist
ADMINUID    := 501

prefix      ?= /usr/local
exec_prefix ?= $(prefix)
sbindir     ?= $(exec_prefix)/sbin
datarootdir ?= $(prefix)/share
datadir     ?= $(datarootdir)/$(NAME)
srcdir      ?= ./src

targetdir   := ./target
target      := $(targetdir)/$(NAME)
destdir     := $(DESTDIR)$(sbindir)
dest        := $(destdir)/$(NAME)
plistfile   := ./plist/$(PLISTNAME)
launchdir   := /Library/LaunchDaemons
launchdest  := $(launchdir)/$(PLISTNAME)

all: build

build:
	go build -o $(target) $(srcdir)/

installdirs:
	install -o $(ADMINUID) -g staff -d $(destdir)/
	install -o $(ADMINUID) -g staff -d $(datadir)/

install: installdirs
	install -o root -g wheel -f uchg $(target) $(destdir)/
	install -m 0644 -o $(ADMINUID) -g staff $(plistfile) $(datadir)/

uninstall:
	chflags nouchg $(dest)
	rm -f $(dest)
	rm -rf $(datadir)/

load:
	install -m 0644 -o root -g wheel $(plistfile) $(launchdir)/
	launchctl load $(launchdest)

unload:
	launchctl unload $(launchdest)
	rm -f $(launchdest)

clean:
	rm -rf $(targetdir)/

NAME        := launch-proxy
ADMINUID    := 501
LAUNCHDIR   := /Library/LaunchDaemons

prefix      ?= /usr/local
exec_prefix ?= $(prefix)
sbindir     ?= $(exec_prefix)/sbin
datarootdir ?= $(prefix)/share
datadir     ?= $(datarootdir)
srcdir      ?= ./src

plistname   := me.lucky.launch-proxy.plist
targetdir   := ./target
target      := $(targetdir)/$(NAME)
sbindestdir := $(DESTDIR)$(sbindir)
datadestdir := $(DESTDIR)$(datadir)/$(NAME)
sbindest    := $(sbindestdir)/$(NAME)
plist       := ./plist/$(plistname)
launchdest  := $(LAUNCHDIR)/$(plistname)

all: build

build:
	go build -o $(target) $(srcdir)/

installdirs:
	install -o $(ADMINUID) -g staff -d $(sbindestdir)/ $(datadestdir)/

install: installdirs
	install -o root -g wheel -f uchg $(target) $(sbindestdir)/
	install -m 0644 -o $(ADMINUID) -g staff $(plist) $(datadestdir)/

uninstall:
	chflags nouchg $(sbindest)
	rm -f $(sbindest)
	rm -rf $(datadestdir)/

load:
	install -m 0644 -o root -g wheel $(plist) $(LAUNCHDIR)/
	launchctl load $(launchdest)

unload:
	launchctl unload $(launchdest)
	rm -f $(launchdest)

clean:
	rm -rf $(targetdir)/

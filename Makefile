NAME        := init-proxyd
ADMINUID    := 501
CONFIGDIR   := /Library/LaunchDaemons
CONFIGNAME  := me.lucky.$(NAME).plist

prefix      ?= /usr/local
exec_prefix ?= $(prefix)
sbindir     ?= $(exec_prefix)/sbin
datarootdir ?= $(prefix)/share
datadir     ?= $(datarootdir)
srcdir      ?= ./src

targetdir   := ./target
target      := $(targetdir)/$(NAME)
sbindestdir := $(DESTDIR)$(sbindir)
datadestdir := $(DESTDIR)$(datadir)/$(NAME)
sbindest    := $(sbindestdir)/$(NAME)
config      := ./config/$(CONFIGNAME)
configdest  := $(CONFIGDIR)/$(CONFIGNAME)

all: build

build:
	go build -o $(target) $(srcdir)/

installdirs:
	install -o $(ADMINUID) -g staff -d $(sbindestdir)/ $(datadestdir)/

install: installdirs
	install -o root -g wheel -f uchg $(target) $(sbindestdir)/
	install -m 0644 -o $(ADMINUID) -g staff -b $(config) $(datadestdir)/

uninstall:
	chflags nouchg $(sbindest)
	rm -f $(sbindest)
	rm -rf $(datadestdir)/

load:
	install -m 0644 -o root -g wheel $(config) $(CONFIGDIR)/
	launchctl load $(configdest)

unload:
	launchctl unload $(configdest)
	rm -f $(configdest)

clean:
	rm -rf $(targetdir)/

NAME        := init-proxyd

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

config      := ./config/launchd/$(CONFIGNAME)
configdest  := $(CONFIGDIR)/$(CONFIGNAME)

all: build

build:
	# ugly fix :(
	(cd $(srcdir); go build -o ../$(target) ".")

installdirs:
	install -d $(sbindestdir)/ $(datadestdir)/

install: installdirs
	install $(target) $(sbindestdir)/
	install -m 0644 -b $(config) $(datadestdir)/

uninstall:
	rm -f $(sbindest)
	rm -rf $(datadestdir)/

load:
	install -m 0644 $(config) $(CONFIGDIR)/
	launchctl load $(configdest)

unload:
	launchctl unload $(configdest)
	rm -f $(configdest)

clean:
	rm -rf $(targetdir)/

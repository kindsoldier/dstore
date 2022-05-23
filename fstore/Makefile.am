
AUTOMAKE_OPTIONS = foreign no-dependencies no-installinfo

SUFFIXES = .go
OBJEXT= none

sbin_PROGRAMS = fstored fstorecli
fstored_SOURCES = fssrv/fsmain.go

EXTRA_fstored_SOURCES = \
	fssrv/fsconf/fsconf.go \
	fssrv/fsconf/fsconf.go.in \
	fssrv/fscont/fscont.go

fstorecli_SOURCES = fscli/fscli.go


GOFLAGS = -ldflags="-s -w"

fstored$(EXEEXT): $(fstored_SOURCES) $(EXTRA_fstored_SOURCES)
	$(GO) build $(GOFLAGS) -o fstored$(EXEEXT) $(fstored_SOURCES)

fstorecli$(EXEEXT): $(fstorecli_SOURCES) $(EXTRA_fstorecli_SOURCES)
	$(GO) build $(GOFLAGS) -o fstorecli$(EXEEXT) $(fstorecli_SOURCES)

EXTRA_DIST = \
	fstore.conf

clean-local:
	rm -rf autom4te.cache
	rm -rf log/ run/ data/

install-data-local:
	test -z $(DESTDIR)$(SRV_CONFDIR) || $(MKDIR_P) $(DESTDIR)$(SRV_CONFDIR)
	test -z $(DESTDIR)$(SRV_LOGDIR) || $(MKDIR_P) $(DESTDIR)$(SRV_LOGDIR)
	test -z $(DESTDIR)$(SRV_RUNDIR) || $(MKDIR_P) $(DESTDIR)$(SRV_RUNDIR)
	test -z $(DESTDIR)$(SRV_DATADIR) || $(MKDIR_P) $(DESTDIR)$(SRV_DATADIR)

PATH := $(PATH):$(HOME)/go/bin

instvis:
	$(GO) get github.com/ofabry/go-callvis
	$(GO) install github.com/ofabry/go-callvis

CVCMD = go-callvis -nostd  -skipbrowser -graphviz
CVIGNORE = nsstore/dslog
CVDIR = .
GOMAIN = ndstore/fstore/fssrv
GOROOT = ndstore

graph:
	$(CVCMD) -ignore $(CVIGNORE) -focus $(GOROOT)/fstore/ndsrv -file $(CVDIR)/fstore $(GOMAIN)
	$(CVCMD) -ignore $(CVIGNORE) -focus $(GOROOT)/dsrpc -file $(CVDIR)/dcrpc $(GOMAIN)

EXTRA_DIST += \
	fstore.svg \
	dsrpc.svg
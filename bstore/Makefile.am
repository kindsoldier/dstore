
AUTOMAKE_OPTIONS = foreign no-dependencies no-installinfo

SUFFIXES = .go
OBJEXT= none

sbin_PROGRAMS = bstored bstorecli
bstored_SOURCES = bssrv/bsmain.go

EXTRA_bstored_SOURCES = \
	bssrv/bsconf/bsconf.go \
	bssrv/bsconf/bsconf.go.in \
	bssrv/bsrec/bsrec.go \
	bssrv/bscont/bscont.go \
	bssrv/bsreg/bsreg.go

bstorecli_SOURCES = bscli/bscli.go


GOFLAGS = -ldflags="-s -w"

bstored$(EXEEXT): $(bstored_SOURCES) $(EXTRA_bstored_SOURCES)
	$(GO) build $(GOFLAGS) -o bstored$(EXEEXT) $(bstored_SOURCES)

bstorecli$(EXEEXT): $(bstorecli_SOURCES) $(EXTRA_bstorecli_SOURCES)
	$(GO) build $(GOFLAGS) -o bstorecli$(EXEEXT) $(bstorecli_SOURCES)

EXTRA_DIST = \
	bstore.conf

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
GOMAIN = ndstore/bstore/bssrv
GOROOT = ndstore

graph:
	$(CVCMD) -ignore $(CVIGNORE) -focus $(GOROOT)/bstore/ndsrv -file $(CVDIR)/bstore $(GOMAIN)
	$(CVCMD) -ignore $(CVIGNORE) -focus $(GOROOT)/dsrpc -file $(CVDIR)/dcrpc $(GOMAIN)

EXTRA_DIST += \
	bstore.svg \
	dsrpc.svg
